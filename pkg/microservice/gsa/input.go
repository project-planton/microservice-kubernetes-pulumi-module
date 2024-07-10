package gsa

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	gcpresourceprojectv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/model/provider/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	microserviceKubernetesName string
	containerClusterProject    *gcpresourceprojectv1.GcpProject
	isWorkloadIdentityEnabled  bool
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		microserviceKubernetesName: contextState.Spec.ResourceName,
		containerClusterProject:    contextState.Spec.ContainerClusterProject,
		isWorkloadIdentityEnabled:  contextState.Spec.IsWorkloadIdentityEnabled,
	}
}
