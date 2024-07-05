package virtualservice

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/controller"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Resources(ctx *pulumi.Context) error {
	input := extractInput(ctx)
	virtualServiceObject := buildVirtualServiceObject(input)
	if err := addVirtualService(ctx, virtualServiceObject, input.workspaceDir, input.kubernetesProvider); err != nil {
		return errors.Wrapf(err, "failed to add virtual service")
	}
	return nil
}

func addVirtualService(ctx *pulumi.Context, virtualServiceObject *v1beta1.VirtualService, workspace string, provider *pulumikubernetes.Provider) error {
	resourceName := fmt.Sprintf("virtual-service-%s", virtualServiceObject.Name)
	manifestPath := filepath.Join(workspace, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, virtualServiceObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{File: manifestPath}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to add virtual-service manifest")
	}
	return nil
}

/*
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:

	name: product-main-dev-planton-live-debug
	namespace: planton-pcs-dev-product

spec:

	gateways:
	- istio-ingress/product-dev-planton-live
	hosts:
	- product-main.dev.planton.live
	tcp:
	- match:
	  - port: 5005
	  route:
	  - destination:
	      host: main.planton-pcs-dev-product.svc.cluster.local
	      port:
	        number: 5005
*/
func buildVirtualServiceObject(input *input) *v1beta1.VirtualService {
	return &v1beta1.VirtualService{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "VirtualService",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      input.virtualServiceName,
			Namespace: input.microserviceNamespaceName,
		},
		Spec: networkingv1beta1.VirtualService{
			Gateways: []string{fmt.Sprintf("%s/%s", input.gatewayNamespaceName, input.gatewayName)},
			Hosts:    input.hostnames,
			Tcp: []*networkingv1beta1.TCPRoute{{
				Match: []*networkingv1beta1.L4MatchAttributes{{
					Port: controller.DebugPort,
				}},
				Route: []*networkingv1beta1.RouteDestination{{
					Destination: &networkingv1beta1.Destination{
						Host: input.kubeServiceFqdn,
						Port: &networkingv1beta1.PortSelector{Number: uint32(controller.DebugPort)},
					},
				}},
			}},
		},
	}
}
