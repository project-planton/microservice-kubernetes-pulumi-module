package gateway

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/controller"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resources adds an istio gateway resource to support remote-debug based on sniHostname.
// tls is configured using a self-signed certificate
func Resources(ctx *pulumi.Context) error {
	input := extractInput(ctx)
	gatewayObject := buildGatewayObject(input)
	resourceName := fmt.Sprintf("gateway-%s", input.gatewayName)
	manifestPath := filepath.Join(input.workspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, gatewayObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName,
		&pulumik8syaml.ConfigFileArgs{
			File: manifestPath,
		}, pulumi.Provider(input.kubernetesProvider))
	if err != nil {
		return errors.Wrap(err, "failed to add ingress-gateway manifest")
	}
	return nil
}

/*
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:

	name: product-dev-planton-live
	namespace: istio-ingress

spec:

	selector:
	  app: istio-ingress
	  istio: ingress
	servers:
	- hosts:
	  - product-main.dev.planton.live
	  name: debug
	  port:
	    name: debug
	    number: 5005
	    protocol: TLS
	  tls:
	    credentialName: cert-dev-planton-live
	    mode: SIMPLE
*/
func buildGatewayObject(input *input) *v1beta1.Gateway {
	return &v1beta1.Gateway{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "Gateway",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      input.gatewayName,
			Namespace: input.gatewayNamespaceName,
			Labels:    input.labels,
		},
		Spec: networkingv1beta1.Gateway{
			Selector: controller.SelectorLabels,
			Servers:  getGatewaySpecServers(input.certSecretName, input.hostnames),
		},
	}
}

func getGatewaySpecServers(certSecretName string, hostnames []string) []*networkingv1beta1.Server {
	servers := make([]*networkingv1beta1.Server, 0)
	servers = append(servers, &networkingv1beta1.Server{
		Name:  "debug",
		Hosts: hostnames,
		Port: &networkingv1beta1.Port{
			Name:     "debug",
			Number:   controller.DebugPort,
			Protocol: "TLS",
		},
		Tls: &networkingv1beta1.ServerTLSSettings{
			Mode:           networkingv1beta1.ServerTLSSettings_SIMPLE,
			CredentialName: certSecretName,
		},
	})
	return servers
}

// GetGatewayName for gateway resource created for remote-debug port
// ex: debug-msi-planton-pcs-dev-product-main
func GetGatewayName(microserviceKubernetesId string) string {
	return fmt.Sprintf("debug-%s", microserviceKubernetesId)
}
