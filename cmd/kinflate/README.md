# kinflate

`kinflate` is a kubernetes cluster configuration utility,
a prototype of ideas discussed in this [proposal for
Declarative Application Management](https://goo.gl/T66ZcD).

## Declarative Application Management (DAM)

[declarations]: https://kubernetes.io/docs/tutorials/object-management-kubectl/declarative-object-management-configuration/

DAM allows one to customize native kubernetes resources
using new modification declarations rather than using
templates or new configuration languages.

It facilitates coupling cluster state changes to version
control commits.  It encourages forking a configuration,
customizing it, and occasionally rebasing to capture
upgrades from the original configuration.

## Installation

This assumes you have [Go](https://golang.org/) installed.

<!-- @installKinflate @test -->
```
tmp=$(mktemp -d)
GOBIN=$tmp go get k8s.io/kubectl/cmd/kinflate
```

## Example

[kubectl]: https://kubernetes.io/docs/user-guide/kubectl-overview/
[YAML]: http://www.yaml.org/start.html

The following commands checkout a _kinflate app_,
customize all the resources in the app with a name prefix,
then pipe the resulting [YAML] directly to a cluster via [kubectl].

```
# Download a demo app
git clone https://github.com/monopole/kapps/demo $tmp

# Use kinflate to customize it
cd $tmp
$GOBIN/kinflate setprefixname acme-

# See the change:
grep prefixName Kube-manifest.yaml

# Send the app to your cluster:
$GOBIN/kinflate inflate -f . | kubectl apply -f -
```
