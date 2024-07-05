package common

import (
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	kubeclustermodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	microservicestatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	MicroserviceContainerName = "microservice"
	DebugPortName             = "debug"
	DebugPortNumber           = 5005
)

type input struct {
	kubernetesProvider *pulumikubernetes.Provider
	labels             map[string]string
	namespaceName      string
	appContainer       *microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpec
	version            string
	envSpec            *microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpecEnvSpec
	appPorts           []*microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpecPortSpec
	sidecars           []*kubeclustermodel.Container
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	return &input{
		kubernetesProvider: contextState.Spec.KubeProvider,
		labels:             contextState.Spec.Labels,
		namespaceName:      contextState.Spec.NamespaceName,
		appContainer:       contextState.Spec.AppContainer,
		version:            contextState.Spec.Version,
		envSpec:            contextState.Spec.EnvSpec,
		appPorts:           contextState.Spec.AppPorts,
		sidecars:           contextState.Spec.Sidecars,
	}
}
