.PHONY: build-linux build-local run migrate clean test

# 编译 Linux 可执行程序
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o adcms-server cmd/server/main.go

# 编译本地可执行程序
build-local:
	go build -o adcms-server cmd/server/main.go

# 运行服务
run:
	go run cmd/server/main.go

# 运行数据库迁移
migrate:
	go run cmd/migrate/main.go

# 清理构建文件
clean:
	rm -f adcms-server adcms-migrate

# 运行测试
test:
	go test -v ./...

# 安装依赖
deps:
	go mod download
	go mod tidy
