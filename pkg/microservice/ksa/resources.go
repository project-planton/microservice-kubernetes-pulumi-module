package ksa

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/workloadidentity"
	pulk8scv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	v12 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	if err := addKsa(ctx); err != nil {
		return errors.Wrap(err, "failed to add kubernetes service account")
	}
	return nil
}

func addKsa(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	_, err := pulk8scv1.NewServiceAccount(ctx, i.namespaceName, &pulk8scv1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: v12.ObjectMetaPtrInput(&v12.ObjectMetaArgs{
			Name:      pulumi.String(i.namespaceName),
			Namespace: i.namespace.Metadata.Name(),
			Annotations: pulumi.StringMap{
				workloadidentity.WorkloadIdentityKubeAnnotationKey: pulumi.Sprintf("%s@%s.iam.gserviceaccount.com",
					pulumi.String(i.containerClusterProject.Id), i.workloadIdentityGsaAccountId.ID()),
			},
		}),
	}, pulumi.Parent(i.namespace))
	if err != nil {
		return errors.Wrap(err, "failed to add service account")
	}
	return nil
}
