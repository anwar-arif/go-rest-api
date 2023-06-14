## Steps to run on local kubernetes cluster
- Create mongodb user and password
```bash
$ kubectl create secret generic mongodb-user-pass \
    --from-literal=MONGO_INITDB_ROOT_USERNAME=your_username \
    --from-literal=MONGO_INITDB_ROOT_PASSWORD='your_password'
```
- Install ingress controller [from here](https://kubernetes.github.io/ingress-nginx/deploy/)
- Go to project root directory and run `kubectl apply -f k8s/`
- Get the address and port by running `kubectl get ingress`
- Open the address in browser and access the server
## Note 
You might see `503 Service Temporarily Unavailable` error in namespaces other than `default`. \
In this case, you have to run in `default` namespace  \
or move ingress controller secret to your preferred namespace 

# Deploy on GKE
- create a cluster in [GKE](https://console.cloud.google.com/)
- create service account for github deployer
- install [helm](https://helm.sh/docs/intro/install/#from-script) in k8s cluster
- create secret for `mongodb-user-pass`
```bash
kubectl create secret generic mongodb-user-pass \
--from-literal=MONGO_INITDB_ROOT_USERNAME=root \
--from-literal=MONGO_INITDB_ROOT_PASSWORD='secret' 
```
- install `ingress-nginx` using [helm](https://kubernetes.github.io/ingress-nginx/deploy/#using-helm)
```bash
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```