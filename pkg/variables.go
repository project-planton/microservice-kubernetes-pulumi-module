package pkg

var vars = struct {
	GatewayIngressClassName    string
	IstioIngressNamespace      string
	IstioIngressSelectorLabels map[string]string
	MainPyConfigMapName        string
	LibFilesConfigMapName      string
}{
	GatewayIngressClassName: "istio",
	IstioIngressNamespace:   "istio-ingress",
	IstioIngressSelectorLabels: map[string]string{
		"app":   "gateway",
		"istio": "gateway",
	},
	MainPyConfigMapName:   "main-py",
	LibFilesConfigMapName: "lib-files",
}
