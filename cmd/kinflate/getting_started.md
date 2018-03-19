# Kinflate: Getting Started

In this `getting started` guide, we will take an off-the-shelf MySQL configuration for Kubernetes and customize it to suit our production scenario.
In production environment, we want:
- MySQL resource names to be prefixed by 'prod-' to make them distinguishable.
- MySQL resources to have 'env: prod' labels so that we can use label selector to query these. 
- MySQL to use persistent disk for storing data.

## Installation

If you have `kinflate` installed already, then you can skip this step. `kinflate` can be installed using `go get`:

<!-- @installKinflate @test -->
```shell
go get k8s.io/kubectl/cmd/kinflate
```
This fetches kinflate and install `kinflate` executable under `GOPATH/bin` - you will want that on your $PATH.


### Get off-the-shelf MySQL configs for Kubernetes
Download a sample MySQL YAML manifest files:

<!-- @makeMySQLDir @test -->
```shell
MYSQL_DIR=$HOME/kinflate_demo/mysql
rm -rf $MYSQL_DIR && mkdir -p $MYSQL_DIR
cd $MYSQL_DIR

# Get MySQL configs
for f in service secret deployment ; do \
	wget https://raw.githubusercontent.com/kinflate/mysql/master/emptyDir/$f.yaml ; \
done

```
This downloads  YAML files `deployment.yaml`, `service.yaml` and `secret.yaml` which are needed to run MySQL in a Kubernetes cluster.

### Initialization
Now that we have base configuration for MySQL in place, we can begin with customization for production environment. Follow the step below:

<!-- @initApp @test -->
```shell
mkdir -p $MYSQL_DIR/prod
cd $MYSQL_DIR/prod

#initialize the customization
kinflate init

cat Kube-manifest.yaml
```

`Kube-manifest.yaml` should contain:

```
apiVersion: manifest.k8s.io/v1alpha1
kind: Manifest
metadata:
  name: helloworld
description: helloworld does useful stuff.
namePrefix: some-prefix
# Labels to add to all objects and selectors.
# These labels would also be used to form the selector for apply --prune
# Named differently than “labels” to avoid confusion with metadata for this object
objectLabels:
  app: helloworld
objectAnnotations:
  note: This is a example annotation
resources:
- deployment.yaml
- service.yaml
# There could also be configmaps in Base, which would make these overlays
configmaps: []
# There could be secrets in Base, if just using a fork/rebase workflow
secrets: []
recursive: true

```

Lets break this down:
- First step create a directory called `prod` and switches to that dir. We will keep all the resources related to production customization in this directory.
- `kinflate init` generates a kinflate manifest file called `Kube-manifest.yaml` that contains metadata about the customizations. You can think of this file as containing instructions which inflate will use to apply to generate the required configuration.

### Add resources

Lets add resource files that we want to `kinflate` to act on. Steps below add the three resources for MySQL.

<!-- @addResources @test -->
```shell

cd $MYSQL_DIR/prod

# add the MySQL resources
kinflate add resource ../secret.yaml
kinflate add resource ../service.yaml
kinflate add resource ../deployment.yaml

cat Kube-manifest.yaml 
```

`Kube-manifest.yaml`'s resources section should contain:

```
apiVersion: manifest.k8s.io/v1alpha1
....
....
resources:
- ../secret.yaml
- ../service.yaml
- ../deployment.yaml
```

Now we are ready to apply our first customization.

### NamePrefix Customization
We want MySQL resources to begin with prefix 'prod' in production environment.  Follow the steps below:

<!-- @customizeLabel @test -->
```shell

cd $MYSQL_DIR/prod

kinflate set nameprefix 'prod-'

cat Kube-manifest.yaml
```

`Kube-manifest.yaml` should have updated value of namePrefix field:

```
apiVersion: manifest.k8s.io/v1alpha1
.....
.....
namePrefix: prod-
objectAnnotations:
  note: This is a example annotation
.....


```

Lets break this down:
- `kinflate set nameprex <prefix>` updates the `namePrefix` directive to `prod` in the manifest file.
- Now if you view the `Kube-manifest.yaml`, you will see `namePrefix` directive updated. Editing the `namePrefix` directive in the file will also achieve the same thing.

At this point you can run `kinflate inflate` to generate name-prefixed configuration as shown below.

<!-- @genNamePrefixConfig @test -->
```shell
cd $MYSQL_DIR/prod

# lets generate name-prefixed resources 
kinflate inflate -f .
```

Output should contain:
```
apiVersion: v1
data:
  password: YWRtaW4=
kind: Secret
metadata:
  ....
  ....
  name: prod-mysql-pass-d2gtcm2t2k
---
apiVersion: v1
kind: Service
metadata:
  ....
  ....
  name: prod-mysql
spec:
  ....
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  ....
  ....
  name: prod-mysql
spec:
  selector:
  	....
	....

```


### Label Customization

We want resources in production environment to have certain labels so that we can query them by label selector. `kinflate` does not have `set label` command to add label, but we can edit `Kube-manifest.yaml` file under `prod` directory and add the production labels under `objectLabels` fields as highlighted below.

```
cd $MYSQL_DIR/prod

edit Kube-manifest.yaml

# Edit the objectLabels 
....
objectLabels:
  app: prod
....

```

At this point, running `kinflate inflate -f .` will generate MySQL configs with name-prefix 'prod-' and labels `app:prod`.

### Storage customization

Off the shelf MySQL uses `emptyDir` type volume, which gets wiped away if the MySQL Pod is recreated, and that is certainly not desirable for production environment. So
we want to use Persistent Disk in production. Kinflate lets you apply `patches` to the resources.

<!-- @customizeOverlay @test -->
```shell
cd $MYSQL_DIR/prod

# create a patch for persistent-disk.yaml
cat <<'EOF' > persistent-disk.yaml
apiVersion: apps/v1beta2 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: mysql
spec:
  template:
    spec:
      volumes:
      - name: mysql-persistent-storage
        emptyDir: null
        gcePersistentDisk:
          pdName: mysql-persistent-storage
EOF

# edit the manifest file to add the above patch or run following command
cat <<'EOF' >> $MYSQL_DIR/prod/Kube-manifest.yaml
patches:
- persistent-disk.yaml
EOF
```

Lets break this down:
- In the first step, we created a YAML file named `persistent-disk.yaml` to patch the resource defined in deployment.yaml
- Then we added `persistent-disk.yaml` to list of `patches` in `Kube-manifest.yaml`. `kinflate inflate` will apply this patch to the deployment resource with the name `mysql` as defined in the patch.

At this point, if you run `kinflate inflate -f .`, it will generate the Kubernetes config for production environement.
