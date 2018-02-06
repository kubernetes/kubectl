# kinflate

`kinflate` is a kubernetes cluster configuration
utility.

It reads a cluster app definition from a
directory of YAML files, and creates new YAML suitable
for directly piping to `kubectl apply`.


## Declarative Application Management

kinflate is a a prototype of ideas discussed in the
[Declarative Application
Management (DAM) proposal](https://goo.gl/T66ZcD).

DAM encourages one to use literal kubernetes resource
declarations to configure a cluster.  DAM understands
new declarations to modify resources, rather than
achieving modification via templates or configuration
languages.

DAM facilitates coupling cluster state changes to version
control commits.  It encourages forking a configuration,
customizing it, and occasionally rebasing to capture
upgrades from the original configuration.

## Installation

<!-- @installKinflate @test -->
```
go get k8s.io/kubectl/cmd/kinflate
```
