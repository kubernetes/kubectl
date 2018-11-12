# The Kubernetes Resource Model

{% panel style="info", title="TL;DR" %}
- A Kubernetes API has 2 parts - a Resource Type and a Controller
- Resources are object declared as json or yaml and written to a cluster
- Controllers asynchronously actuate Resources after they are stored
{% endpanel %}

## Resources

Instances of Kubernetes objects such as
[Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/),
[StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/),
[Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/),
[CronJobs](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) and
[DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) are called Resources.

It is important to understand the structure of Resources, as Resources are how users interact
with Kubernetes.

Users work with Resource APIs by declaring the desired state of Kubernetes Resources in
files called Resource Config.  *After* Resource Config is Applied to a cluster and the
request completes, a Controller actuates the API.

Resources are keyed by:

- **apiVersion**
- **kind**
- (metadata) **namespace**
- (metadata) **name**

{% panel style="info", title="Default Namespace" %}
If namespace is omitted from the Resource Config, the *default* namespace is used.
{% endpanel %}

{% method %}
### Resources Structure

Resources have the following components.

**TypeMeta:** Resource Type **apiVersion** and **kind**.

**ObjectMeta:** Resource **name** and **namespace** + other metadata (labels, annotations, etc).

**Spec:** the desired state of the Resource - declared by the user.

**Status:** the observed state of the object - recorded by the Controller.

Resource Config omits the Status.

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
Resources such as ConfigMaps and Secrets do not have a Status written by a Controller,
and as a result their Spec is implicit (i.e. they don't have a spec field).
{% endpanel %}

## Controllers

Controllers actuate Kubernetes APIs.  They observe the state of the system and look for
changes either to desired state of Resources (create, update, delete) or the system
(Pod or Node dies).

Controllers then make changes to the cluster to fulfill the intent specified by the user
(e.g. in Resource Config) or automation (e.g. changes from Autoscalers).

{% panel style="info", title="Asynchronous Actuation" %}
Because Controllers run asynchronously, issues such as a bad
Container Image or unschedulable Pods will not be present in the CRUD response.
Tools must facilitate watching the state of the system until changes are
completely actuated by Controllers.
{% endpanel %}

### Controller Structure

**Reconcile**

Controllers actuate Resources by reading the Resource they are Reconciling + related Resources.
Controllers **do not** Reconcile events, instead they compare the expected
cluster state to the observed cluster state, and make changes.

- Deployment Controller creates/deletes ReplicaSets
- ReplicaSet Controller creates/delete Pods
- Scheduler (Controller) writes Nodes to Pods
- Node (Controller) runs Containers specifid in Pods on the Node

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
of the system at the time the Controller is run, several different changes may be observed
and Reconciled together.  This is referred to as a **Level Based** system, whereas a system that
responds to each requested state would be an **Edge Based** system.
{% endpanel %}