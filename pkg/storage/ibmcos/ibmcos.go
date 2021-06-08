package ibmcos

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"

	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go/aws/session"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	configapiv1 "github.com/openshift/api/config/v1"
	imageregistryv1 "github.com/openshift/api/imageregistry/v1"
	operatorapi "github.com/openshift/api/operator/v1"

	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
	"github.com/openshift/cluster-image-registry-operator/pkg/defaults"
	"github.com/openshift/cluster-image-registry-operator/pkg/envvar"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/util"
)

const IAMEndpoint = "https://iam.cloud.ibm.com/identity/token"

type driver struct {
	Context    context.Context
	Config     *imageregistryv1.ImageRegistryConfigStorageIBMCOS
	KubeConfig *rest.Config
	Listers    *regopclient.Listers

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
	)
	return
}

// CreateStorage attempts to create a COS bucket and apply any provided
// configuration.
func (d *driver) CreateStorage(cr *imageregistryv1.Config) error {
	fmt.Println("[WIP] ibmcos.CreateStorage")
	return nil
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

	_, err = client.HeadBucket(&s3.HeadBucketInput{
		Bucket: &bucketName,
	})

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
	fmt.Println("[WIP] ibmcos.VolumeSecrets")
	return nil, nil
}

// Volumes returns configuration for mounting credentials data as a Volume for
// image-registry Pods.
func (d *driver) Volumes() ([]corev1.Volume, []corev1.VolumeMount, error) {
	fmt.Println("[WIP] ibmcos.Volumes")
	return []corev1.Volume{}, []corev1.VolumeMount{}, nil
}
