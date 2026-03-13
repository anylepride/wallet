FROM golang:1.26.1 AS builder
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /server
COPY . .

# 下载依赖库
RUN go mod download

# 构建成功后程序存放路径
RUN mkdir /server/bin
ARG SERVERNAME

# 编译构建钱包服务程序
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" \
    -trimpath -buildvcs=false \
    -o /server/bin/$SERVERNAME ./cmd/$SERVERNAME

FROM alpine:3.23
RUN apk update && apk add --update tzdata

ARG SERVERNAME
ENV TZ=Asia/Shanghai
ENV SERVERNAME=$SERVERNAME

WORKDIR /server
COPY --from=builder /server .

EXPOSE 8000 8020
ENTRYPOINT ["sh", "-c", "/server/bin/$SERVERNAME"]