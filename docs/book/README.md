{% panel style="info", title="TL;DR" %}
- Kubernetes runs Containerized Workloads in a cluster
- The Kubectl Book explains Kubernetes tools and workflows
{% endpanel %}

# Introduction

The goal of this book is to document how to configure, deploy and manage their containerized
Workloads in Kubernetes using Kubectl.

It covers the following topics:

- Introduction to Kubernetes Workload APIs & Kubectl
- Declarative Configuration
- Deployment Techniques
- Printing information about Workloads
- Debugging Workloads
- Imperative Porcelain Commands

## Overview

Kubernetes is a set of APIs to run containerized Workloads in a cluster.

Users define API objects (i.e. Resources) in files which are typically checked into source control.
They then use kubectl to Apply (i.e. create, update, delete) to update cluster state.

### Pods
 
Containers are run in [*Pods*](https://kubernetes.io/docs/concepts/workloads/pods/pod-overview/) which are
scheduled to run on *Nodes* (i.e. worker machines) in a cluster.

Pods run a *single replica* of an Application and provide:

- Compute Resources (cpu, memory, disk)
- Environment Variables
- Readiness and Health Checking
- Network (IP address shared by containers in the Pod)
- Mounting Shared Configuration and Secrets
- Mounting Storage Volumes
- Initialization

{% panel style="warning", title="Multi Container Pods" %}
Multiple replicas of an Application should be created using a Workload API to manage
creation and deletion of Pod replicas using a PodTemplate.

In some cases a Pod may contain multiple Containers forming a single instance of an Application.  These
containers may coordinate with one another through shared network (IP) and storage.
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

{% panel style="success", title="API Abstraction Layers" %}
High-level Workload APIs may manage lower-level Workload APIs instead of directly managing Pods
(e.g. Deployments manage ReplicaSets).
{% endpanel %}

### Service Discovery and Load Balancing

Service discovery and Load Balancing may be managed by a *Service* object.  Services provide a single
virtual IP address and dns name load balanced to a collection of Pods matching Labels.

{% panel style="info", title="Internal vs External Services" %}
- [Services Resources](https://kubernetes.io/docs/concepts/services-networking/service/)
  (L4) may expose Pods internally within a cluster or externally through an HA proxy.
- [Ingress Resources](https://kubernetes.io/docs/concepts/services-networking/ingress/) (L7)
  may expose URI endpoints and route them to Services.
{% endpanel %}

### Configuration and Secrets

Shared Configuration and Secret data may be provided by ConfigMaps and Secrets.  This allows
Environment Variables, Commandline Arguments and Files to be loosely injected into
the Pods and Containers that consume them.

{% panel style="info", title="ConfigMaps vs Secrets" %}
- [ConfigMaps](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)
  are for providing non-sensitive data to Pods.
- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
  are for providing sensitive data to Pods.
{% endpanel %}