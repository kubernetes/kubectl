{% panel style="info", title="TL;DR" %}
- Resource Config is stored in one or more git repositories
- Directory hierarchy, git branches and git repositories may be used for loose coupling
{% endpanel %}


# Resource Config Structure

The chapters in this section cover how to structure Resource Config using git.

Users may start with a pure Directory Hierarchy approach, and migrate to include Branches
and / or Repositories as part of the structure.

## Background

Tools:

- *Bases:* provide **common or shared Resource Config to be factored out** that can be
  imported into multiple projects.
- *Overlays and Customizations:* tailor **common or shared Resource Config to be modified** to
  a specific application, environment or purpose.

Applied Concepts:

- Resource Config may be initially structured using Directory Hierarchy for organization.
  - Use Bases with Overlays / Customizations for factoring across Directories
- Different Deployment environments for the same app may be loosely coupled, using separate **Branches for separate
  environments**.
  - Use Bases with Overlays / Customization for factoring across Branches
- Teams owning separate (shared) Config may be loosely coupled using separate **Repositories for
  separate teams**.
  - Use Bases with Overlays / Customization for factoring across Repositories


Table:

| Technique                                   | Decouple Changes               | Best For                                           | Workflow |
|---------------------------------------------|-----------------------|----------------------------------------------------|----------|
| [Directories](structure_directories.md)     | NA                    | Simple organizational and deployment structure.    | Changes are immediately propagated globally.  |
| [Branches](structure_branches.md)           | *Across Environments*       | Promoting changes across Environments. | Changes to Bases are pushed across multiple linear stages by the owners of Bases. |
| [Repositories](structure_repositories.md)   | *Across Teams*              | Fetching changes across config shared across Teams. | Changes to Bases are pulled by individual consumers into their projects. |
