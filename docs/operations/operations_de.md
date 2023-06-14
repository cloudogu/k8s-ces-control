# k8s-ces-control betreiben

## Installation


## Abgelaufenes Zertifikat aktualisieren
1. Zertifikat löschen
   - `kubectl exec etcd-client etcdctl rm /config/_global/certificate/cesappd/server.crt`
   - `kubectl exec etcd-client etcdctl rm /config/_global/certificate/k8s-ces-control/server.crt`
   - `kubectl exec etcd-client etcdctl rm /config/_global/certificate/k8s-ces-control/server.key` 
2. k8s-ces-control neu starten
   - `kubectl delete pod k8s-ces-control`