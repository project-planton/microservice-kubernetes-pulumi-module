package outputs

import (
	microservicekubernetesmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/microservicekubernetes/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	Namespace               = "namespace"
	Service                 = "service"
	KubePortForwardCommand  = "port-forward-command"
	KubeEndpoint            = "kube-endpoint"
	IngressExternalHostname = "ingress-external-hostname"
	IngressInternalHostname = "ingress-internal-hostname"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *microservicekubernetesmodel.MicroserviceKubernetesStackInput) *microservicekubernetesmodel.MicroserviceKubernetesStackOutputs {
	return &microservicekubernetesmodel.MicroserviceKubernetesStackOutputs{
		Namespace:          autoapistackoutput.GetVal(pulumiOutputs, Namespace),
		KubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, KubeEndpoint),
		Service:            autoapistackoutput.GetVal(pulumiOutputs, Service),
		PortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, KubePortForwardCommand),
		ExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressExternalHostname),
		InternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressInternalHostname),
	}
}
