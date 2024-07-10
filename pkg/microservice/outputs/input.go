package outputs

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	microservicekubernetesstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	resourceId                      string
	resourceName                    string
	version                         string
	environmentName                 string
	endpointDomainName              string
	namespaceName                   string
	internalHostname                string
	externalHostname                string
	kubeServiceName                 string
	kubeLocalEndpoint               string
	kubernetesImagePullSecretInputs []*microservicekubernetesstackmodel.KubernetesImagePullSecretInput
	forwardServicePort              int32
	listenerServicePort             int32
	gsaEmailId                      pulumi.StringOutput
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		resourceId:                      contextState.Spec.ResourceId,
		resourceName:                    contextState.Spec.ResourceName,
		version:                         contextState.Spec.Version,
		environmentName:                 contextState.Spec.EnvironmentInfo.EnvironmentName,
		endpointDomainName:              contextState.Spec.EndpointDomainName,
		namespaceName:                   contextState.Spec.NamespaceName,
		kubeServiceName:                 contextState.Spec.KubeServiceName,
		kubeLocalEndpoint:               contextState.Spec.KubeLocalEndpoint,
		kubernetesImagePullSecretInputs: contextState.Spec.KubernetesImagePullSecretInputs,
		forwardServicePort:              contextState.Spec.ForwardServicePort,
		listenerServicePort:             contextState.Spec.ListenerServicePort,
		gsaEmailId:                      contextState.Status.AddedResources.GsaEmailId,
		internalHostname:                contextState.Spec.InternalHostname,
		externalHostname:                contextState.Spec.ExternalHostname,
	}
}
