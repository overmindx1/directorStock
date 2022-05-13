package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/overmind/go_test/controller"
	"github.com/overmind/go_test/middlewares"
	"github.com/overmind/go_test/rootVar"
	"github.com/overmind/go_test/task"
)

func main() {
	// 啟用Sql
	sqlConnect := rootVar.InitDatabaseOrm()
	defer sqlConnect.Close()
	os.Setenv("TZ", "Asia/Taipei")
	f, _ := os.Create("./log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	//判斷Domain
	domains := []string{"stock.hapopo.com", "localhost"}
	r.Use(middlewares.AllowDomains(domains))

	// 限制每個人的讀取量
	r.Use(middlewares.SetRequestLimit())

	//壓縮Response
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.LoadHTMLFiles("view/index.html")
	r.Static("/css", "view/css")
	r.Static("/js", "view/js")
	r.Static("/fonts", "view/fonts")
	r.Static("/icons", "view/icons")
	r.StaticFile("/favicon.ico", "view/favicon.ico")
	r.StaticFile("/sw.js", "view/sw.js")
	r.StaticFile("/manifest.json", "view/manifest.json")
	// index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 抓初始會長資料
	r.GET("/stockitems", controller.GetDirectorSelection)

	// 開始比對會員資料
	r.POST("/mss", controller.GetMemberSpreadsheetCompare)

	// 開始排程工作
	task.InitTask()
	// ACME Server - For CertBot
	//ACME := gin.Default()
	//ACME.GET("/.well-known/acme-challenge/:acmename", controller.AcmeHandle)
	//ACME.Run(":80")

	// r.Run(":9527") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.RunTLS(":443", "./ssl/fullchain.pem", "./ssl/privkey.pem")
}
