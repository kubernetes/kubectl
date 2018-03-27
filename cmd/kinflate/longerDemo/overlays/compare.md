# Compare them

[original]: https://github.com/kinflate/example-hello

`DEMO_HOME` now contains a _base_ directory - your
slightly customized clone of the [original]
configuration, and an _overlays_ directory, that
contains all one needs to create a _staging_ and
_production_ instance in a cluster.

Review the directory structure:

<!-- @listFiles @test -->
```
tree $DEMO_HOME
```

<!-- @compareOutput -->
```
diff \
  <(kinflate inflate -f $OVERLAYS/staging) \
  <(kinflate inflate -f $OVERLAYS/production) |\
  more
```

__Next:__ [Deploy](deploy)
