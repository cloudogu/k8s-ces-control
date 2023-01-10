MONITORING_NAMESPACE=monitoring

.PHONY: install-monitoring-stack
install-monitoring-stack: install-loki-stack install-prometheus install-event-exporter

.PHONY: delete-monitoring-stack
delete-monitoring-stack: delete-loki-stack delete-prometheus delete-event-exporter

.PHONY: expose-grafana
expose-grafana:
	@kubectl --namespace=${MONITORING_NAMESPACE} port-forward svc/grafana 8080:80 &

.PHONY: install-loki-stack
install-loki-stack: install-promtail install-minio install-loki install-grafana

.PHONY: delete-loki-stack
delete-loki-stack: delete-promtail delete-minio delete-loki delete-grafana

.PHONY: monitoring-namespace
monitoring-namespace:
	@kubectl create namespace monitoring || true

.PHONY: install-promtail
install-promtail: monitoring-namespace loki-example-credentials
	@kubectl apply -f ${WORKDIR}/monitoring-stack/promtail --namespace=${MONITORING_NAMESPACE}

.PHONY: loki-example-credentials
loki-example-credentials: monitoring-namespace
	@kubectl create secret generic loki-credentials --from-literal=username=loki --from-literal=password=loki  --namespace=${MONITORING_NAMESPACE} || true

.PHONY: delete-promtail
delete-promtail: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/monitoring-stack/promtail --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-minio
install-minio: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/monitoring-stack/minio --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-minio
delete-minio: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/monitoring-stack/minio --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-loki
install-loki: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/monitoring-stack/loki/ --namespace=${MONITORING_NAMESPACE}
	@kubectl apply -f ${WORKDIR}/monitoring-stack/loki/write --namespace=${MONITORING_NAMESPACE}
	@kubectl apply -f ${WORKDIR}/monitoring-stack/loki/read --namespace=${MONITORING_NAMESPACE}
	@kubectl apply -f ${WORKDIR}/monitoring-stack/loki/gateway --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-loki
delete-loki: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/monitoring-stack/loki/ --namespace=${MONITORING_NAMESPACE} || true
	@kubectl delete -f ${WORKDIR}/monitoring-stack/loki/write --namespace=${MONITORING_NAMESPACE} || true
	@kubectl delete -f ${WORKDIR}/monitoring-stack/loki/read --namespace=${MONITORING_NAMESPACE} || true
	@kubectl delete -f ${WORKDIR}/monitoring-stack/loki/gateway --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-grafana
install-grafana: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/monitoring-stack/grafana --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-grafana
delete-grafana: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/monitoring-stack/grafana --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-prometheus
install-prometheus: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/monitoring-stack/prometheus --recursive --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-prometheus
delete-prometheus: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/monitoring-stack/prometheus --recursive --namespace=${MONITORING_NAMESPACE} || true

.PHONY: install-event-exporter
install-event-exporter: monitoring-namespace
	@kubectl apply -f ${WORKDIR}/monitoring-stack/event --namespace=${MONITORING_NAMESPACE}

.PHONY: delete-event-exporter
delete-event-exporter: monitoring-namespace
	@kubectl delete -f ${WORKDIR}/monitoring-stack/event --namespace=${MONITORING_NAMESPACE} || true
