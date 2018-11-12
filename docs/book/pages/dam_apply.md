{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

{% panel style="info", title="TL;DR" %}
- Apply Creates, Updates and Deletes Resources in a cluster through running `kubectl apply` on Resource Config.
- Apply manages complexity such as ordering of operations and merging user defined and cluster defined state.
{% endpanel %}

# Apply

## Motivation

Apply will update a Kubernetes cluster to match state defined locally in files.

- Fully declarative - don't need to specify create, update or delete - just manage files
- Merges user owned state (e.g. Service `selector`) with state owned by the cluster (e.g. Service `clusterIp`)

## Definitions

- **Resources**: *Objects* in a cluster - e.g. Deployments, Services, etc.
- **Resource Config**: *Files* declaring the desired state for Resources - e.g. deployment.yaml.
  Resources are created and updated using Apply with these files.

*kubectl apply* Creates, Updates and Deletes Resources using Resource Config with a *kustomization.yaml*.

## Usage

{% method %}
kustomization.yaml defines and transforms a collection of Resources Configs.  `kubectl apply`
takes a list of directories containing `kustomization.yaml` files.

This `kustomization.yaml` file contains a list of Resource Configs that will be applied.
The file paths are relative to the `kustomization.yaml` file.
{% sample lang="yaml" %}

```yaml
# kustomization.yaml
apiVersion: v1beta1
kind: Kustomization
resources:
- deployment.yaml
```

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.15.4
```

{% endmethod %}

{% method %}
Users run Apply on directories containing `kustomization.yaml` files.

{% sample lang="yaml" %}
```bash
# Apply the Resource Config
$ kubectl apply -f .

# View the Resources
$ kubectl get -f .
```
{% endmethod %}

{% panel style="info", title="Multi-Resource Configs" %}
A single Resource Config file may declare multiple Resources separated by `\n---\n`.
{% endpanel %}

## CRUD Operations

### Creating Resources

Any Resources that do not exist and are declared in Resource Config when Apply is run will be Created.

### Updating Resources

Any Resources that already exist and are declared in Resource Config when Apply is run may be Updated.

**Added Fields**

Any fields that have been added to the Resource Config will be set on the Resource.

**Updated Fields** 
 
Any fields that contain different values for the fields specified locally in the Resource Config from what is
in the Resource will be updated by merging the Resource Config into the live Resource.  See [merging](dam_merge.md)
for more details.

**Deleted Fields**

Fields that were in the Resource Config the last time Apply was run, will be deleted from the Resource, and return to their default values.

**Unmanaged Fields**

Fields that were not specified in the Resource Config but are set on the Resource will be left unmodified.

### Deleting Resources

{% method %}
Any Resources that exist and are not declared in the Resource Config when Apply is run will be deleted
if and only if all of the following are true:

- A `pruneSelector` field is set in the `kustomization.yaml`
- The Resource matches the `pruneSelector` label selector
- If the Resource is Namespaced and the `namespace` field is set in the `kustomization.yaml` and the namespace matches
- The Resource is not a ConfigMap or Secret that is currently in use

ConfigMaps and Secrets are not deleted immediately because they may be in use by Workloads even if they are not
directly referenced in the newest Resource Config.

{% sample lang="yaml" %}
```yaml
# kustomization.yaml
apiVersion: v1beta1
kind: Kustomization
namespace: foo
pruneSelector:
  app: my-app
labels:
  app: my-app
```
{% endmethod %}

## Waiting for Rollouts to Complete

{% method %}

By default, Apply ends immediately after all changes have been sent to the apiserver, but before
the changes have been fully rolled out by the Controllers.

Using `--wait` will cause Apply to wait until the Controllers have rolled out all changes before
exiting.  When used with `--wait`, Apply will output the status of roll outs.
 
 
It is possible to Apply updates to Resources *sequentially* by using `--wait` and structuring the
Resource Config so Apply is called separately for each sequential collection of Resources.

{% sample lang="yaml" %}

```bash
kubectl apply -f /path/to/project1/ --wait
```

{% endmethod %}

## Continuously Applying a Project

{% method %}

A common pattern is to continuously run Apply on a Project so that Resource Config is automatically
Applied shortly after any changes are made.

The unix`watch` command allow users to periodically invoke Apply.

{% sample lang="yaml" %}

```bash
watch -n 60 kubectl apply -f https://github.com/myorg/myrepo
```

{% endmethod %}

## Resource Creation Ordering

Certain Resource Types may be dependent on other Resource Types being created first.  e.g. Namespaced
Resources on the Namespaces, RoleBindings on Roles, CustomResources on the CRDs, etc.

Apply sorts the Resources by Resource type to ensure Resources with these dependencies
are created in the correct order.
