# Kinflate Demo

Goal:

 1. Clone a simple off-the-shelf example as a base configuration.
 1. Customize it.
 1. Create two different instances based on the customization.

First install the tool, and define a place to work locally:

<!-- @install @test -->
```
go get k8s.io/kubectl/cmd/kinflate
```
<!-- @makeWorkplace @test -->
```
DEMO_HOME=$(mktemp -d)
```

Alternatively, use

> ```
> DEMO_HOME=~/hello
> ```



__Next:__ [Clone an Example](clone)
