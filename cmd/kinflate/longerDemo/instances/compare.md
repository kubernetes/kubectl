# Compare them

Review the instance differences:

<!-- @reviewDiffs @test -->
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

<!-- @deployStaging @test -->
```
kinflate inflate -f $TUT_APP/staging |\
    kubectl apply -f -
```

<!-- @deployProduction @test -->
```
kinflate inflate -f $TUT_APP/production |\
    kubectl apply -f -
```

Query them:

<!-- @queryStaging @demo -->
```
tut_query staging-acme-tut-service pear
```
<!-- @queryProduction @demo -->
```
tut_query production-acme-tut-service apple
```
