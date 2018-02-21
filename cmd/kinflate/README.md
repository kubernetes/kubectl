# Overview

`kinflate` is a command line interface to manage Kubernetes Application configuration. 

## Installation

`kinflate` is written in Go and can be installed using `go get`:

<!-- @installKinflate @test -->
```shell
go get k8s.io/kubectl/cmd/kinflate
```

After running the above command, `kinflate` should get installed in your `GOPATH/bin` directory.

## Usage

This section explains sub-commands of `kinflate` using examples. 

### init 
`kinflate init` initializes a directory by creating an application manifest file named `Kube-manifest.yaml`. This command should be run in the same directory where application config resides.

```shell
kinflate init
```

### set nameprefix

`kinflate set nameprefix` updates the nameprefix field in the [manifest].
`kinflate inflate` uses the prefix to generate prefixed names for Kubernetes resources.

Examples
```shell
# set nameprefix to 'staging'
kinflate set nameprefix staging

# set nameprefix to 'john`
kinflate set nameprefix john
```

### add 
`kinflate add` adds a resource, configmap or secret to the [manifest]. 

#### examples
```shell
# adds a resource to the manifest
kinflate add resource <path/to/resource-file-or-dir>

# adds a configmap from a literal or file-source to the manifest
kinflate add configmap my-configmap --from-file=my-key=file/path --from-literal=my-literal=12345

# Adds a generic secret to the Manifest (with a specified key)
kinflate add secret generic my-secret --from-file=my-key=file/path --from-literal=my-literal=12345

# Adds a TLS secret to the Manifest (with a specified key)
kinflate add secret tls my-tls-secret --cert=cert/path.cert --key=key/path.key        
```

### inflate
`kinflate inflate` generates application configuration by applying
transformations (e.g. name-prefix, label, annotations) specified in the [manifest].

#### examples
```shell
# generate application configuration and output to stdout
kinflate inflate -f <path/to/application/dir>

# generate application configuration and apply it to a Kubernetes cluster
kinflate inflate -f <path/to/application/dir> | kubectl apply -f -
```


## Demos and tutorial

 * [Getting Started](getting_started.md)
 * [Short demo](shortDemo.md)
 * [Longer demo / tutorial](longerDemo/README.md)
