#!/usr/bin/env bash
MYSQL=mysql.yaml
NAMESPACE=namespace.yaml
PERSISTENT_VOLUME=persistent-volume.yaml
PERSISTENT_VOLUME_CLAIM=persistent-volume-claim.yaml
SUMWHERE=sumwhere-application.yaml
REDIS=redis-single.yaml
RABBITMQ=rabbitmq.yaml
RABBITMQ_RBAC=rabbitmq_rbac.yaml


if [ $# != 1 ]; then
    echo argument not 1
    exit 1
fi

if [ $1 == start ]; then
    kubectl taint nodes --all node-role.kubernetes.io/master-
    kubectl apply -f ${NAMESPACE}
    kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=qkrqjadn --docker-password=1q2w3e4r --docker-email=qjadn0914@naver.com -n sumwhere
    kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml
    kubectl apply -f layer2config.yaml
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
    kubectl apply -f default-http-backend.yaml
    kubectl apply -f nginx-ingress-service.yaml
    kubectl apply -f ${PERSISTENT_VOLUME}
    kubectl apply -f ${PERSISTENT_VOLUME_CLAIM}
    kubectl apply -f ${SUMWHERE}
    kubectl apply -f ${REDIS}
    kubectl apply -f ${RABBITMQ_RBAC}
    kubectl apply -f ${RABBITMQ}
    kubectl apply -f ${MYSQL}
    kubectl create configmap firebase-json --from-file ./galmal-8f900-firebase-adminsdk-zhjsl-f6d034ad3b.json -n sumwhere

elif [ $1 == stop ]; then
    kubectl delete -f nginx-ingress-service.yaml
    kubectl delete -f default-http-backend.yaml
    kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
    kubectl delete -f ingress.yaml
    kubectl delete -f layer2config.yaml
    kubectl delete -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml
    kubectl delete -f ${RABBITMQ}
    kubectl delete -f ${RABBITMQ_RBAC}
    kubectl delete -f ${REDIS}
    kubectl delete -f ${SUMWHERE}
    kubectl delete -f ${MYSQL}
    kubectl delete -f ${PERSISTENT_VOLUME_CLAIM}
    kubectl delete -f ${PERSISTENT_VOLUME}
    kubectl delete secret regcred -n sumwhere
    kubectl delete -f ${NAMESPACE}
fi




