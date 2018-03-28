# Demo: hello world

[kubectl]: https://kubernetes.io/docs/user-guide/kubectl-overview/

Steps:

 1. Clone an off-the-shelf configuration.
 1. Customize its resources with a name prefix.
 1. Apply the result to a cluster via [kubectl].

First make a place to work:

<!-- @makeDemoDir @test -->
```
DEMO_HOME=$(mktemp -d)
```

[example-hello]: https://github.com/kinflate/example-hello
Clone an example configuration ([example-hello]):

<!-- @cloneExample @test -->
```
cd $DEMO_HOME
git clone https://github.com/kinflate/example-hello
cd example-hello
```

Customize the base application by specifying a prefix
that will be applied to all resource names:

<!-- @customizeApp @test -->
```
kinflate set nameprefix acme-
```

Confirm that the edit happened:

<!-- @confirmEdit @test -->
```
grep --context=3 namePrefix Kube-manifest.yaml
```

Confirm that the prefix appears in the output:

<!-- @confirmResourceNames @test -->
```
kinflate inflate -f . | grep --context=3 acme-
```

Optionally apply the modified configuration to a cluster

<!-- @applyToCluster -->
```
kinflate inflate -f . | kubectl apply -f -
```

This fork of [example-hello] could be commited to a
private repo, and evolve independently of its upstream.

One could rebase from upstream as desired.

__Next:__ [hello world (with instances)](demoHelloWorldLong/README.md).
