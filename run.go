package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-minio-sync/config"
	"go-minio-sync/minio"
	"go-minio-sync/sync"
)

func run() {
	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听系统信号（优雅退出）
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// 加载配置文件
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Println("加载配置成功")

	// 初始化 MinIO 客户端
	minioClient, err := minio.NewClient(cfg)
	if err != nil {
		log.Fatalf("初始化 MinIO 客户端失败: %v", err)
	}
	log.Println("初始化 MinIO 客户端成功")

	// 确保存储桶存在
	err = minioClient.MakeBucketIfNotExists(ctx, cfg)
	if err != nil {
		log.Fatalf("创建存储桶失败: %v", err)
	}
	log.Println("创建存储桶成功")

	// 启动文件监听
	err = sync.StartFileWatcher(cfg, sync.Callback)
	if err != nil {
		log.Fatalf("启动文件监听失败: %v", err)
	}
	log.Println("启动文件监听成功")

	// 等待退出信号
	<-stop
	log.Println("收到终止信号，正在退出...")
	cancel()
	log.Println("退出完成。")
}
