package secret

import (
	"github.com/pkg/errors"
	microservicesecretdata "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/secret/data"
	code2cloudv1deploymsistackk8smodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/stack/model"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	pk8smv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	for _, imagePullSecretInput := range i.kubernetesImagePullSecretInputs {
		_, err := addImagePullSecret(ctx, imagePullSecretInput, i.namespace, i.labels)
		if err != nil {
			return errors.Wrap(err, "failed to add deployment")
		}
	}
	return nil
}

func addImagePullSecret(ctx *pulumi.Context, secretInput *code2cloudv1deploymsistackk8smodel.KubernetesImagePullSecretInput,
	addedNamespace *v1.Namespace, labels map[string]string) (*v1.Secret, error) {
	secretStringData, err := microservicesecretdata.Get(secretInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get secret string microservicesecretdata")
	}
	ns, err := v1.NewSecret(ctx, secretInput.ImagePullSecretName, &v1.SecretArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Secret"),
		Metadata: pk8smv1.ObjectMetaArgs{
			Name:      pulumi.String(secretInput.ImagePullSecretName),
			Namespace: addedNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(labels),
		},
		Type:       pulumi.String(KubernetesSecretTypeName),
		StringData: pulumi.ToStringMap(secretStringData),
	}, pulumi.Parent(addedNamespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}
	return ns, nil
}
