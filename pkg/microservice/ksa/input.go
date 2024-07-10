package ksa

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	namespaceName             string
	labels                    map[string]string
	namespace                 *kubernetescorev1.Namespace
	isWorkloadIdentityEnabled bool
	gsaEmailId                pulumi.StringOutput
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		isWorkloadIdentityEnabled: contextState.Spec.IsWorkloadIdentityEnabled,
		namespaceName:             contextState.Spec.NamespaceName,
		namespace:                 contextState.Status.AddedResources.Namespace,
		labels:                    contextState.Spec.Labels,
		gsaEmailId:                contextState.Status.AddedResources.GsaEmailId,
	}
}
