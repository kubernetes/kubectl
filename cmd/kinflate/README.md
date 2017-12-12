# kinflate

_TODO: flesh out this placeholder documentation_
 
`kinflate` is a kubernetes cluster configuration utility,
a prototype of ideas discussed in [this doc][DAM].

It accepts one or more file system path arguments,
reads the content thereof, and emits kubernetes
[resource] yaml to stdout, suitable for piping 
into [kubectl] for [declarative] application to a
kubernetes cluster.

For example, if your current working directory contained
> ```
> mycluster/
>   Kube-Manifest.yaml
>   deployment.yaml
>   service.yaml
>   instances/
>     testing/
>       Kube-Manifest.yaml
>       deployment.yaml
>     prod/
>       Kube-Manifest.yaml
>       deployment.yaml
>  ...
> ```

then the command 
```
kinflate ./mycluster/instances | kubectl apply -f -
```
would modify your cluster per the resources
generated from the manifests and associated
resource files found in `mycluster`.

[resource]: https://kubernetes.io/docs/api-reference/v1.7/
[DAM]: https://goo.gl/T66ZcD
[yaml]: http://www.yaml.org/start.html
[kubectl]: https://kubernetes.io/docs/user-guide/kubectl-overview/
[declarative]: https://kubernetes.io/docs/tutorials/object-management-kubectl/declarative-object-management-configuration/
