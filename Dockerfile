FROM golang:1.24.3-alpine AS backendbuilder
RUN apk add --no-cache build-base git bash

ENV WORKDIR=/k8s-ces-control
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY go.mod go.sum ${WORKDIR}/
RUN go mod download

## Copy golang source code
COPY main.go ${WORKDIR}/
COPY interfaces.go ${WORKDIR}/
COPY packages ${WORKDIR}/packages
COPY .git ${WORKDIR}/.git

## Copy makefiles
COPY Makefile ${WORKDIR}/
COPY makefiles ${WORKDIR}/makefiles
COPY build ${WORKDIR}/build

RUN go mod vendor
RUN make compile-generic

FROM alpine:20221110
RUN apk add --no-cache git bash

ENV WORKDIR=/k8s-ces-control
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

LABEL maintainer="hello@cloudogu.com" \
    NAME="k8s-ces-control" \
    VERSION="1.7.3"

COPY --from=backendbuilder /k8s-ces-control/target/k8s-ces-control ${WORKDIR}/k8s-ces-control

ENV USER=k8s-ces-control \
    GRPC_HEALTH_PROBE_VERSION=0.4.36

RUN set -eux -o pipefail \
    && apk update \
    && apk upgrade \
    && addgroup -S -g 1000 "${USER}" \
    && adduser -S -h ${WORKDIR} -G "${USER}" -u 1000 -s /bin/bash "${USER}" \
    && chown -R ${USER}:${USER} ${WORKDIR} /etc/ssl/certs \
    && rm -rf /var/cache/apk/*

# Install grpc health probe cli tool to be used to verify the readiness of our app
RUN wget -O /grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
    && chmod +x /grpc_health_probe

# Create folder for k8s-ces-control files.
RUN mkdir /etc/k8s-ces-control \
    && chown -R ${USER}:${USER} /etc/k8s-ces-control

EXPOSE 50051
#HEALTHCHECK CMD /grpc_health_probe -tls -tulpn | grep LISTEN | grep 50051

WORKDIR ${WORKDIR}
USER k8s-ces-control

CMD ./k8s-ces-control start
