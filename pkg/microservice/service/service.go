package service

import (
	"fmt"
	"github.com/pkg/errors"
	kubedns "github.com/plantoncloud-inc/go-commons/kubernetes/network/dns"
	microservicestatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	_, err := addService(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add kubernetes service")
	}
	return nil
}

func addService(ctx *pulumi.Context) (*corev1.Service, error) {
	i := extractInput(ctx)
	svc, err := corev1.NewService(ctx, i.version, &corev1.ServiceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(i.version),
			Namespace: pulumi.String(i.namespaceName),
			Labels:    pulumi.ToStringMap(i.labels),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type:     pulumi.String("ClusterIP"),
			Selector: pulumi.ToStringMap(i.labels),
			Ports:    getServicePorts(i.appPorts),
		},
	}, pulumi.Provider(i.kubernetesProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add service")
	}
	return svc, nil
}

func getServicePorts(appContainerPorts []*microservicestatemodel.MicroserviceKubernetesSpecContainerSpecAppSpecPortSpec) corev1.ServicePortArray {
	portsArray := make(corev1.ServicePortArray, 0)
	for _, p := range appContainerPorts {
		portsArray = append(portsArray, &corev1.ServicePortArgs{
			Name:        pulumi.String(p.Name),
			Protocol:    pulumi.String(p.NetworkProtocol),
			Port:        pulumi.Int(p.ServicePort),
			TargetPort:  pulumi.Int(p.ContainerPort),
			AppProtocol: pulumi.String(p.AppProtocol),
		})
	}
	return portsArray
}

// GetFqdn for the service created for the microservice deployment
// ex: main.planton-pcs-dev-product.svc.cluster.local
func GetFqdn(microserviceInstanceVersion, namespaceName string) string {
	return fmt.Sprintf("%s.%s.%s", microserviceInstanceVersion, namespaceName, kubedns.DefaultDomain)
}
