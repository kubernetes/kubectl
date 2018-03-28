# kinflate

[_kubectl apply_]: https://goo.gl/UbCRuf
[Declarative Application Management]: https://goo.gl/T66ZcD


`kinflate` is a command line tool supporting
template-free customization of declarative
configuration targetted to kubernetes.

It's an implementation of ideas described in Brian
Grant's [Declarative Application Management] proposal.


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

Note: golang 1.10 or newer is required.

## Demos

 * [hello world one-pager](demoHelloWorldShort.md)
 * [hello world, with instances, slide format](demoHelloWorldLong/README.md)
 * [mysql one-pager](demoMySql.md)
