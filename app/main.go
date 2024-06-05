package main

import (
	"LittleVideo/app/route"
	"LittleVideo/app/store"
	"LittleVideo/config"
	"LittleVideo/middleware"
	"LittleVideo/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	config.InitEnvConf()
	util.InitLogger()

	config.LoadConfig()

	//store.ConnectRedis()
	store.ConnectMysql()
	store.InitMemoryCache()
	//启动http服务
	log.Println("服务正在启动，监听端口:", util.GetLocalIp()+":504", ",PID:", strconv.Itoa(os.Getpid()))

	if !config.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.Recovery())
	route.SetupRouter(r)
	server := &http.Server{
		Addr:         ":504",
		WriteTimeout: 20 * time.Second,
		Handler:      r,
	}
	//err := gracehttp.Serve(server)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("服务器启动失败:", err.Error())
	}
}
