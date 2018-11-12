{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

# Introduction

This section of the book covers how to build complex Projects that may contain different subcomponents
owned by multiple teams or organizations.  Apply facilitates defining projects as abstractions that
may be consumed and extended by other projects.

- Projects run in multiple Environments with slightly different configurations
- Projects run in multiple Clusters with slightly different configurations
- Projects that configure Ready-Made Applications published as Resource Config
- Meta-Projects that compose multiple Projects together.

Chapters:

* **[Bases and Variations](pages/project_variants.md):** Defining and Reusing Resource Config Abstractions
* **[Customizing Whitebox Bases](pages/project_whitebox.md):** Customizing Public Ready-Made Applications
* **[Composing Multiple Bases](pages/project_composition.md):** Advanced Resource Config Reuse
