# Apply

{% panel style="info", title="TL;DR" %}
- Apply Creates, Updates and Deletes Resources in a cluster through declarative files.
- Apply merges user owned state with state owned by the cluster when Updating Resources.
- Apply ensures the correct ordering when creating Resources.
{% endpanel %}

## Motivation

Apply will update a Kubernetes cluster to match state defined locally in files.

- Fully declarative - don't need to specify create, update or delete - just manage files
- Merges user owned state (e.g. Service `selector`) with state owned by the cluster (e.g. Service `clusterIp`)

## Definitions

- **Resources**: *Objects* in a cluster - e.g. Deployments, Services, etc.
- **Resource Config**: *Files* declaring the desired state for Resources - e.g. deployment.yaml.
  Resources are created and updated using Apply with these files.

*kubectl apply* Creates, Updates and Deletes Resources using Resource Config with a *apply.yaml*.

## Applying a Project

{% method %}
apply.yaml defines and transforms a collection of Resources Configs.  `kubectl apply`
takes a list of directories containing `apply.yaml` files.

This `apply.yaml` file contains a list of Resource Configs that will be applied.
The file paths are relative to the `apply.yaml` file.
{% sample lang="yaml" %}
```yaml
# apply.yaml
resources:
- some-service.yaml
- ../some-dir/some-deployment.yaml
```
{% endmethod %}

{% method %}
Users run Apply on directories containing `apply.yaml` files.

{% sample lang="yaml" %}
```bash
$ kubectl apply -f /path/to/project1/ -f /path/to/project2/
```
{% endmethod %}

{% panel style="info", title="Multi-Resource Configs" %}
A single Resource Config file may declare multiple Resources separated by `\n---\n`.
{% endpanel %}

## CRUD Operations

### Creating Resources

Any Resources that do not exist and are declared in Resource Config when Apply is run will be Created.

### Updating Resources

Any Resources that already exist and are declaraed in Resource Config when Apply is run may be updated.

**Added Fields**

Any fields that have been added to the Resource Config will be set on the Resource.

**Updated Fields** 
 
Any fields that contain different values for the fields specified locally in the Resource Config from what is
in the Resource will be updated by Apply.


Apply will overwrite any fields that are defined in the Resource Config, and whose values do not match.  For
example, if a Resource was modified by another source (e.g. an Autoscaler), it will be changed back to have the
values that are specified in the Resource Config.

**Deleted Fields**

If a field is deleted from the Resource Config:

- Fields that were *most recently* set by Apply, will be deleted from the Resource, and return to their default values.
- Fields that were set by Apply, but have been set by another source more recently (e.g. an Autoscaler), will
  be left unmodified.

**Unmanaged Fields**

Fields that have not been specified in the Resource Config but are set on the Resource will be left unmodified.


### Deleting Resources

{% method %}
Any Resources that exist and are not declared in the Resource Config when Apply is run will be deleted
if and only if all of the following are true:

- A `pruneSelector` field is set in the `apply.yaml`
- The Resource matches the `pruneSelector` label selector
- If the Resource is Namespaced and the `namespace` field is set in the `apply.yaml` and the namespace matches
- The Resource is not a ConfigMap or Secret that is currently in use

ConfigMaps and Secrets are not deleted immediately because they may be in use by Workloads even if they are not
directly referenced in the newest Resource Config.

{% sample lang="yaml" %}
```yaml
# apply.yaml
namespace: foo
pruneSelector:
  app: my-app
labels:
  app: my-app
```
{% endmethod %}

## Resource Creation Ordering

Certain Resource Types may be dependent on other Resource Types being created first.  e.g. Namespaced
Resources on the Namespaces, RoleBindings on Roles, CustomResources on the CRDs, etc.

Apply sorts the Resources by Resource type to ensure Resources with these dependencies
are created in the correct order.

## Blocking on Completion

{% method %}

By default, Apply exists immediately after all changes have been sent to the apiserver, but before
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


