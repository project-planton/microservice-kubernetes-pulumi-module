package pkg

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-module/pkg/outputs"
	microservicekubernetesmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/microservicekubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/kubernetesdockercredential/enums/dockerrepoprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	IngressCertClusterIssuerName string
	IngressCertSecretName        string
	IngressExternalHostname      string
	IngressHostnames             []string
	IngressInternalHostname      string
	KubePortForwardCommand       string
	KubeServiceFqdn              string
	KubeServiceName              string
	Namespace                    string
	MicroserviceKubernetes       *microservicekubernetesmodel.MicroserviceKubernetes
	ImagePullSecretData          map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *microservicekubernetesmodel.MicroserviceKubernetesStackInput) (*Locals, error) {
	locals := &Locals{}
	//assign value for the locals variable to make it available across the project
	locals.MicroserviceKubernetes = stackInput.ApiResource

	if stackInput.KubernetesDockerCredential != nil &&
		dockerrepoprovider.DockerRepoProvider_gcp_artifact_registry == stackInput.KubernetesDockerCredential.Spec.DockerRepoProvider {
		decodedStringBytes, err := b64.StdEncoding.DecodeString(stackInput.KubernetesDockerCredential.Spec.GcpArtifactRegistryCredentialSpec.GcpServiceAccountKeyBase64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode gcp service account key base64")
		}
		dockerConfigAuth := fmt.Sprintf("_json_key:%s", string(decodedStringBytes))

		dockerConfigAuth = b64.StdEncoding.EncodeToString([]byte(dockerConfigAuth))
		locals.ImagePullSecretData = map[string]string{".dockerconfigjson": fmt.Sprintf(`
			{
  				"auths": {
    				"%s": {
      					"username": "_json_key",
						"auth": "%s"
					}
  				}
			}`, stackInput.KubernetesDockerCredential.Spec.GcpArtifactRegistryCredentialSpec.DockerRepoHostname, dockerConfigAuth)}
	}

	microserviceKubernetes := stackInput.ApiResource

	//decide on the namespace
	locals.Namespace = microserviceKubernetes.Metadata.Id
	ctx.Export(outputs.Namespace, pulumi.String(locals.Namespace))

	locals.KubeServiceName = microserviceKubernetes.Spec.Version

	//export kubernetes service name
	ctx.Export(outputs.Service, pulumi.String(locals.KubeServiceName))

	locals.KubeServiceFqdn = fmt.Sprintf("%s.%s.svc.cluster.local", locals.KubeServiceName, locals.Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.KubeEndpoint, pulumi.String(locals.KubeServiceFqdn))

	locals.KubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s 8080:8080",
		locals.Namespace, locals.KubeServiceName)

	//export kube-port-forward command
	ctx.Export(outputs.KubePortForwardCommand, pulumi.String(locals.KubePortForwardCommand))

	if microserviceKubernetes.Spec.Ingress == nil ||
		!microserviceKubernetes.Spec.Ingress.IsEnabled ||
		microserviceKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return locals, nil
	}

	locals.IngressExternalHostname = fmt.Sprintf("%s.%s", microserviceKubernetes.Metadata.Id,
		microserviceKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressInternalHostname = fmt.Sprintf("%s-internal.%s", microserviceKubernetes.Metadata.Id,
		microserviceKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHostnames = []string{
		locals.IngressExternalHostname,
		locals.IngressInternalHostname,
	}

	//export ingress hostnames
	ctx.Export(outputs.IngressExternalHostname, pulumi.String(locals.IngressExternalHostname))
	ctx.Export(outputs.IngressInternalHostname, pulumi.String(locals.IngressInternalHostname))

	//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
	//this is typically taken care of by the kubernetes cluster administrator.
	//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
	//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
	//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
	locals.IngressCertClusterIssuerName = microserviceKubernetes.Spec.Ingress.EndpointDomainName

	locals.IngressCertSecretName = microserviceKubernetes.Metadata.Id

	return locals, nil
}