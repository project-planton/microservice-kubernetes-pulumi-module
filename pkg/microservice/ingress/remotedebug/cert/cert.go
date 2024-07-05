package cert

import (
	"fmt"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/certmanager/clusterissuer"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"path/filepath"

	k8sapimachineryv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Resources(ctx *pulumi.Context) error {
	if err := addCert(ctx); err != nil {
		return errors.Wrap(err, "failed to add cert")
	}
	return nil
}

func addCert(ctx *pulumi.Context) error {
	input := extractInput(ctx)
	certObj := buildCertObject(input.certName, input.certNamespaceName, input.hostnames, input.labels)
	resourceName := fmt.Sprintf("cert-%s", certObj.Name)
	manifestPath := filepath.Join(input.workspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, certObj); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{
		File: manifestPath,
	}, pulumi.Provider(input.kubernetesProvider))
	if err != nil {
		return errors.Wrap(err, "failed to add cert manifest")
	}
	return nil
}

/*
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:

	name: debug-planton-pcs-dev-product-main
	namespace: planton-pcs-dev-postgres-apr

spec:

	secretName: cert-debug-planton-pcs-dev-product-main
	dnsNames:
	  - product-main.dev.planton.cloud
	  - product-main.dev.planton.live
	privateKey:
	  algorithm: ECDSA
	  size: 256
	issuerRef:
	  name: self-signed
	  kind: ClusterIssuer
	  group: cert-manager.io
*/
func buildCertObject(certName string, namespaceName string, hostnames []string, labels map[string]string) *v1.Certificate {
	return &v1.Certificate{
		TypeMeta: k8sapimachineryv1.TypeMeta{
			APIVersion: "cert-manager.io/v1",
			Kind:       "Certificate",
		},
		ObjectMeta: k8sapimachineryv1.ObjectMeta{
			Name:      certName,
			Namespace: namespaceName,
			Labels:    labels,
		},
		Spec: v1.CertificateSpec{
			SecretName: GetCertSecretName(certName),
			DNSNames:   hostnames,
			PrivateKey: &v1.CertificatePrivateKey{
				Algorithm: "ECDSA",
				Size:      256,
			},
			IssuerRef: cmmeta.ObjectReference{
				Kind:  "ClusterIssuer",
				Name:  clusterissuer.SelfSignedIssuerName,
				Group: "cert-manager.io",
			},
		},
	}
}

func GetCertSecretName(certName string) string {
	return fmt.Sprintf("cert-%s", certName)
}

// GetCertName for cert resource created for remote-debug port
// ex: debug-msi-planton-pcs-dev-product-main
func GetCertName(microserviceKubernetesId string) string {
	return fmt.Sprintf("debug-%s", microserviceKubernetesId)
}
