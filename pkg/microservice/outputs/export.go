package outputs

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/kubernetes/pulumikubernetesprovider"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/pulumi/pulumicustomoutput"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Export(ctx *pulumi.Context) error {
	var i = extractInput(ctx)
	var kubePortForwardCommand = getPortForwardCommand(i.version, i.namespaceName, i.listenerServicePort, i.forwardServicePort)

	ctx.Export(GetExternalClusterHostnameOutputName(), pulumi.String(i.externalHostname))
	ctx.Export(GetInternalClusterHostnameOutputName(), pulumi.String(i.internalHostname))

	ctx.Export(GetServiceNameOutputName(), pulumi.String(i.kubeServiceName))
	ctx.Export(GetFqdnOutputName(), pulumi.String(i.kubeLocalEndpoint))

	ctx.Export(GetPortForwardCommandOutputName(), pulumi.String(kubePortForwardCommand))
	ctx.Export(GetNamespaceNameOutputName(), pulumi.String(i.namespaceName))

	for _, imagePullSecretInput := range i.kubernetesImagePullSecretInputs {
		ctx.Export(GetSecretNameOutputName(imagePullSecretInput.ImagePullSecretName), pulumi.String(imagePullSecretInput.ImagePullSecretName))
	}

	return nil
}

func GetExternalClusterHostnameOutputName() string {
	return pulumicustomoutput.Name("external-hostname")
}

func GetInternalClusterHostnameOutputName() string {
	return pulumicustomoutput.Name("internal-hostname")
}

func GetSecretNameOutputName(imagePullSecretName string) string {
	return pulumikubernetesprovider.PulumiOutputName(kubernetescorev1.Secret{}, imagePullSecretName, englishword.EnglishWord_name.String())
}

func GetFqdnOutputName() string {
	return pulumikubernetesprovider.PulumiOutputName(kubernetescorev1.Service{}, englishword.EnglishWord_endpoint.String())
}

func GetPortForwardCommandOutputName() string {
	return pulumikubernetesprovider.PulumiOutputName(kubernetescorev1.Service{}, "port-forward-command")
}

func GetServiceNameOutputName() string {
	return pulumikubernetesprovider.PulumiOutputName(kubernetescorev1.Service{}, englishword.EnglishWord_name.String())
}

func GetNamespaceNameOutputName() string {
	return pulumikubernetesprovider.PulumiOutputName(kubernetescorev1.Namespace{}, englishword.EnglishWord_namespace.String())
}

// GetPortForwardCommand for the service created for the microservice deployment
// ex: kubectl port-forward service/main 80:8080 -n planton-pcs-dev-product
func getPortForwardCommand(microserviceInstanceVersion, namespaceName string, containerPort, servicePort int32) string {
	return fmt.Sprintf("kubectl port-forward service/%s %d:%d -n %s",
		microserviceInstanceVersion,
		servicePort,
		containerPort,
		namespaceName)
}

func GetGsaEmailOutputName(namespaceName string) string {
	return pulumikubernetesprovider.PulumiOutputName(serviceaccount.Account{}, namespaceName)
}
