# helm acr

[![CircleCI](https://circleci.com/gh/AliyunContainerService/helm-acr.svg?style=svg)](https://circleci.com/gh/AliyunContainerService/helm-acr)
[![Go Report Card](https://goreportcard.com/badge/github.com/AliyunContainerService/helm-acr)](https://goreportcard.com/report/github.com/AliyunContainerService/helm-acr)

Helm plugin to push chart package to [ChartMuseum](https://github.com/helm/chartmuseum).

This project is forked from [chartmuseum/helm-push](https://github.com/chartmuseum/helm-push). 

Some modifications has been made to meet the security requirements on Alibaba Cloud:
* the plugin is able to talk to auth server to gain a Bearer Token.
* the plugin is able to use the Bearer Token to download/upload charts to Chartmuseum.
* the plugin registers `acr`(short for Alibaba Cloud Container Registry) as protocol name in `plugin.yaml`.

### Installation

```bash
# make sure you have git installed
yum install -y git

# install plugin
helm plugin install https://github.com/AliyunContainerService/helm-acr
```

### Usage

Before you use Alibaba Cloud Container Registry's hosted Helm charts service, you should:
* purchase an ACR Enterprise Edition instance and activate its Helm charts service
* have a Kubernetes cluster and have `helm init` done
* make sure you have Internet access to GitHub to download plugin
* create a Helm chart namespace in your ACR Enterprise Edition

```bash
# add namespace/repo to your local repository
# please change username/password/namespace/repo/url below
export HELM_REPO_USERNAME=username; export HELM_REPO_PASSWORD=password;
helm repo add demo acr://hello-acr-helm.cn-hangzhou.cr.aliyuncs.com/foo/bar --username ${HELM_REPO_USERNAME} --password ${HELM_REPO_PASSWORD}

# create an empty chart locally
helm create hello-acr

# push the chart
helm push hello-acr demo

# delete local chart
rm -r hello-acr

# update charts index from remote
helm repo update

# show all remote charts
helm search

# fetch the chart we uploaded
helm fetch demo/hello-acr

# delete local repository
helm repo remove demo
```