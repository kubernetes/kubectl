# Glossary

[Resource]: #resource
[application]: #application
[base]: #base
[bases]: #base
[bespoke]: #bespoke-configuration
[kinflate]: #kinflate
[manifest]: #manifest
[overlay]: #overlay
[overlays]: #overlay
[off-the-shelf]: #off-the-shelf
[OTS]: #off-the-shelf
[patch]: #patch
[patches]: #patch
[proposal]: https://github.com/kubernetes/community/pull/1629
[resource]: #resource
[resources]: #resource
[target]: #target
[TypeMeta]: https://github.com/kubernetes/kubernetes/blob/master/pkg/api/unversioned/types.go
[apt]: https://en.wikipedia.org/wiki/APT_(Debian)
[rpm]: https://en.wikipedia.org/wiki/Rpm_(software)
[YAML]: http://www.yaml.org/start.html
[JSON]: https://www.json.org/
[workflow]: #workflow


## application

An _application_ is a group of k8s resources that
server some common purpose, e.g.  a webserver backed by
a database.

[Resource] labelling, naming and metadata schemes have
historically served to group resources together for
collective operations like _list_ and _remove_.

This [proposal] describes a new k8s _application_
resource, which describe a more formal grouping and
official type to support application operations and
dashboards.  [kinflate] configures k8s resources, and
the proposed application resource is just another
resource.

There's some conceptual overlap between the proposal’s
application resource and the kinflate [manifest].  The
application resource has a `Components` field that
serves a purpose similar to the `resources` field in
the kinflate manifest.  This overlap can be resolved in
various ways, the simplest being that kinflate does
nothing special with the application resource.


## base

A _base_ is a [target] that some [overlay] modifies.

Any target, including an overlay, can be a base to
another target.

A typical base would be a set of resources common to
some [application], e.g. mysql.  A base might be a cloned
fork of some set of canonical resources describing a
mysql deployment, maintained somewhere on the web.

## bespoke configuration

A _bespoke_ configuration is a manifest and some
resources maintained internally by some organization
for their own purposes.

The [workflow] associated with a _bespoke_ config is
simpler than the workflow associated with an
[off-the-shelf] (OTS) config, because there's no notion
of periodically capturing upgrades from some OTS config
that someone else maintains.


## instance

An _instance_ is the outcome, in a cluster, of applying
an [overlay] to a [base].

> E.g. a _staging_ and _production_ overlay both modify
> some common base to create distinct instances.
>
> The _staging_ instance is exposed to quality assurance
> testing, or to some external users who'd like to see
> what the next version of production will look like.
>
> The _production_ instance is the set of resources
> exposed to all production traffic, and thus may employ
> deployments with a large number of replicas and higher
> cpu and memory requests.

In a best practices end-user target layout, a directory
called _overlays_ (or _instances_) contains overlays in
its sub-directories.

Roughly synonymous with [overlay].


## kinflate

_kinflate_ is a command line tool supporting template-free
customization of declarative configuration targetted to
k8s.

Targetted to k8s means that kinflate reserves the right
to understand (to whatever level needed) k8s API
resources, k8s concepts like names, labels, namespaces,
etc. and the semantics of resource patching.

## manifest

A _manifest_ is a file called `Kube-manifest.yaml` (and
nothing else) that describes a configuration.

A manifest contains fields falling into these categories:

 * (_TBD_) Standard k8s API kind-version fields, e.g. [TypeMeta].
 * Immediate customization instructions - _nameprefix_, _labelprefix_, etc.
 * Resource _generators_ for configmaps and secrets.
 * Cargo - _names of external files_ in these categories:
   * [resources] - completely specified k8s API objects,
      e.g. `deployment.yaml`, `configmap.yaml`, etc.
   * [patches] - _partial_ resources that modify full
     resources defined in a [base] (only meaningful in an [overlay]).
   * [bases] - path to a directory containing a [manifest]
      (only meaningful in an [overlay]).

## off-the-shelf configuration

An _off-the-shelf_ configuration is a manifest and
resources intentionally published somewhere for others
to use.

One can offer some user a simple workflow by placing
the manifest and its associated resources in one git
repository, e.g.

> ```
> github.com/kinflate/ldap/
>  Kube-manifest.yaml
>  deployment.yaml
>  configmap.yaml
>  README.md
>  LICENSE.md
> ```

A consumer could then _fork_ this repo (on github) and
_clone_ their fork to their local disk for
customization.

This clone could act as a [base] for the user's
own [overlays] to do further customization.

## overlay

An _overlay_ is a [target] that modifies (and thus
depends on) another target.

The [manifest] in an overlay refers to (via file path,
URI or other method) to _some other manifest_, known as
its [base].  An overlay is unusable without its base.

An overlay supports the typical notion of a
_development_, _QA_, _staging_ and _production_
environment instances.

The configuration of these environments is specified in
individual overlays (one per environment) that all
refer to a common base that holds common configuration.
One configures the cluser like this:

> ```
>  kinflate inflate -f ldap/overlays/staging | kubectl apply -f -
>  kinflate inflate -f ldap/overlays/production | kubectl apply -f -
> ```

etc.

Usage of the base is implicit (one would need to
examine the manifests to see it).

An overlay may act as a base to another overlay.

## package

The word _package_ has no meaning in kinflate, as
kinflate is not to be confused with a package
management tool in the tradition of, say, [apt] or
[rpm].

## patch

A _patch_ is a partially defined k8s resource with a
name that must match a resource already known per
traversal rules built into [kinflate].

_Patch_ is a field in the manifest, distinct from
resources, because a patch file looks like a resource
file, but has different semantics.  A patch depends on
(modifies) a resource, whereas a resourse has no
dependencies.  Since any resource file can be used as a
patch, one cannot reliably distinguish a resource from
a patch just by looking at the file's [YAML].

## resource

A _resource_ is a path to a [YAML] or [JSON] file that
completely defines a functional k8s API object.

## sub-target / sub-application / sub-package

A _sub-whatever_ is not a thing. There are only [bases] and [overlays].

## target

The _target_ is the argument to `inflate`, e.g.:

> ```
>  kinflate inflate -f $target | kubectl apply -f -
> ```

`$target` must be

 * a file path ending with `Kube-manifest.yaml`,
 * a directory that immediately contains a file with that name,
 * a URI that resolves to some file or directory on the web
   (e.g. a github repo) meeting the aforemention conditions.

The target contains all the information needs for kinflate
to create the customized resources that will be applied to
one's cluster.


## workflow

A _workflow_ is the steps one takes to maintain and use
a configuration.

Common workflows:

#### No local files

> ```
> kinflate inflate -f https://github.com/kinflate/ldap
> ```

You just install some configuration from the web.

#### local, bare manifest

> ```
> kinflate inflate -f ~/Kube-manifest.yaml
> ```

The manifest is a simple overlay of some web-based
customization target, e.g. specifiying a name prefix.
It references no other files in ~, and the user doesn’t
maintain it in a repo.

#### one local overlay

> ```
> kinflate inflate -f ~/myldap
> ```

The myldap dir contains a `Kube-manifest.yaml`
referencing the base via URL, and a `deployment.yaml`
that increase the replica count specified in the base.

#### multiple instances

> ```
> # Make a workspace
> mkdir ldap
>
> # Clone a target you wish to customize (or fork it on github,
> # then clone your fork, and make the original your upstream):
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
> # To instantiate the customizations in a cluster, enter:
> kinflate inflate -f ldap/overlays/staging
> kinflate inflate -f ldap/overlays/production
>
> ```

The [overlays] are siblings to each other and to the
[base] they depend on.  The overlays directory is
maintained in its own repo.

The [base] directory is maintained in another repo whose
upstream is an [OTS] configuration, in this case
https://github.com/kinflate/ldap.  The user can rebase
at will.

### bad practice

> ```
> git clone https://github.com/kinflate/ldap
> mkdir ldap/staging
> mkdir ldap/production # ...
> ```

This nests kinflate targets, confusing one’s ability to
maintain them in distinct git repos, and and increases
the chance of a cycle.
