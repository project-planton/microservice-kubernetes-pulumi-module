package podmanager

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/integration/v1/kubernetes/apiresources/enums/podmanagertype"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	podManagerType podmanagertype.PodManagerType
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		podManagerType: contextState.Spec.PodManagerType,
	}
}
