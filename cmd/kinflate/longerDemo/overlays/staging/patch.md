# Staging patch

Add a configmap customization to change the server
greeting from _Good Morning!_ to _Have a pineapple!_

Also, enable the _risky_ flag.

<!-- @stagingMap @test -->
```
cat <<EOF >$OVERLAYS/staging/map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: the-map
data:
  altGreeting: "Have a pineapple!"
  enableRisky: "true"
EOF
```

__Next:__ [Production Manifest](../production)
