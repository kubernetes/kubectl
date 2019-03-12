{% panel style="warning", title="Warning: Alpha Recommendations" %}
This chapter contains recommendations that are still being actively evaluated, and may
be changed in the future.
{% endpanel %}


{% panel style="info", title="TL;DR" %}
- Use separate branches for separate Environments to 
  - **Decouple release specific and operational specific changes** to Resource Config
  - Clean audit log of changes to the Environment
  - Facilitate Rollbacks to an Environment
{% endpanel %}

# Branch Structure Based Composition

The are several techniques for users to structure their Resource Config files.

| Type                                   | Summary               | Benefits                                           |
|----------------------------------------|-----------------------|----------------------------------------------------|
| [Directories](structure_directories.md)   | *Simplest approach*   | Easy to get started and understand               |
| **[Branches](structure_branches.md)**   | **More flexible**       | **Loose coupling between release specific and operation changes** |
| [Repositories](structure_repositories.md) | *Fine grain control*  | Isolated permissions model                         |

## Motivation

This chapter describes conventions for using **Branches** with Directories.

**Advantages:**

- Flexibility
  - Decouple release specific changes (e.g. images) from operational specific changes (e.g. cpu reservations)
  - Trigger Webhooks scoped to specific Environments
  - Simple to view changes for a specific Environment in *git*
  - Simple to rollback changes for a specific Environment using *git*

**Drawbacks:**

- More complicated to setup and configure
- Additional steps for merging release specific changes into and Environment branches

{% method %}

## Branch Structure

### Resource Config

The convention shown here should be changed and adapted as needed.

Structure:

- Use a Base (e.g. master, release-version, etc) branch for configuration tightly coupled to releasing new code
  - Looks like [Directories](structure_directories.md) approach
- Create separate branches for deploying to different Environments
  - Create a **new Directory for the operational overlays** - e.g. `release-<env>`
  - Uses Bases that point to Directories merged from the Base branch

Techniques:

- Add new required flags and environment variables to the Resource Config in the base branch at the
  time they are added to the code.
  - Will be rolled out when the code is rolled out.
- Adjust flags and configuration to the Resource Config in the Operational branch in the release directory
  - Will be rolled out immediately independent of releases
- Merge code from the Base branch to the Operational branches to perform a Rollout

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

**Operational Branches:**

```bash
$ tree
.
├── bases # From Base Branch
│   ├── ...
├── prod # From Base Branch
│   ├── ... 
└── release-prod # Prod release folder
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
├── release-staging # Staging release folder
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
├──release-test # Test release folder
│   └── us-west 
│       └── kustomization.yaml # Uses bases: ["../../test/us-west"]
└── test # From Base Branch
    └── ...
```

{% endmethod %}

{% panel style="success", title="Better Merging" %}

- If a user's application source code and Resource Config are both in the Base branch, user's may want
  to only merge the Resource Config.  This could be done using `git checkout` - e.g.
  `git checkout <base-branch> bases/ prod/`
  
- Instead of merging from the Base branch directly, users can create release branches of the Base.
  Alternatively, users can tag the Base branch commits as releases and check these out.

{% endpanel %}

{% method %}

## Alternative Branch Structure

An alternative to the above structure is to use branches similar to how *GitHub Pages* branches
functions - where code is not merged between branches and is similar to having a new repository.

This approach looks very similar to the [Repository Based Structure](structure_repositories.md), but
using branches instead of Repositories.

- Use a Base (e.g. master, release-version, etc) branch for configuration tightly coupled to releasing new code
  - Looks like [Directories](structure_directories.md) approach
- Create separate branches for deploying to different Environments
  - Create a **new Directory for the operational overlays** - e.g. `release-<env>`
  - Base Branch is never merge.  Operational overlays refer to Bases as remote urls.

Techniques:

- Add new required flags and environment variables to the Resource Config in the base branch at the
  time they are added to the code.
  - Will be rolled out when the code is rolled out.
- Adjust flags and configuration to the Resource Config in the Operational branch in the release directory
  - Will be rolled out immediately independent of releases
- Tag the base branch with releases
  - Operational branches use tagged references as their bases

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

**Operational Branches:**

```bash
$ tree
.
└── release-prod
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
└── release-staging
    └── us-west 
        └── kustomization.yaml
```

```bash
$ tree
.
└── release-test
    └── us-west 
        └── kustomization.yaml
```

{% endmethod %}