# Resources

Make three resource files.  The resources have name
fields with these values:

* tut-deployment
* tut-map
* tut-service

_tut-_ just stands for _tutorial_.


<!-- @mkAppDir @test -->
```
TUT_APP=$TUT_DIR/tuthello
/bin/rm -rf $TUT_APP
mkdir -p $TUT_APP
```

<!-- @writeDeploymentYaml @test -->
```
cat <<EOF >$TUT_APP/deployment.yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: tut-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        deployment: tuthello
    spec:
      containers:
      - name: tut-container
        image: monopole/tuthello:1
        command: ["/tuthello",
                  "--port=8080",
                  "--enableRiskyFeature=\$(ENABLE_RISKY)"]
        ports:
        - containerPort: 8080
        env:
        - name: ALT_GREETING
          valueFrom:
            configMapKeyRef:
              name: tut-map
              key: altGreeting
        - name: ENABLE_RISKY
          valueFrom:
            configMapKeyRef:
              name: tut-map
              key: enableRisky
EOF
```

This deployment takes some parameters from a map:

<!-- @writeMapYaml @test -->
```
cat <<EOF >$TUT_APP/configMap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tut-map
data:
  altGreeting: "Good Morning!"
  enableRisky: "false"
EOF
```

A named label selector (a service) is used to
query the deployment.

<!-- @writeServiceYaml @test -->
```
cat <<EOF >$TUT_APP/service.yaml
kind: Service
apiVersion: v1
metadata:
  name: tut-service
spec:
  selector:
    deployment: tuthello
  type: LoadBalancer
  ports:
  - protocol: TCP
    port: 8666
    targetPort: 8080
EOF
```

Review the app definition so far:

<!-- @listFiles @test -->
```
find $TUT_APP
```
