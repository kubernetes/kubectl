#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

function exit_with {
  local msg=$1
  echo >&2 ${msg}
  exit 1
}

base_dir="$( cd "$(dirname "$0")/../../.." && pwd )"
cd "$base_dir" || {
  echo "Cannot cd to '$base_dir'. Aborting." >&2
  exit 1
}

# Install kinflate to $GOPATH/bin and export PATH
go install ./cmd/kinflate || { exit_with "Failed to install kinflate"; }
export PATH=$GOPATH/bin:$PATH

home=`pwd`
example_dir="some/default/dir/for/examples"
if [ $# -eq 1 ]; then
    example_dir=$1
fi
if [ ! -d ${example_dir} ]; then
    exit_with "directory ${example_dir} doesn't exist"
fi

test_targets=$(ls ${example_dir})

for t in ${test_targets}; do
    cd ${example_dir}/${t}
    if [ -x "tests/test.sh" ]; then
        tests/test.sh .
        if [ $? -eq 0 ]; then
            echo "testing ${t} passed."
        else
            exit_with "testing ${t} failed."
        fi
    fi
    cd ${home}
done