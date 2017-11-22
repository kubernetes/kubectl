#!/usr/bin/env bash
set -eu

# Use DEBUG=1 ./set-pipeline.sh to get debug output
[[ -z "${DEBUG:-""}" ]] || set -x

# Use CONCOURSE_TARGET=my-concourse ./set-pipeline.sh to connect to your local concourse
: "${CONCOURSE_TARGET:="wings"}"
# Use PIPELINE_NAME=my-name ./set-pipeline.sh to give your pipeline a different name
: "${PIPELINE_NAME:="kubectl"}"

# Use PAIR1_LASTPASS=my-lastpass-key ./set-pipeline.sh to get your github keys and URL from your lastpass entry
: "${PAIR1_LASTPASS:="oss-k8s-github-gds-keypair"}"
: "${PAIR2_LASTPASS:="oss-k8s-github-hhorl-keypair"}"

github_pair1_key="$(lpass show "${PAIR1_LASTPASS}" --field "Private Key")"
github_pair2_key="$(lpass show "${PAIR2_LASTPASS}" --field "Private Key")"
github_pair1_url="$(lpass show "${PAIR1_LASTPASS}" --notes)"
github_pair2_url="$(lpass show "${PAIR2_LASTPASS}" --notes)"

script_dir="$(cd "$(dirname "$0")" ; pwd)"

# Create/Update the pipline
fly set-pipeline \
  --target="${CONCOURSE_TARGET}" \
  --pipeline="${PIPELINE_NAME}" \
  --config="${script_dir}/pipeline.yml" \
  --var=git-dev-url="${github_pair1_url}" \
  --var=git-pair1-url="${github_pair1_url}" \
  --var=git-pair2-url="${github_pair2_url}" \
  --var=git-dev-private-key="${github_pair1_key}" \
  --var=git-pair1-private-key="${github_pair1_key}" \
  --var=git-pair2-private-key="${github_pair2_key}"

# Make the pipeline publicly available
fly expose-pipeline \
  --target="${CONCOURSE_TARGET}" \
  --pipeline="${PIPELINE_NAME}"
