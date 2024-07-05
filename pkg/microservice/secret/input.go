package secret

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	microservicekubernetesstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	KubernetesSecretTypeName = "kubernetes.io/dockerconfigjson"
)

type input struct {
	labels                          map[string]string
	namespace                       *kubernetescorev1.Namespace
	kubernetesImagePullSecretInputs []*microservicekubernetesstackmodel.KubernetesImagePullSecretInput
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		namespace:                       contextState.Status.AddedResources.Namespace,
		labels:                          contextState.Spec.Labels,
		kubernetesImagePullSecretInputs: contextState.Spec.KubernetesImagePullSecretInputs,
	}
}
