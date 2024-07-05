package service

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	microservicestatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	kubernetesProvider *pulumikubernetes.Provider
	version            string
	labels             map[string]string
	namespaceName      string
	appPorts           []*microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpecPortSpec
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		kubernetesProvider: contextState.Spec.KubeProvider,
		version:            contextState.Spec.Version,
		labels:             contextState.Spec.Labels,
		namespaceName:      contextState.Spec.NamespaceName,
		appPorts:           contextState.Spec.AppPorts,
	}
}
