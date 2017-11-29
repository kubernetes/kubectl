#!/usr/bin/env bash
set -eu

# Use DEBUG=1 ./scripts/run-tests.sh to get debug output
[[ -z "${DEBUG:-""}" ]] || set -x

GINKGO="ginkgo"
if [[ -n "${GINKGO_WATCH:-""}" ]] ; then
  GINKGO="$GINKGO watch"
fi

if [[ -z ${GINKGO_PERFORMANCE:-""} ]] ; then
  GINKGO="$GINKGO -skipMeasurements"
fi

test_framework_dir="$(cd "$(dirname "$0")/.." ; pwd)"

export KUBE_ASSETS_DIR="${test_framework_dir}/assets/bin"

$GINKGO -r "${test_framework_dir}"
