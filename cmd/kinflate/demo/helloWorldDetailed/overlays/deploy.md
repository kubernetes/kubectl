# Deploy

The individual resource sets are:

<!-- @runKinflateStaging @test -->
```
kinflate inflate $OVERLAYS/staging
```

<!-- @runKinflateProduction @test -->
```
kinflate inflate $OVERLAYS/production
```

To deploy, pipe the above commands to kubectl apply:

> ```
> kinflate inflate $OVERLAYS/staging |\
>     kubectl apply -f -
> ```

> ```
> kinflate inflate $OVERLAYS/production |\
>    kubectl apply -f -
> ```

__Next:__ [Editting](../editor.md)
