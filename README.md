<p align="center">
  <img src="https://go.dev/blog/gopher/header.jpg" height="80"/>
  &nbsp;&nbsp;&nbsp;
  <img src="https://upload.wikimedia.org/wikipedia/commons/3/39/Kubernetes_logo_without_workmark.svg" height="70"/>
</p>

<h1 align="center">simple-operator — MiniApp</h1>

<p align="center">
  A production-style Kubernetes Operator built with Kubebuilder & Go.<br/>
  Introduces a custom resource <code>MiniApp</code> and reconciles it into a native Kubernetes Deployment.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21-blue?logo=go"/>
  <img src="https://img.shields.io/badge/Kubernetes-Operator-326CE5?logo=kubernetes"/>
  <img src="https://img.shields.io/badge/Kubebuilder-v3-orange"/>
  <img src="https://img.shields.io/github/stars/amirhosssein0/simple-operator?style=social"/>
</p>

---

## What it does

`MiniApp` is a custom Kubernetes resource (CRD) that abstracts away raw Deployment complexity.
You define your app — the operator handles the rest.

- Watches `MiniApp` custom resources via a reconciliation loop
- Creates or updates a matching `Deployment` automatically
- `spec.image` is required — returns a descriptive error if missing
- Sensible defaults: `replicas=1`, `port=8080`
- Uses `OwnerReference` — deleting a `MiniApp` cascades to its `Deployment`
- Idempotent reconciler — safe to re-run at any time

---

## Architecture

```
kubectl apply MiniApp CR
        │
        ▼
┌───────────────────┐
│  MiniApp CRD      │  Custom Resource Definition
└────────┬──────────┘
         │ watch
         ▼
┌───────────────────┐
│  Reconciler Loop  │  controller-runtime
│  (Go controller)  │
└────────┬──────────┘
         │ creates/updates
         ▼
┌───────────────────┐
│  Deployment       │  Native K8s resource
│  (OwnerRef set)   │
└───────────────────┘
```

---

## MiniApp Spec

```yaml
apiVersion: apps.amir.local/v1alpha1
kind: MiniApp
metadata:
  name: my-app
spec:
  image: nginx:latest   # required
  replicas: 2           # optional, default: 1
  port: 8080            # optional, default: 8080
```

---

## Run locally

```bash
git clone https://github.com/amirhosssein0/simple-operator.git
cd simple-operator
make generate && make manifests
make install
make run
```

Then in another terminal:

```bash
kubectl apply -f examples/miniapp.yaml
kubectl get miniapp
kubectl get deploy
kubectl get pods
```

---

## Key concepts demonstrated

| Concept | Implementation |
|---|---|
| Custom Resource Definition | `MiniApp` CRD via Kubebuilder markers |
| Reconciliation loop | `controller-runtime` Reconciler interface |
| OwnerReference | Cascade delete MiniApp → Deployment |
| Error handling | Validation on required `spec.image` |
| Idempotency | Create-or-update pattern |

---

## Why I built this

Most DevOps engineers use Kubernetes — few understand how to extend it.
This operator was built to deeply understand the controller-runtime reconciliation model,
CRD design, and how tools like ArgoCD and cert-manager work under the hood.
