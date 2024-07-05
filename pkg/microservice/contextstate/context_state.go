package contextstate

import (
	gcpresourceprojectv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/model/provider/gcp"
	code2cloudenvironmentmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/environment/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	kubeclustermodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	microservicestatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	microservicekubernetesstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/integration/v1/kubernetes/apiresources/enums/podmanagertype"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
)

const (
	Key = "ctx-state"
)

type ContextState struct {
	Spec   *Spec
	Status *Status
}

type Spec struct {
	KubeProvider        *kubernetes.Provider
	ResourceId          string
	ResourceName        string
	Labels              map[string]string
	WorkspaceDir        string
	NamespaceName       string
	EnvironmentInfo     *code2cloudenvironmentmodel.ApiResourceEnvironmentInfo
	IsIngressEnabled    bool
	IngressType         kubernetesworkloadingresstype.KubernetesWorkloadIngressType
	ExternalHostname    string
	InternalHostname    string
	ForwardServicePort  int32
	ListenerServicePort int32
	EndpointDomainName  string
	EnvDomainName       string
	KubeServiceName     string
	KubeLocalEndpoint   string

	KubernetesImagePullSecretInputs []*microservicekubernetesstackmodel.KubernetesImagePullSecretInput
	ContainerClusterProject         *gcpresourceprojectv1.GcpProject
	GcpProvider                     *pulumigcp.Provider
	IsWorkloadIdentityEnabled       bool
	WorkloadIdentityGsaAccountId    string
	PodManagerType                  podmanagertype.PodManagerType
	AppContainer                    *microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpec
	Version                         string
	EnvSpec                         *microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpecEnvSpec
	AppPorts                        []*microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpecPortSpec
	Sidecars                        []*kubeclustermodel.Container
	AvailabilitySpec                *microservicestatemodel.MicroserviceKubernetesSpecAvailabilitySpec
}

type Status struct {
	AddedResources *AddedResources
}

type AddedResources struct {
	Namespace *kubernetescorev1.Namespace
}
