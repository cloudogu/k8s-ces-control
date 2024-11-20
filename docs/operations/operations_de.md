# k8s-ces-control betreiben

## Installation

1. Installieren Sie die folgenden Komponenten
   - [k8s-loki](https://github.com/cloudogu/k8s-loki)
   - [k8s-minio](https://github.com/cloudogu/k8s-minio)
2. Installieren Sie k8s-ces-control


Dies kann z. B. so erreicht werden:

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
