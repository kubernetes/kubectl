# workflows

A _workflow_ is the steps one takes to maintain and use
a configuration.


### No local files

> ```
> kinflate inflate -f https://github.com/kinflate/ldap
> ```

You just install some configuration from the web.

### local, bare manifest

> ```
> kinflate inflate -f ~/Kube-manifest.yaml
> ```

The manifest is a simple overlay of some web-based
customization target, e.g. specifiying a name prefix.
It references no other files in ~, and the user doesn’t
maintain it in a repo.

### one local overlay

> ```
> kinflate inflate -f ~/myldap
> ```

The myldap dir contains a `Kube-manifest.yaml`
referencing the base via URL, and a `deployment.yaml`
that increase the replica count specified in the base.

### multiple instances

> ```
> # Make a workspace
> mkdir ldap
>
> # Clone your fork of some target you wish to customize:
> git clone https://github.com/kinflate/ldap ldap/base
>
> # Create a directory to hold overlays.
> mkdir ldap/overlays
>
> # Create an overlay, in this case called “staging”
> mkdir ldap/overlays/staging
>
> # To “staging” add a Kube-manifest.yaml file,
> # and optionally some resources and/or patches,
> # e.g. a configmap that turns on an experiment flag.
>
> # Create another overlay.
> mkdir ldap/overlays/production
> # And add customization to this directory as was done in staging,
> # e.g. a patch that increases a replica count.
>
> # Apply the instances to a cluster:
> kinflate inflate -f ldap/overlays/staging | kubectl apply -f -
> kinflate inflate -f ldap/overlays/production | kubectl apply -f -
>
> ```

[overlays]: glossary.md#overlay
[base]: glossary.md#base
[off-the-shelf]: glossary.md#off-the-shelf
[rebase]: https://git-scm.com/docs/git-rebase

The [overlays] are siblings to each other and to the
[base] they depend on.  The overlays directory is
maintained in its own repo.

The [base] directory is maintained in another repo whose
upstream is an [off-the-shelf] configuration, in this case
https://github.com/kinflate/ldap.  The user can [rebase]
this [base] at will to capture upgrades.

### bad practice

> ```
> git clone https://github.com/kinflate/ldap
> mkdir ldap/staging
> mkdir ldap/production # ...
> ```

This nests kinflate targets, confusing one’s ability to
maintain them in distinct git repos, and and increases
the chance of a cycle.
