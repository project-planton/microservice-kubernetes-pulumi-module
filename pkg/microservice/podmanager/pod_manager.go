package podmanager

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/podmanager/deployment"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/podmanager/statefulset"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/integration/v1/kubernetes/apiresources/enums/podmanagertype"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	if i.podManagerType == podmanagertype.PodManagerType_deployment {
		if err := deployment.Resources(ctx); err != nil {
			return errors.Wrap(err, "failed to add deployment")
		}
	}
	if i.podManagerType == podmanagertype.PodManagerType_stateful_set {
		if err := statefulset.Resources(ctx); err != nil {
			return errors.Wrap(err, "failed to add stateful-set")
		}
	}
	return nil
}
