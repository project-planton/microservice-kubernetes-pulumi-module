package pkg

var vars = struct {
	IstioIngressNamespace      string
	IstioIngressSelectorLabels map[string]string
	MainPyConfigMapName        string
	LibFilesConfigMapName      string
}{
	IstioIngressNamespace: "istio-ingress",
	IstioIngressSelectorLabels: map[string]string{
		"app":   "istio-ingress",
		"istio": "ingress",
	},
	MainPyConfigMapName:   "main-py",
	LibFilesConfigMapName: "lib-files",
}
