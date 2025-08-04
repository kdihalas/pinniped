---
title: Install the Pinniped command-line tool
description: Download and set up the `pinniped` command-line tool on macOS, Linux, or Windows clients.
cascade:
  layout: docs
menu:
  docs:
    name: Install CLI
    weight: 10
    parent: howtos
---
The `pinniped` command-line tool is used to generate Pinniped-compatible kubeconfig files, and is also an important part of the Pinniped-based login flow.

It must be installed by administrators setting up a Pinniped cluster as well as by users accessing a Pinniped-enabled cluster.

## Install using Homebrew on macOS or Linux

Use [Homebrew](https://brew.sh/) to install from the Pinniped [tap](https://github.com/vmware/homebrew-pinniped):

- `brew install vmware/pinniped/pinniped-cli`

## Download binaries

Find the appropriate binary for your platform from the [latest release](https://github.com/vmware/pinniped/releases/latest):

{{< buttonlink filename="pinniped-cli-darwin-amd64" >}}Download {{< latestversion >}} for macOS/amd64{{< buttonicon "download.png" >}}{{< /buttonlink >}}
{{< buttonlink filename="pinniped-cli-darwin-arm64" >}}Download {{< latestversion >}} for macOS/arm64{{< buttonicon "download.png" >}}{{< /buttonlink >}}

{{< buttonlink filename="pinniped-cli-linux-amd64" >}}Download {{< latestversion >}} for Linux/amd64{{< buttonicon "download.png" >}}{{< /buttonlink >}}
{{< buttonlink filename="pinniped-cli-linux-arm64" >}}Download {{< latestversion >}} for Linux/arm64{{< buttonicon "download.png" >}}{{< /buttonlink >}}

{{< buttonlink filename="pinniped-cli-windows-amd64.exe" >}}Download {{< latestversion >}} for Windows/amd64{{< buttonicon "download.png" >}}{{< /buttonlink >}}
{{< buttonlink filename="pinniped-cli-windows-arm64.exe" >}}Download {{< latestversion >}} for Windows/arm64{{< buttonicon "download.png" >}}{{< /buttonlink >}}

You should put the command-line tool somewhere on your `$PATH`, such as `/usr/local/bin` on macOS/Linux.
You'll also need to mark the file as executable, e.g. `chmod +x pinniped` on macOS/Linux.

To find specific versions or view all available platforms and architectures, visit the [releases page](https://github.com/vmware/pinniped/releases/).

### Gatekeeper

If you are using macOS, you may get an error dialog when you first run `pinniped` that says `“pinniped” cannot be opened because the developer cannot be verified`.
Cancel this dialog, open System Preferences, click Security & Privacy, and click the Allow Anyway button next to the Pinniped message.

Run the command again and another dialog appears saying `macOS cannot verify the developer of “pinniped”. Are you sure you want to open it?`.
Click Open to allow the command to proceed.

## Install a specific version via script

Choose your preferred [release](https://github.com/vmware/pinniped/releases) and use it to replace the version number in the URL below.

For example, to install {{< latestversion >}} on Linux/amd64:

```sh
curl -Lso pinniped https://get.pinniped.dev/{{< latestversion >}}/pinniped-cli-linux-amd64 \
  && chmod +x pinniped \
  && sudo mv pinniped /usr/local/bin/pinniped
```

## Next steps

Next, [install the Supervisor]({{< ref "install-supervisor.md" >}}) and/or [install the Concierge]({{< ref "install-concierge.md" >}})!
