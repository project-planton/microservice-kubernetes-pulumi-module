package workloadidentity

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	pulk8scv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	if !i.isWorkloadIdentityEnabled {
		return nil
	}
	addedGsa, err := addGsa(ctx, i)
	if err != nil {
		return errors.Wrap(err, "failed to add google service account")
	}
	if err := addWorkloadIdentityBinding(ctx, i, addedGsa); err != nil {
		return errors.Wrap(err, "failed to add workload identity binding")
	}
	return nil
}

func addGsa(ctx *pulumi.Context, i *input) (*serviceaccount.Account, error) {
	gsa, err := serviceaccount.NewAccount(ctx, i.microserviceKubernetesId,
		&serviceaccount.AccountArgs{
			Project:     pulumi.String(i.containerClusterProject.Id),
			Description: pulumi.Sprintf("workload identity for %s", i.namespaceName),
			AccountId:   i.workloadIdentityGsaAccountId.ID(),
			DisplayName: pulumi.Sprintf("workload identity for %s kubernetes namespace", i.namespaceName),
		}, pulumi.Provider(i.gcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed add google service account")
	}
	return gsa, nil
}

func addWorkloadIdentityBinding(ctx *pulumi.Context, i *input, gsa *serviceaccount.Account) error {
	_, err := serviceaccount.NewIAMBinding(ctx,
		fmt.Sprintf("%s-workload-identity", i.microserviceKubernetesId), &serviceaccount.IAMBindingArgs{
			ServiceAccountId: gsa.Name,
			Role:             pulumi.String(standard.Iam_workloadIdentityUser),
			Members:          pulumi.StringArray(getMembers(i.containerClusterProject.Id, i.microserviceKubernetesId, i.namespace)),
		}, pulumi.Parent(gsa))
	if err != nil {
		return errors.Wrapf(err, "failed to add workload identity binding")
	}
	return nil
}

func getMembers(gcpProjectId, ksaName string, addedNamespace *pulk8scv1.Namespace) []pulumi.StringInput {
	return []pulumi.StringInput{
		pulumi.Sprintf("serviceAccount:%s.svc.id.goog[%s/%s]", gcpProjectId, addedNamespace.Metadata.Name().Elem(), ksaName),
	}
}
