MYSQL=mysql.yaml
NAMESPACE=namespace.yaml
PERSISTENT_VOLUME=persistent-volume.yaml
PERSISTENT_VOLUME_CLAIM=persistent-volume-claim.yaml
SUMWHERE=sumwhere-application.yaml
REDIS=redis-single.yaml


clean:
	kubectl delete secret regcred -n sumwhere;
	kubectl delete -f $(REDIS)
	kubectl delete -f $(SUMWHERE)
	kubectl delete -f $(MYSQL)
	kubectl delete -f $(PERSISTENT_VOLUME_CLAIM)
	kubectl delete -f $(PERSISTENT_VOLUME)
	kubectl delete -f $(NAMESPACE)
	kubectl delete -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml;
	kubectl delete -f layer2config.yaml
	kubectl delete -f ingress.yaml
	kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
	kubectl delete -f default-http-backend.yaml
	kubectl delete -f nginx-ingress-service.yaml


.PHONY: metallb
metallb:
	kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.7.3/manifests/metallb.yaml
	kubectl apply -f layer2config.yaml

.PHONY: ingress
ingress:
	kubectl apply -f ingress.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
	kubectl apply -f default-http-backend.yaml
	kubectl apply -f nginx-ingress-service.yaml

.PHONY: sumwhere
sumwhere:
	kubectl apply -f $(NAMESPACE)
	kubectl create secret docker-registry regcred --docker-server=https://index.docker.io/v1/ --docker-username=qkrqjadn --docker-password=1q2w3e4r --docker-email=qjadn0914@naver.com -n sumwhere
	kubectl apply -f $(PERSISTENT_VOLUME)
	kubectl apply -f $(PERSISTENT_VOLUME_CLAIM)
	kubectl apply -f $(MYSQL)
	kubectl apply -f $(SUMWHERE)
	kubectl apply -f $(REDIS)

.PHONY: all
all: sumwhere metallb ingress

