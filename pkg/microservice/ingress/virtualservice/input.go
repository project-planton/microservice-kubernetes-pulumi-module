package virtualservice

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	kubernetesProvider       *pulumikubernetes.Provider
	labels                   map[string]string
	workspaceDir             string
	microserviceKubernetesId string
	namespaceName            string
	microserviceVersion      string
	servicePort              int32
	kubeServiceFqdn          string
	environmentName          string
	endpointDomainName       string
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		kubernetesProvider:       contextState.Spec.KubeProvider,
		labels:                   contextState.Spec.Labels,
		workspaceDir:             contextState.Spec.WorkspaceDir,
		microserviceKubernetesId: contextState.Spec.ResourceId,
		namespaceName:            contextState.Spec.NamespaceName,
		microserviceVersion:      contextState.Spec.Version,
		servicePort:              contextState.Spec.ForwardServicePort,
		kubeServiceFqdn:          contextState.Spec.KubeLocalEndpoint,
		environmentName:          contextState.Spec.EnvironmentInfo.EnvironmentName,
		endpointDomainName:       contextState.Spec.EndpointDomainName,
	}
}
