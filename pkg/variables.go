package pkg

var vars = struct {
	GatewayIngressClassName                    string
	GatewayInternalLoadBalancerServiceHostname string
	GatewayExternalLoadBalancerServiceHostname string
	IstioIngressNamespace                      string
	MainPyConfigMapName                        string
	LibFilesConfigMapName                      string
}{
	GatewayIngressClassName:                    "istio",
	GatewayInternalLoadBalancerServiceHostname: "ingress-internal.istio-ingress.svc.cluster.local",
	GatewayExternalLoadBalancerServiceHostname: "ingress-external.istio-ingress.svc.cluster.local",
	IstioIngressNamespace:                      "istio-ingress",
	MainPyConfigMapName:                        "main-py",
	LibFilesConfigMapName:                      "lib-files",
}
