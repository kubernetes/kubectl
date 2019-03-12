{% panel style="info", title="TL;DR" %}
- A Kubernetes API has 2 parts - a Resource Type and a Controller
- Resources are objects declared as json or yaml and written to a cluster
- Controllers asynchronously actuate Resources after they are stored
{% endpanel %}

# The Kubernetes Resource Model

## Resources

Instances of Kubernetes objects such as
[Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/),
[CronJobs](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) and
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) are called **Resources**.

**Users work with Resource APIs by declaring them in files called Resource Config.**  Resource Config is
*Applied* (declarative Create/Update/Delete) to a Kubernetes cluster, and actuated by a *Controller*.

Resources are keyed by:

- **apiVersion** (API Type Group and Version)
- **kind** (API Type Name)
- **metadata.namespace** (Instance namespace)
- **metadata.name** (Instance name)

{% panel style="warning", title="Default Namespace" %}
If namespace is omitted from the Resource Config, the *default* namespace is used.  Users
should almost always explicitly specify the namespace for their Application using a
`kustomization.yaml`.
{% endpanel %}

{% method %}

### Resources Structure

Resources have the following components.

**TypeMeta:** Resource Type **apiVersion** and **kind**.

**ObjectMeta:** Resource **name** and **namespace** + other metadata (labels, annotations, etc).

**Spec:** the desired state of the Resource - intended state the user provides to the cluster.

**Status:** the observed state of the object - recorded state the cluster provides to the user.

Resource Config written by the user omits the Status field.

**Example Deployment Resource Config**
{% sample lang="yaml" %}

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
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

{% panel style="info", title="Spec and Status" %}
Resources such as ConfigMaps and Secrets do not have a Status,
and as a result their Spec is implicit (i.e. they don't have a spec field).
{% endpanel %}

## Controllers

Controllers actuate Kubernetes APIs.  They observe the state of the system and look for
changes either to desired state of Resources (create, update, delete) or the system
(Pod or Node dies).

Controllers then make changes to the cluster to fulfill the intent specified by the user
(e.g. in Resource Config) or automation (e.g. changes from Autoscalers).

**Example:** After a user creates a Deployment, the Deployment Controller will see
that the Deployment exists and verify that the corresponding ReplicaSet it expects
to find exists.  The Controller will see that the ReplicaSet does not exist and will
create one.

{% panel style="warning", title="Asynchronous Actuation" %}
Because Controllers run asynchronously, issues such as a bad
Container Image or unschedulable Pods will not be present in the CRUD response.
Tooling must facilitate processes for watching the state of the system until changes are
completely actuated by Controllers.  Once the changes have been fully actuated such
that the desired state matches the observed state, the Resource is considered *Settled*.
{% endpanel %}

### Controller Structure

**Reconcile**

Controllers actuate Resources by reading the Resource they are Reconciling + related Resources,
such as those that they create and delete.

**Controllers *do not* Reconcile events, rather they Reconcile the expected
cluster state to the observed cluster state at the time Reconcile is run.**

1. Deployment Controller creates/deletes ReplicaSets
1. ReplicaSet Controller creates/delete Pods
1. Scheduler (Controller) writes Nodes to Pods
1. Node (Controller) runs Containers specifid in Pods on the Node

**Watch**

Controllers actuate Resources *after* they are written by Watching Resource Types, and then
triggering Reconciles from Events.  After a Resource is created/updated/deleted, Controllers
Watching the Resource Type will receive a notification that the Resource has been changed,
and they will read the state of the system to see what has changed (instead of relying on
the Event for this information).

- Deployment Controller watches Deployments + ReplicaSets (+ Pods)
- ReplicaSet Controller watches ReplicaSets + Pods
- Scheduler (Controller) watches Pods
- Node (Controller) watches Pods (+ Secrets + ConfigMaps)

{% panel style="info", title="Level vs Edge Based Reconciliation" %}
Because Controllers don't respond to individual Events, but instead Reconcile the state
of the system at the time that Reconcile is run, **changes from several different events may be observed
and Reconciled together.**  This is referred to as a *Level Based* system, whereas a system that
responds to each event individually would be referred to as an *Edge Based* system.
{% endpanel %}