{% panel style="warning", title="Warning: Alpha Recommendations" %}
This chapter contains recommendations that are **still being actively evaluated, and are
expected to evolve.**

The intent of this chapter is to share the way kubectl developers are thinking about solving
this problem as they develop more solutions.

Before using these recommendations, carefully evaluate if they are right for your organization.
{% endpanel %}



{% panel style="info", title="TL;DR" %}
- Use separate branches for separate Environments to 
  - **Decouple version specific and live operational specific changes** to Resource Config
  - Clean audit log of changes to the Environment
  - Facilitate Rollbacks to an Environment by reverting commits
{% endpanel %}

# Branch Structure Based Composition

The are several techniques for users to structure their Resource Config files.

| Type                                   | Summary               | Benefits                                           |
|----------------------------------------|-----------------------|----------------------------------------------------|
| [Directories](structure_directories.md)   | *Simplest approach*   | Easy to get started and understand               |
| **[Branches](structure_branches.md)**   | **More flexible**       | **Loose coupling between version specific and live operational changes** |
| [Repositories](structure_repositories.md) | *Fine grain control*  | Isolated permissions model                         |

## Motivation

This chapter describes conventions for using **Branches** with Directories.

**Advantages:**

- Flexibility
  - Loose coupling between:
    - Version dependent changes made when changes are introduced to the code (e.g. adding flags, adding StatefulSets)
    - Operational dependent changes for tuning a production environment (e.g. cpu reservations, flag enabling features)
  - Simple to view changes for a specific Environment in *git*
  - Simple to rollback changes for a specific Environment using *git*
  - Trigger Webhooks scoped to specific Environments

**Drawbacks:**

- More complicated to setup and configure.  Requires and understanding of using git branches to manage Resource Config.

{% method %}

## Branch Structure

### Resource Config

The convention shown here should be changed and adapted as needed.


| Branch Type Name                                   | Purpose               | Examples |
|----------------------------------------|-----------------------|----|
| Base   | Contains shared Bases for all deploy environments and version dependent configuration.  When new code is added that requires additional configuration, this branch is updated.  **This Resource Config is never deployed directly.** | `master`, `release-1.14`, `i1026` |
| Deploy   | Contains relevant Config from the Base to be referred to as Bases from `kustomization.yaml`s.  Rather than directly modifying the directories from the Base branch, Deploy branches contain separate directories with customizations overlayed on the Base branch directories. **Resource Config only ever gets deployed from these branches.** | `deploy-test`, `deploy-staging`, `deploy-prod` |


Structure:

- Create (e.g. `master`, `app-version`, etc) a Base branch for version dependent Config changes which
  will be used as a Base for deployment dependent Config.
  - May have a similar structure to [Directories](structure_directories.md) approach
- Create separate Deploy branches for separate deployment Environments
  - Create a **new Directory in each branch containing only the deployment specific overlays** - e.g. `deploy-<env>`
  - Create `kustomization.yaml`'s and refer to the version dependent

Techniques:

- Add new required flags and environment variables to the Resource Config in the Base branch at the
  time they are added to the code.
  - Will be rolled out when the code is rolled out.
- Adjust flags and configuration to the Resource Config in the Deploy branch in the deploy directory
  - Will be rolled out immediately independent of versions
- Merge code from the Base branch to the Deploy branches to perform a Rollout

{% sample lang="yaml" %}

**Base Branch:**

```bash
$ tree
.
├── bases
│   ├── ...
├── prod
│   ├── bases 
│   │   ├── ...
│   ├── us-central
│   │   ├── kustomization.yaml
│   │   └── backend
│   │       └── deployment-patch.yaml
│   ├── us-east 
│   │   └── kustomization.yaml
│   └── us-west 
│       └── kustomization.yaml
├── staging
│   ├── bases 
│   │   ├── ...
│   └── us-west 
│       └── kustomization.yaml
└── test
    ├── bases 
    │   ├── ...
    └── us-west 
        └── kustomization.yaml
```

**Deploy Branches:**

```bash
$ tree
.
├── bases # From Base Branch
│   ├── ...
├── prod # From Base Branch
│   ├── ... 
└── deploy-prod # Prod deploy folder
    ├── us-central
    │   ├── kustomization.yaml # Uses bases: ["../../prod/us-central"]
    ├── us-east 
    │   └── kustomization.yaml # Uses bases: ["../../prod/us-east"]
    └── us-west 
        └── kustomization.yaml # Uses bases: ["../../prod/us-west"]
```

```bash
$ tree
.
├── bases # From Base Branch
│   ├── ...
├── deploy-staging # Staging deploy folder
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["../../staging/us-west"]
└── staging # From Base Branch
    └── ...
```

```bash
$ tree
.
├── bases # From Base Branch
│   ├── ...
├──deploy-test # Test deploy folder
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["../../test/us-west"]
└── test # From Base Branch
    └── ...
```

{% endmethod %}

{% panel style="success", title="Promotion from the Base branch to Deploy branches" %}

- If a user's application source code and Resource Config are both in the Base branch, user's may want
  to only merge the Resource Config.  This could be done using `git checkout` - e.g.
  `git checkout <base-branch> bases/ prod/`
  
- Instead of merging from the Base branch directly, users can create Deploy branches of the Base.
  Alternatively, users can tag the Base branch commits as deploys and check these out.

{% endpanel %}

{% method %}

## Alternative Branch Structure

An alternative to the above structure is to use branches similar to how *GitHub Pages* branches
functions - where code is not merged between branches and is similar to having a new repository.

This approach looks very similar to the [Repository Based Structure](structure_repositories.md), but
using branches instead of Repositories.

- Use a Base (e.g. master, deploy-version, etc) branch for configuration tightly coupled to releasing new code
  - Looks like [Directories](structure_directories.md) approach
- Create separate branches for deploying to different Environments
  - Create a **new Directory for the Deploy overlays** - e.g. `deploy-<env>`
  - Base Branch is never merge.  Deploy overlays refer to Bases as remote urls.

Techniques:

- Add new required flags and environment variables to the Resource Config in the base branch at the
  time they are added to the code.
  - Will be rolled out when the code is rolled out.
- Adjust flags and configuration to the Resource Config in the Deploy branch in the deploy directory
  - Will be rolled out immediately independent of deploys
- Tag the base branch with deploys
  - Deploy branches use tagged references as their bases

{% sample lang="yaml" %}

**Base Branch:**

```bash
$ tree
.
├── bases
│   ├── ...
├── prod
│   ├── bases 
│   │   ├── ...
│   ├── us-central
│   │   ├── kustomization.yaml
│   │   └── backend
│   │       └── deployment-patch.yaml
│   ├── us-east 
│   │   └── kustomization.yaml
│   └── us-west 
│       └── kustomization.yaml
├── staging
│   ├── bases 
│   │   ├── ...
│   └── us-west 
│       └── kustomization.yaml
└── test
    ├── bases 
    │   ├── ...
    └── us-west 
        └── kustomization.yaml
```

**Deploy Branches:**

```bash
$ tree
.
└── deploy-prod
    ├── us-central
    │   ├── kustomization.yaml
    ├── us-east 
    │   └── kustomization.yaml
    └── us-west 
        └── kustomization.yaml
```

```bash
$ tree
.
└── deploy-staging
    └── us-west 
        └── kustomization.yaml
```

```bash
$ tree
.
└── deploy-test
    └── us-west 
        └── kustomization.yaml
```

{% endmethod %}