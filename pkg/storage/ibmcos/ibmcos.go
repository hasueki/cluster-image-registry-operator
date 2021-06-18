package ibmcos

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	configapiv1 "github.com/openshift/api/config/v1"
	imageregistryv1 "github.com/openshift/api/imageregistry/v1"
	operatorapi "github.com/openshift/api/operator/v1"

	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
	"github.com/openshift/cluster-image-registry-operator/pkg/defaults"
	"github.com/openshift/cluster-image-registry-operator/pkg/envvar"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/util"
)

const (
	IAMEndpoint                   = "https://iam.cloud.ibm.com/identity/token"
	imageRegistrySecretDataKey    = "ibmcloud_api_key"
	imageRegistrySecretMountpoint = "/var/run/secrets/cloud"
)

type driver struct {
	Context context.Context
	Config  *imageregistryv1.ImageRegistryConfigStorageIBMCOS
	Listers *regopclient.Listers

	// httpClient is used only during tests.
	// httpClient *http.Client
}

// NewDriver creates a new IBM COS storage driver.
// Used during bootstrapping.
func NewDriver(ctx context.Context, c *imageregistryv1.ImageRegistryConfigStorageIBMCOS, listers *regopclient.Listers) *driver {
	return &driver{
		Context: ctx,
		Config:  c,
		Listers: listers,
	}
}

// ConfigEnv configures the environment variables that will be
// used in the image registry deployment.
func (d *driver) ConfigEnv() (envs envvar.List, err error) {
	envs = append(envs,
		envvar.EnvVar{Name: "REGISTRY_STORAGE", Value: "ibmcos"},
		envvar.EnvVar{Name: "REGISTRY_STORAGE_IBMCOS_BUCKET", Value: d.Config.Bucket},
		envvar.EnvVar{Name: "REGISTRY_STORAGE_IBMCOS_LOCATION", Value: d.Config.Location},
		envvar.EnvVar{Name: "REGISTRY_STORAGE_IBMCOS_RESOURCEGROUPNAME", Value: d.Config.ResourceGroupName},
		envvar.EnvVar{Name: "REGISTRY_STORAGE_IBMCOS_SERVICEINSTANCECRN", Value: d.Config.ServiceInstanceCRN},
		envvar.EnvVar{Name: "REGISTRY_STORAGE_IBMCOS_CREDENTIALSCONFIGPATH", Value: filepath.Join(imageRegistrySecretMountpoint, imageRegistrySecretDataKey)},
	)
	return
}

// CreateStorage attempts to create an IBM COS service instance and bucket.
func (d *driver) CreateStorage(cr *imageregistryv1.Config) error {
	// Get Infrastructure spec
	infra, err := util.GetInfrastructure(d.Listers)
	if err != nil {
		return err
	}

	// Set configs from Infrastructure
	d.Config.Location = infra.Status.PlatformStatus.IBMCloud.Location
	d.Config.ResourceGroupName = infra.Status.PlatformStatus.IBMCloud.ResourceGroupName

	// Get resource controller service
	rc, err := d.getResouceControllerService()
	if err != nil {
		return err
	}

	// Get resource manager service
	rm, err := d.getResourceManagerService()
	if err != nil {
		return err
	}

	// Check if service instance exists
	if len(d.Config.ServiceInstanceCRN) != 0 {
		instance, resp, err := rc.GetResourceInstanceWithContext(
			d.Context,
			&resourcecontrollerv2.GetResourceInstanceOptions{
				ID: &d.Config.ServiceInstanceCRN,
			},
		)
		if err != nil {
			return fmt.Errorf("unable to get resource instance: %s with resp code: %d", err.Error(), resp.StatusCode)
		}

		switch *instance.State {
		case resourcecontrollerv2.ListResourceInstancesOptionsStateActiveConst:
			// Service instance exists and is active
			if *instance.ResourceGroupID != "" {
				// Get resource group name
				rg, resp, err := rm.GetResourceGroupWithContext(
					d.Context,
					&resourcemanagerv2.GetResourceGroupOptions{
						ID: instance.ResourceGroupID,
					},
				)
				if err != nil {
					return fmt.Errorf("unable to get resource group: %s with resp code: %d", err.Error(), resp.StatusCode)
				}
				// Set resource group name
				d.Config.ResourceGroupName = *rg.Name
			}
			cr.Status.Storage = imageregistryv1.ImageRegistryConfigStorage{
				IBMCOS: d.Config.DeepCopy(),
			}
			cr.Spec.Storage.IBMCOS = d.Config.DeepCopy()
			util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, "IBM COS Instance Active", "IBM COS service instance is active")
		case resourcecontrollerv2.ListResourceInstancesOptionsStateProvisioningConst:
			// Service instance exists and is provisioning
			util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, "IBM COS Instance Provisioning", "IBM COS service instance is provisioning")
			return fmt.Errorf("waiting for IBM COS service instance to finish provisioning")
		default:
			// Service instance does not exist
			d.Config.ServiceInstanceCRN = ""
			util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, "IBM COS Instance Gone", "IBM COS service instance is inactive or has been removed.")
		}
	}

	// Attempt to create a new service instance
	if len(d.Config.ServiceInstanceCRN) == 0 {
		// Get resource group details
		resourceGroups, resp, err := rm.ListResourceGroupsWithContext(
			d.Context,
			&resourcemanagerv2.ListResourceGroupsOptions{
				Name: &d.Config.ResourceGroupName,
			},
		)
		if len(resourceGroups.Resources) == 0 || err != nil {
			return fmt.Errorf("unable to get resource groups: %s with resp code: %d", err.Error(), resp.StatusCode)
		}

		// Define instance options
		serviceInstanceName := fmt.Sprintf("%s-%s", infra.Status.InfrastructureName, defaults.ImageRegistryName)
		serviceTarget := "bluemix-global"
		resourceGroupID := *resourceGroups.Resources[0].ID
		resourcePlanID := "744bfc56-d12c-4866-88d5-dac9139e0e5d"

		// Check if service instance with name already exists
		instances, resp, err := rc.ListResourceInstancesWithContext(
			d.Context,
			&resourcecontrollerv2.ListResourceInstancesOptions{
				Name:            &serviceInstanceName,
				ResourceGroupID: &resourceGroupID,
				ResourcePlanID:  &resourcePlanID,
			},
		)
		if instances == nil || err != nil {
			return fmt.Errorf("unable to get resource instances: %s with resp code: %d", err.Error(), resp.StatusCode)
		}

		var instance *resourcecontrollerv2.ResourceInstance
		if len(instances.Resources) != 0 {
			// Service instance found
			instance = &instances.Resources[0]
		} else {
			// Create COS service instance
			instance, resp, err = rc.CreateResourceInstanceWithContext(
				d.Context,
				&resourcecontrollerv2.CreateResourceInstanceOptions{
					Name:           &serviceInstanceName,
					Target:         &serviceTarget,
					ResourceGroup:  &resourceGroupID,
					ResourcePlanID: &resourcePlanID,
					Tags:           []string{fmt.Sprintf("kubernetes.io_cluster_%s:owned", infra.Status.InfrastructureName)},
				},
			)
			if instance == nil || err != nil {
				return fmt.Errorf("unable to create resource instance: %s with resp code: %d", err.Error(), resp.StatusCode)
			}

			if cr.Spec.Storage.ManagementState == "" {
				cr.Spec.Storage.ManagementState = imageregistryv1.StorageManagementStateManaged
			}
		}

		d.Config.ServiceInstanceCRN = *instance.CRN
		cr.Status.Storage = imageregistryv1.ImageRegistryConfigStorage{
			IBMCOS: d.Config.DeepCopy(),
		}
		cr.Spec.Storage.IBMCOS = d.Config.DeepCopy()
		util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, "IBM COS Instance Creation Successful", "IBM COS service instance was successfully created")
	}

	// Check if bucket already exists
	var bucketExists bool
	if len(d.Config.Bucket) != 0 {
		if err := d.bucketExists(d.Config.Bucket, d.Config.ServiceInstanceCRN); err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket, "Forbidden", "NotFound":
					// If the bucket doesn't exist that's ok, we'll try to create it
					util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, aerr.Code(), aerr.Error())
				default:
					util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionUnknown, "Unknown Error Occurred", err.Error())
					return err
				}
			} else {
				util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionUnknown, "Unknown Error Occurred", err.Error())
				return err
			}
		} else {
			bucketExists = true
		}
	}

	// Create new bucket if required
	if len(d.Config.Bucket) != 0 && bucketExists {
		if cr.Spec.Storage.ManagementState == "" {
			cr.Spec.Storage.ManagementState = imageregistryv1.StorageManagementStateUnmanaged
		}

		cr.Status.Storage = imageregistryv1.ImageRegistryConfigStorage{
			IBMCOS: d.Config.DeepCopy(),
		}
		util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionTrue, "IBM COS Bucket Exists", "User supplied IBM COS bucket exists and is accessible")
	} else {
		// Attempt to create new bucket
		if len(d.Config.Bucket) == 0 {
			if d.Config.Bucket, err = util.GenerateStorageName(d.Listers, d.Config.Location); err != nil {
				return err
			}
		}

		// Get COS client
		client, err := d.getIBMCOSClient(d.Config.ServiceInstanceCRN)
		if err != nil {
			return err
		}

		// Create COS bucket
		_, err = client.CreateBucketWithContext(
			d.Context,
			&s3.CreateBucketInput{
				Bucket: aws.String(d.Config.Bucket),
				CreateBucketConfiguration: &s3.CreateBucketConfiguration{
					LocationConstraint: aws.String(fmt.Sprintf("%s-smart", d.Config.Location)),
				},
			},
		)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, aerr.Code(), aerr.Error())
			}
			return err
		}

		if cr.Spec.Storage.ManagementState == "" {
			cr.Spec.Storage.ManagementState = imageregistryv1.StorageManagementStateManaged
		}
		cr.Status.Storage = imageregistryv1.ImageRegistryConfigStorage{
			IBMCOS: d.Config.DeepCopy(),
		}
		cr.Spec.Storage.IBMCOS = d.Config.DeepCopy()
		util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionTrue, "Creation Successful", "IBM COS bucket was successfully created")

		// Wait until the bucket exists
		if err := client.WaitUntilBucketExistsWithContext(
			d.Context,
			&s3.HeadBucketInput{
				Bucket: aws.String(d.Config.Bucket),
			},
		); err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, aerr.Code(), aerr.Error())
			}
			return err
		}
	}

	return nil
}

func (d *driver) getResouceControllerService() (*resourcecontrollerv2.ResourceControllerV2, error) {
	IAMAPIKey, err := d.getCredentialsConfigData()
	if err != nil {
		return nil, err
	}

	service, err := resourcecontrollerv2.NewResourceControllerV2(
		&resourcecontrollerv2.ResourceControllerV2Options{
			Authenticator: &core.IamAuthenticator{
				ApiKey: IAMAPIKey,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (d *driver) getResourceManagerService() (*resourcemanagerv2.ResourceManagerV2, error) {
	IAMAPIKey, err := d.getCredentialsConfigData()
	if err != nil {
		return nil, err
	}

	service, err := resourcemanagerv2.NewResourceManagerV2(
		&resourcemanagerv2.ResourceManagerV2Options{
			Authenticator: &core.IamAuthenticator{
				ApiKey: IAMAPIKey,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// ID return the underlying storage identificator, in this case the bucket name.
func (d *driver) ID() string {
	fmt.Println("[WIP] ibmcos.ID")
	return d.Config.Bucket
}

// RemoveStorage deletes the storage medium that was created.
// The COS bucket must be empty before it can be removed
func (d *driver) RemoveStorage(cr *imageregistryv1.Config) (bool, error) {
	fmt.Println("[WIP] ibmcos.RemoveStorage")
	return false, nil
}

// StorageChanged checks to see if the name of the storage medium
// has changed
func (d *driver) StorageChanged(cr *imageregistryv1.Config) bool {
	if !reflect.DeepEqual(cr.Status.Storage.IBMCOS, cr.Spec.Storage.IBMCOS) {
		util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionUnknown, "IBMCOS Configuration Changed", "IBMCOS storage is in an unknown state")
		return true
	}
	return false
}

// StorageExists checks if an IBM COS bucket with the given name exists
// and we can access it
func (d *driver) StorageExists(cr *imageregistryv1.Config) (bool, error) {
	if len(d.Config.Bucket) == 0 || len(d.Config.ServiceInstanceCRN) == 0 {
		return false, nil
	}

	err := d.bucketExists(d.Config.Bucket, d.Config.ServiceInstanceCRN)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket, "Forbidden", "NotFound":
				util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionFalse, aerr.Code(), aerr.Error())
				return false, nil
			}
		}
		util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionUnknown, "Unknown Error Occurred", err.Error())
		return false, err
	}

	util.UpdateCondition(cr, defaults.StorageExists, operatorapi.ConditionTrue, "IBM COS Bucket Exists", "")
	return true, nil
}

// bucketExists checks whether or not the IBM COS bucket exists
func (d *driver) bucketExists(bucketName string, serviceInstanceCRN string) error {
	client, err := d.getIBMCOSClient(serviceInstanceCRN)
	if err != nil {
		return err
	}

	_, err = client.HeadBucketWithContext(
		d.Context,
		&s3.HeadBucketInput{
			Bucket: &bucketName,
		},
	)

	return err
}

// getIBMCOSClient returns a client that allows us to interact
// with the IBM COS service
func (d *driver) getIBMCOSClient(serviceInstanceCRN string) (*s3.S3, error) {

	infra, err := util.GetInfrastructure(d.Listers)
	if err != nil {
		return nil, err
	}

	IBMCOSLocation := imageregistryv1.ImageRegistryConfigStorageIBMCOS{}.Location
	if infra.Status.PlatformStatus != nil && infra.Status.PlatformStatus.Type == configapiv1.IBMCloudPlatformType {
		IBMCOSLocation = infra.Status.PlatformStatus.IBMCloud.Location
	}

	if IBMCOSLocation == "" {
		return nil, fmt.Errorf("unable to get location from infrastructure")
	}

	serviceEndpoint := fmt.Sprintf("s3.%s.cloud-object-storage.appdomain.cloud", IBMCOSLocation)
	IAMAPIKey, err := d.getCredentialsConfigData()
	if err != nil {
		return nil, err
	}

	conf := aws.NewConfig().
		WithEndpoint(serviceEndpoint).
		WithCredentials(ibmiam.NewStaticCredentials(aws.NewConfig(), IAMEndpoint, IAMAPIKey, serviceInstanceCRN)).
		WithS3ForcePathStyle(true)

	sess := session.Must(session.NewSession())

	return s3.New(sess, conf), nil
}

// getCredentialsConfigData reads credential data for IBM Cloud.
func (d *driver) getCredentialsConfigData() (string, error) {
	// Look for a user defined secret to get the IBM Cloud credentials from first
	sec, err := d.Listers.Secrets.Get(defaults.ImageRegistryPrivateConfigurationUser)
	if err != nil && errors.IsNotFound(err) {
		// Fall back to those provided by the credential minter if nothing is provided by the user
		sec, err = d.Listers.Secrets.Get(defaults.CloudCredentialsName)
		if err != nil {
			return "", fmt.Errorf("unable to get cluster minted credentials %q: %v", fmt.Sprintf("%s/%s", defaults.ImageRegistryOperatorNamespace, defaults.CloudCredentialsName), err)
		}
		if v, ok := sec.Data["ibmcloud_api_key"]; ok {
			return string(v), nil
		} else {
			return "", fmt.Errorf("secret %q does not contain required key \"ibmcloud_api_key\"", fmt.Sprintf("%s/%s", defaults.ImageRegistryOperatorNamespace, defaults.CloudCredentialsName))
		}
	} else if err != nil {
		return "", err
	} else {
		if v, ok := sec.Data["REGISTRY_STORAGE_IBMCOS_IAMAPIKEY"]; ok {
			return string(v), nil
		} else {
			return "", fmt.Errorf("secret %q does not contain required key \"REGISTRY_STORAGE_IBMCOS_IAMAPIKEY\"", fmt.Sprintf("%s/%s", defaults.ImageRegistryOperatorNamespace, defaults.ImageRegistryPrivateConfigurationUser))
		}
	}
}

// VolumeSecrets returns the same credentials data that the image-registry-operator
// is using so that it can be stored in the image-registry Pod's Secret.
func (d *driver) VolumeSecrets() (map[string]string, error) {
	IAMAPIKey, err := d.getCredentialsConfigData()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		imageRegistrySecretDataKey: IAMAPIKey,
	}, nil
}

// Volumes returns configuration for mounting credentials data as a Volume for
// image-registry Pods.
func (d *driver) Volumes() ([]corev1.Volume, []corev1.VolumeMount, error) {
	optional := false

	volume := corev1.Volume{
		Name: defaults.ImageRegistryPrivateConfiguration,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: defaults.ImageRegistryPrivateConfiguration,
				Optional:   &optional,
			},
		},
	}

	mount := corev1.VolumeMount{
		Name:      volume.Name,
		MountPath: imageRegistrySecretMountpoint,
		ReadOnly:  true,
	}

	return []corev1.Volume{volume}, []corev1.VolumeMount{mount}, nil
}
