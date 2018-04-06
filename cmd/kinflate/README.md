# kinflate

[_kubectl apply_]: docs/glossary.md#apply
[base]: docs/glossary.md#base
[declarative configuration]: docs/glossary.md#declarative-application-management
[imageBase]: docs/base.jpg
[imageOverlay]: docs/overlay.jpg
[manifest]: docs/glossary.md#manifest
[overlay]: docs/glossary.md#overlay
[resources]: docs/glossary.md#resource
[workflows]: docs/workflows.md

`kinflate` is a command line tool supporting
template-free customization of declarative
configuration targetted to kubernetes.

## Usage

#### 1) Make a base

A [base] configuration is a [manifest] listing a set of
k8s [resources] - deployments, services, configmaps,
secrets that serve some common purpose.

![base image][imageBase]

#### 2) Customize it with overlays

An [overlay] customizes your base along different dimensions
for different purposes or different teams, e.g. for
_development, staging, production_.

![overlay image][imageOverlay]

#### 3) Run kinflate

Run kinflate on your overlay.  The result, a set of
complete resources, is printed to stdout, suitable for
sending to your cluster.

## Installation

Assumes [Go](https://golang.org/) is installed
and your `PATH` contains `$GOPATH/bin`:

<!-- @installKinflate @test -->
```
go get k8s.io/kubectl/cmd/kinflate
```

## Demos

 * [hello world one-pager](demo/helloWorldOnePager.md)
 * [hello world detailed](demo/helloWorldDetailed/README.md) (with instances, in slide format)
 * [mysql one-pager](demo/mySql.md)
