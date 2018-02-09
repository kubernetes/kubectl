# Short Demo

Prerequisites: [Go](https://golang.org/), [git](https://git-scm.com),
a `PATH` containing `$GOPATH/bin`, [kubectl] set up
to talk to a cluster.

[kubectl]: https://kubernetes.io/docs/user-guide/kubectl-overview/

The following clones an app, customizes the resources in
the app with a name prefix, then pipes the resulting
YAML directly to a cluster via [kubectl].

<!-- @makeWorkDir @test -->
```
TUT_DIR=$HOME/kinflate_demo
/bin/rm -rf $TUT_DIR
mkdir -p $TUT_DIR
```

Download the app:

<!-- @downloadTutorialHelloApp @test -->
```
cd $TUT_DIR
git clone https://github.com/kinflate/tuthello
cd tuthello
```

Customize the base application by specifying a prefix that will be applied to
all resource names:

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

Optionally apply to a cluster:

<!-- @applyToCluster -->
```
kinflate inflate -f . | kubectl apply -f -
```

This fork of the app repository could be committed to a private repo.  Rebase as
desired.
