# Compare them

Before running kinflate on the two different instance
directories, review the directory
structure:

<!-- @listFiles @test -->
```
find $TUT_APP
```

<!-- @compareKinflateOutput -->
```
diff \
  <(kinflate inflate -f $TUT_APP/staging) \
  <(kinflate inflate -f $TUT_APP/production) |\
  more
```

Look at the output individually:

<!-- @runKinflateStaging @test -->
```
kinflate inflate -f $TUT_APP/staging
```

<!-- @runKinflateProduction @test -->
```
kinflate inflate -f $TUT_APP/production
```

Deploy them:

<!-- @deployStaging -->
```
kinflate inflate -f $TUT_APP/staging |\
    kubectl apply -f -
```

<!-- @deployProduction -->
```
kinflate inflate -f $TUT_APP/production |\
    kubectl apply -f -
```

<!-- @getAll -->
```
kubectl get all
```

Delete the resources:

<!-- @deleteStaging -->
```
kubectl delete configmap staging-acme-tut-map
kubectl delete service staging-acme-tut-service
kubectl delete deployment staging-acme-tut-deployment
```

<!-- @deleteProduction -->
```
kinflate inflate -f $TUT_APP/production |\
    kubectl delete -f -
```

__Next:__ [Lifecycle](../lifecycle.md)
