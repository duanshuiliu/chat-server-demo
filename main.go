package main

import (
	"os"
	"go.uber.org/zap"

	config    "chat/pkg/conf"
	logger    "chat/pkg/logger"
	websocket "chat/internal/conn/websocket"
)

func init() {
	// 加载配置文件
	if err := config.Register("./config"); err != nil {
		logger.NewLogger().Error("配置文件初始化失败", zap.String("info", err.Error()))
		os.Exit(0)
	}
}

func main() {
	appConf, err := config.New("app")
	if err != nil {
		logger.NewLogger().Error("获取配置文件失败", zap.String("info", err.Error()))
		return
	}

	// 创建websocket服务器
	port,_ := appConf.String("websocket::port")
	addr := "0.0.0.0:"+port
	logger.NewLogger().Info("Start Server", zap.String("server", addr))

	websocket.Run(addr)
}
