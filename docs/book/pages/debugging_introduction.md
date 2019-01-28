{% panel style="danger", title="Proposal Only" %}
Many of the features and workflows.  The features that must be implemented
are tracked [here](https://github.com/kubernetes/kubectl/projects/7)
{% endpanel %}

# Debugging Workloads

This section covers the following tasks for view and debugging Workloads:

- Printing information about Resources
- Working with Container Logs, Files, Shells, etc
- Connecting to Services or Pods

Chapters:

- **[Summarizing Resource](debugging_summary.md):** Print Summaries of Resources
- **[Print Raw Resource](debugging_complete.md):** Printing Raw Resources
- **[Print Resource Fields](debugging_fields.md):** Print Resource Fields
- **[Describe Resources](debugging_debug.md):** Print Verbose Resource Information
- **[Queries and Options](debugging_queries.md):** Queries for Get and Describe
- **[Watching Resources](debugging_watch.md):** Continuously Watch for changes to Resources
- **[Container Logs](debugging_logs.md):** Print container logs
- **[Copying Container Files](debugging_files.md):** Copy files to / from containers
- **[Executing a Command in a Container](debugging_shell.md):** Execute a Command or get a Shell
  in a Container
- **[Accessing Services](debugging_proxy.md):** Connect to a Service running in the cluster
- **[Port Forward to Pods](debugging_forward.md):** Port forward to Pods in the cluster
  traffic through the apiserver.
