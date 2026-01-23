FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags '-s -w' \
    -o codexec ./cmd/codexec

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git make g++ pkg-config \
    autoconf bison flex libtool \
    protobuf-compiler libprotobuf-dev \
    libnl-route-3-dev \
    time



ARG ADD_NODE=false
RUN if [ "$ADD_NODE" = "true" ]; then \
    curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y --no-install-recommends nodejs; \
    fi

ARG ADD_PYTHON=false
RUN if [ "$ADD_PYTHON" = "true" ]; then \
    apt-get install -y --no-install-recommends python3 python3-minimal; \
    fi

RUN rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/codexec /app/codexec
RUN chmod +x /app/codexec

RUN git clone --depth 1 --branch 3.4 https://github.com/google/nsjail.git /opt/nsjail \
    && make -C /opt/nsjail \
    && ln -s /opt/nsjail/nsjail /usr/local/bin/nsjail

RUN useradd -m -u 1001 -s /bin/bash runner \
    && mkdir -p /jobs \
    && chown -R runner:runner /jobs

RUN mkdir -p /opt/nsjail/rootfs/usr/bin \
    /opt/nsjail/rootfs/usr/lib \
    /opt/nsjail/rootfs/usr/lib64 \
    /opt/nsjail/rootfs/lib \
    /opt/nsjail/rootfs/lib64 \
    /opt/nsjail/rootfs/work \
    /opt/nsjail/rootfs/tmp \
    /opt/nsjail/rootfs/dev

WORKDIR /

# CMD ["/app/codexec"]