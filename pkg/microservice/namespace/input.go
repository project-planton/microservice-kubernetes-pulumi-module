package namespace

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	namespaceName string
	labels        map[string]string
	kubeProvider  *pulumikubernetes.Provider
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		namespaceName: contextState.Spec.NamespaceName,
		labels:        contextState.Spec.Labels,
		kubeProvider:  contextState.Spec.KubeProvider,
	}
}
