# Deploy

The individual resource sets are:

<!-- @runKinflateStaging @test -->
```
kinflate inflate -f $OVERLAYS/staging
```

<!-- @runKinflateProduction @test -->
```
kinflate inflate -f $OVERLAYS/production
```

To deploy, pipe the above commands to kubectl apply:

> ```
> kinflate inflate -f $OVERLAYS/staging |\
>     kubectl apply -f -
> ```

> ```
> kinflate inflate -f $OVERLAYS/production |\
>    kubectl apply -f -
> ```
