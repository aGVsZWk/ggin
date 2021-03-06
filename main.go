package main

import (
	"fmt"
	"net/http"

	"ggin/routers"
	"ggin/pkg/setting"
	"log"
	"os"
	"os/signal"
	"context"
	"time"
	"ggin/models"
	"ggin/pkg/logging"
	"ggin/pkg/gredis"
)

func main() {

	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)	// 对 Interrupt 信号进行捕捉至 quit 通道
	<- quit		// 阻塞，通道等待

	log.Println("Shutdown Server ...")
	ctx, cannel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cannel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
