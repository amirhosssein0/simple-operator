<p align="center">
  <img src="https://go.dev/images/go-logo-blue.svg" height="60"/>
  &nbsp;&nbsp;&nbsp;
  <img src="https://upload.wikimedia.org/wikipedia/commons/3/39/Kubernetes_logo_without_workmark.svg" height="60"/>
</p>

# A tiny Kubebuilder operator that introduces a new CRD `MiniApp` and reconciles it into a Kubernetes `Deployment`.

What it does:
- Watches MiniApp resources and keeps a matching Deployment in sync
- Deployment name = MiniApp name (same namespace)
- Required: spec.image (if missing, reconcile returns an error)
- Defaults: replicas=1, port=8080
- Sets OwnerReference so deleting MiniApp deletes the Deployment

Prereqs:
- Ubuntu + Kubernetes (k3s/minikube/kind)
- Go, kubectl, kubebuilder
(Optional network fix)
go env -w GOPROXY=https://goproxy.io,direct ; go env -w GOSUMDB=off

Install & run (local controller):
git clone https://github.com/amirhosssein0/simple-operator.git
cd simple-operator
make generate && make manifests
make install
make run

Example:
cat > miniapp.yaml <<'YAML'
apiVersion: apps.amir.local/v1alpha1
kind: MiniApp
metadata:
  name: hello
spec:
  image: nginx:latest
  replicas: 2
  port: 80
YAML
kubectl apply -f miniapp.yaml
kubectl get deploy
kubectl get pods -l app=hello

Cleanup:
kubectl delete -f miniapp.yaml
make uninstall
