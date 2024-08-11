package pkg

import (
	"github.com/pkg/errors"
	certmanagerv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/certmanager/certmanager/v1"
	gatewayv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/gatewayapis/gateway/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func istioIngress(ctx *pulumi.Context, locals *Locals, kubernetesProvider *kubernetes.Provider,
	createdNamespace *kubernetescorev1.Namespace, labels map[string]string) error {
	//crate new certificate
	addedCertificate, err := certmanagerv1.NewCertificate(ctx,
		"ingress-certificate",
		&certmanagerv1.CertificateArgs{
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
		}, pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrap(err, "error creating certificate")
	}

	//create gateway
	createdGateway, err := gatewayv1.NewGateway(ctx,
		locals.MicroserviceKubernetes.Metadata.Id,
		&gatewayv1.GatewayArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name: pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
				//all gateway resources should be created in the istio-ingress deployment namespace
				Namespace: pulumi.String(vars.IstioIngressNamespace),
				Labels:    pulumi.ToStringMap(labels),
			},
			Spec: gatewayv1.GatewaySpecArgs{
				//the selector labels map should match the desired istio-ingress deployment.
				GatewayClassName: pulumi.String(vars.GatewayIngressClassName),
				Listeners: gatewayv1.GatewaySpecListenersArray{
					&gatewayv1.GatewaySpecListenersArgs{
						Name:     pulumi.String("https-external"),
						Hostname: pulumi.String(locals.IngressExternalHostname),
						Port:     pulumi.Int(443),
						Protocol: pulumi.String("HTTPS"),
						Tls: &gatewayv1.GatewaySpecListenersTlsArgs{
							Mode: pulumi.String("Terminate"),
							CertificateRefs: gatewayv1.GatewaySpecListenersTlsCertificateRefsArray{
								gatewayv1.GatewaySpecListenersTlsCertificateRefsArgs{
									Name: pulumi.String(locals.IngressCertSecretName),
								},
							},
						},
					},
					&gatewayv1.GatewaySpecListenersArgs{
						Name:     pulumi.String("https-internal"),
						Hostname: pulumi.String(locals.IngressInternalHostname),
						Port:     pulumi.Int(443),
						Protocol: pulumi.String("HTTPS"),
						Tls: &gatewayv1.GatewaySpecListenersTlsArgs{
							Mode: pulumi.String("Terminate"),
							CertificateRefs: gatewayv1.GatewaySpecListenersTlsCertificateRefsArray{
								gatewayv1.GatewaySpecListenersTlsCertificateRefsArgs{
									Name: pulumi.String(locals.IngressCertSecretName),
								},
							},
						},
					},
					&gatewayv1.GatewaySpecListenersArgs{
						Name:     pulumi.String("http-external"),
						Hostname: pulumi.String(locals.IngressExternalHostname),
						Port:     pulumi.Int(80),
						Protocol: pulumi.String("HTTP"),
					},
					&gatewayv1.GatewaySpecListenersArgs{
						Name:     pulumi.String("http-internal"),
						Hostname: pulumi.String(locals.IngressInternalHostname),
						Port:     pulumi.Int(80),
						Protocol: pulumi.String("HTTP"),
					},
				},
			},
		}, pulumi.Provider(kubernetesProvider), pulumi.DependsOn([]pulumi.Resource{addedCertificate}))
	if err != nil {
		return errors.Wrap(err, "error creating gateway")
	}

	var destinationServicePort = pulumi.Int(80)
	for _, p := range locals.MicroserviceKubernetes.Spec.Container.App.Ports {
		if p.IsIngressPort {
			destinationServicePort = pulumi.Int(p.ServicePort)
		}
	}

	//create http-route
	_, err = gatewayv1.NewHTTPRoute(ctx,
		"main",
		&gatewayv1.HTTPRouteArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:      pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(labels),
			},
			Spec: gatewayv1.HTTPRouteSpecArgs{
				Hostnames: pulumi.ToStringArray(locals.IngressHostnames),
				ParentRefs: gatewayv1.HTTPRouteSpecParentRefsArray{
					gatewayv1.HTTPRouteSpecParentRefsArgs{
						Name:      pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
						Namespace: createdGateway.Metadata.Namespace(),
					},
				},
				Rules: gatewayv1.HTTPRouteSpecRulesArray{
					gatewayv1.HTTPRouteSpecRulesArgs{
						Matches: gatewayv1.HTTPRouteSpecRulesMatchesArray{
							gatewayv1.HTTPRouteSpecRulesMatchesArgs{
								Path: gatewayv1.HTTPRouteSpecRulesMatchesPathArgs{
									Type:  pulumi.String("PathPrefix"),
									Value: pulumi.String("/"),
								},
							},
						},
						BackendRefs: gatewayv1.HTTPRouteSpecRulesBackendRefsArray{
							gatewayv1.HTTPRouteSpecRulesBackendRefsArgs{
								Name:      pulumi.String(locals.KubeServiceName),
								Namespace: createdNamespace.Metadata.Name(),
								Port:      destinationServicePort,
							},
						},
					},
				},
			},
		}, pulumi.Parent(createdNamespace))
	return nil
}
