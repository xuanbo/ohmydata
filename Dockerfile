# 基于 base 镜像编译
FROM golang:1.14-alpine as builder

# 工作路径
WORKDIR /apps/ohmydata

# 国内代理
ENV GOPROXY=https://goproxy.cn
# 启用模块化
ENV GO111MODULE on

# 代码复制到容器中
COPY . .

# 编译成可执行程序
RUN go build -mod=vendor -o bin/app cmd/ohmydata/main.go

# 基于 alpine 镜像运行
FROM alpine:3.12 as runner

# 工作路径
WORKDIR /apps/ohmydata

# 国内源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# Install root certificates, they are needed for email validator to work
# with the TLS SMTP servers like Gmail or Mailjet. Also add bash and grep.
RUN apk update && apk add --no-cache ca-certificates bash grep

# 设置时区为上海
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 复制可执行文件、配置文件
COPY --from=builder /apps/ohmydata/bin/app bin/
COPY --from=builder /apps/ohmydata/config config

# 暴露端口
EXPOSE 9090

# 执行
ENTRYPOINT ["./bin/app"]