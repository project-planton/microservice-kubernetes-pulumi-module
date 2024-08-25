package pkg

import (
	"github.com/pkg/errors"
	microservicekubernetesmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/microservicekubernetes"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *microservicekubernetesmodel.MicroserviceKubernetesStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals, err := initializeLocals(ctx, s.Input)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesClusterCredential(ctx,
		s.Input.KubernetesClusterCredential, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	//create namespace resource
	createdNamespace, err := kubernetescorev1.NewNamespace(ctx,
		locals.Namespace,
		&kubernetescorev1.NamespaceArgs{
			Metadata: kubernetesmetav1.ObjectMetaPtrInput(
				&kubernetesmetav1.ObjectMetaArgs{
					Name:   pulumi.String(locals.Namespace),
					Labels: pulumi.ToStringMap(s.Labels),
				}),
		},
		pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s namespace", locals.Namespace)
	}

	//create kubernetes deployment resources
	if err := deployment(ctx, locals, createdNamespace, s.Labels); err != nil {
		return errors.Wrap(err, "failed to create microservice deployment")
	}

	//create kubernetes service resources
	if err := service(ctx, locals, createdNamespace, s.Labels); err != nil {
		return errors.Wrap(err, "failed to create microservice kubernetes service resource")
	}

	if err := externalSecret(ctx, locals, createdNamespace, s.Labels); err != nil {
		return errors.Wrap(err, "failed to create external secret")
	}

	//create istio-ingress resources if ingress is enabled.
	if locals.MicroserviceKubernetes.Spec.Ingress.IsEnabled {
		if err := istioIngress(ctx, locals, kubernetesProvider, createdNamespace, s.Labels); err != nil {
			return errors.Wrap(err, "failed to create istio ingress resources")
		}
	}

	return nil
}
