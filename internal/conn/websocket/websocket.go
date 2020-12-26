package websocket

import (
	"net/http"
	"time"

	controllers "chat/internal/conn/websocket/controllers"
)

func Run(address string) error {
	// 初始化客户端 - websocket
	controllers.WebsocketClients = &controllers.ClientsPool{
		RegisterChan: make(chan *controllers.Client, 100),
		LeaveChan   : make(chan *controllers.Client, 100),
	}

	go controllers.WebsocketClients.Run()

	// 定义路由
	defaultServeMux := http.NewServeMux()
	defaultServeMux.Handle("/chat", &controllers.ChatController{})

	// 服务器状态监控
	//defaultServeMux.Handle("/monitor", &controllers.MonitorController{})
	//defaultServeMux.HandleFunc("/debug/pprof/", pprof.Index)
	//defaultServeMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//defaultServeMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//defaultServeMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//defaultServeMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// 搭建服务器
	s := &http.Server{
		Addr: address,
		Handler: defaultServeMux,
		ReadTimeout: 1000*time.Second,
		WriteTimeout: 1000*time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 监听服务
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
