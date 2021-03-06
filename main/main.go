/*
@Time : 2019/9/16 18:46
@Author : zxr
@File : main
@Software: GoLand
*/
package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"poetry/app/bootstrap"
	"poetry/config/define"
	"poetry/libary/server"
	"runtime"
)
var (
	confFile string
	HttpEnv string
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	HttpEnv = os.Getenv("http_env")
	if HttpEnv == define.DevEnvStr{
		HttpEnv = define.DevEnvStr
	}else{
		HttpEnv = define.ProductEnvStr
	}
}

func initConfFile() {
	confFile = define.BaseDir + "/config/"+HttpEnv+"/config.json"
}

func init()  {
	initEnv()
	initConfFile()
}

func main() {
	var err error
	if err = bootstrap.InitBootstrap(confFile);err!=nil{
		fmt.Println(err)
		return
	}
    //启动prometheus  metrics端口
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(bootstrap.G_Conf.MetricsPortStr, nil))
	}()
   server.StartHttp()
   return
}
