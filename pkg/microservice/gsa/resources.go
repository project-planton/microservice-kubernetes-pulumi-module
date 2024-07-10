package gsa

import (
	"fmt"
	microservicecontextstate "github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/contextstate"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	i := extractInput(ctx)
	prefix := i.microserviceKubernetesName

	if len(prefix) > 17 {
		prefix = prefix[:17]
	}

	gsaAccountId, err := random.NewRandomId(ctx, "generate-gsa-account-id", &random.RandomIdArgs{
		ByteLength: pulumi.Int(5),
		Prefix:     pulumi.String(fmt.Sprintf("%s-", prefix)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate gsa account id value: %w", err)
	}

	var contextState = ctx.Value(microservicecontextstate.Key).(microservicecontextstate.ContextState)

	addWorkLoadIdentityGsaAccountIdToContext(&contextState, gsaAccountId)
	ctx = ctx.WithValue(microservicecontextstate.Key, contextState)
	return ctx, nil
}

func addWorkLoadIdentityGsaAccountIdToContext(existingConfig *microservicecontextstate.ContextState, workLoadIdentityGsaAccountId *random.RandomId) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &microservicecontextstate.AddedResources{
			WorkLoadIdentityGsaAccountId: workLoadIdentityGsaAccountId,
		}
		return
	}
	existingConfig.Status.AddedResources.WorkLoadIdentityGsaAccountId = workLoadIdentityGsaAccountId
}
