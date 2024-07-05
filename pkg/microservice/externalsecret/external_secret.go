package externalsecret

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/datatypes/maps"
	"path/filepath"
	"time"

	externalsecretsv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/externalsecrets/clustersecretstore"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Resources(ctx *pulumi.Context) error {
	input := extractInput(ctx)
	if input.envSecrets == nil || len(input.envSecrets) == 0 {
		return nil
	}
	externalSecretObject := buildExternalSecretObject(input, input.envSecrets)
	if err := addExternalSecret(ctx, externalSecretObject, input.workspaceDir, input.kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to add external secret")
	}
	return nil
}

func addExternalSecret(ctx *pulumi.Context, externalSecretsObject *externalsecretsv1beta1.ExternalSecret, workspace string, provider *pulumikubernetes.Provider) error {
	resourceName := fmt.Sprintf("external-secret-%s", externalSecretsObject.Name)
	manifestPath := filepath.Join(workspace, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, externalSecretsObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	if _, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{File: manifestPath}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to add kubernetes config file")
	}
	return nil
}

/*
---
# Source: pcs-project/templates/external-secret.yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:

	name: gql-master-dev
	namespace: pcs-stream
	labels:
	  product: planton
	  productCode: pcs

spec:

	refreshInterval: 1m
	secretStoreRef:
	  kind: ClusterSecretStore
	  name: gcp-backend
	target:
	  name: gql-master-dev  # secret name in kubernetes
	data:
	  - secretKey: pcs-kfk-sasl-username    # name of the gcp secret
	    remoteRef:
	      key: pcs-kfk-sasl-username # name of the gcp secret
	  - secretKey: pcs-kfk-sasl-password    # name of the gcp secret
	    remoteRef:
	      key: pcs-kfk-sasl-password # name of the gcp secret
*/
func buildExternalSecretObject(input *input, secrets map[string]string) *externalsecretsv1beta1.ExternalSecret {
	return &externalsecretsv1beta1.ExternalSecret{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "external-secrets.io/v1beta1",
			Kind:       "ExternalSecret",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      input.version,
			Namespace: input.namespaceName,
			Labels:    input.labels,
		},
		Spec: externalsecretsv1beta1.ExternalSecretSpec{
			SecretStoreRef: externalsecretsv1beta1.SecretStoreRef{
				Kind: clustersecretstore.Kind,
				Name: clustersecretstore.Name,
			},
			Target: externalsecretsv1beta1.ExternalSecretTarget{
				Name: input.version,
			},
			RefreshInterval: &k8smetav1.Duration{
				Duration: 1 * time.Minute,
			},
			Data: buildSecretData(secrets),
		},
	}
}

/*
data:
  - secretKey: pcs-kfk-sasl-username # name of the key inside the kubernetes secret (spec.data.<key>)
    remoteRef:
    key: pcs-kfk-sasl-username # id of the gcp secret
*/
func buildSecretData(secrets map[string]string) []externalsecretsv1beta1.ExternalSecretData {
	resp := make([]externalsecretsv1beta1.ExternalSecretData, 0)
	sortedSecretKeys := maps.SortMapKeys(secrets)
	for _, sortedSecretKey := range sortedSecretKeys {
		resp = append(resp, externalsecretsv1beta1.ExternalSecretData{
			SecretKey: sortedSecretKey,
			RemoteRef: externalsecretsv1beta1.ExternalSecretDataRemoteRef{
				Key:     secrets[sortedSecretKey],
				Version: "latest",
			},
		})
	}
	return resp
}
