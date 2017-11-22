#!/usr/bin/env bash
set -eu

# Use DEBUG=1 ./scripts/run-tests.sh to get debug output
[[ -z "${DEBUG:-""}" ]] || set -x

ginkgo_args=''
[[ -z "${GINKGO_WATCH:-""}" ]] || ginkgo_args="${ginkgo_args} watch"

test_framework_dir="$(cd "$(dirname "$0")/.." ; pwd)"

export KUBE_ASSETS_DIR="${test_framework_dir}/assets/bin"

ginkgo $ginkgo_args -r "${test_framework_dir}"
