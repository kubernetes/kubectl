{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

{% panel style="info", title="TL;DR" %}
- Kubernetes runs Containerized Workloads in a cluster
- The Kubectl Book explains Kubernetes tools and workflows
{% endpanel %}

# Introduction

The goal of this book is to document how users should configure, deploy and manage their 
containerized Workloads in Kubernetes.

It is broken into the following sections:

- Introduction to Kubernetes & Kubectl
- How to configure Applications through Resource Config
- How to debug Applications using Kubectl
- How to configure Projects composed of multiple Applications
- How to rollout changes with CICD
- How to perform cluster maintenance operations using Kubectl

## Background

Kubernetes is a set of APIs to run containerized Workloads in a cluster.

Users define API objects (i.e. Resources) in files checked into source control, and use kubectl
to Apply (i.e. create, update, delete) configuration files to cluster Resources.

### Pods
 
Containers are run in [*Pods*](https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/) which are scheduled to *Nodes* (i.e. worker machines) in a cluster.

Pods provide the following a Pod running a *single instance* of an Application:

- Compute Resources (cpu, memory, disk)
- Environment Variables
- Readiness and Health Checking
- Network (IP address shared by containers in the Pod)
- Mounting Shared Configuration and Secrets

{% panel style="warning", title="Multi Container Pods" %}
Multiple identical instances of an Application should be run by creating multiple copies of
the same Pod using a Workload API.

A Pod may contain multiple Containers which are a single instance of an Application.  These
containers may coordinate with one another through shared network (IP) and files.
{% endpanel %}

### Workloads

Pods are typically managed by higher level abstractions that handle concerns such as
replication, identity, persistent storage, custom scheduling, rolling updates, etc.

The most common out-of-the-box Workload APIs (manage Pods) are:

- [Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) (Stateless Applications)
  - replication + roll outs
- [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) (Stateful Applications)
  - replication + roll outs + persistent storage + identity
- [Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/) (Batch Work)
  - run to completion
- [CronJobs](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) (Scheduled Batch Work)
  - scheduled run to completion
- [DaemonSets](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/) (Per-Machine)
  - per-Node scheduling

{% panel style="info", title="Abstraction Layers" %}
High-level Workload APIs may manage lower-level Workload APIs instead of directly managing Pods
(e.g. Deployments manage ReplicaSets).
{% endpanel %}

### Service Discovery and Load Balancing

Service discovery and Load Balancing is managed by a *Service* object.  Services provide a single
IP address and dns name to talk to a collection of Pods.

{% panel style="info", title="Internal vs External Services" %}
- [Services Resources](https://kubernetes.io/docs/concepts/services-networking/service/)
  (L4) may expose Pods internally within a cluster or externally through an HA proxy.
- [Ingress Resources](https://kubernetes.io/docs/concepts/services-networking/ingress/) (L7)
  may expose URI endpoints and route them to Services.
{% endpanel %}

### Configuration and Secrets

Shared Configuration and Secret data may be provided by Secrets and ConfigMaps.  This allows
Environment Variables, Commandline Arguments and Files to be setup and decoupled from
the Pods and Containers.

{% panel style="info", title="ConfigMaps vs Secrets" %}
- [ConfigMaps](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
  are for providing non-sensitive data to Pods.
- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
  are for providing sensitive data to Pods.
{% endpanel %}