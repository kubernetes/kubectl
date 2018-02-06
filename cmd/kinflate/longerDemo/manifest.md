# Manifest

Add a manifest, to turn a directory of resources into a kinflate
application.

A manifest is always called

> `Kube-manifest.yaml`

<!-- @defineManifest @demo -->
```
export TUT_APP_MANIFEST=$TUT_APP/Kube-manifest.yaml
```

<!-- @makeManifest @demo -->
```
cat <<'EOF' >$TUT_APP_MANIFEST

apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: tuthello

description: Tuthello offers greetings as a service.
keywords: [hospitality, kubernetes]
appVersion: 0.1.0

home: https://github.com/TODO_MakeAHomeForThis

sources:
- https://github.com/TODO_MakeAHomeForThis

icon: https://www.pexels.com/photo/some-icon.png

maintainers:
- name: Miles Mackentunnel
  email: milesmackentunnel@gmail.com
  github: milesmackentunnel

resources:
- deployment.yaml
- configMap.yaml
- service.yaml

EOF
```

### K8s API fields

[k8s API style]: https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/#required-fields

The manifest follows [k8s API style], defining

 * _apiVersion_ - which version of the k8s API is known to work
   with this object
 * _kind_ - the object type
 * _metadata_ - data that helps uniquely identify the
   object, e.g. a name.


### The Usual Suspects

The manifest has the fields one would expect:
_description_, _version_, _home page_, _maintainers_, etc.

### Resources field

This is a _manifest_, it must list the cargo.

## Save a copy

Before going any further, save a copy of the manifest:

<!-- @copyManifest @test -->
```
export TUT_ORG_MANIFEST=$TUT_TMP/original-manifest.yaml
cp $TUT_APP_MANIFEST $TUT_ORG_MANIFEST
```
