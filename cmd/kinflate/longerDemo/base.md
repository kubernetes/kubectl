# The base app

These resources work as is. To optionally confirm this,
apply it to your cluster.

<!-- @runKinflate -->
```
kubectl apply -f $TUT_APP
```

<!-- @showResources -->
```
kubectl get deployments
```

Define some functions to query the server directly:

<!-- @funcGetAddress -->
```
function tut_getServiceAddress {
  local name=$1
  local tm='{{range .spec.ports -}}{{.nodePort}}{{end}}'
  local nodePort=$(\
    kubectl get -o go-template="$tm" service $name)
  echo $($MINIKUBE_HOME/minikube ip):$nodePort
}

function tut_query {
  local addr=$(tut_getServiceAddress $1)
  curl --fail --silent --max-time 3 $addr/$2
}
```

Query it:

<!-- @query -->
```
tut_query tut-service peach
```

All done.  Clear the cluster for the next example.

<!-- @query -->
```
kubectl delete deployment tut-deployment
kubectl delete service tut-service
kubectl delete configmap tut-map
```

__Next:__ [Describe the app with a Manifest](manifest.md)
