# ktl (kubectl2)

## Metadata

## Table of Contents

- [Title](#title)
   - [Metadata](#metadata)
   - [Table of Contents](#table-of-contents)
   - [Summary](#summary)
   - [Motivation](#motivation)
   - [Reference-level explanation](#reference-level-explanation)
   - [Graduation Criteria](#graduation-criteria)
   - [Alternatives](#alternatives-optional)
   - [Unresolved Questions](#unresolved-questions-optional)

## Summary

- Stop trying to wholesale move kubectl out of kubernetes/kubernetes
  - we have been working on this for a year with limited success
- Instead, create a new cli command called `ktl`.
  - `ktl` lives in its own repo and vendors commands developed in other repos (kubernetes/kubernetes, kubernetes/kubectl)
- `ktl` is built and released frequently
  - vendored sub commands follow their own release cycles and updated by released the vendored code

## Motivation

The core Kubernetes cli is published as a single binary called `kubectl`.
`kubectl` is a statically linked go binary in the kubernetes/kubernetes
repo that contains many sub commands to perform operations against a Kubernetes cluster, such as:

- creating & updating objects defined by commandline args
- managing objects using configuration files
- debugging objects in a cluster

Benefits of this approach:

- having a single binary facilitates discovery for the suite of Kubernetes cli commands
- static binary makes it simple to distribute the cli across all OS platforms
- static binary makes it simple for users to install the cli
- single binary limits the size of the install (MB) vs distributing ~50 separate go commands (GB)
- centralizing the development facilitates the construction of shared libraries and infrastructure
- leveraging the kubernetes/kubernetes repo infrastructure provides a process to build and release

Challenges of this approach:

- the cli cannot be released at a different intervals from the cluster itself
  - slows down velocity
  - slows down feedback on alpha and beta features - longer time to GA
  - each release is heavier weight due to its length
  - obfuscates version skew as a first class problem
- individual commands cannot be released at different intervals from one another
  - some commands should be released every 3 months with the kubernetes cluster
  - others could be released daily
- shared cli infrastructure cannot easily be publish to be used by repos / commands outside kubernetes/kubernetes
  - requires vendoring kubernetes/kubernetes which is painful
- cli cannnot be owned and maintained independent from unrelated components
  - submit / merge queue blocks on unrelated tests
  - GitHub permissions apply to whole repo - cannot add collaborators or maintainers for just the cli code
  - GitHub notifications for kubernetes/kubernetes are a firehose
  - hard to manage PRs and issues because they are not scoped to CLI

To address these challenges, sig-cli maintainers have worked toward moving kubectl out of the main repo
for the past year.  This process has been slow and challenging due to the following reasons:

- many kubectl commands depend on internal kubernetes/kubernetes libraries
- most commands should not have these dependencies, but removing them requires large rewrites of the commands
- *some* commands do actually need these dependencies (e.g. convert)
- continuing to develop in the kubernetes/kubernetes repo results in more *bad* dependencies being added
  even as we try to remove old ones
- many commands depend on test infrastructure bespoke to kubernetes/kubernetes, which would need to be moved as well

Additionally, many kubectl commands are laden with technical debt, using anti-patterns for working with APIs
that do not work with version skew or support extensibility.  We have since grown out of using these patterns,
but they are still pervasive.  Frequently, is it much faster and effective to rewrite large pieces
instead of trying to refactor them into different designs.

## Goals for ktl

Thing we need.

Keep the advantages of:

- easy distribution
- reasonably sized (~100MB)
- easy installation
- use common cli / client infrastructure
- discoverable commands

And also:

- allow the cli to be released independently from the cluster
- allow sub command groups within the cli to be released independently from one another
- allow a decentralized ecosystem of tools to leverage centralized maintained cli / client infrastructure
- facilitate end-to-end ownership of the cli, and in some cases sub command groups with the cli
- facilitate decentralized development of extensions for the cli

## Anti-goals

Things we want to avoid.

- block on moving existing kubectl commands out of the kubernetes/kubernetes repo
- rewrite kubectl from the ground up in a new repo

## Non-goals

Even if these are good ideas, don't let them distract us from meeting our goals will simpler solutions.

- build solution for discovering installable plugins and installing them
  - rely on existing package management solutions for this until we need something more
- invent new build and distribution infrastructure
- fix issues with the existing kubectl commands
- dog fooding the plugin mechanism for core commands

## Reference-level explanation

- build a new cli binary `ktl` (kubectl2) under the `kubernetes/ktl` repo that dispatches to commands developed in other repos
- keep old cli commands in `kubernetes/kubernetes/cmd/kubectl` and vendor them into `kubectl/ktl`
- build new cli commands in `kubernetes/kubectl` and vendor them into `kurnetes/ktl`
- build common cli Kubernetes client infrastructure and libraries that can be used to develop a decentralized cli ecosystem

### ktl

In a new repo (`kubernetes/ktl`), create the `ktl` binary that dispatches (both statically and dynamically)
to commands developed in other repos.

#### Dispatch

Static dispatch:

- vendor in kubernetes/kubernetes/cmd/kubectl cobra commands, and register them under `ktl kube`
- vendor in kubernetes/kubectl/cmd cobra commands, and register them under `ktl`
- vendor in commands from other sources as needed over time

Dynamic dispatch:

- support git-style plugins that allow registering new commands under `ktl`
  - use simplified version of the kubectl plugin implementation - make configuration files optional
- plugins only purpose is discovery of kubernetes related commands
  - plugins can leverage shared client cli libraries (whether they are installed as plugins or not)
- by default throw an error if plugin names conflict with existing commands
  - this is configurable
- plugins can be disabled through modifying `~/.ktlconfig`

Overriding existing commands:

- support command alias' to allow overriding one command with another
  - allows plugins to extend and override built-in commands
  - allows `ktl kube *` commands to be aliased to `ktl *`

#### Configuration

Use [viper](https://github.com/spf13/viper) to configure dispatch precedence

### Cli commands

All new cli commands should be built outside of kubernetes/kubernetes and vendored into the kubernetes/ktl.  Command repos
implement an interface that returns a list of cobra (sp13) commands, which ktl can register under the root.

Initially new cli commands can be built in the kubernetes/kubectl, but if needed, development maybe moved to other repos.

#### Conventions

- ktl top level commands are one of
  - command groups that work against a specific API group/versions (e.g. isto, svcat, kube)  e.g. `create`
    for a given resource.
  - generic commands that are agnostic to specific APIs, but instead discover API metadata and work against
    all APIs following published conventions.  e.g. a `scale` that works against anything with a scale sub resource


### Library infrastructure

Develop shared set of client / cli libraries that handle registering flags / environment variables / config files
and provide common functionality to cli clients of Kubernetes.

- reading and parsing kubeconfig
- printing libraries
- reading and parsing config into objects from files, directories, stdin
- indexing + querying discovery API and openapi data
- manipulating unstructured objects using openapi schema
- merging objects
- generating patches from objects
- writing objects to the server - optimistic locking failure retries with exp backoff
- providing tty for subresources that accept - e.g. exec, attach
- defining shared exit codes (e.g. exit code when not doing a conditional write)
- test infrastructure for unit and integration tests
  - support version skew testing

### Build and release infrastructure

Develop build triggers to automatically cut and publish builds based on the presence of GitHub tags.  Aggregate
release notes from vendored commands.

- use GCP container builder + mirrored GCP source repo
- publish binary to gs:// bucket
- publish binary + release notes to GitHub releases page

## Graduation Criteria

None

## Alternatives [optional]

### Keeping the kubectl name instead of rebranding

Alternatively we could call the new command `kubectl` and attach all of the legacy kubectl commands
at the root level.  This would make it look and feel exactly like `kubectl`, but allow for new pieces
to be built out of kubernetes/kubernetes.

Steps:

- rename `kubernetes/kubernets/cmd/kubectl` to `kubernetes/kubernets/cmd/legacykubectl`
- rename `kubernetes/kubernets/pkg/kubectl` to `kubernetes/kubernets/pkg/legacykubectl`
- create new `kubernetes/kubectl/cmd/kubectl` command
- vendor kubectl subcommands directly under `kubernetes/kubectl/cmd/kubectl` root command

Tradeoffs:

Keeping kubectl as name: Need to make it continue to look and act like kubectl

- Minimal change from users perspective
- Can easily swap with existing kubectl without updating docs, blog posts, etc

Renaming to ktl: Need to get everyone to use the new command

- Easier to phase out old commands with new ones while changing behavior
- Easier to restructure command layout
- Easier to redefine command wide things - such as version skew support


### Don't vendor in commands, make them plugins instead*

This would be a more complicated approach that can be considered in a later iteration.

## Unresolved Questions [optional]

- Should we rename / restructure the command, or keep kubectl and all the bagage that comes with?