[manifest]: ../docs/glossary.md#manifest
[base]: ../docs/glossary.md#base
[overlay]: ../docs/glossary.md#overlay
[overlays]: ../docs/glossary.md#overlay
[instance]: ../docs/glossary.md#instance
[instances]: ../docs/glossary.md#instance
[hello]: https://github.com/monopole/hello
[config]: https://github.com/kinflate/example-hello
[original]: https://github.com/kinflate/example-hello

# Demo: hello world with instances

Steps:

 1. Clone an existing configuration as a [base].
 1. Customize it.
 1. Create two different [overlays] (_staging_ and _production_)
    from the customized base.
 1. Run kinflate and kubectl to deploy staging and production.

First define a place to work:

<!-- @makeWorkplace @test -->
```
DEMO_HOME=$(mktemp -d)
```

Alternatively, use

> ```
> DEMO_HOME=~/hello
> ```

## Clone an example

Let's run the [hello] service.
Here's an existing [config] for it.

Clone this config to a directory called `base`:

<!-- @cloneIt @test -->
```
git clone \
    https://github.com/kinflate/example-hello \
    $DEMO_HOME/base
```

<!-- @runTree @test -->
```
tree $DEMO_HOME
```

Expecting something like:
> ```
> /tmp/tmp.IyYQQlHaJP
> └── base
>     ├── configMap.yaml
>     ├── deployment.yaml
>     ├── Kube-manifest.yaml
>     ├── LICENSE
>     ├── README.md
>     └── service.yaml
> ```


One could immediately apply these resources to a
cluster:

> ```
> kubectl apply -f $DEMO_HOME/base
> ```

to instantiate the _hello_ service.  `kubectl`
would only recognize the resource files.

## The Manifest

The `base` directory has a [manifest] file:

<!-- @manifest @test -->
```
BASE=$DEMO_HOME/base
more $BASE/Kube-manifest.yaml
```

Run `kinflate` on the base to emit customized resources
to `stdout`:

<!-- @manifest @test -->
```
kinflate build $BASE
```

## Customize the base

A first customization step could be to change the _app
label_ applied to all resources:

<!-- @manifest @test -->
```
sed -i 's/app: hello/app: my-hello/' \
    $BASE/Kube-manifest.yaml
```

See the effect:
<!-- @manifest @test -->
```
kinflate build $BASE | grep -C 3 app:
```

## Create Overlays

Create a _staging_ and _production_ [overlay]:

 * _Staging_ enables a risky feature not enabled in production.
 * _Production_ has a higher replica count.
 * Web server greetings from these cluster
   [instances] will differ from each other.

<!-- @overlayDirectories @test -->
```
OVERLAYS=$DEMO_HOME/overlays
mkdir -p $OVERLAYS/staging
mkdir -p $OVERLAYS/production
```

#### Staging Manifest

In the `staging` directory, make a manifest
defining a new name prefix, and some different labels.

<!-- @makeStagingManifest @test -->
```
cat <<'EOF' >$OVERLAYS/staging/Kube-manifest.yaml
apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: makes-staging-hello
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

#### Staging Patch

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

#### Production Manifest

In the production directory, make a manifest
with a different name prefix and labels.

<!-- @makeProductionManifest @test -->
```
cat <<EOF >$OVERLAYS/production/Kube-manifest.yaml
apiVersion: manifest.k8s.io/v1alpha1
kind: Package
metadata:
  name: makes-production-tuthello
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


#### Production Patch

Make a production patch that increases the replica
count (because production takes more traffic).

<!-- @productionDeployment @test -->
```
cat <<EOF >$OVERLAYS/production/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: the-deployment
spec:
  replicas: 10
EOF
```

## Compare overlays


`DEMO_HOME` now contains:

 - a _base_ directory - a slightly customized clone
   of the original configuration, and

 - an _overlays_ directory, containing the manifests
   and patches required to create distinct _staging_
   and _production_ instances in a cluster.

Review the directory structure and differences:

<!-- @listFiles @test -->
```
tree $DEMO_HOME
```

Expecting something like:
> ```
> /tmp/tmp.IyYQQlHaJP1<
> ├── base
> │   ├── configMap.yaml
> │   ├── deployment.yaml
> │   ├── Kube-manifest.yaml
> │   ├── LICENSE
> │   ├── README.md
> │   └── service.yaml
> └── overlays
>     ├── production
>     │   ├── deployment.yaml
>     │   └── Kube-manifest.yaml
>     └── staging
>         ├── Kube-manifest.yaml
>         └── map.yaml
> ```

<!-- @compareOutput -->
```
diff \
  <(kinflate build $OVERLAYS/staging) \
  <(kinflate build $OVERLAYS/production) |\
  more
```


## Deploy

The individual resource sets are:

<!-- @buildStaging @test -->
```
kinflate build $OVERLAYS/staging
```

<!-- @buildProduction @test -->
```
kinflate build $OVERLAYS/production
```

To deploy, pipe the above commands to kubectl apply:

> ```
> kinflate build $OVERLAYS/staging |\
>     kubectl apply -f -
> ```

> ```
> kinflate build $OVERLAYS/production |\
>    kubectl apply -f -
> ```

## Rollback

To rollback, one would undo whatever edits were made to
the configuation in source control, then rerun kinflate
on the reverted configuration and apply it to the
cluster.
