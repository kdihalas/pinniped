---
title: Install the Pinniped Concierge
description: Install the Pinniped Concierge service in a Kubernetes cluster.
cascade:
  layout: docs
menu:
  docs:
    name: Install Concierge
    weight: 20
    parent: howtos      
---
This guide shows you how to install the Pinniped Concierge.
You should have a [supported Kubernetes cluster]({{< ref "../reference/supported-clusters" >}}).

In the examples below, you can replace *{{< latestversion >}}* with your preferred version number.
You can find a list of Pinniped releases [on GitHub](https://github.com/vmware/pinniped/releases).

## With default options

**Warning:** the default Concierge configuration may create a public LoadBalancer Service on your cluster if that is the default on your cloud provider.
If you'd prefer to customize the annotations or load balancer IP address, see the "With custom options" section below.

### Using kapp

1. Install the latest version of the Concierge into the `pinniped-concierge` namespace with default options using [kapp](https://carvel.dev/kapp/):

   - `kapp deploy --app pinniped-concierge --file https://get.pinniped.dev/{{< latestversion >}}/install-pinniped-concierge.yaml`

### Using kubectl

1. Install the latest version of the Concierge CustomResourceDefinitions:

   - `kubectl apply -f https://get.pinniped.dev/{{< latestversion >}}/install-pinniped-concierge-crds.yaml`

   This step is required so kubectl can validate the custom resources deployed in the next step.

1. Install the latest version of the Concierge into the `pinniped-concierge` namespace with default options:

   - `kubectl apply -f https://get.pinniped.dev/{{< latestversion >}}/install-pinniped-concierge-resources.yaml`

## With custom options

Pinniped uses [ytt](https://carvel.dev/ytt/) from [Carvel](https://carvel.dev/) as a templating system.

1. Install the `ytt` and `kapp` command-line tools using the instructions from the [Carvel documentation](https://carvel.dev/#whole-suite).

1. Clone the Pinniped GitHub repository and visit the `deploy/concierge` directory:

   - `git clone git@github.com:vmware/pinniped.git`
   - `cd pinniped/deploy/concierge`

1. Decide which release version you would like to install. All release versions are [listed on GitHub](https://github.com/vmware/pinniped/releases).

1. Checkout your preferred version tag, e.g. `{{< latestversion >}}`.

   - `git checkout {{< latestversion >}}`

1. Customize configuration parameters:

    - See the [default values](http://github.com/vmware/pinniped/tree/main/deploy/concierge/values.yaml) for documentation about individual configuration parameters.
      For example, you can change the number of Concierge pods by setting `replicas` or apply custom annotations to the impersonation proxy service using `impersonation_proxy_spec`.

    - In a different directory, create a new YAML file to contain your site-specific configuration. For example, you might call this file `site/dev-env.yaml`.

      In the file, add the special ytt comment for a values file and the YAML triple-dash which starts a new YAML document.
      Then add custom overrides for any of the parameters from [`values.yaml`](http://github.com/vmware/pinniped/tree/main/deploy/concierge/values.yaml).

      Override the `image_tag` value to match your preferred version tag, e.g. `{{< latestversion >}}`,
      to ensure that you use the version of the server which matches these templates.

      Here is an example which overrides the image tag, the default logging level, and the number of replicas:
      ```yaml
      #@data/values
      ---
      image_tag: {{< latestversion >}}
      log_level: debug
      replicas: 1
      ```
    - Parameters for which you would like to use the default value should be excluded from this file.

    - If you are using a GitOps-style workflow to manage the installation of Pinniped, then you may wish to commit this new YAML file to your GitOps repository.

1. Render templated YAML manifests:

   - `ytt --file . --file site/dev-env.yaml`

   By putting the override file last in the list of `--file` options, it will override the default values.

1. Deploy the templated YAML manifests:

   - `ytt --file . --file site/dev-env.yaml | kapp deploy --app pinniped-concierge --file -`

## Supported Node Architectures

The Pinniped Concierge can be installed on Kubernetes clusters with available `amd64` or `arm64` linux nodes.

## Other notes

_Important:_ Configure Kubernetes authorization policies (i.e. RBAC) to prevent non-admin users from reading the
resources, especially the Secrets, in the Concierge's namespace.

## Next steps

Next, configure the Concierge for
[JWT]({{< ref "configure-concierge-jwt.md" >}}) or [webhook]({{< ref "configure-concierge-webhook.md" >}}) authentication,
or [configure the Concierge to use the Supervisor for authentication]({{< ref "configure-concierge-supervisor-jwt" >}}).
