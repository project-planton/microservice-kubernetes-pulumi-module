package statefulset

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/podmanager/common"
	code2cloudv1deploymsimodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	if err := addStatefulSet(ctx); err != nil {
		return errors.Wrap(err, "failed to add deployment")
	}
	return nil
}

func addStatefulSet(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	_, err := appsv1.NewStatefulSet(ctx, i.version, &appsv1.StatefulSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(i.resourceName),
			Namespace: pulumi.String(i.namespaceName),
			Labels:    pulumi.ToStringMap(i.labels),
			Annotations: pulumi.StringMap{
				"pulumi.com/patchForce": pulumi.String("true"),
			}},
		Spec: &appsv1.StatefulSetSpecArgs{
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
			VolumeClaimTemplates: getPersistentVolumeClaimTemplates(i.appContainer.VolumeMounts, i.labels),
		},
	}, pulumi.Provider(i.kubernetesProvider),
		pulumi.IgnoreChanges([]string{
			"status",
		}))
	if err != nil {
		return errors.Wrap(err, "failed to add deployment")
	}
	return nil
}

func getPersistentVolumeClaimTemplates(volumeMounts []*code2cloudv1deploymsimodel.MicroserviceKubernetesSpecContainerSpecAppSpecVolumeMountSpec, labels map[string]string) corev1.PersistentVolumeClaimTypeArray {
	resp := make(corev1.PersistentVolumeClaimTypeArray, 0)
	for _, v := range volumeMounts {
		resp = append(resp, corev1.PersistentVolumeClaimTypeArgs{
			ApiVersion: pulumi.String("v1"),
			Kind:       pulumi.String("PersistentVolumeClaim"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(v.Name),
				Labels: pulumi.ToStringMap(labels),
			},
			Spec: corev1.PersistentVolumeClaimSpecArgs{
				AccessModes: pulumi.ToStringArray([]string{"ReadWriteOnce"}),
				Resources: corev1.VolumeResourceRequirementsArgs{
					Requests: pulumi.ToStringMap(map[string]string{"storage": v.Size}),
				},
			},
		})
	}
	return resp
}
