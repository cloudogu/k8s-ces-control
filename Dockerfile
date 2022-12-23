FROM golang:1.19.3-alpine3.17 AS backendBuilder
RUN apk add --no-cache build-base git bash

ENV WORKDIR=/k8s-ces-control
RUN mkdir -p ${WORKDIR}
WORKDIR ${WORKDIR}

COPY go.mod go.sum ${WORKDIR}/
RUN go mod download

## Copy golang source code
COPY main.go ${WORKDIR}/
COPY packages ${WORKDIR}/packages
COPY generated ${WORKDIR}/generated
COPY .git ${WORKDIR}/.git

## Copy makefiles
COPY Makefile ${WORKDIR}/
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
    VERSION="0.0.0"

COPY --from=backendBuilder /k8s-ces-control/target/k8s-ces-control ${WORKDIR}/k8s-ces-control

ENV USER=k8s-ces-control
RUN set -eux -o pipefail \
    && apk update \
    && apk upgrade \
    && addgroup -S -g 1000 "${USER}" \
    && adduser -S -h ${WORKDIR} -G "${USER}" -u 1000 -s /bin/bash "${USER}" \
    && chown -R ${USER}:${USER} ${WORKDIR} /etc/ssl/certs \
    && rm -rf /var/cache/apk/*

EXPOSE 50051
HEALTHCHECK CMD netstat -tulpn | grep LISTEN | grep 50051

WORKDIR ${WORKDIR}
USER ${USER}

CMD LOG_LEVEL=DEBUG ./k8s-ces-control start