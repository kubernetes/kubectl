# Staging Instance

In the staging subdirectory, define a shorter manifest,
with a new name prefix, and some different labels.

<!-- @makeStagingManifest @test -->
```
cat <<'EOF' >$TUT_APP/staging/Kube-manifest.yaml
apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: makes-staging-tuthello

description: Tuthello configured for staging

namePrefix: staging-

objectLabels:
  instance: staging
  org: acmeCorporation

objectAnnotations:
  note: Hello, I am staging!

packages:
- ..

patches:
- map.yaml

EOF
```

Add a configmap customization to change the
server greeting from _Good Morning!_ to _Have a
pineapple!_.

Also, enable the _risky_ flag.

<!-- @stagingMap @test -->
```
cat <<EOF >$TUT_APP/staging/map.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tut-map
data:
  altGreeting: "Have a pineapple!"
  enableRisky: "true"
EOF
```

__Next:__ [Production](production.md)
