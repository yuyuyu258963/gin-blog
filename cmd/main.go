package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin_example/models"
	"gin_example/pkg/gredis"
	log "gin_example/pkg/logging"
	"gin_example/pkg/setting"
	"gin_example/routers"
)

func init() {
	// 1. 抽离出来统一加载的好处就是可以控制加载的流程
	// 比如必须先让setting配置了，然后再进行别的操作
	// 2. 增加可读性和可见性，不然还得去找那些出现了问题
	setting.Setup()
	log.Setup()
	models.Setup()
	gredis.Setup()
}

// endless 热更新采用创建子进程后，将原进程退出的方式
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	router := routers.InitRouter()
	ServerSetting := setting.ServerSetting
	srv := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    ServerSetting.ReadTimeout,
		WriteTimeout:   ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.InfoF("starting server at port: %d ...", ServerSetting.HTTPPort)
		if err := srv.ListenAndServe(); err != nil {
			log.InfoF("Failed Listen: %v \n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	reload := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(reload, syscall.SIGUSR2)

	// 分别等待退出信号和重启型号
	select {
	case <-quit:
		log.Info("Shutdown Server ...")
		shutdown(srv)
	case <-reload:
		shutdown(srv)
		log.Info("Triggering graceful restart ...")

		// 启动一个新进程
		if err := startNewProcess(); err != nil {
			log.ErrorF("Failed to start new process: %v", err)
			// 添加重试逻辑
			for i := 0; i < 3; i++ {
				time.Sleep(time.Second)
				if err := startNewProcess(); err == nil {
					break
				}
			}
		}
	}
}

// startNewProcess 启动新进程
func startNewProcess() error {

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path")
	}

	// 设置新进程的环境变量和参数
	env := os.Environ()
	args := os.Args[1:]

	// 创建新进程
	process, err := os.StartProcess(executable, args, &os.ProcAttr{
		Dir:   ".",
		Env:   env,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})

	if err != nil {
		return fmt.Errorf("failed to start new process: %v", err)
	}

	// 释放这个新创建的子进程，让它能独立运行
	log.InfoF("new server serve pid=%d", syscall.Getpid())
	err = process.Release()
	if err != nil {
		return fmt.Errorf("failed to release new process: %v", err)
	}

	return nil
}

// shutdown 处理服务关闭
func shutdown(srv *http.Server) {
	// 父进程，等待当前请求处理完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server ShutDown:", err)
	}
	<-ctx.Done()
	log.Info("Server exiting")
}
