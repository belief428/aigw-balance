.PHONY: start build

RELEASE_VERSION = v1.0.1

ROOT = release
NAME = app.server

# Basic go commands
GOCMD = go
GORUN = $(GOCMD) run
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

all: run

run:
	@$(GORUN) main.go

migrate:
	@$(GORUN) main.go migrate

test:
	@$(GOTEST) -v ./...

clean:
	@$(GOCLEAN)
	@rm -f $(ROOT)/linux/$(NAME)

# Cross compilation
build-windows:
	$(GOBUILD) -ldflags "-w -s" -o $(ROOT)/windows/$(NAME).exe main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "-w -s" -o $(ROOT)/linux/$(NAME) main.go

build-arm:
	@echo "build-arm "$(TAGS)
ifeq ($(TAGS), rk3399)
	CC="/usr/local/src/gcc-linaro-7.4.1-2019.02-x86_64_aarch64-linux-gnu/bin/aarch64-linux-gnu-gcc" CGO_ENABLED="1" GOARCH=arm64 $(GOBUILD) -ldflags "-w -s" -o $(ROOT)/arm/$(NAME) main.go
else
	CC="/usr/local/src/arm-2014.05/bin/arm-none-linux-gnueabi-gcc" CGO_ENABLED="1" GOARCH=arm GOARM=5 $(GOBUILD) -ldflags "-w -s" -o $(ROOT)/arm/$(NAME) main.go
endif
build-docker:
	docker build --rm -t hanson/app_server:$(RELEASE_VERSION) .
	docker run --name app -p 8088:8088 -d hanson/app_server:$(RELEASE_VERSION)
	#docker push hanson/app_server:$(RELEASE_VERSION)

push:
	@echo "发布版本"
	curl -T $(ROOT)/arm/$(NAME) ftp://root:suxx@192.168.137.2/home/suxx/ai/
	ssh root@192.168.137.2 'cd /home/suxx/ai && rm -rf aigw && mv app.server aigw && chmod 755 aigw && killall aigw'

help:
	@echo "make 直接运行程序"
	@echo "make build 编译go代码生成二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试case"
	@echo "make docker 构建docker镜像"
