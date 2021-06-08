package ibmcos

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"

	imageregistryv1 "github.com/openshift/api/imageregistry/v1"
	operatorapi "github.com/openshift/api/operator/v1"

	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
	"github.com/openshift/cluster-image-registry-operator/pkg/defaults"
	"github.com/openshift/cluster-image-registry-operator/pkg/envvar"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/util"
)

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

// StorageExists checks if an S3 bucket with the given name exists
// and we can access it
func (d *driver) StorageExists(cr *imageregistryv1.Config) (bool, error) {
	fmt.Println("[WIP] ibmcos.StorageExists")
	return false, nil
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
