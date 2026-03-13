# wallet
wallet service

# 构建wallet_server程序镜像
docker build -t wallet/server --build-args SERVERNAME=server .

# 构建gateway程序镜像
docker build -t wallet/gateway --build-args SERVERNAME=gateway .