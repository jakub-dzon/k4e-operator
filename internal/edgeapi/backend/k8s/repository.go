package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/project-flotta/flotta-operator/api/v1alpha1"
	"github.com/project-flotta/flotta-operator/internal/common/repository/edgedevice"
	"github.com/project-flotta/flotta-operator/internal/common/repository/edgedeviceset"
	"github.com/project-flotta/flotta-operator/internal/common/repository/edgedevicesignedrequest"
	"github.com/project-flotta/flotta-operator/internal/common/repository/edgeworkload"
	"github.com/project-flotta/flotta-operator/internal/edgeapi/k8sclient"
)

type EdgeDeviceRepository interface {
	GetEdgeDevice(ctx context.Context, name, namespace string) (*v1alpha1.EdgeDevice, error)
	PatchEdgeDeviceStatus(ctx context.Context, edgeDevice *v1alpha1.EdgeDevice, patch *client.Patch) error
	UpdateEdgeDeviceLabels(ctx context.Context, device *v1alpha1.EdgeDevice, labels map[string]string) error
	PatchEdgeDevice(ctx context.Context, old, new *v1alpha1.EdgeDevice) error
	RemoveEdgeDeviceFinalizer(ctx context.Context, edgeDevice *v1alpha1.EdgeDevice, finalizer string) error
}

type EdgeDeviceSignedRequestRepository interface {
	GetEdgeDeviceSignedRequest(ctx context.Context, name string, namespace string) (*v1alpha1.EdgeDeviceSignedRequest, error)
	CreateEdgeDeviceSignedRequest(ctx context.Context, edgeDeviceSignedRequest *v1alpha1.EdgeDeviceSignedRequest) error
}

type EdgeWorkloadRepository interface {
	GetEdgeWorkload(ctx context.Context, name string, namespace string) (*v1alpha1.EdgeWorkload, error)
}

type EdgeDeviceSetRepository interface {
	GetEdgeDeviceSet(ctx context.Context, name string, namespace string) (*v1alpha1.EdgeDeviceSet, error)
}

type CoreRepository interface {
	GetSecret(ctx context.Context, name string, namespace string) (*v1.Secret, error)
	GetConfigMap(ctx context.Context, name string, namespace string) (*v1.ConfigMap, error)
}

//go:generate mockgen -package=k8s -destination=mock_repository_facade.go . RepositoryFacade
type RepositoryFacade interface {
	EdgeDeviceRepository
	EdgeDeviceSignedRequestRepository
	EdgeWorkloadRepository
	EdgeDeviceSetRepository
	CoreRepository
}
type repositoryFacade struct {
	deviceSignedRequestRepository edgedevicesignedrequest.Repository
	deviceRepository              edgedevice.Repository
	workloadRepository            edgeworkload.Repository
	deviceSetRepository           edgedeviceset.Repository

	client k8sclient.K8sClient
}

func NewRepository(deviceSignedRequestRepository edgedevicesignedrequest.Repository,
	deviceRepository edgedevice.Repository,
	workloadRepository edgeworkload.Repository,
	deviceSetRepository edgedeviceset.Repository,
	client k8sclient.K8sClient) RepositoryFacade {
	return &repositoryFacade{
		deviceSignedRequestRepository: deviceSignedRequestRepository,
		deviceRepository:              deviceRepository,
		deviceSetRepository:           deviceSetRepository,
		workloadRepository:            workloadRepository,
		client:                        client,
	}
}

func (b *repositoryFacade) GetEdgeDevice(ctx context.Context, name, namespace string) (*v1alpha1.EdgeDevice, error) {
	return b.deviceRepository.Read(ctx, name, namespace)
}

func (b *repositoryFacade) PatchEdgeDeviceStatus(ctx context.Context, edgeDevice *v1alpha1.EdgeDevice, patch *client.Patch) error {
	return b.deviceRepository.PatchStatus(ctx, edgeDevice, patch)
}

func (b *repositoryFacade) UpdateEdgeDeviceLabels(ctx context.Context, device *v1alpha1.EdgeDevice, labels map[string]string) error {
	return b.deviceRepository.UpdateLabels(ctx, device, labels)
}

func (b *repositoryFacade) PatchEdgeDevice(ctx context.Context, old, new *v1alpha1.EdgeDevice) error {
	return b.deviceRepository.Patch(ctx, old, new)
}

func (b *repositoryFacade) RemoveEdgeDeviceFinalizer(ctx context.Context, edgeDevice *v1alpha1.EdgeDevice, finalizer string) error {
	return b.deviceRepository.RemoveFinalizer(ctx, edgeDevice, finalizer)
}

func (b *repositoryFacade) GetEdgeDeviceSignedRequest(ctx context.Context, name string, namespace string) (*v1alpha1.EdgeDeviceSignedRequest, error) {
	return b.deviceSignedRequestRepository.Read(ctx, name, namespace)
}

func (b *repositoryFacade) CreateEdgeDeviceSignedRequest(ctx context.Context, edgedeviceSignedRequest *v1alpha1.EdgeDeviceSignedRequest) error {
	return b.deviceSignedRequestRepository.Create(ctx, edgedeviceSignedRequest)
}

func (b *repositoryFacade) GetEdgeWorkload(ctx context.Context, name string, namespace string) (*v1alpha1.EdgeWorkload, error) {
	return b.workloadRepository.Read(ctx, name, namespace)
}

func (b *repositoryFacade) GetEdgeDeviceSet(ctx context.Context, name string, namespace string) (*v1alpha1.EdgeDeviceSet, error) {
	return b.deviceSetRepository.Read(ctx, name, namespace)
}

func (b *repositoryFacade) GetSecret(ctx context.Context, name string, namespace string) (*v1.Secret, error) {
	secret := v1.Secret{}
	err := b.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func (b *repositoryFacade) GetConfigMap(ctx context.Context, name string, namespace string) (*v1.ConfigMap, error) {
	configMap := v1.ConfigMap{}
	err := b.client.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &configMap)
	if err != nil {
		return nil, err
	}
	return &configMap, nil
}
