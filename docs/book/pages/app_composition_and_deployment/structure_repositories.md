{% panel style="warning", title="Warning: Alpha Recommendations" %}
This chapter contains recommendations that are **still being actively evaluated, and are
expected to evolve.**

The intent of this chapter is to share the way kubectl developers are thinking about solving
this problem as they develop more solutions.

Before using these recommendations, carefully evaluate if they are right for your organization.
{% endpanel %}


{% panel style="info", title="TL;DR" %}
- Finer grain management using separate repos for separate Team 
  - Separate permissions for committing changes to separate environments
  - Separate Issue, Project and PR tracking
{% endpanel %}

# Repository Structure Based Composition

The are several techniques for users to structure their Resource Config files.

| Type                                        | Summary               | Benefits                                           |
|---------------------------------------------|-----------------------|----------------------------------------------------|
| [Directories](structure_directories.md)        | *Simplest approach*   | Easy to get started and understand               |
| [Branches](structure_branches.md)   | *More flexible*       | Loose coupling between version specific and live operational changes |
| **[Repositories](structure_repositories.md)** | **Fine grain control**  | **Isolated permissions model**                         |

## Motivation

This chapter describes conventions for using **Repositories** with Directories.

**Advantages:**

- **Isolation between teams** managing separate Environments
  - Permissions
- **Fine grain control** over
  - PRs
  - Issues
  - Projects
  - Automation
   
**Drawbacks:**

- Tools designed to work with files and directories don't work across Repositories
- Complicated to setup and manage
- **Harder to reason about the system as a whole**
  - State spread across multiple Repositories

## Directory Structure

{% method %}

### Resource Config

The convention shown here should be changed and adapted as needed.

| Repo Type Name                                   | Purpose               | Examples |
|----------------------------------------|-----------------------|----|
| Base   | Contains shared Bases for all deploy environments and version dependent configuration.  When new code is added that requires additional configuration, this repository is updated.  **This Resource Config is never deployed directly.** | `app-name` |
| Deploy   | Does not contain Config from the Base, rather refers to the Base Config remotely through the git url.  Deploy repositories contain directories with similar structure to the Base directories, but instead contain customizations overlayed on the remote Bases. **Resource Config only ever gets deployed from these Repositories.** | `app-name-test`, `app-name-staging`, `app-name-prod` |


Structure:

- Create a Base Repository for shared configuration
  - Looks like [Directories](structure_directories.md) approach
- For each **separate Environment, create a separate Deploy Repository**
  - Remotely reference the Base Repository in from the Deploy Repository

Techniques:

- Use techniques described in [Directories](structure_directories.md) and [Branches](structure_branches.md)

{% sample lang="yaml" %}


**Base Repository:**

```bash
$ tree
.
├── bases # Used as a Base only
│   ├── kustomization.yaml
│   ├── backend
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   ├── frontend
│   │   ├── deployment.yaml
│   │   ├── ingress.yaml
│   │   └── service.yaml
│   └── storage
│       ├── service.yaml
│       └── statefulset.yaml
├── prod # Production
│   ├── bases 
│   │   ├── kustomization.yaml # Uses bases: ["../../bases"]
│   │   ├── backend
│   │   │   └── deployment-patch.yaml # Production Env specific backend overrides
│   │   ├── frontend
│   │   │   └── deployment-patch.yaml # Production Env specific frontend overrides
│   │   └── storage
│   │       └── statefulset-patch.yaml # Production Env specific storage overrides
│   ├── us-central
│   │   ├── kustomization.yaml # Uses bases: ["../bases"]
│   │   └── backend
│   │       └── deployment-patch.yaml # us-central cluster specific backend overrides
│   ├── us-east 
│   │   └── kustomization.yaml # Uses bases: ["../bases"]
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["../bases"]
├── staging # Staging
│   ├── bases 
│   │   ├── kustomization.yaml # Uses bases: ["../../bases"]
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["../bases"]
└── test # Test
    ├── bases 
    │   ├── kustomization.yaml # Uses bases: ["../../bases"]
    └── us-west 
        └── kustomization.yaml # Uses bases: ["../bases"]
```

**Team Repositories:**

```bash
# sre team
$ tree
.
├── prod
│   ├── us-central
│   │   ├── kustomization.yaml # Uses bases: ["https://<your-repo>/prod/us-central?ref=<prod-release>"]
│   ├── us-east 
│   │   └── kustomization.yaml # Uses bases: ["https://<your-repo>/prod/us-east?ref=<prod-release>"]
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["https://<your-repo>/prod/us-west?ref=<prod-release>"]
```

```bash
# qa team
$ tree
.
├── staging # Staging
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["https://<your-repo>/staging/us-west?ref=<staging-release>"]
```

```bash
# dev team
$ tree
.
└── test # Test
    └── us-west 
        └── kustomization.yaml # Uses bases: ["https://<your-repo>/test/us-west?ref=<test-release>"]
```

{% endmethod %}
