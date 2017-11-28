#!/usr/bin/env bash
set -eu

# Use DEBUG=1 ./scripts/run-tests.sh to get debug output
[[ -z "${DEBUG:-""}" ]] || set -x

declare -a ginkgo_args

if [[ -n "${GINKGO_WATCH:-""}" ]] ; then
  ginkgo_args=( "${ginkgo_args[@]}" "watch" )
fi

if [[ -z ${GINKGO_PERFORMANCE:-""} ]] ; then
  ginkgo_args=( "${ginkgo_args[@]}" "-skipMeasurements" )
fi

test_framework_dir="$(cd "$(dirname "$0")/.." ; pwd)"

export KUBE_ASSETS_DIR="${test_framework_dir}/assets/bin"

ginkgo "${ginkgo_args[@]}" -r "${test_framework_dir}"
