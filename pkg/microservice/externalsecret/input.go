package externalsecret

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	kubernetesProvider *pulumikubernetes.Provider
	workspaceDir       string
	version            string
	labels             map[string]string
	namespaceName      string
	envSecrets         map[string]string
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		kubernetesProvider: contextState.Spec.KubeProvider,
		workspaceDir:       contextState.Spec.WorkspaceDir,
		version:            contextState.Spec.Version,
		labels:             contextState.Spec.Labels,
		namespaceName:      contextState.Spec.NamespaceName,
		envSecrets:         contextState.Spec.EnvSpec.Secrets,
	}
}
