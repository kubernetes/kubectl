# The base app

These resources work as is. To optionally confirm this,
apply it to your cluster.

<!-- @runKinflate @demo -->
```
kubectl apply -f $TUT_APP
```

Define some functions to query the server directly:

<!-- @funcGetAddress @env @test -->
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

<!-- @query @demo -->
```
tut_query tut-service peach
```

All done.  Clear the cluster for the next example.

<!-- @query @demo -->
```
kubectl delete deployment tut-deployment
kubectl delete service tut-service
kubectl delete configmap tut-map
```
