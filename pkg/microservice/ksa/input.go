package ksa

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	gcpresourceprojectv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/model/provider/gcp"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	namespaceName                string
	labels                       map[string]string
	namespace                    *kubernetescorev1.Namespace
	isWorkloadIdentityEnabled    bool
	containerClusterProject      *gcpresourceprojectv1.GcpProject
	workloadIdentityGsaAccountId *random.RandomId
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		isWorkloadIdentityEnabled:    contextState.Spec.IsWorkloadIdentityEnabled,
		namespaceName:                contextState.Spec.NamespaceName,
		namespace:                    contextState.Status.AddedResources.Namespace,
		labels:                       contextState.Spec.Labels,
		containerClusterProject:      contextState.Spec.ContainerClusterProject,
		workloadIdentityGsaAccountId: contextState.Status.AddedResources.WorkLoadIdentityGsaAccountId,
	}
}
