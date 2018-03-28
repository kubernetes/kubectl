# Demo: hello world with instances

Steps:

 1. Clone an off-the-shelf configuration as your base.
 1. Customize it.
 1. Create two different instances (_staging_ and _production_)
    from your customized base.

First define a place to work:

<!-- @makeWorkplace @test -->
```
DEMO_HOME=$(mktemp -d)
```

Alternatively, use

> ```
> DEMO_HOME=~/hello
> ```

__Next:__ [Clone an Example](clone.md)
