#! /usr/bin/env bash

KIND_CLUSTER_1_NAME="${1}"
KIND_CLUSTER_2_NAME="${2}"

if [ "$(kind get clusters | grep -c "${KIND_CLUSTER_1_NAME}")" == 1 ]; then
	echo "Kind test cluster 1 already exists!"
else
	kind create cluster --config internal/e2e/kind/kind-config.yaml --name "${KIND_CLUSTER_1_NAME}" --kubeconfig ~/.kube/config
fi

if [ "$(kind get clusters | grep -c "${KIND_CLUSTER_2_NAME}")" == 1 ]; then
	echo "Kind test cluster 2 already exists!"
else
	kind create cluster --config internal/e2e/kind/kind-config.yaml --name "${KIND_CLUSTER_2_NAME}" --kubeconfig ~/.kube/config
fi
