#!/usr/bin/env bash
set -eu

# Use DEBUG=1 ./scripts/download-binaries.sh to get debug output
quiet="--quiet"
[[ -z "${DEBUG:-""}" ]] || {
  set -x
  quiet=""
}

# Use BASE_URL=https://my/binaries/url ./scripts/download-binaries to download
# from a different bucket
: "${BASE_URL:="https://storage.googleapis.com/k8s-c10s-test-binaries"}"

test_framework_dir="$(cd "$(dirname "$0")/.." ; pwd)"
os="$(uname -s)"
arch="$(uname -m)"

echo "About to download a couple of binaries. This might take a while..."
wget $quiet "${BASE_URL}/etcd-${os}-${arch}" -O "${test_framework_dir}/assets/bin/etcd"
wget $quiet "${BASE_URL}/kube-apiserver-${os}-${arch}" -O "${test_framework_dir}/assets/bin/kube-apiserver"
chmod +x "${test_framework_dir}/assets/bin/etcd"
chmod +x "${test_framework_dir}/assets/bin/kube-apiserver"
echo "Done!"
