# Staging Manifest

In the `staging` directory, make a manifest
defining a new name prefix, and some different labels.

<!-- @makeStagingManifest @test -->
```
cat <<'EOF' >$OVERLAYS/staging/Kube-manifest.yaml
apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: makes-staging-hello
description: hello configured for staging
namePrefix: staging-
objectLabels:
  instance: staging
  org: acmeCorporation
objectAnnotations:
  note: Hello, I am staging!
bases:
- ../../base
patches:
- map.yaml
EOF
```
__Next:__ [Staging Patch](patch.md)
