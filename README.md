## How to make a call to Minikube Kubernetes API

run Minikube with command

```bash

minikube start
```

1) first, you need to find IP of your local Minikube instance via command

```bash
minikube ip
```

2) secondly, give a system user permissions to read the list of roles

```bash
kubectl create clusterrolebinding permissive-binding --clusterrole=cluster-admin --user=admin --user=kubelet --group=system:serviceaccounts
```


## How to build docker image in local

go on project root

```bash
eval $(minikube docker-env)
docker build -t demo-kubernate-api:dev .
```

## How to deploy app to the Minikube

navigate to the folder resource-manifests and call two commands:

```bash
kubectl apply -f deployment-enumerator.yaml
kubectl apply -f service-enumerator.yaml
```

and then you can open the service in your browser like so:

```bash
minikube service front-enumerator-api
```

but, probably, the page will be empty. Please, modify the URL and call like so:

```bash
curl http://192.168.64.2:32262/ping
``` 

## How to update version

```bash
kubectl rollout status deployment enumerator
```

## Example request

In order to send a request to get filtered list of roles, please do the:

```bash
curl -X POST http://192.168.64.2:32262/enumerate --data '{"filter_by": "cloud"}'
```

