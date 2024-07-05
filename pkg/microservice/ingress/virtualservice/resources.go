package virtualservice

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/pulumi/pulumicustomoutput"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud-inc/go-commons/network/dns/zone"
	ingressnamespace "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Resources include on virtual service object for each endpoint domain name.
// it is not possible to create a single virtual-service resource for all endpoint domain names as
// istio-gateway resource is different for each endpoint domain name and one virtual service allows configuring
// only one gateway resource name
func Resources(ctx *pulumi.Context) error {
	input := extractInput(ctx)
	msiIngressDomains := make([]string, 0)
	gatewayIngressDomainName := getGatewayEndpointDomain(input.endpointDomainName, input.environmentName)
	msiIngressDomain := getMicroserviceKubernetesEndpointDomain(
		input.microserviceKubernetesId,
		gatewayIngressDomainName,
	)

	msiIngressDomains = append(msiIngressDomains, msiIngressDomain)

	virtualServiceObject := buildVirtualServiceObject(input, gatewayIngressDomainName, msiIngressDomain)
	if err := addVirtualService(ctx, virtualServiceObject, input.workspaceDir, input.kubernetesProvider); err != nil {
		return errors.Wrapf(err, "failed to add virtual service for %s domain", msiIngressDomain)
	}
	exportOutputs(ctx, msiIngressDomains)
	return nil
}

func exportOutputs(ctx *pulumi.Context, domains []string) {
	ctx.Export(GetIngressDomainsOutputName(), pulumi.ToStringArray(domains))
}

func GetIngressDomainsOutputName() string {
	return pulumicustomoutput.Name("ingress-domains")
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

// getGatewayEndpointDomain returns the env domain of the microservice deployment endpoint.
// ex: dev.planton.cloud. this domain is only used for adding the correct ingress gateway resource and is not included in the hostnames in the virtual resource.
func getGatewayEndpointDomain(domainName string, productEnvName string) string {
	return fmt.Sprintf("%s.%s", productEnvName, domainName)
}

// getMicroserviceKubernetesEndpointDomain returns the versioned env domain of the microservice deployment endpoint.
// ex: msi-planton-pcs-prod-console-main.dev.planton.cloud
func getMicroserviceKubernetesEndpointDomain(microserviceKubernetesId, envDomainName string) string {
	return fmt.Sprintf("%s.%s", microserviceKubernetesId, envDomainName)
}

/*
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:

	name: testing
	namespace: planton-pcs-dev-kubernetes

spec:

	gateways:
	- istio-ingress/dev-planton-cloud
	hosts:
	- testing-main.dev.planton.cloud
	http:
		- name: default
		  route:
		  - destination:
			  host: testing.planton-pcs-dev-kubernetes.svc.cluster.local
			  port:
				number: 80
*/
func buildVirtualServiceObject(input *input, gatewayIngressDomainName, msiIngressDomain string) *v1beta1.VirtualService {
	return &v1beta1.VirtualService{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "VirtualService",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      zone.GetZoneName(msiIngressDomain),
			Namespace: input.namespaceName,
		},
		Spec: networkingv1beta1.VirtualService{
			Gateways: []string{fmt.Sprintf("%s/%s", ingressnamespace.Name, zone.GetZoneName(gatewayIngressDomainName))},
			Hosts:    []string{msiIngressDomain},
			Http: []*networkingv1beta1.HTTPRoute{{
				Name: input.microserviceVersion,
				Route: []*networkingv1beta1.HTTPRouteDestination{
					{
						Destination: &networkingv1beta1.Destination{
							Host: input.kubeServiceFqdn,
							Port: &networkingv1beta1.PortSelector{Number: uint32(input.servicePort)},
						},
					},
				},
			}},
		},
	}
}
