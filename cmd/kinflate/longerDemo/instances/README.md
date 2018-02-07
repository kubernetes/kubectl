# Instances

_Instances_ of a cluster app represent a common configuration problem.

Their configuration is mostly the same.

We'd like to focus on managing _differences_.

The DAM approach is to create _overlays_.

An overlay is just a sub-directory with another manifest file,
and optionally more (or fewer) resources.

## Example

Create a _staging_ and _production_ instance.

 * The greetings from the _hello world_ web servers will differ.
 * _Staging_ enables a risky feature (for testing).
 * _Production_ has a higher replica count.

<!-- @makePatchDirectories @test -->
```
mkdir -p $TUT_APP/staging
mkdir -p $TUT_APP/production
```
