package deployment

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/podmanager/common"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	if err := addDeployment(ctx); err != nil {
		return errors.Wrap(err, "failed to add deployment")
	}
	return nil
}

func addDeployment(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	_, err := appsv1.NewDeployment(ctx, i.version, &appsv1.DeploymentArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(i.resourceName),
			Namespace: pulumi.String(i.namespaceName),
			Labels:    pulumi.ToStringMap(i.labels),
			Annotations: pulumi.StringMap{
				"pulumi.com/patchForce": pulumi.String("true"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(i.availabilitySpec.MinReplicas),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.ToStringMap(i.labels),
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.ToStringMap(i.labels),
				},
				Spec: &corev1.PodSpecArgs{
					ServiceAccountName: pulumi.String(i.namespaceName),
					ImagePullSecrets: corev1.LocalObjectReferenceArray{corev1.LocalObjectReferenceArgs{
						Name: pulumi.String(i.appContainer.Image.PullSecretName),
					}},
					Containers: corev1.ContainerArray(common.BuildContainers(ctx)),
					//wait for 60 seconds before sending the termination signal to the processes in the pod
					TerminationGracePeriodSeconds: pulumi.IntPtr(60),
				},
			},
		},
	}, pulumi.Provider(i.kubernetesProvider), pulumi.IgnoreChanges([]string{
		"status",
	}))
	if err != nil {
		return errors.Wrap(err, "failed to add deployment")
	}
	return nil
}
