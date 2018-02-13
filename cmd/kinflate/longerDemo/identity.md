# Identity transformation

Run it

<!-- @noCustomization @test -->
```
kinflate inflate -f $TUT_APP >$TUT_TMP/original_out
```

kinflate expects to find `Kube-manifest.yaml` in `$TUT_APP`.

The above command discovers the resources, processes them,
and emits the result to `stdout`.

<!-- @showOutput -->
```
more $TUT_TMP/original_out
```

As the app now stands, this command spits out
_unmodified resources_.

The output could be piped directly to kubectl:

> ```
> kinflate inflate -f $TUT_APP | kubectl apply -f -
> ```

The resulting change to the cluster would be no
different than using kubectl directly:

> ```
> kubectl apply -f $TUT_APP
> ```

__Next:__ [Customization](customization.md)
