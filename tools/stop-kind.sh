#! /usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

KIND_CLUSTER_1_NAME="${1}"
KIND_CLUSTER_2_NAME="${2}"

kind delete cluster --name "${KIND_CLUSTER_1_NAME}"
kind delete cluster --name "${KIND_CLUSTER_2_NAME}"
