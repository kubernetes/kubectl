{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

{% panel style="info", title="TL;DR" %}
- Port Forward connections to Pods running in a cluster 
{% endpanel %}

# Port Forward

## Motivation

Connect to ports of Pods running a cluster by port forwarding local ports.

{% method %}
## Forward Multiple Ports

Listen on ports 5000 and 6000 locally, forwarding data to/from ports 5000 and 6000 in the pod
{% sample lang="yaml" %}

```bash
$ kubectl port-forward pod/mypod 5000 6000
```

{% endmethod %}

---

{% method %}
## Pod in a Workload

Listen on ports 5000 and 6000 locally, forwarding data to/from ports 5000 and 6000 in a pod selected by the
deployment
{% sample lang="yaml" %}

```bash
$ kubectl port-forward deployment/mydeployment 5000 6000
```

{% endmethod %}

---

{% method %}
## Different Local and Remote Ports

Listen on port 8888 locally, forwarding to 5000 in the pod
{% sample lang="yaml" %}

```bash
$ kubectl port-forward pod/mypod 8888:5000
```

{% endmethod %}

---

{% method %}
## Random Local Port

Listen on a random port locally, forwarding to 5000 in the pod
{% sample lang="yaml" %}

```bash
$ kubectl port-forward pod/mypod :5000
```

{% endmethod %}

---

{% method %}
## Specify the Conainer

Specify the Container within a Pod running multiple containers.

- `-c <container-name>`
{% sample lang="yaml" %}

```bash
$ kubectl cp /tmp/foo <some-pod>:/tmp/bar -c <specific-container>
```

{% endmethod %}
  
---

{% method %}
## Namespaces

Set the Pod namespace by prefixing the Pod name with `<namespace>/` .

- `<pod-namespace>/<pod-name>:<path>`
{% sample lang="yaml" %}

```bash
$ kubectl cp /tmp/foo <some-namespace>/<some-pod>:/tmp/bar
```

{% endmethod %}
