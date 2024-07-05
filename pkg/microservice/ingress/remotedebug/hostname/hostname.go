package hostname

import "fmt"

// GetHostnames to be included in the self-signed certificate, gateway and virtual-service resources
func GetHostnames(microserviceKubernetesId, productEnvName, endpointDomainName string) []string {
	hostnames := make([]string, 0)
	hostnames = append(hostnames,
		fmt.Sprintf("%s.%s.%s", microserviceKubernetesId, productEnvName, endpointDomainName))
	return hostnames
}

func GetInternalHostname(microserviceKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s-internal.%s", microserviceKubernetesId, environmentName, endpointDomainName)
}

func GetExternalHostname(microserviceKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s.%s", microserviceKubernetesId, environmentName, endpointDomainName)
}
