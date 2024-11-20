# Operating k8s-ces-control 

## Installation

1. Install the following components
   - [k8s-loki](https://github.com/cloudogu/k8s-loki)
   - [k8s-minio](https://github.com/cloudogu/k8s-minio)
2. Install k8s-ces-control

This can be done f. i. like this:

```bash
echo '
apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-minio
  labels:
    app: ces
spec:
  name: k8s-minio
  namespace: k8s
  version: 2023.9.23-7

---

apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-loki
  labels:
    app: ces
spec:
  name: k8s-loki
  namespace: k8s

---

apiVersion: k8s.cloudogu.com/v1
kind: Component
metadata:
  name: k8s-ces-control
  labels:
    app: ces
spec:
  name: k8s-ces-control
  namespace: k8s
' > install-k8s-ces-control-and-deps.yaml

kubectl apply -f install-k8s-ces-control-and-deps.yaml --namespace ecosystem
```
