/*

Package test implements a integration testing framework for kubernetes.

It provides a kubernetes API you can connect to and test your
kubernetes client implementations against. The lifecycle of the components
needed to provide this API is managed by this framework.

Components

Currently the framework provides the following components:

ControlPlane: The ControlPlane is wrapping Etcd & APIServer and provides some
convenience as it creates instances of the components and wires them together
correctly. A ControlPlane can be new'd up, stopped & started and provide the
URL to connect to the API.
The ControlPlane is supposed to be the default entry point for you.

Etcd: Manages an Etcd binary, which can be started, stopped and connected to.
Etcd will listen on a random port for http connections and will create a
temporary directory for it'd data; unless configured differently.

APIServer: Manages an Kube-APIServer binary, which can be started, stopped and
connected to. APIServer will listen on a random port for http connections and
will create a temporary directory to store the (auto-generated) certificates;
unless configured differently.

Binaries

Both Etcd and APIServer use the same mechanism to determine which binaries to
use when they get started.

1. If the component is configured with a `Path` the framework tries to run that
binary. (Note: To overwrite the `Path` of a component you cannot use the
ControlPlane, but you need to create, wire and start & stop all the components
yourself.)

2. If a environment variable named
`TEST_ASSET_KUBE_APISERVER` or `TEST_ASSET_ETCD` is set, this value is used as a
path to the binary for the APIServer or Etcd.

3. If neither an environment variable is set nor the `Path` is overwritten
explicitly, the framework tries to use the binaries apiserver or etcd in the
directory ${FRAMEWORK_DIRECTORY}/assets/bin/ .

For convenience this framework ships with
${FRAMEWORK_DIR}/scripts/download-binaries.sh which can be used to download
pre-compiled versions of the needed binaries and place them in the default
location.

*/
package test
