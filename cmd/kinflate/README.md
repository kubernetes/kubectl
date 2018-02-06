# kinflate

`kinflate` is a kubernetes cluster configuration utility,
a prototype of ideas discussed in this [proposal for
Declarative Application Management](https://goo.gl/T66ZcD) (DAM).

[declarations]: https://kubernetes.io/docs/tutorials/object-management-kubectl/declarative-object-management-configuration/

DAM allows one to customize native kubernetes resources
using new modification declarations rather than using
template or configuration languages.

It facilitates coupling cluster state changes to version
control commits.  It encourages forking a configuration,
customizing it, and occasionally rebasing to capture
upgrades from the original configuration.

## Installation

Prerequisites: [Go](https://golang.org/), [git](https://git-scm.com)

<!-- @installKinflate @test -->
```
tmp=$(mktemp -d)
GOBIN=$tmp go get k8s.io/kubectl/cmd/kinflate
```

## Example

[kubectl]: https://kubernetes.io/docs/user-guide/kubectl-overview/
[YAML]: http://www.yaml.org/start.html

Clone an app, customize the resources in the app with a
name prefix, then pipe the resulting [YAML] directly to
a cluster via [kubectl].

```
# Download a demo app
cd $tmp
git clone https://github.com/monopole/damapps
cd damapps/demo

# Use kinflate to customize it
$tmp/kinflate setprefixname acme-

# See the added line:
grep namePrefix Kube-manifest.yaml
git diff Kube-manifest.yaml

# Send the app to your cluster:
$tmp/kinflate inflate -f . | kubectl apply -f -

```

At this point, you could commit your fork
to your own repository, and rebase as desired.
