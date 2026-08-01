# Git 备份管理系统 - 构建入口
# 常用命令：
#   make            # 构建前端 + 编译内嵌前端的后端二进制
#   make build-web  # 仅构建前端
#   make run        # 构建并直接运行（:8080）
#   make clean      # 清理构建产物

APP_NAME := git-backup-server

.PHONY: all build build-web build-server run clean

all: build

build: build-web build-server

build-web:
	cd web && npm ci && npm run build

build-server:
	rm -rf server/webroot/dist
	mkdir -p server/webroot/dist
	cp -r web/dist/. server/webroot/dist/
	cd server && go build -trimpath -ldflags="-s -w" -o ../$(APP_NAME) .

run: build
	./$(APP_NAME)

clean:
	rm -rf server/webroot/dist $(APP_NAME)
	rm -rf web/dist
