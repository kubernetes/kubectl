# Edit tool

> _kinflate is a text editor / linter!_

Kinflate's basic function is to read manifests and resources to create new YAML.

It also offers some basic manifest file operations.

Make a new workspace:

<!-- @workspace @test -->
```
TUT_EDITS=$TUT_TMP/edits
/bin/rm -rf $TUT_EDITS
mkdir -p $TUT_EDITS
pushd $TUT_EDITS
```

Create a manifest:

<!-- @init @test -->
```
kinflate init
```

<!-- @showIt @test -->
```
clear
more Kube-manifest.yaml
```

Write a resource file:

<!-- @writeResource @test -->
```
cat <<EOF >configMap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: tut-map
data:
  altGreeting: "Good Morning!"
  enableRisky: "false"
EOF
```

<!-- @ls @test -->
```
ls
```

Add it to the manifest:

<!-- @addResource @test -->
```
kinflate add resource configMap.yaml
```

<!-- @showItAgain @test -->
```
more Kube-manifest.yaml
```

Add one that isn't there:

<!-- @addNoResource @test -->
```
kinflate add resource does_not_exist.yaml
```

Try to reinit
<!-- @initAgain @test -->
```
kinflate init
```


<!-- @allDone @test -->
```
popd
```
