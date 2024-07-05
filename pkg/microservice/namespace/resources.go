package namespace

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system/istiod"
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	namespace, err := addNamespace(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}

	var ctxConfig = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	addNamespaceToContext(&ctxConfig, namespace)
	ctx = ctx.WithValue(microservicecontextstate.Key, ctxConfig)
	return ctx, nil
}

func addNamespace(ctx *pulumi.Context) (*kubernetescorev1.Namespace, error) {
	var i = extractInput(ctx)

	ns, err := kubernetescorev1.NewNamespace(ctx, i.namespaceName, &kubernetescorev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Namespace"),
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:   pulumi.String(i.namespaceName),
			Labels: pulumi.ToStringMap(addIstioInjectionLabel(i.labels)),
		}),
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(i.kubeProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s namespace", i.namespaceName)
	}
	return ns, nil
}

func addIstioInjectionLabel(inputLabels map[string]string) map[string]string {
	labels := inputLabels
	labels[istiod.SidecarInjectionLabelKey] = istiod.SidecarInjectionLabelValue
	return labels
}

func addNamespaceToContext(existingConfig *microservicecontextstate.ContextState, namespace *kubernetescorev1.Namespace) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &microservicecontextstate.AddedResources{
			Namespace: namespace,
		}
		return
	}
	existingConfig.Status.AddedResources.Namespace = namespace
}
