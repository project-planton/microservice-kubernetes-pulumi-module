package microservice

import (
	"github.com/pkg/errors"
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/externalsecret"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/gsa"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ksa"
	microservicenamespace "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/namespace"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/outputs"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/podmanager"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/secret"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/service"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/workloadidentity"
	code2cloudv1deploymsistackk8smodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	WorkspaceDir     string
	Input            *code2cloudv1deploymsistackk8smodel.MicroserviceKubernetesStackInput
	KubernetesLabels map[string]string
}

func (resourceStack *ResourceStack) Resources(ctx *pulumi.Context) error {
	//load context config
	var ctxConfig, err = loadConfig(ctx, resourceStack)
	if err != nil {
		return errors.Wrap(err, "failed to initiate context config")
	}
	ctx = ctx.WithValue(microservicecontextstate.Key, *ctxConfig)

	// Create the namespace resource
	ctx, err = microservicenamespace.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace resource")
	}

	// Create the gcp service account id resource
	ctx, err = gsa.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create gcp service account id")
	}

	if err := secret.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add secret resources")
	}

	err = ksa.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add kubernetes service account resources")
	}
	err = workloadidentity.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add workload identity resources")
	}
	if err := podmanager.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add pod controller resources")
	}
	err = service.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add kubernetes service resources")
	}
	if err := externalsecret.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add external secret resources")
	}

	err = ingress.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add istio virtual-service resources")
	}

	err = outputs.Export(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to export microservice kubernetes outputs")
	}

	return nil
}
