package remotedebug

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/remotedebug/cert"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/remotedebug/gateway"
	"github.com/plantoncloud/microservice-kubernetes-pulumi-blueprint/pkg/microservice/ingress/remotedebug/virtualservice"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources required for supporting remote-debug for java microservices.
// jwdp protocol only supports tcp protocol. all ingress controllers rely on host header in http or sniHostname
// in tls protocol to route the incoming requests to the correct kubernetes service.
// ingress controllers can not route incoming tcp requests based on hostname as that information is not included
// in the tcp request.
// to address this gap, developers will use stunnel on client side to wrap jwdp initiated tcp connection to tls
// connection, which will be terminated by istio and route the request to the correct java microservice based on
// sniHostname in the terminated tls request. java microservice would only receive terminated tls connections as tcp
// connections.
// important: istio can only support routing terminated tls connection to kubernetes services when gateways
// resource contains the exact hostnames instead of wildcard hosts. so, a separate gateway resource with
// exact hostname is created for each versioned microservice deployment.
func Resources(ctx *pulumi.Context) error {
	err := cert.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add self-signed cert")
	}

	err = gateway.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add gateway resources")
	}

	err = virtualservice.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add virtual-service resources")
	}

	return nil
}
