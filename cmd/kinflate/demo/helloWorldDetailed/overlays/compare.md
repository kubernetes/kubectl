# Compare overlays

[original]: https://github.com/kinflate/example-hello

`DEMO_HOME` now contains:

 - a _base_ directory - your slightly customized clone of the [original]
   configuration, and

 - an _overlays_ directory, containing the manifests and patches required to
   create distinct _staging_ and _production_ instance in a cluster.

Review the directory structure and differences:

<!-- @listFiles @test -->
```
tree $DEMO_HOME
```

<!-- @compareOutput -->
```
diff \
  <(kinflate inflate $OVERLAYS/staging) \
  <(kinflate inflate $OVERLAYS/production) |\
  more
```

__Next:__ [Deploy](deploy.md)
