# Instances

_Instances_ of a cluster app represent a common configuration problem.

Their configuration is mostly the same.

We only want to manage _differences_.

The DAM approach is to create _overlays_.

An overlay is just a sub-directory with another manifest file,
and optionally more (or fewer) resources.

## Example

Create a _staging_ and _production_ instance.

The differences:

 * In staging we'll enable a risky feature.
 * The greetings will differ.
 * Production will have a higher replica count because
   it takes public traffic.

<!-- @makePatchDirectories @test -->
```
mkdir -p $TUT_APP/staging
mkdir -p $TUT_APP/production
```
