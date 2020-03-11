#!/usr/bin/env bash

# Install outdated helm.
curl -L https://get.helm.sh/helm-v2.16.0-linux-amd64.tar.gz | tar xvz --strip-components 1 linux-amd64/helm
chmod +x helm

kubectl create ns giantswarm

# Install Tiller with outdated retagged image.
helm init --tiller-namespace giantswarm --tiller-image quay.io/giantswarm/tiller:v2.16.0 --wait