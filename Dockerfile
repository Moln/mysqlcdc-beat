ARG GO_VERSION=1.24

FROM golang:${GO_VERSION}-alpine as base

# 配置go代理
ENV GOPROXY=https://goproxy.cn

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

FROM base AS build

ARG TARGETOS
ARG TARGETARCH

COPY . .

# 编译go程序
RUN --mount=type=cache,target=/go/pkg/mod \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -v -o /bin/cdcbeat

FROM alpine:3.20

# 安装mysql客户端
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache mariadb-connector-c mariadb-client && \
    rm /usr/lib/perl5 /usr/share/perl5/ -rf

WORKDIR /usr/share/cdcbeat

# 挂载配置目录
VOLUME /usr/share/cdcbeat/data/

# 从build阶段拷贝go程序
COPY --from=build /bin/cdcbeat /usr/share/cdcbeat/cdcbeat
# 从build阶段拷贝配置
COPY ./_meta/config/cdcbeat.yml ./cdcbeat.yml
COPY ./_meta/config/cdcbeat.reference.yml.tmpl ./cdcbeat.reference.yml
COPY --chmod=755 ./_meta/docker/docker-entrypoint /usr/local/bin/docker-entrypoint

RUN ln -s /usr/share/cdcbeat/cdcbeat /usr/local/bin/cdcbeat

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]

CMD ["/usr/share/cdcbeat/cdcbeat"]
