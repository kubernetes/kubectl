{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

# Introduction

This section covers how to automate and structure rolling out Workload changes through Resource Config.

A typical rollout may include:

1. Building source artifacts into Containers
1. Modifying Resource Config
1. Auditing and Reviewing changes before they are actuated in a specific cluster or environment
1. Actuating the Resource Config for a given cluster or environment
