# 在docker中重新构建
# FROM golang:latest

# ENV GOPROXY=https://goproxy.cn,direct
# ENV GO111MODULE=on

# # 设置工作目录
# WORKDIR /app
# COPY . .

# # 下载依赖
# RUN go mod tidy
# # 构建（修改构建命令）
# RUN go build -o main ./cmd/main.go

# # 暴露端口
# EXPOSE 8000

# # 运行
# ENTRYPOINT ["./main"]

# 直接构建好后放到docker中运行
FROM scratch

WORKDIR /app
COPY . .
EXPOSE 8000
ENTRYPOINT [ "./main" ]