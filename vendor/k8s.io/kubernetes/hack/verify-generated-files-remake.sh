#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${KUBE_ROOT}/hack/lib/init.sh"

# $1 = filename pattern as in "zz_generated.$1.go"
function find_genfiles() {
    find .                         \
        \(                         \
          -not \(                  \
            \(                     \
                -path ./_\* -o     \
                -path ./.\*        \
            \) -prune              \
          \)                       \
        \) -name "zz_generated.$1.go"
}

# $1 = filename pattern as in "zz_generated.$1.go"
# $2 timestamp file
function newer() {
    find_genfiles "$1" | while read F; do
        if [[ "${F}" -nt "$2" ]]; then
            echo "${F}"
        fi
    done
}

# $1 = filename pattern as in "zz_generated.$1.go"
# $2 timestamp file
function older() {
    find_genfiles "$1" | while read F; do
        if [[ "$2" -nt "${F}" ]]; then
            echo "${F}"
        fi
    done
}

function assert_clean() {
    make generated_files >/dev/null
    touch "${STAMP}"
    make generated_files >/dev/null
    X=($(newer deepcopy "${STAMP}"))
    if [[ "${#X[*]}" != 0 ]]; then
        echo "Generated files changed on back-to-back 'make' runs:"
        echo "  ${X[@]:-(none)}"
        return 1
    fi
    true
}

STAMP=/tmp/stamp.$RANDOM

#
# Test when we touch a file in a package that needs codegen.
#

assert_clean

DIR=staging/src/k8s.io/sample-apiserver/pkg/apis/wardle/v1alpha1
touch "$DIR/types.go"
touch "${STAMP}"
make generated_files >/dev/null
X=($(newer deepcopy "${STAMP}"))
if [[ "${#X[*]}" != 1 || ! ( "${X[0]}" =~ "${DIR}/zz_generated.deepcopy.go" ) ]]; then
    echo "Wrong generated deepcopy files changed after touching src file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi
X=($(newer defaults "${STAMP}"))
if [[ "${#X[*]}" != 1 || ! ( "${X[0]}" =~ "${DIR}/zz_generated.defaults.go" ) ]]; then
    echo "Wrong generated defaults files changed after touching src file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi
X=($(newer conversion "${STAMP}"))
if [[ "${#X[*]}" != 1 || ! ( "${X[0]}" =~ "${DIR}/zz_generated.conversion.go" ) ]]; then
    echo "Wrong generated conversion files changed after touching src file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

#
# Test when the codegen tool itself changes: deepcopy
#

assert_clean

touch staging/src/k8s.io/code-generator/cmd/deepcopy-gen/main.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older deepcopy "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated deepcopy files did not change after touching code-generator file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch staging/src/k8s.io/code-generator/cmd/deepcopy-gen/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older deepcopy "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated deepcopy files did not change after touching code-generator dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/gengo/examples/deepcopy-gen/generators/deepcopy.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older deepcopy "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated deepcopy files did not change after touching code-generator dep file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/gengo/examples/deepcopy-gen/generators/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older deepcopy "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated deepcopy files did not change after touching code-generator dep dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

#
# Test when the codegen tool itself changes: defaults
#

assert_clean

touch staging/src/k8s.io/code-generator/cmd/defaulter-gen/main.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older defaults "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated defaults files did not change after touching code-generator file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch staging/src/k8s.io/code-generator/cmd/defaulter-gen/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older defaults "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated defaults files did not change after touching code-generator dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/gengo/examples/defaulter-gen/generators/defaulter.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older defaults "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated defaults files did not change after touching code-generator dep file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/gengo/examples/defaulter-gen/generators/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older defaults "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated defaults files did not change after touching code-generator dep dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

#
# Test when the codegen tool itself changes: conversion
#

assert_clean

touch staging/src/k8s.io/code-generator/cmd/conversion-gen/main.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older conversion "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated conversion files did not change after touching code-generator file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch staging/src/k8s.io/code-generator/cmd/conversion-gen/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older conversion "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated conversion files did not change after touching code-generator dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/code-generator/cmd/conversion-gen/generators/conversion.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older conversion "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated conversion files did not change after touching code-generator dep file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/code-generator/cmd/conversion-gen/generators/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older conversion "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated conversion files did not change after touching code-generator dep dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

#
# Test when we touch a file in a package that needs codegen.
#

assert_clean

touch "staging/src/k8s.io/api/core/v1/types.go"
touch "${STAMP}"
make generated_files >/dev/null
X=($(newer openapi "${STAMP}"))
if [[ "${#X[*]}" != 1 || ! ( "${X[0]}" =~ "pkg/generated/openapi/zz_generated.openapi.go" ) ]]; then
    echo "Wrong generated openapi files changed after touching src file:"
    echo "${X[@]:-(none)}"
    exit 1
fi

#
# Test when the codegen tool itself changes: openapi
#

assert_clean

touch vendor/k8s.io/kube-openapi/cmd/openapi-gen/openapi-gen.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older openapi "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated openapi files did not change after touching code-generator file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/kube-openapi/cmd/openapi-gen/
touch "${STAMP}"
make generated_files >/dev/null
X=($(older openapi "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated openapi files did not change after touching code-generator dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/kube-openapi/pkg/generators/openapi.go
touch "${STAMP}"
make generated_files >/dev/null
X=($(older openapi "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated openapi files did not change after touching code-generator dep file:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi

assert_clean

touch vendor/k8s.io/kube-openapi/pkg/generators
touch "${STAMP}"
make generated_files >/dev/null
X=($(older openapi "${STAMP}"))
if [[ "${#X[*]}" != 0 ]]; then
    echo "Generated openapi files did not change after touching code-generator dep dir:"
    echo "  ${X[@]:-(none)}"
    exit 1
fi
