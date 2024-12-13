package main

import (
	"context"
	"fmt"
	"net"
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

// 添加文件描述符环境变量的常量
const LISTEN_FDS_START = 3
const ENV_LISTEN_FDS = "LISTEN_FDS"
const ENV_LISTEN_PID = "LISTEN_PID"

// endless 热更新采用创建子进程后，将原进程退出的方式
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
func main() {
	ServerSetting := setting.ServerSetting
	fd := LISTEN_FDS_START
	var listener net.Listener
	var err error
	log.Info("uintptr(fd)", uintptr(fd))

	log.InfoF("os.Getenv(ENV_LISTEN_PID) %s, fmt.Sprint(os.Getpid()) %s", os.Getenv(ENV_LISTEN_PID), fmt.Sprint(os.Getpid()))
	// 判断是否为子进程
	if os.Getenv(ENV_LISTEN_PID) != "" {
		log.InfoF("starting server reuse listener at new process pid : %d ...", syscall.Getpid())

		file := os.NewFile(uintptr(fd), "")
		listener, err = net.FileListener(file)
		if err != nil {
			log.FatalF("Failed to create listener from file descriptor: %v", err)
		}
	} else {
		log.InfoF("starting server first time with pid : %d ...", os.Getpid())
		listener, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", ServerSetting.HTTPPort))
		if err != nil {
			log.FatalF("Failed to create listener: %v", err)
		}
	}

	router := routers.InitRouter()
	srv := &http.Server{
		// Addr:           fmt.Sprintf("0.0.0.0:%d", ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    ServerSetting.ReadTimeout,
		WriteTimeout:   ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.InfoF("starting server at port: %d, pid : %d ...", ServerSetting.HTTPPort, syscall.Getpid())
		if err := srv.Serve(listener); err != nil {
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
		listener.Close()
		os.Exit(0)
	case <-reload:
		log.Info("Triggering graceful restart ...")

		// 启动一个新进程
		if err := startNewProcess(listener); err != nil {
			log.ErrorF("Failed to start new process: %v", err)
			// 添加重试逻辑
			for i := 0; i < 3; i++ {
				time.Sleep(time.Second)
				if err := startNewProcess(listener); err == nil {
					break
				}
			}
		}
	}
}

// startNewProcess 启动新进程
func startNewProcess(listener net.Listener) error {
	// 获取 listener 的文件描述符
	listenerFile, err := listener.(*net.TCPListener).File()
	if err != nil {
		return fmt.Errorf("failed to get listener file:%v", listener)
	}
	defer listenerFile.Close()

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path")
	}

	// 设置新进程的环境变量和参数
	env := os.Environ()
	env = append(env,
		fmt.Sprintf("%s=%d", ENV_LISTEN_FDS, 1),
		fmt.Sprintf("%s=%d", ENV_LISTEN_PID, os.Getpid()))

	args := os.Args[1:]

	// 创建新进程
	process, err := os.StartProcess(executable, args, &os.ProcAttr{
		Dir:   ".",
		Env:   env,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, listenerFile}, // 额外传递 listener 的文件描述符
	})

	if err != nil {
		return fmt.Errorf("failed to start new process: %v", err)
	}

	// 检查新创建线程的装阿嚏
	if err := process.Signal(syscall.Signal(0)); err != nil {
		return fmt.Errorf("new process closed err:%v", err)
	}

	// 释放这个新创建的子进程，让它能独立运行
	return process.Release()
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
