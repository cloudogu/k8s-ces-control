# Operating k8s-ces-control 

## Installation

1. Install the following components
   - [k8s-minio](https://github.com/cloudogu/k8s-minio)
   - [k8s-loki](https://github.com/cloudogu/k8s-loki)
2. Install k8s-ces-control

Additionally, the [Admin dogu documentation](https://github.com/cloudogu/admin/blob/develop/k8s-samples/k8s-admin-dependencies.yaml) offers a file with all dependencies which need to be installed in the cluster with `kubectl apply -f k8s-admin-dependencies.yaml --namespace ecosystem`.