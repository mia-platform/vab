#! /usr/bin/env bash

KIND_IMAGE_VERSION="${1}"
KIND_CLUSTER_1_NAME="${2}"
KIND_CLUSTER_2_NAME="${3}"

if [ "$(kind get clusters | grep -c "${KIND_CLUSTER_1_NAME}")" == 1 ]; then
	echo "Kind test cluster 1 already exists!"
else
	kind create cluster \
    --config internal/e2e/kind/kind-config.yaml \
    --name "${KIND_CLUSTER_1_NAME}" \
    --kubeconfig ~/.kube/config \
    --image "${KIND_IMAGE_VERSION}"
fi

if [ "$(kind get clusters | grep -c "${KIND_CLUSTER_2_NAME}")" == 1 ]; then
	echo "Kind test cluster 2 already exists!"
else
	kind create cluster \
    --config internal/e2e/kind/kind-config.yaml \
    --name "${KIND_CLUSTER_2_NAME}" \
    --kubeconfig ~/.kube/config \
    --image "${KIND_IMAGE_VERSION}"
fi
