package ingress

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/remotedebug"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/virtualservice"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	if !i.isIngressEnabled || i.endpointDomainName == "" {
		return nil
	}
	err := virtualservice.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add virtual service resource")
	}
	err = remotedebug.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add remote-debug resources")
	}
	return nil
}
