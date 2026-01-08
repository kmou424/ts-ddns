# 构建阶段
FROM golang:latest AS builder

WORKDIR /build

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译静态二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o ts-ddns ./cmd/ts-ddns

# 运行阶段
FROM scratch

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/ts-ddns /app/ts-ddns

# 从构建阶段复制 CA 证书（用于 HTTPS 请求）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/ts-ddns"]
