# Creating Bases and Variations

{% panel style="info", title="TL;DR" %}
- Create Variants of a Project for different Environments.
- Customize Resource Config shared across multiple Projects.
{% endpanel %}

## Motivation

It is common for users to deploy several variants of the same project.

Examples:

- a project may be deployed to dev, test, staging, canary and production environments,
  but with variants between the environments.
- a project may be deployed to different clusters that are tuned differently or running
  different versions of the project.
  
Apply allows users to refer to another project as *Base*, and then apply additional customizations
to it.

Examples of changes between variants:

- Change replica count and resource
- Change image tag
- Change Environment Variables and Command Args 

## Referring to a Base

A project can refer by adding a path (relative to the `apply.yaml`) to `base` that
points to a directory containing another `apply.yaml` file.  This will automatically
add all of the Resources from the base project to the current project.

Bases can be:

- Relative paths from the `apply.yaml` - e.g. `../base`
- Urls - e.g. `github.com/kubernetes-sigs/kustomize/examples/multibases?ref=v1.0.6`

{% panel style="info", title="URL Syntax" %}
The Base URLs should follow
[hashicorp/go-getter URL format](https://github.com/hashicorp/go-getter#url-format).
{% endpanel %}

{% method %}
**Example:** Add the Resource Config from a base.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base

# ../base/apply.yaml
configMapGenerator:
- name: myJavaServerEnvVars
  literals:	
  - JAVA_HOME=/opt/java/jdk
  - JAVA_TOOL_OPTIONS=-agentlib:hprof
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      volumes:
      - configMap:
          name: myJavaServerEnvVars
        name: config-volume
```

**Applied:** The Resource that is Applied to the cluster

```yaml
# Unmodified Generated Base Resource
apiVersion: v1
kind: ConfigMap
metadata:
  name: myJavaServerEnvVars-k44mhd6h5f
data:
  JAVA_HOME: /opt/java/jdk
  JAVA_TOOL_OPTIONS: -agentlib:hprof
---
# Unmodified  Config Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      volumes:
      - configMap:
          name: myJavaServerEnvVars-k44mhd6h5f
        name: config-volume
```
{% endmethod %}

## Customizing Each Variant

When users have multiple similar projects with a shared base, they will want
to create variants that customize the original base.

### Customizing Pod Environment Variables

{% method %}
Customizing Pod Command arguments may be performed by generating different ConfigMaps
in each Variant and using the ConfigMap values in the Pod Environment Variables.

- Base uses ConfigMap data in Pods as Environment Variables
- Each Variant defines different ConfigMap data

**Use Case:** Different Environments (test, dev, staging, canary, prod) provide different Environment
Variables to a Pod.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
configMapGenerator:
- name: special-config
  literals:
  - special.how=very
  - special.type=charm

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        env:
        - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: special.how
```

**Applied:** The Resources that are Applied to the cluster

```yaml
# Generated Variant Resource
apiVersion: v1
kind: ConfigMap
metadata:
  name: special-config-82tc88cmcg
data:
  special.how: very
  special.type: charm
---
# Unmodified Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - env:
        - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              key: special.how
              name: special-config-82tc88cmcg
        image: nginx
        name: nginx
```
{% endmethod %}

See [ConfigMaps and Secrets](dam_generators.md).


### Customizing Pod Command Arguments

{% method %}
Customizing Pod Command arguments may be performed by generating different ConfigMaps
in each Variant and using the ConfigMap values in the Pod Command Arguments.

- Base uses ConfigMap data in Pods as Command Arguments
- Each Variant defines different ConfigMap data

**Use Case:** Different Environments (test, dev, staging, canary, prod) provide different Commandline
Arguments to a Pod.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
configMapGenerator:
- name: special-config
  literals:
  - special.how=very
  - special.type=charm

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: test-container
        image: k8s.gcr.io/busybox
        # Use the ConfigMap Environment Variables in the Command
        command: [ "/bin/sh", "-c", "echo $(SPECIAL_LEVEL_KEY) $(SPECIAL_TYPE_KEY)" ]
        env:
        - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: SPECIAL_LEVEL
        - name: SPECIAL_TYPE_KEY
          valueFrom:
            configMapKeyRef:
              name: special-config
              key: SPECIAL_TYPE
```

**Applied:** The Resources that are Applied to the cluster

```yaml
# Generated Variant Resource
apiVersion: v1
kind: ConfigMap
metadata:
  name: special-config-82tc88cmcg
data:
  special.how: very
  special.type: charm
---
# Unmodified Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - command:
        - /bin/sh
        - -c
        # Use the ConfigMap Environment Variables in the Command
        - echo $(SPECIAL_LEVEL_KEY) $(SPECIAL_TYPE_KEY)
        env:
        - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              key: SPECIAL_LEVEL
              name: special-config-82tc88cmcg
        - name: SPECIAL_TYPE_KEY
          valueFrom:
            configMapKeyRef:
              key: SPECIAL_TYPE
              name: special-config-82tc88cmcg
        image: k8s.gcr.io/busybox
        name: test-container
```

{% endmethod %}
See [ConfigMaps and Secrets](dam_generators.md).

### Customizing Image Tags

{% method %}
Customizing the Image Tag run in each Variant can be performed by specifying `imageTags`
in each Variant `apply.yaml`.

**Use Case:** Different Environments (test, dev, staging, canary, prod) can use images with different tags.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
imageTags:
  - name: nginx
    newTag: 1.8.0

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
```

**Applied:** The Resource that is Applied to the cluster

```yaml
# Modified Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      # The image has been changed to include the tag
      - image: nginx:1.8.0
        name: nginx
```
{% endmethod %}

See [Image Tags](dam_images.md).

### Customizing Namespace

{% method %}
Customizing the Namespace in each Variant can be performed by specifying `namespace` in each
Variant `apply.yaml`.

**Use Case:** Different Environments (test, dev, staging, canary, prod) run in different Namespaces.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
namespace: test

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
```

**Applied:** The Resource that is Applied to the cluster

```yaml
# Modified Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
  # Namespace has been set
  namespace: test
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx

```
{% endmethod %}

See [Namespaces and Names](dam_namespaces.md).

### Customizing Resource Name Prefixes

{% method %}
Customizing the Name by adding a prefix in each Variant can be performed by specifying `namePrefix` in each
Variant `apply.yaml`.

**Use Case:** Different Environments (test, dev, staging, canary, prod) have different Naming conventions.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
namePrefix: test-

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
```

**Applied:** The Resource that is Applied to the cluster

```yaml
# Modified Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  # Name has been prefixed with the environment
  name: test-nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
```
{% endmethod %}

See [Namespaces and Names](dam_namespaces.md).

### Customizing Arbitrary Fields with Overlays

{% method %}
Arbitrary fields may be added, changed, or deleted by supplying *Overlays* against the
Resources provided by the base.  Overlays are sparse Resource definitions that
allow arbitrary customizations to be performed without requiring a base to expose
the customization as a template.

Overlays require the *Group, Version, Kind* and *Name* of the Resource to be specified, as
well as any fields that should be set on the base Resource.  Overlays are applied using
*StrategicMergePatch*.

**Use Case:** Different Environments (test, dev, staging, canary, prod) require fields such as
replicas or resources to be overridden.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
patchesStrategicMerge:
- overlay.yaml

# overlay.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  # override replicas
  replicas: 3
  template:
    spec:
      containers:
      - name: nginx
        # override resources
        resources:
          limits:
            cpu: "1"
          requests:
            cpu: "0.5"

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        resources:
          limits:
            cpu: "0.2"
          requests:
            cpu: "0.1"        
```

**Applied:** The Resource that is Applied to the cluster

```yaml
# Overlayed Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  # replicas field has been added
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
        # resources have been overridden
        resources:
          limits:
            cpu: "1"
          requests:
            cpu: "0.5"
```
{% endmethod %}

{% panel style="info", title="Overlay URLs" %}
Like Bases, Overlays may also be URLs and should follow the
[hashicorp/go-getter URL format](https://github.com/hashicorp/go-getter#url-format).
{% endpanel %}


### Customizing Arbitrary Fields with JsonPatch

{% method %}
Arbitrary fields may be added, changed, or deleted by supplying *Json 6902 Patches* against the
Resources provided by the base.

**Use Case:** Different Environments (test, dev, staging, canary, prod) require fields such as
replicas or resources to be overridden.

Json 6902 Patches are [rfc6902](https://tools.ietf.org/html/rfc6902) patches that are applied
to resources.  Patches require the *Group, Version, Kind* and *Name* of the Resource to be
specified in addition to the Patch.  Patches offer a number of powerful imperative operations
for modifying the base Resources.

{% sample lang="yaml" %}
**Input:** The apply.yaml file

```yaml
# apply.yaml
bases:
- ../base
patchesJson6902:
- target:
    group: apps
    version: v1
    kind: Deployment
    name: nginx-deployment
  path: patch.yaml

# patch.yaml
- op: add
  path: /spec/replicas
  value: 3

# ../base/apply.yaml
resources:
- deployment.yaml

# ../base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
```

**Applied:** The Resource that is Applied to the cluster

```yaml
# Patched Base Resource
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: nginx
  name: nginx-deployment
spec:
  # replicas field has been added
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
```
{% endmethod %}

### Adding Resources to a Base

Additional Resources not specified in the Base may be added to Variants by
Variants specifying them as `resources` in their `apply.yaml`.
