package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	externalsecretsv1beta1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/externalsecrets/externalsecrets/v1beta1"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/datatypes/maps"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func externalSecret(ctx *pulumi.Context, locals *Locals,
	createdNamespace *kubernetescorev1.Namespace, labels map[string]string) error {
	var secrets = locals.MicroserviceKubernetes.Spec.Container.App.Env.Secrets
	if secrets == nil || len(secrets) == 0 {
		return nil
	}

	var secretData = externalsecretsv1beta1.ExternalSecretSpecDataArray{}

	sortedSecretKeys := maps.SortMapKeys(secrets)
	for _, sortedSecretKey := range sortedSecretKeys {
		secretData = append(secretData, externalsecretsv1beta1.ExternalSecretSpecDataArgs{
			SecretKey: pulumi.String(sortedSecretKey),
			RemoteRef: externalsecretsv1beta1.ExternalSecretSpecDataRemoteRefArgs{
				Key:     pulumi.String(secrets[sortedSecretKey]),
				Version: pulumi.String("latest"),
			},
		})
	}

	_, err := externalsecretsv1beta1.NewExternalSecret(ctx,
		fmt.Sprintf("external-secret-%s", locals.MicroserviceKubernetes.Spec.Version),
		&externalsecretsv1beta1.ExternalSecretArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(locals.MicroserviceKubernetes.Spec.Version),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(labels),
			},
			Spec: &externalsecretsv1beta1.ExternalSecretSpecArgs{
				Data:            secretData,
				RefreshInterval: pulumi.String("1m"),
				SecretStoreRef: &externalsecretsv1beta1.ExternalSecretSpecSecretStoreRefArgs{
					Kind: pulumi.String("ClusterSecretStore"),
					Name: pulumi.String(vars.GcpSecretsManagerClusterSecretStoreName),
				},
				Target: &externalsecretsv1beta1.ExternalSecretSpecTargetArgs{
					Name: pulumi.String(locals.MicroserviceKubernetes.Spec.Version),
				},
			},
		}, pulumi.Parent(createdNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to add external secret")
	}

	return nil
}