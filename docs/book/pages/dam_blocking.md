# Rollout Status

## Apply

{% method %}
{% sample lang="yaml" %}
```bash
$ kubectl apply -f dir/ --wait
```
{% endmethod %}


## Checking on the Status of an existing Rollout

{% method %}
{% sample lang="yaml" %}
```bash
$ kubectl rollout status -f dir/
```
{% endmethod %}

## Conditions and Fields

### Rollout Completion

### Rollout Health