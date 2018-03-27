# Production Patch

Make a production patch that increases the replica count (because production
takes more traffic).

<!-- @productionDeployment @test -->
```
cat <<EOF >$OVERLAYS/production/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  replicas: 6
EOF
```

__Next:__ [Compare](../compare)
