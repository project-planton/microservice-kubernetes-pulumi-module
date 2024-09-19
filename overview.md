# Overview

The Pulumi module provided is designed to automate the deployment of a microservice onto a Kubernetes cluster using
Pulumi and Go. It reads the microservice specifications from a standardized API resource model, which includes details
like the API version, kind, metadata, spec, and status. The module utilizes this structured input to create and
configure the necessary Kubernetes resources that represent the desired state of the microservice deployment.

Key functionalities of the module include setting up a Kubernetes namespace, creating image pull secrets based on
provided Docker credentials, and deploying the microservice using a Kubernetes Deployment resource. It configures
environment variables, secrets, ports, and resource limits as specified. The module also creates a Kubernetes Service to
expose the microservice internally within the cluster. If ingress is enabled, it sets up ingress resources using the
Gateway API and Cert-Manager to handle external and internal HTTPS traffic, including the provisioning of TLS
certificates and defining routing rules for hostnames.
