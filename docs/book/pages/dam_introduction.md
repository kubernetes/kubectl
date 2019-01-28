{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

# Declarative Application Management

This section covers how to define a Project for running Workloads in a Kubernetes cluster.

Workloads in a cluster are configured through files called *Resource Config*.  These files are
typically checked into source control, and allow cluster state changes to be reviewed before they
are Applied, and audited.

Chapters:

- **[Apply](dam_apply.md):** Create, Update, Delete Kubernetes Resources
- **[Secret and ConfigMap Generation](dam_generators.md):** Generate Secrets and ConfigMaps from files and commands
- **[Container Image Tags](dam_images.md):** Set tags for Container Images across a Project
- **[Namespaces and Names](dam_namespaces.md):** Set namespaces and namePrefixes for Resources across a Project
- **[Labels and Annotations](dam_labels.md):** Set labels and annotations for Resources across a Project
- **[Config Reflection](dam_variables.md):** Populate Container Command Arguments and Environment Variables from
  Resource Config values from other Resources in the Project.
- **[Field Merge Semantics](dam_merge.md):** How Resources are updated from Resource Config.