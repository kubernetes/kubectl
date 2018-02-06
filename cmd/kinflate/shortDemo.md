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

<!-- @downloadTutorialHelloApp @test -->
```
cd $TUT_DIR
git clone https://github.com/kinflate/tuthello
cd tuthello
```

<!-- @customizeApp @test -->
```
kinflate setnameprefix acme-
```

<!-- @confirmEdit @test -->
```
grep namePrefix Kube-manifest.yaml
```

<!-- @applyToCluster @test -->
```
kinflate inflate -f . | kubectl apply -f -
```

This fork of the app repository could be
commited to a private repo.  Rebase as desired.
