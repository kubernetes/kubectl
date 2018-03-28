# Customize

A first customization step could be to change the _app label_
applied to all resources:

<!-- @manifest @test -->
```
sed -i 's/app: hello/app: my-hello/' \
    $BASE/Kube-manifest.yaml
```

See the effect:
<!-- @manifest @test -->
```
kinflate inflate -f $BASE | grep -C 3 app:
```

__Next:__ [Overlays](overlays)
