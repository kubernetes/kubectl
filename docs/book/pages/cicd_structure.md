{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

{% panel style="info", title="TL;DR" %}
- Structure Resource Config directories using a hierarchy that matches your environments, zones, versions, etc.
{% endpanel %}

# Project Structure

When creating a Kubernetes Project for your Application, there many possible was to structure it.  This
Chapter provides direction and common conventions for directory structure.

## Definitions

- **Application:** Collection of containerized Workloads that are run together.
- **Project:** One or more bundles of Resource Config that configure an Application for Kubernetes.
- **Bespoke:** Written by an end-user for themselves (e.g. run by the author).
- **Ready-Made:** Published for consumption by other users (e.g. run by users besides the author).

## Directory Structure

{% method %}

### Bespoke Projects

Bespoke Resource Config is for Projects that are run by the same organization that develops them.

Bespoke Resource Config must be structured so that it can be rolled out across multiple environments and availability
zones.

**Organizing Resource Config:** While the convention show here is purely optional, it is recommended for
across Kubernetes Projects consistency.

- Resource Configs `<project-name>/<environment>/<zone>/<component>/<resource-type>.yaml`
- Apply targets under `<project-name>/<environment>/<cluster>/kustomization.yaml`
- Reusable bases should be put under `<project-name>/bases/<component>` and
  `<project-name>/<environment>/bases/<component>`

**Best Practices:**
 
- Each *Base* should add a `namePrefix` and `commonLabels` to build up well structured Resources.
- Each *Environment Base* should set a `namespace` unique to that Project + Environment 

{% sample lang="yaml" %}

```bash
$ tree
.
├── bases # Shared Across all Environments
│   ├── kustomization.yaml # Used as a Base by Environment Bases
│   ├── backend # Backend Resource Config
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   ├── frontend # Frontend Resource Config
│   │   ├── deployment.yaml
│   │   ├── ingress.yaml
│   │   └── service.yaml
│   └── storage # Storage Resource Config
│       ├── service.yaml
│       └── statefulset.yaml
├── prod # Production Resource Config
│   ├── bases # Production specfic configuration
│   │   ├── kustomization.yaml # Used as a Base by Zones
│   │   ├── backend
│   │   │   └── deployment-patch.yaml
│   │   ├── frontend
│   │   │   └── deployment-patch.yaml
│   │   └── storage
│   │       └── statefulset-patch.yaml
│   ├── us-central # us-central specific configuration
│   │   ├── kustomization.yaml
│   │   ├── configmap-patch.yaml
│   │   └── backend
│   │       └── deployment-patch.yaml
│   ├── us-east # us-east specific configuration
│   │   └── kustomization.yaml 
│   └── us-west # us-west specific configuration
│       └── kustomization.yaml
├── staging
│   ├── bases # Staging specific configuration
│   │   ├── kustomization.yaml # Used as a Base by Zones
│   │   ├── backend
│   │   │   └── deployment-patch.yaml
│   │   ├── frontend
│   │   │   └── deployment-patch.yaml
│   │   └── storage
│   │       └── statefulset-patch.yaml
│   └── us-west
│       └── kustomization.yaml
└── test
    ├── bases # Test specific configuration
    │   ├── kustomization.yaml
    │   ├── configmap-patch.yaml
    └── us-west
        └── kustomization.yaml
```

{% endmethod %}

{% panel style="success", title="Applying Environment + Cluster" %}
While the directory structure contains the cluster, Apply won't read use this to determine the kubeconfig
context.  To Apply a specific cluster, add that cluster to the `kubectl config`, and
specify the corresponding context when running Apply.

```bash
$ kubectl apply -f myproject/prod/us-west --context us-west --wait
```

{% endpanel %}

{% method %}

### Ready-Made Projects

Ready-Made Resource Config is for Projects that are run by a different organization than develops them.

Ready-Made Config should be structured to support multiple concurrent versions / stability releases of a Project.
This structure may take the form of directories or branches (if using git).

**Organizing Resource Config:** While the convention show here is purely optional, it is recommended for
across Kubernetes Projects consistency.

- Resource Configs `<project-name>/<version | stability>/<resource-type>.yaml`
- Resource Configs `<project-name>/<version | stability>/kustomization.yaml`
- Reusable bases should be put under `<project-name>/bases/<component>` and
  `<project-name>/<version>/bases/<component>`

{% sample lang="yaml" %}

```bash
$ tree
.
├── bases
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
├── v1.0
│   ├── kustomization.yaml
│   └── bases
│       ├── kustomization.yaml
│       ├── backend
│       │   └── deployment-patch.yaml
│       ├── frontend
│       │   └── deployment-patch.yaml
│       └── storage
│           └── statefulset-patch.yaml
└── v1.1
    ├── kustomization.yaml
    └── bases
        ├── kustomization.yaml
        ├── backend
        │   └── deployment-patch.yaml
        ├── frontend
        │   └── deployment-patch.yaml
        └── storage
            └── statefulset-patch.yaml
```

{% endmethod %}

### Forking and Consuming Ready-Made Projects

Bespoke Projects may consume Ready-Made Projects by referencing them as bases.  These Projects may either
directly reference the Ready-Made Projects using their URLs, or may fork/clone them (e.g. using git).

When forking/cloning a Ready-Made Project, it may be put in the Bespoke Project bases, or in a location
shared by multiple Bespoke Projects.