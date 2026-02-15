<p align="center">
  <img src="https://go.dev/images/go-logo-blue.svg" height="60"/>
  &nbsp;&nbsp;&nbsp;
  <img src="https://upload.wikimedia.org/wikipedia/commons/3/39/Kubernetes_logo_without_workmark.svg" height="60"/>
</p>

# simple-operator — MiniApp

A minimal Kubebuilder + Go Kubernetes Operator that introduces a custom resource MiniApp and reconciles it into a native Kubernetes Deployment.

## What it does
- Watches MiniApp
- Creates/updates a matching Deployment
- spec.image is required (returns error if missing)
- Defaults: replicas=1, port=8080
- Uses OwnerReference (deleting MiniApp deletes Deployment)

## Run locally

```bash
git clone https://github.com/amirhosssein0/simple-operator.git
cd simple-operator
make generate && make manifests
make install
make run
```

## Example

```bash
kubectl apply -f examples/miniapp.yaml
kubectl get deploy
kubectl get pods
```
