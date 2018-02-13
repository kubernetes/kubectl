# Production Instance

Again, different name prefix and labels. This time patch the deployment to
increase the replica count.


<!-- @makeProductionManifest @test -->
```
cat <<'EOF' >$TUT_APP/production/Kube-manifest.yaml
apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: makes-production-tuthello
description: Tuthello configured for production

namePrefix: production-

objectLabels:
  instance: production
  org: acmeCorporation

objectAnnotations:
  note: Hello, I am production!

packages:
- ..

patches:
- deployment.yaml

EOF
```

<!-- @productionDeployment @test -->
```
cat <<EOF >$TUT_APP/production/deployment.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: tut-deployment
spec:
  replicas: 6
EOF
```

__Next:__ [Compare](compare.md)
