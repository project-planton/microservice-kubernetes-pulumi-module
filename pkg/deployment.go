package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/datatypes/maps"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployment(ctx *pulumi.Context, locals *Locals, createdNamespace *kubernetescorev1.Namespace, labels map[string]string) (*appsv1.Deployment, error) {

	// create image pull secret resources
	_, err := kubernetescorev1.NewSecret(ctx, locals.MicroserviceKubernetes.Spec.KubernetesDockerCredentialId, &kubernetescorev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.MicroserviceKubernetes.Spec.KubernetesDockerCredentialId),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(labels),
		},
		Type:       pulumi.String("kubernetes.io/dockerconfigjson"),
		StringData: pulumi.ToStringMap(locals.ImagePullSecretData),
	}, pulumi.Parent(createdNamespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add image pull secret")
	}

	// create service account
	createdServiceAccount, err := kubernetescorev1.NewServiceAccount(ctx, locals.MicroserviceKubernetes.Metadata.Id, &kubernetescorev1.ServiceAccountArgs{
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.MicroserviceKubernetes.Metadata.Id),
			Namespace: createdNamespace.Metadata.Name(),
		}),
	}, pulumi.Parent(createdNamespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add service account")
	}

	envVarInputs := make([]kubernetescorev1.EnvVarInput, 0)
	//add HOSTNAME env var
	envVarInputs = append(envVarInputs, kubernetescorev1.EnvVarInput(kubernetescorev1.EnvVarArgs{
		Name: pulumi.String("HOSTNAME"),
		ValueFrom: &kubernetescorev1.EnvVarSourceArgs{
			FieldRef: &kubernetescorev1.ObjectFieldSelectorArgs{
				FieldPath: pulumi.String("status.podIP"),
			},
		},
	}))
	//add K8S_POD_ID env var
	envVarInputs = append(envVarInputs, kubernetescorev1.EnvVarInput(kubernetescorev1.EnvVarArgs{
		Name: pulumi.String("K8S_POD_ID"),
		ValueFrom: &kubernetescorev1.EnvVarSourceArgs{
			FieldRef: &kubernetescorev1.ObjectFieldSelectorArgs{
				ApiVersion: pulumi.String("v1"),
				FieldPath:  pulumi.String("metadata.name"),
			},
		},
	}))

	sortedEnvVariableKeys := maps.SortMapKeys(locals.MicroserviceKubernetes.Spec.Container.App.Env.Variables)

	for _, environmentVariableKey := range sortedEnvVariableKeys {
		envVarInputs = append(envVarInputs, kubernetescorev1.EnvVarInput(kubernetescorev1.EnvVarArgs{
			Name:  pulumi.String(environmentVariableKey),
			Value: pulumi.String(locals.MicroserviceKubernetes.Spec.Container.App.Env.Variables[environmentVariableKey]),
		}))
	}

	sortedEnvironmentSecretKeys := maps.SortMapKeys(locals.MicroserviceKubernetes.Spec.Container.App.Env.Secrets)

	for _, environmentSecretKey := range sortedEnvironmentSecretKeys {
		envVarInputs = append(envVarInputs, kubernetescorev1.EnvVarInput(kubernetescorev1.EnvVarArgs{
			Name: pulumi.String(environmentSecretKey),
			ValueFrom: &kubernetescorev1.EnvVarSourceArgs{
				SecretKeyRef: &kubernetescorev1.SecretKeySelectorArgs{
					Name: pulumi.String(locals.MicroserviceKubernetes.Spec.Version),
					Key:  pulumi.String(environmentSecretKey),
				},
			},
		}))
	}

	portsArray := make(kubernetescorev1.ContainerPortArray, 0)
	for _, p := range locals.MicroserviceKubernetes.Spec.Container.App.Ports {
		portsArray = append(portsArray, &kubernetescorev1.ContainerPortArgs{
			Name:          pulumi.String(p.Name),
			ContainerPort: pulumi.Int(p.ContainerPort),
		})
	}

	containerInputs := make([]kubernetescorev1.ContainerInput, 0)
	//add main container
	containerInputs = append(containerInputs, kubernetescorev1.ContainerInput(
		&kubernetescorev1.ContainerArgs{
			Name:  pulumi.String("microservice"),
			Image: pulumi.String(fmt.Sprintf("%s:%s", locals.MicroserviceKubernetes.Spec.Container.App.Image.Repo, locals.MicroserviceKubernetes.Spec.Container.App.Image.Tag)),
			Env:   kubernetescorev1.EnvVarArray(envVarInputs),
			Ports: portsArray,
			Resources: kubernetescorev1.ResourceRequirementsArgs{
				Limits: pulumi.ToStringMap(map[string]string{
					englishword.EnglishWord_cpu.String():    locals.MicroserviceKubernetes.Spec.Container.App.Resources.Limits.Cpu,
					englishword.EnglishWord_memory.String(): locals.MicroserviceKubernetes.Spec.Container.App.Resources.Limits.Memory,
				}),
				Requests: pulumi.ToStringMap(map[string]string{
					englishword.EnglishWord_cpu.String():    locals.MicroserviceKubernetes.Spec.Container.App.Resources.Requests.Cpu,
					englishword.EnglishWord_memory.String(): locals.MicroserviceKubernetes.Spec.Container.App.Resources.Requests.Memory,
				}),
			},
			Lifecycle: kubernetescorev1.LifecycleArgs{
				PreStop: kubernetescorev1.LifecycleHandlerArgs{
					Exec: kubernetescorev1.ExecActionArgs{
						//wait for 60 seconds before killing the main process
						//this is particularly useful and required when deploying build and stack microservices of planton cloud service in production.
						Command: pulumi.ToStringArray([]string{"/bin/sleep", "60"}),
					},
				},
			},
		}))

	//create deployment
	createdDeployment, err := appsv1.NewDeployment(ctx,
		locals.MicroserviceKubernetes.Spec.Version,
		&appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(locals.MicroserviceKubernetes.Metadata.Name),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(labels),
				Annotations: pulumi.StringMap{
					"pulumi.com/patchForce": pulumi.String("true"),
				},
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(locals.MicroserviceKubernetes.Spec.Availability.MinReplicas),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.ToStringMap(labels),
				},
				Template: &kubernetescorev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.ToStringMap(labels),
					},
					Spec: &kubernetescorev1.PodSpecArgs{
						ServiceAccountName: createdServiceAccount.Metadata.Name(),
						ImagePullSecrets: kubernetescorev1.LocalObjectReferenceArray{kubernetescorev1.LocalObjectReferenceArgs{
							Name: pulumi.String(locals.MicroserviceKubernetes.Spec.KubernetesDockerCredentialId),
						}},
						Containers: kubernetescorev1.ContainerArray(containerInputs),
						//wait for 60 seconds before sending the termination signal to the processes in the pod
						TerminationGracePeriodSeconds: pulumi.IntPtr(60),
					},
				},
			},
		}, pulumi.Parent(createdNamespace), pulumi.IgnoreChanges([]string{
			//WARNING: adding metdata.managedFields to ignoreChanges is rejected from kubernetes api-server for some reason
			//although the issue must have been resolved by now,per, https://github.com/pulumi/pulumi-kubernetes/issues/1075,
			//apparently it is not.
			//error from the api-server is "metadata.managedFields must be nil"
			//"metadata.managedFields", "status",
		}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add deployment")
	}

	return createdDeployment, nil
}
