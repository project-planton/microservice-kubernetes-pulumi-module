package gsa

import (
	"fmt"
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	i := extractInput(ctx)
	if !i.isWorkloadIdentityEnabled {
		return ctx, nil
	}
	gsaAccountId, err := generateGsaAccountId(ctx, i)
	if err != nil {
		return nil, fmt.Errorf("failed to generate gsa account id value: %w", err)
	}

	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	addWorkLoadIdentityGsaAccountIdToContext(&contextState, gsaAccountId, i.containerClusterProject.Id)
	ctx = ctx.WithValue(microservicecontextstate.Key, contextState)
	return ctx, nil
}

func generateGsaAccountId(ctx *pulumi.Context, i *input) (*random.RandomId, error) {
	prefix := i.microserviceKubernetesName
	if len(prefix) > 18 {
		prefix = prefix[:18]
	}

	gsaAccountId, err := random.NewRandomId(ctx, "generate-gsa-account-id", &random.RandomIdArgs{
		ByteLength: pulumi.Int(5),
		Prefix:     pulumi.String(fmt.Sprintf("%s-", prefix)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate gsa account id value: %w", err)
	}
	return gsaAccountId, nil
}

func addWorkLoadIdentityGsaAccountIdToContext(existingConfig *microservicecontextstate.ContextState,
	workLoadIdentityGsaAccountId *random.RandomId,
	containerClusterProjectId string) {
	var gsaEmailId = pulumi.Sprintf("%s@%s.iam.gserviceaccount.com", pulumi.String(containerClusterProjectId), workLoadIdentityGsaAccountId.Hex)
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &microservicecontextstate.AddedResources{
			GsaEmailId:                   gsaEmailId,
			WorkloadIdentityGsaAccountId: workLoadIdentityGsaAccountId,
		}
		return
	}
	existingConfig.Status.AddedResources.GsaEmailId = gsaEmailId
	existingConfig.Status.AddedResources.WorkloadIdentityGsaAccountId = workLoadIdentityGsaAccountId
}
