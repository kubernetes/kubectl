# Production Manifest

In the production directory, make a manifest
with a different name prefix and labels.

<!-- @makeProductionManifest @test -->
```
cat <<EOF >$OVERLAYS/production/Kube-manifest.yaml
apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: makes-production-tuthello
# description: hello configured for production
namePrefix: production-
objectLabels:
  instance: production
  org: acmeCorporation
objectAnnotations:
  note: Hello, I am production!
bases:
- ../../base
patches:
- deployment.yaml
EOF
```

__Next:__ [Production Patch](patch.md)
