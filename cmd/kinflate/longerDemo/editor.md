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
cat Kube-manifest.yaml
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

<!-- @confirmIt @test -->
```
grep configMap.yaml Kube-manifest.yaml
```

Attempt to add a missing resource; kinflate should complain.

<!-- @addNoResource -->
```
kinflate add resource does_not_exist.yaml
```

Try to reinit; kinflate should complain.
<!-- @initAgain -->
```
kinflate init
```


<!-- @allDone @test -->
```
popd
```
