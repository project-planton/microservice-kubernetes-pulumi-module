package workloadidentity

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	gcpresourceprojectv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/model/provider/gcp"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	microserviceKubernetesId     string
	gcpProvider                  *pulumigcp.Provider
	isWorkloadIdentityEnabled    bool
	namespaceName                string
	namespace                    *kubernetescorev1.Namespace
	labels                       map[string]string
	containerClusterProject      *gcpresourceprojectv1.GcpProject
	workloadIdentityGsaAccountId *random.RandomId
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		microserviceKubernetesId:     contextState.Spec.ResourceId,
		gcpProvider:                  contextState.Spec.GcpProvider,
		isWorkloadIdentityEnabled:    contextState.Spec.IsWorkloadIdentityEnabled,
		namespaceName:                contextState.Spec.NamespaceName,
		namespace:                    contextState.Status.AddedResources.Namespace,
		labels:                       contextState.Spec.Labels,
		containerClusterProject:      contextState.Spec.ContainerClusterProject,
		workloadIdentityGsaAccountId: contextState.Status.AddedResources.WorkloadIdentityGsaAccountId,
	}
}
