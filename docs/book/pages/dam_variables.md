{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

{% panel style="info", title="TL;DR" %}
- Inject the values of other Resource Config fields into Pod Env Vars and Command Args with `vars`.
{% endpanel %}

# Config Reflection

## Motivation

Pods may need to refer to the values of Resource Config fields.  For example
a Pod may take the name of Service defined in the Project as a command argument.
Instead of hard coding the value directly into the PodSpec, it is preferable
to use a `vars` entry to reference the value by path.  This will ensure
if the value is updated or transformed by the `kustomization.yaml` file, the
value will be propagated to where it is referenced in the PodSpec. 

## Vars

The `vars` section contains variable references to Resource Config fields within the project.  They require
the following to be defined:

- Resource Kind
- Resource Version
- Resource name
- Field path

{% method %}

**Example:** Set the Pod command argument to the value of a Service name.

Apply will set the resolve $(BACKEND_SERVICE_NAME) to a value using the path
specified in `vars`.

{% sample lang="yaml" %}
**Input:** The kustomization.yaml, deployment.yaml and service.yaml files

```yaml
# kustomization.yaml
vars:
- name: BACKEND_SERVICE_NAME
  objref:
    kind: Service
    name: backend-service
    apiVersion: v1
  fieldref:
    fieldpath: metadata.name
resources:
- deployment.yaml
- service.yaml

# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: curl-deployment
  labels:
    app: curl
spec:
  selector:
    matchLabels:
      app: curl
  template:
    metadata:
      labels:
        app: curl
    spec:
      containers:
      - name: curl
        image: ubuntu
        command: ["curl", "$(BACKEND_SERVICE)"]
        
# service.yaml
kind: Service
apiVersion: v1
metadata:
  name: backend-service
spec:
  selector:
    app: backend
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
```

**Applied:** The Resources that are Applied to the cluster

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 9376
  selector:
    app: backend
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: curl
  name: curl-deployment
spec:
  selector:
    matchLabels:
      app: curl
  template:
    metadata:
      labels:
        app: curl
    spec:
      containers:
      - command:
        - curl
        # $(BACKEND_SERVICE_NAME) has been resolved to
        # backend-service
        - backend-service
        image: ubuntu
        name: curl
```
{% endmethod %}

{% panel style="info", title="Referencing Variables" %}
Variables are intended to allow Pods to access Resource Config values from.  They are
**not** intended as a general templating mechanism.  Overriding values should be done with
patches instead of variables.  See [Bases and Variations](project_variants.md).
{% endpanel %}
