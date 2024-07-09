package microservice

import (
	"fmt"
	"github.com/pkg/errors"
	environmentblueprinthostnames "github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/hostnames"
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/remotedebug/hostname"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/service"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

func loadConfig(ctx *pulumi.Context, resourceStack *ResourceStack) (*microservicecontextstate.ContextState, error) {

	kubernetesProvider, err := pulumikubernetesprovider.GetWithStackCredentials(ctx, resourceStack.Input.CredentialsInput.Kubernetes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup kubernetes provider")
	}

	gcpProvider, err := pulumigoogleprovider.Get(ctx, resourceStack.Input.CredentialsInput.Google)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup gcp provider")
	}

	var resourceId = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Metadata.Id
	var resourceName = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Metadata.Name
	var environmentInfo = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.EnvironmentInfo
	var isIngressEnabled = false
	var externalHostname = ""
	var internalHostname = ""
	var forwardServicePort = int32(0)
	var listenerServicePort = int32(0)

	if resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Ingress != nil {
		isIngressEnabled = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Ingress.IsEnabled
	}

	var endpointDomainName = ""
	var envDomainName = ""
	var ingressType = kubernetesworkloadingresstype.KubernetesWorkloadIngressType_unspecified

	if isIngressEnabled {
		endpointDomainName = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Ingress.EndpointDomainName
		envDomainName = environmentblueprinthostnames.GetExternalEnvHostname(environmentInfo.EnvironmentName, endpointDomainName)
		ingressType = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Ingress.IngressType

		forwardServicePort = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Ingress.ForwardServicePort
		listenerServicePort = resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Ingress.ListenerIngressPort

		externalHostname = hostname.GetExternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
		internalHostname = hostname.GetInternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
	}

	return &microservicecontextstate.ContextState{
		Spec: &microservicecontextstate.Spec{
			KubeProvider:        kubernetesProvider,
			ResourceId:          resourceId,
			ResourceName:        resourceName,
			Labels:              resourceStack.KubernetesLabels,
			WorkspaceDir:        resourceStack.WorkspaceDir,
			NamespaceName:       resourceId,
			EnvironmentInfo:     resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.EnvironmentInfo,
			IsIngressEnabled:    isIngressEnabled,
			IngressType:         ingressType,
			ExternalHostname:    externalHostname,
			InternalHostname:    internalHostname,
			ForwardServicePort:  forwardServicePort,
			ListenerServicePort: listenerServicePort,
			EndpointDomainName:  endpointDomainName,
			EnvDomainName:       envDomainName,
			KubeServiceName:     resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Version,
			KubeLocalEndpoint:   service.GetFqdn(resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Version, resourceId),

			KubernetesImagePullSecretInputs: resourceStack.Input.ResourceInput.KubernetesImagePullSecrets,
			GcpProvider:                     gcpProvider,
			ContainerClusterProject:         resourceStack.Input.ResourceInput.ContainerClusterProject,
			IsWorkloadIdentityEnabled:       resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.IsWorkloadIdentityEnabled,
			WorkloadIdentityGsaAccountId:    getGsaAccountId(environmentInfo.EnvironmentId, resourceName),
			PodManagerType:                  resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.PodManagerType,
			AppContainer:                    resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Container.App,
			Version:                         resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Version,
			EnvSpec:                         resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Container.App.Env,
			AppPorts:                        resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Container.App.Ports,
			Sidecars:                        resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Container.Sidecars,
			AvailabilitySpec:                resourceStack.Input.ResourceInput.MicroserviceKubernetes.Spec.Availability,
		},
		Status: &microservicecontextstate.Status{},
	}, nil
}

func getGsaAccountId(environmentId, microserviceKubernetesName string) string {
	return fmt.Sprintf("%s-%s", removeCompanyId(environmentId), microserviceKubernetesName)
}

func removeCompanyId(environmentId string) string {
	delimiter := "-"
	parts := strings.Split(environmentId, delimiter)
	if len(parts) > 1 {
		return strings.Join(parts[1:], delimiter)
	}
	return environmentId
}
