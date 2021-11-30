---
title: Configure the Pinniped Supervisor as an OIDC issuer
description: Set up the Pinniped Supervisor to provide seamless login flows across multiple clusters.
cascade:
  layout: docs
menu:
  docs:
    name: Configure Supervisor as an OIDC Issuer
    weight: 70
    parent: howtos
---

The Supervisor is an [OpenID Connect (OIDC)](https://openid.net/connect/) issuer that supports connecting a single
"upstream" identity provider to many "downstream" cluster clients. When a user authenticates, the Supervisor can issue
[JSON Web Tokens (JWTs)](https://tools.ietf.org/html/rfc7519) that can be [validated by the Pinniped Concierge]({{< ref "configure-concierge-jwt" >}}).

This guide explains how to expose the Supervisor's REST endpoints to clients.

## Prerequisites

This how-to guide assumes that you have already [installed the Pinniped Supervisor]({{< ref "install-supervisor" >}}).

## Exposing the Supervisor app's endpoints outside the cluster

The Supervisor app's endpoints should be exposed as HTTPS endpoints with proper TLS certificates signed by a
certificate authority (CA) which is trusted by your end user's web browsers.

It is recommended that the traffic to these endpoints should be encrypted via TLS all the way into the
Supervisor pods, even when crossing boundaries that are entirely inside the Kubernetes cluster.
The credentials and tokens that are handled by these endpoints are too sensitive to transmit without encryption.

In all versions of the Supervisor app so far, there are both HTTP and HTTPS ports available for use by default.
These ports each host all the Supervisor's endpoints. Unfortunately, this has caused some confusion in the community
and some blog posts have been written which demonstrate using the HTTP port in such a way that a portion of the traffic's
path is unencrypted. **Anything which exposes the non-TLS HTTP port outside the Pod should be considered deprecated**.
A future version of the Supervisor app may include a breaking change to adjust the default behavior of the HTTP port
to only listen on localhost (or perhaps even to be disabled) to make it more clear that the Supervisor app is not intended to receive
non-TLS HTTP traffic from outside the Pod (expect the `/healthz` endpoint).

Because there are many ways to expose TLS services from a Kubernetes cluster, the Supervisor app leaves this up to the user.
The most common ways are:

- Define a [TCP LoadBalancer Service](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer).

  In this case, the Service is a layer 4 load balancer which does not terminate TLS, so the Supervisor app needs to be
  configured with TLS certificates and will terminate the TLS connection itself (see the section about FederationDomain
  below). The LoadBalancer Service should be configured to use the HTTPS port 443 of the Supervisor pods as its `targetPort`.

  *Warning:* Never expose the Supervisor's HTTP port 8080 to the public. It would not be secure for the OIDC protocol
  to use HTTP, because the user's secret OIDC tokens would be transmitted across the network without encryption.

- Or, define an [Ingress resource](https://kubernetes.io/docs/concepts/services-networking/ingress/).

   In this case, the [Ingress typically terminates TLS](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls)
   and then talks plain HTTP to its backend,
   which would be a NodePort or LoadBalancer Service in front of the HTTP port 8080 of the Supervisor pods.
   However, because the Supervisor's endpoints deal with sensitive credentials, it is much better if the
   traffic is encrypted using TLS all the way into the Supervisor's Pods. Some Ingress implementations
   may support re-encrypting the traffic before sending it to the backend. If your Ingress controller does not
   support this, then consider using one of the other configurations described here instead of using an Ingress.

   The required configuration of the Ingress is specific to your cluster's Ingress Controller, so please refer to the
   documentation from your Kubernetes provider. If you are using a cluster from a cloud provider, then you'll probably
   want to start with that provider's documentation. For example, if your cluster is a Google GKE cluster, refer to
   the [GKE documentation for Ingress](https://cloud.google.com/kubernetes-engine/docs/concepts/ingress).
   Otherwise, the Kubernetes documentation provides a list of popular
   [Ingress Controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/), including
   [Contour](https://projectcontour.io/) and many others.

- Or, expose the Supervisor app using a Kubernetes service mesh technology (e.g. [Istio](https://istio.io/)).

   In this case, the setup would be similar to the previous description
   for defining an Ingress, except the service mesh would probably provide both the ingress with TLS termination
   and the service. Please see the documentation for your service mesh.

   If your service mesh is capable of transparently encrypting traffic all the way into the
   Supervisor Pods, then you should use that capability. In this case, it may make sense to configure the Supervisor's
   HTTP port to only listen on localhost, such as when the service mesh injects a sidecar container that can securely
   access the HTTP port via localhost networking from within the same Pod.
   See the `container_http_listener` option in [deploy/supervisor/values.yml](https://github.com/vmware-tanzu/pinniped/blob/main/deploy/supervisor/values.yaml)
   for more information.
   This would prevent any unencrypted traffic from accidentally being transmitted from outside the Pod into the
   Supervisor app's HTTP port.

## Creating a Service to expose the Supervisor app's endpoints within the cluster

Now that you've selected a strategy to expose the endpoints outside the cluster, you can choose how to expose
the endpoints inside the cluster in support of that strategy.

If you've decided to use a LoadBalancer Service then you'll need to create it. On the other hand, if you've decided to
use an Ingress then you'll need to create a Service which the Ingress can use as its backend. Either way, how you
create the Service will depend on how you choose to install the Supervisor:

- If you installed using `ytt` then you can use
the related `service_*` options from [deploy/supervisor/values.yml](https://github.com/vmware-tanzu/pinniped/blob/main/deploy/supervisor/values.yaml)
to create a Service. 
- If you installed using the pre-rendered manifests attached to the Pinniped GitHub releases, then you can create
the Service separately after installing the Supervisor app.

There is no Ingress included in either the `ytt` templates or the pre-rendered manifests,
so if you choose to use an Ingress then you'll need to create the Ingress separately after installing the Supervisor app.

### Example: Creating a LoadBalancer Service

This is an example of creating a LoadBalancer Service to expose port 8443 of the Supervisor app outside the cluster.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: pinniped-supervisor-loadbalancer
  # Assuming that this is the namespace where the Supervisor was installed.
  # This is the default.
  namespace: pinniped-supervisor
spec:
  type: LoadBalancer
  selector:
    # Assuming that this is how the Supervisor Pods are labeled.
    # This is the default.
    app: pinniped-supervisor
  ports:
  - protocol: TCP
    port: 443
    targetPort: 8443 # 8443 is the TLS port. Do not expose port 8080.
```

### Example: Creating a NodePort Service

A NodePort Service exposes the app as a port on the nodes of the cluster. For example, a NodePort Service could also be
used as the backend of an Ingress.

This is also convenient for use with Kind clusters, because kind can
[expose node ports as localhost ports on the host machine](https://kind.sigs.k8s.io/docs/user/configuration/#extra-port-mappings)
without requiring an Ingress, although
[Kind also supports several Ingress Controllers](https://kind.sigs.k8s.io/docs/user/ingress).

```yaml
apiVersion: v1
kind: Service
metadata:
  name: pinniped-supervisor-nodeport
  # Assuming that this is the namespace where the Supervisor was installed.
  # This is the default.
  namespace: pinniped-supervisor
spec:
  type: NodePort
  selector:
    # Assuming that this is how the Supervisor Pods are labeled.
    # This is the default.
    app: pinniped-supervisor
  ports:
  - protocol: TCP
    port: 443
    targetPort: 8443
    # This is the port that you would forward to the kind host.
    # Or omit this key for a random port on the node.
    nodePort: 31234
```

## Configuring the Supervisor to act as an OIDC provider

The Supervisor can be configured as an OIDC provider by creating FederationDomain resources
in the same namespace where the Supervisor app was installed. At least one FederationDomain must be configured
for the Supervisor to provide its functionality.

Here is an example of a FederationDomain.

```yaml
apiVersion: config.supervisor.pinniped.dev/v1alpha1
kind: FederationDomain
metadata:
  name: my-provider
  # Assuming that this is the namespace where the supervisor was installed.
  # This is the default.
  namespace: pinniped-supervisor
spec:
  # The hostname would typically match the DNS name of the public ingress
  # or load balancer for the cluster.
  # Any path can be specified, which allows a single hostname to have
  # multiple different issuers. The path is optional.
  issuer: https://my-issuer.example.com/any/path
  # Optionally configure the name of a Secret in the same namespace,
  # of type `kubernetes.io/tls`, which contains the TLS serving certificate
  # for the HTTPS endpoints served by this OIDC Provider.
  tls:
    secretName: my-tls-cert-secret
```

You can create multiple FederationDomains as long as each has a unique issuer string.
Each FederationDomain can be used to provide access to a set of Kubernetes clusters for a set of user identities.

### Configuring TLS for the Supervisor OIDC endpoints

If you have terminated TLS outside the app, for example using service mesh which handles encrypting the traffic for you,
then you do not need to configure TLS certificates on the FederationDomain.  Otherwise, you need to configure the
Supervisor app to terminate TLS.

There are two places to optionally configure TLS certificates:

1. Each FederationDomain can be configured with TLS certificates, using the `spec.tls.secretName` field.

1. The default TLS certificate for all FederationDomains can be configured by creating a Secret called
`pinniped-supervisor-default-tls-certificate` in the same namespace in which the Supervisor was installed.

Each incoming request to the endpoints of the Supervisor may use TLS certificates that were configured in either
of the above ways. The TLS certificate to present to the client is selected dynamically for each request
using Server Name Indication (SNI):
- When incoming requests use SNI to specify a hostname, and that hostname matches the hostname
  of a FederationDomain, and that FederationDomain specifies `spec.tls.secretName`, then the TLS certificate from the
  `spec.tls.secretName` Secret will be used.
- Any other request will use the default TLS certificate, if it is specified. This includes any request to a host
  which is an IP address, because SNI does not work for IP addresses. If the default TLS certificate is not specified,
  then these requests will fail.

It is recommended that you have a DNS entry for your load balancer or Ingress, and that you configure the
OIDC provider's `issuer` using that DNS hostname, and that the TLS certificate for that provider also
covers that same hostname.

You can create the certificate Secrets however you like, for example you could use [cert-manager](https://cert-manager.io/)
or `kubectl create secret tls`. They must be Secrets of type `kubernetes.io/tls`.
Keep in mind that your end users must load some of these endpoints in their web browsers, so the TLS certificates
should be signed by a certificate authority that is trusted by their browsers.

## Next steps

Next, configure an OIDCIdentityProvider, ActiveDirectoryIdentityProvider, or an LDAPIdentityProvider for the Supervisor
(several examples are available in these guides),
and [configure the Concierge to use the Supervisor for authentication]({{< ref "configure-concierge-supervisor-jwt" >}})
on each cluster!
