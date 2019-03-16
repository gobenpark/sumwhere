#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

#if [ $# != 1 ]; then
#    echo argument not 1
#    exit 1
#fi

function sumwhere::certmanager() {
    kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.6/deploy/manifests/00-crds.yaml
    kubectl create namespace cert-manager
    kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true
    helm install --name cert-manager --namespace cert-manager stable/cert-manager
}

function sumwhere::database() {
    kubectl apply -f mysql.yaml
}

function sumwhere::nfs() {
    helm install -n nfs-claim --namespace sumwhere stable/nfs-client-provisioner --set nfs.server=192.168.0.10 --set nfs.path=/HDD2/claim
}

function sumwhere::metallb() {
    kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml
    kubectl apply -f layer2config.yaml
}

function sumwhere::firebase() {
    kubectl create configmap firebase-json --from-file ./galmal-8f900-firebase-adminsdk-zhjsl-f6d034ad3b.json -n sumwhere
}

function sumwhere::volume() {
    kubectl apply -f persistent-volume.yaml
    kubectl apply -f persistent-volume-claim.yaml
}

function sumwhere::weavenet() {
    kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
}

function sumwhere::ingress() {
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
    kubectl apply -f default-http-backend.yaml
    kubectl apply -f nginx-ingress-service.yaml
}

function sumwhere::helm() {
    kubectl create serviceaccount --namespace kube-system tiller
    kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
    helm init --service-account tiller
}

function sumwhere() {
    kubectl apply -f sumwhere-application.yaml
}


kubeadm init
kubectl taint nodes --all node-role.kubernetes.io/master- # 마스터도 스케쥴링 되도록
sumwhere::weavenet
sumwhere::helm
kubectl create namespace sumwhere
kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=qkrqjadn --docker-password=1q2w3e4r --docker-email=qjadn0914@naver.com -n sumwhere
sumwhere::firebase
sumwhere::metallb
sumwhere::ingress
sumwhere::volume
sumwhere::nfs
sumwhere::database
sumwhere::certmanager
sumwhere

#sysctl -q -w net.ipv6.conf.all.disable_ipv6=1 # curl -i 안되는 문제




