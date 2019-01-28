{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

{% panel style="info", title="TL;DR" %}
- Kubectl is the Kubernetes CLI.
- Kubectl has different command groups for different types of user workflows.
{% endpanel %}

# Kubectl

Kubectl is the Kubernetes CLI and used to manage Resources.

## Command Families

While Kubectl has many different commands, they fall into only a few categories.

- Declaratively Creating, Updating, Deleting Resources (Apply)
- Debugging Workloads and Reading Cluster State
- Managing the cluster itself
- Porcelain commands for working with Resources

## Declaratively Creating, Updating, Deleting Resources (Apply)

Creating, Updating and Deleting Resources is done through declarative files called Resource Config
and the Kubectl *Apply* command.  This command reads a local (or remote) file structure and modifies
cluster state to reflect the declared intent.

{% panel style="info", title="Apply" %}
Apply is the preferred mechanism for managing Resources in a Kubernetes cluster.
{% endpanel %}

## Debugging Workloads and Reading Cluster State

Users will need to debug and view Workloads running in a cluster.  Kubectl supports debugging
by providing commands for:

- printing state and information about Resources
- printing Container logs
- printing cluster events
- exec or attaching to a Container
- copying files from Containers in the cluster to a user's filesystem

## Cluster Management

On occasion, users may need to perform operations to the Nodes of cluster.  Kubectl supports
commands to drain Workloads from a Node so that it can be decommission or debugged.

## Porcelain

Users may find using Resource Config overly verbose for *Development* and prefer to work with
the cluster *imperatively* with a shell-like workflow.  Kubectl offers porcelain commands for
generating and modifying Resources.

- generating + creating Resources such as Deployments, StatefulSets, Services, ConfigMaps, etc
- setting fields on Resources
- editing (live) Resources in a text editor

{% panel style="danger", title="Porcelain For Dev Only" %}
Porcelain commands are time saving for experimenting with workloads in a dev cluster, but
shouldn't be used for production.
{% endpanel %}
