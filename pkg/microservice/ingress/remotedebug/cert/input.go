package cert

import (
	ingressnamespace "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/remotedebug/hostname"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	kubernetesProvider *pulumikubernetes.Provider
	labels             map[string]string
	workspaceDir       string
	certNamespaceName  string
	certName           string
	hostnames          []string
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	hostnames := hostname.GetHostnames(contextState.Spec.ResourceId, contextState.Spec.EnvironmentInfo.EnvironmentName, contextState.Spec.EndpointDomainName)
	certName := GetCertName(contextState.Spec.ResourceId)

	return &input{
		kubernetesProvider: contextState.Spec.KubeProvider,
		labels:             contextState.Spec.Labels,
		workspaceDir:       contextState.Spec.WorkspaceDir,
		certNamespaceName:  ingressnamespace.Name,
		certName:           certName,
		hostnames:          hostnames,
	}
}
