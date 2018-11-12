{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

# Rolling out Across Environments

## Automatic Sequential Rollouts

**Immediate:** Rollout one environment immediately after the previous

**Delayed:** Pause between rollouts

## Manual Sequential Rollouts

Require humans to push

## Incorporating Application Metrics

Use Application metrics to identify issues and pause or rollback.