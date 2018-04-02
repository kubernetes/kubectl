# kinflate

[_kubectl apply_]: glossary.md#apply
[DAM]: glossary.md#declarative-application-management
[workflows]: workflows.md

`kinflate` is a command line tool supporting
template-free customization of declarative
configuration targetted to kubernetes.  It's an
implementation of ideas described in Brian Grant's
[DAM] proposal.

kinflate plays a role in various configuration
management [workflows].

## Design tenets

 * __composable__: do one thing, use plain text, work
   with pipes, usable in other tools (e.g. helm,
   kubernetes-deploy, etc.).

 * __simple__: no templating, no logic, no inheritance,
   no new API obscuring the k8s API.

 * __accountable__: gitops friendly, diff against
   declaration in version control, diff against
   cluster.

 * __k8s style__: recognizable k8s resources,
   extensible (openAPI, CRDs),
   patching, intended for [_kubectl apply_].


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
