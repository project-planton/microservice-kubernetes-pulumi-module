package ingress

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	isIngressEnabled   bool
	endpointDomainName string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		isIngressEnabled:   ctxConfig.Spec.IsIngressEnabled,
		endpointDomainName: ctxConfig.Spec.EndpointDomainName,
	}
}
