# Clone

[hello]: https://github.com/monopole/hello

Assume you want to run the [hello] service.

[off-the-shelf config]: https://github.com/kinflate/example-hello

Find an [off-the-shelf config] for it, and clone that
config into a directory called `base`:

<!-- @cloneIt @test -->
```
git clone \
    https://github.com/kinflate/example-hello \
    $DEMO_HOME/base
```

<!-- @runTree @test -->
```
tree $DEMO_HOME
```

One could immediately apply these resources to a
cluster:

> ```
> kubectl apply -f $DEMO_HOME/base
> ```

to instantiate the _hello_ service in off-the-shelf form.

__Next:__ [The Base Manifest](manifest)
