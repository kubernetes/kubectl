# Compare them

To deploy staging (or production) one would run a command like

> ```
> kinflate inflate -f $TUT_APP/staging | kubectl apply -f -
> ```

Review the instance differences:

<!-- @reviewDiffs @test -->
```
diff \
  <(kinflate inflate -f $TUT_APP/staging) \
  <(kinflate inflate -f $TUT_APP/production) | more
```

Look out output individually:

<!-- @runKinflateStaging @test -->
```
kinflate inflate -f $TUT_APP/staging
```

<!-- @runKinflateProduction @test -->
```
kinflate inflate -f $TUT_APP/production
```
