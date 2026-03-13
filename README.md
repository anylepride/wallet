# wallet
wallet service

# 程序运行
wallet_server: 
go run cmd/server

wallet_gateway:
go run cmd/gateway

# 构建wallet_server程序镜像
docker build -t wallet/server --build-args SERVERNAME=server .

# 构建gateway程序镜像
docker build -t wallet/gateway --build-args SERVERNAME=gateway .