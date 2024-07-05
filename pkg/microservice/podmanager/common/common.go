package common

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/datatypes/maps"
	"strings"

	kubernetesv1model "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	code2cloudv1deploymsimodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/microservicekubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// BuildContainers returns list of containers that are to be added to the deployment
// todo: the design needs to be reviewed as there is too much marshalling an unmarshalling between planton, kubernetes and pulumi objects.
func BuildContainers(ctx *pulumi.Context) []corev1.ContainerInput {
	i := extractInput(ctx)
	resp := make([]corev1.ContainerInput, 0)
	//add main container
	resp = append(resp, corev1.ContainerInput(
		&corev1.ContainerArgs{
			Name:  pulumi.String(MicroserviceContainerName),
			Image: pulumi.String(getContainerImage(i.appContainer.Image)),
			Env:   corev1.EnvVarArray(buildEnvVars(i.version, i.envSpec)),
			Ports: getContainerPorts(i.appPorts),
			Resources: corev1.ResourceRequirementsArgs{
				Limits: pulumi.ToStringMap(map[string]string{
					englishword.EnglishWord_cpu.String():    i.appContainer.Resources.Limits.Cpu,
					englishword.EnglishWord_memory.String(): i.appContainer.Resources.Limits.Memory,
				}),
				Requests: pulumi.ToStringMap(map[string]string{
					englishword.EnglishWord_cpu.String():    i.appContainer.Resources.Requests.Cpu,
					englishword.EnglishWord_memory.String(): i.appContainer.Resources.Requests.Memory,
				}),
			},
			//note: VolumeMounts should only be added for stateful set containers
			VolumeMounts: getVolumeMounts(i.appContainer.VolumeMounts),
			Lifecycle: corev1.LifecycleArgs{
				PreStop: corev1.LifecycleHandlerArgs{
					Exec: corev1.ExecActionArgs{
						//wait for 60 seconds before killing the main process
						//this is particularly useful and required when deploying build and stack microservices of planton cloud service in production.
						Command: pulumi.ToStringArray([]string{"/bin/sleep", "60"}),
					},
				},
			},
		}))
	//add sidecar containers
	for _, s := range i.sidecars {
		sideCarEnvVars := make([]corev1.EnvVarInput, 0)
		//this does not currently support reading from kubernetes secrets
		for _, envVar := range s.Env {
			sideCarEnvVars = append(sideCarEnvVars, corev1.EnvVarInput(corev1.EnvVarArgs{
				Name:  pulumi.String(envVar.Name),
				Value: pulumi.String(envVar.Value),
			}))
		}
		sideCarPorts := make([]corev1.ContainerPortInput, 0)
		//this does not currently support reading from kubernetes secrets
		for _, port := range s.Ports {
			sideCarPorts = append(sideCarPorts, corev1.ContainerPortInput(corev1.ContainerPortArgs{
				Name:          pulumi.String(port.Name),
				Protocol:      pulumi.String(port.Protocol),
				ContainerPort: pulumi.Int(port.ContainerPort),
			}))
		}
		resp = append(resp, corev1.ContainerInput(
			&corev1.ContainerArgs{
				Name:  pulumi.String(s.Name),
				Image: pulumi.String(s.Image),
				Ports: corev1.ContainerPortArray(sideCarPorts),
				Env:   corev1.EnvVarArray(sideCarEnvVars),
			}))
	}
	return resp
}

func getContainerPorts(ports []*code2cloudv1deploymsimodel.MicroserviceKubernetesSpecContainerSpecAppSpecPortSpec) corev1.ContainerPortArray {
	portsArray := make(corev1.ContainerPortArray, 0)
	for _, p := range ports {
		portsArray = append(portsArray, &corev1.ContainerPortArgs{
			Name:          pulumi.String(p.Name),
			ContainerPort: pulumi.Int(p.ContainerPort),
		})
	}
	return portsArray
}

func getContainerImage(image *kubernetesv1model.ContainerImage) string {
	return fmt.Sprintf("%s:%s", image.Repo, image.Tag)
}

func getVolumeMounts(volumeMounts []*code2cloudv1deploymsimodel.MicroserviceKubernetesSpecContainerSpecAppSpecVolumeMountSpec) corev1.VolumeMountArray {
	resp := make(corev1.VolumeMountArray, 0)
	for _, v := range volumeMounts {
		resp = append(resp, corev1.VolumeMountArgs{
			Name:      pulumi.String(v.Name),
			MountPath: pulumi.String(v.MountPath),
		})
	}
	return resp
}

/*
  - name: KFK_SASL_USERNAME
    valueFrom:
    secretKeyRef:
    name: gql-master-dev
    key: pcs-kfk-sasl-username

  - name: KFK_SASL_PASSWORD
    valueFrom:
    secretKeyRef:
    name: gql-master-dev
    key: pcs-kfk-sasl-password

  - name: HOSTNAME
    valueFrom:
    fieldRef:
    fieldPath: status.podIP

  - name: K8S_POD_ID
    valueFrom:
    fieldRef:
    apiVersion: v1
    fieldPath: metadata.name
*/
func buildEnvVars(version string, containerEnv *code2cloudv1deploymsimodel.MicroserviceKubernetesSpecContainerSpecAppSpecEnvSpec) []corev1.EnvVarInput {
	resp := make([]corev1.EnvVarInput, 0)
	//add HOSTNAME env var
	resp = append(resp, corev1.EnvVarInput(corev1.EnvVarArgs{
		Name: pulumi.String("HOSTNAME"),
		ValueFrom: &corev1.EnvVarSourceArgs{
			FieldRef: &corev1.ObjectFieldSelectorArgs{
				FieldPath: pulumi.String("status.podIP"),
			},
		},
	}))
	//add K8S_POD_ID env var
	resp = append(resp, corev1.EnvVarInput(corev1.EnvVarArgs{
		Name: pulumi.String("K8S_POD_ID"),
		ValueFrom: &corev1.EnvVarSourceArgs{
			FieldRef: &corev1.ObjectFieldSelectorArgs{
				ApiVersion: pulumi.String("v1"),
				FieldPath:  pulumi.String("metadata.name"),
			},
		},
	}))

	sortedEnvVariableKeys := maps.SortMapKeys(containerEnv.Variables)

	for _, environmentVariableKey := range sortedEnvVariableKeys {
		resp = append(resp, corev1.EnvVarInput(corev1.EnvVarArgs{
			Name:  pulumi.String(environmentVariableKey),
			Value: pulumi.String(containerEnv.Variables[environmentVariableKey]),
		}))
	}

	sortedEnvironmentSecretKeys := maps.SortMapKeys(containerEnv.Secrets)

	for _, environmentSecretKey := range sortedEnvironmentSecretKeys {
		resp = append(resp, corev1.EnvVarInput(corev1.EnvVarArgs{
			Name: pulumi.String(environmentSecretKey),
			ValueFrom: &corev1.EnvVarSourceArgs{
				SecretKeyRef: &corev1.SecretKeySelectorArgs{
					Name: pulumi.String(version),
					Key:  pulumi.String(environmentSecretKey),
				},
			},
		}))
	}

	return resp
}

func GetNormalizedMountPath(mountPath string) string {
	return strings.ReplaceAll(mountPath, "/", "-")
}
