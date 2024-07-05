package statefulset

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	microservicestatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	resourceName       string
	version            string
	namespaceName      string
	labels             map[string]string
	kubernetesProvider *pulumikubernetes.Provider
	availabilitySpec   *microservicestatemodel.MicroserviceKubernetesSpecAvailabilitySpec
	appContainer       *microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpec
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		resourceName:       contextState.Spec.ResourceName,
		version:            contextState.Spec.Version,
		namespaceName:      contextState.Spec.NamespaceName,
		labels:             contextState.Spec.Labels,
		kubernetesProvider: contextState.Spec.KubeProvider,
		availabilitySpec:   contextState.Spec.AvailabilitySpec,
		appContainer:       contextState.Spec.AppContainer,
	}
}
