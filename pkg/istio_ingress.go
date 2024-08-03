package pkg

import (
	"github.com/pkg/errors"
	certmanagerv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/certmanager/certmanager/v1"
	istiov1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/istio/networking/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	v1 "istio.io/api/networking/v1"
)

func istioIngress(ctx *pulumi.Context, locals *Locals, createdNamespace *kubernetescorev1.Namespace, labels map[string]string) error {
	//crate new certificate
	addedCertificate, err := certmanagerv1.NewCertificate(ctx, "ingress-certificate", &certmanagerv1.CertificateArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
			Namespace: pulumi.String(vars.IstioIngressNamespace),
			Labels:    pulumi.ToStringMap(labels),
		},
		Spec: certmanagerv1.CertificateSpecArgs{
			DnsNames:   pulumi.ToStringArray(locals.IngressHostnames),
			SecretName: pulumi.String(locals.IngressCertSecretName),
			IssuerRef: certmanagerv1.CertificateSpecIssuerRefArgs{
				Kind: pulumi.String("ClusterIssuer"),
				Name: pulumi.String(locals.IngressCertClusterIssuerName),
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "error creating certificate")
	}

	//create gateway
	_, err = istiov1.NewGateway(ctx, locals.MicroserviceKubernetes.Metadata.Id, &istiov1.GatewayArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
			//all gateway resources should be created in the istio-ingress deployment namespace
			Namespace: pulumi.String(vars.IstioIngressNamespace),
			Labels:    pulumi.ToStringMap(labels),
		},
		Spec: istiov1.GatewaySpecArgs{
			//the selector labels map should match the desired istio-ingress deployment.
			Selector: pulumi.ToStringMap(vars.IstioIngressSelectorLabels),
			Servers: istiov1.GatewaySpecServersArray{
				&istiov1.GatewaySpecServersArgs{
					Name: pulumi.String("https"),
					Port: &istiov1.GatewaySpecServersPortArgs{
						Number:   pulumi.Int(443),
						Name:     pulumi.String("https"),
						Protocol: pulumi.String("HTTPS"),
					},
					Hosts: pulumi.ToStringArray(locals.IngressHostnames),
					Tls: &istiov1.GatewaySpecServersTlsArgs{
						CredentialName: addedCertificate.Spec.SecretName(),
						Mode:           pulumi.String(v1.ServerTLSSettings_SIMPLE.String()),
					},
				},
				&istiov1.GatewaySpecServersArgs{
					Name: pulumi.String("http"),
					Port: &istiov1.GatewaySpecServersPortArgs{
						Number:   pulumi.Int(80),
						Name:     pulumi.String("http"),
						Protocol: pulumi.String("HTTP"),
					},
					Hosts: pulumi.ToStringArray(locals.IngressHostnames),
					Tls: &istiov1.GatewaySpecServersTlsArgs{
						HttpsRedirect: pulumi.Bool(true),
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "error creating gateway")
	}

	var destinationServicePort = pulumi.Int(80)
	for _, p := range locals.MicroserviceKubernetes.Spec.Container.App.Ports {
		if p.IsIngressPort {
			destinationServicePort = pulumi.Int(p.ServicePort)
		}
	}

	//create virtual-service
	_, err = istiov1.NewVirtualService(ctx, locals.MicroserviceKubernetes.Metadata.Id, &istiov1.VirtualServiceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(labels),
		},
		Spec: istiov1.VirtualServiceSpecArgs{
			Gateways: pulumi.StringArray{
				pulumi.Sprintf("%s/%s", vars.IstioIngressNamespace, locals.MicroserviceKubernetes.Metadata.Id),
			},
			Hosts: pulumi.ToStringArray(locals.IngressHostnames),
			Http: istiov1.VirtualServiceSpecHttpArray{
				&istiov1.VirtualServiceSpecHttpArgs{
					Name: pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
					Route: istiov1.VirtualServiceSpecHttpRouteArray{
						&istiov1.VirtualServiceSpecHttpRouteArgs{
							Destination: istiov1.VirtualServiceSpecHttpRouteDestinationArgs{
								Host: pulumi.String(locals.KubeServiceFqdn),
								Port: istiov1.VirtualServiceSpecHttpRouteDestinationPortArgs{
									Number: destinationServicePort,
								},
							},
						},
					},
				},
			},
		},
		Status: nil,
	})
	if err != nil {
		return errors.Wrap(err, "error creating virtual-service")
	}
	return nil
}
