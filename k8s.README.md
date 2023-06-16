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
- set region `gcloud config set compute/region your_compute_region`
- set compute zone `gcloud config set compute/zone your_compute_zone`
- set cluster `gcloud config set container/cluster your_cluster_name`
- get credentials `gcloud container clusters get-credentials your_cluster_name`
- create secret for `mongodb-user-pass`
```bash
kubectl create secret generic mongodb-user-pass \
--from-literal=MONGO_INITDB_ROOT_USERNAME=root \
--from-literal=MONGO_INITDB_ROOT_PASSWORD='secret' 
```
- create service account in `IAM` section with role `Kubernetes Engine Admin` to deploy from github action 
- create a key for the service account and download the key
- encode the key with `base64` and add the encoded key in GitHub keys for corresponding project
- add these secrets in `github action secrets` section
- also add `cluster name`, `project id` `compute zone` in github secrets
- names of these secrets must match the secret names in github workflows.yml file
- install [helm](https://helm.sh/docs/intro/install/#from-script) in k8s cluster
```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh 
```
- install `ingress-nginx` using [helm](https://kubernetes.github.io/ingress-nginx/deploy/#using-helm)
```bash
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```
- commit any changes and result should trigger the deployment job
- install certificate manager using helm in your cluster
```bash
helm repo add jetstack https://charts.jetstack.io
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.12.0 \
  --set installCRDs=true
```
- now you should be able to see the working server