# 定义一个函数，用于杀死指定进程
kill_process = $(shell ps -ef | grep $(1) | awk '{print $2}' | xargs kill -9 2>/dev/null)

# 定义一个目标，用于停止 gin-exmaple 进程
stop:
	@$(call kill_process, gin-exmaple)
	@echo "Stopped gin-exmaple process."


# 生成自动API文档
# https://www.lixueduan.com/posts/go/swagger/
swagInit:
# 注意要在项目的根目录下执行，并指定服务的入口地址，这样swag才能根据当前目录生成其他服务的代码
	@echo "===== swag init to generate auto api ====="
	@swag init -g ./cmd/main.go

# 构建docker镜像
buildDockerImage:
# -t 指定名称为gin-blog-docker, . 构建的内容为当前上下文
	@echo "===== create docker image named gin-blog-docker ====="
	@docker build -t gin-blog-docker .

buildProject: swagInit
# 直接将生成的可执行文件静态链接到所依赖的库
	@echo "===== Static linked executable file ====="
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

runDocker: buildProject buildDockerImage
	@echo "===== start run in docker container gin-example-app ====="
	@mkdir -p runtime/logs
	@chmod -R 777 runtime/logs
	@docker run -d \
		--name gin-example-app	\
		-p 8000:8000	\
		-v $(PWD)/runtime:/app/runtime \
		gin-blog-docker

clean:
	@docker rm -f gin-example-app
	@docker rmi gin-blog-docker