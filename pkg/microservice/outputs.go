package microservice

import (
	"context"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/pkg/errors"
	microservicestatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	microservicestackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
)

func Outputs(ctx context.Context, input *microservicestackmodel.MicroserviceKubernetesStackInput) (*microservicestatemodel.MicroserviceKubernetesStatusStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return OutputMapTransformer(stackOutput, input), nil
}

// OutputMapTransformer transforms untyped map of outputs into strictly typed proto object
func OutputMapTransformer(stackOutput map[string]interface{}, input *microservicestackmodel.MicroserviceKubernetesStackInput) *microservicestatemodel.MicroserviceKubernetesStatusStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &microservicestatemodel.MicroserviceKubernetesStatusStackOutputs{}
	}
	return &microservicestatemodel.MicroserviceKubernetesStatusStackOutputs{
		Namespace:          backend.GetVal(stackOutput, outputs.GetNamespaceNameOutputName()),
		Service:            backend.GetVal(stackOutput, outputs.GetServiceNameOutputName()),
		PortForwardCommand: backend.GetVal(stackOutput, outputs.GetPortForwardCommandOutputName()),
		KubeEndpoint:       backend.GetVal(stackOutput, outputs.GetFqdnOutputName()),
		ExternalHostname:   backend.GetVal(stackOutput, outputs.GetExternalClusterHostnameOutputName()),
		InternalHostname:   backend.GetVal(stackOutput, outputs.GetInternalClusterHostnameOutputName()),
		WorkloadIdentityAccountId: backend.GetVal(stackOutput, outputs.GetGsaEmailOutputName(
			input.ResourceInput.MicroserviceKubernetes.Metadata.Id,
		)),
	}
}
