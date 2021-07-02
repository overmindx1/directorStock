package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/overmind/go_test/controller"
	"github.com/overmind/go_test/rootVar"
	"github.com/overmind/go_test/task"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"
)

func main() {
	// 啟用Sql
	sqlConnect := rootVar.InitDatabaseOrm()
	defer sqlConnect.Close()
	os.Setenv("TZ", "Asia/Taipei")
	f, _ := os.Create("./log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	// //判斷Domain
	// r.Use(func(c *gin.Context) {
	// 	if strings.Index(c.Request.Host, "stock.hapopo.com") == -1 {
	// 		c.AbortWithStatus(401)
	// 	}
	// 	c.Next()
	// })

	// 限制每個人的讀取量
	r.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(100*time.Millisecond), 5), time.Hour // limit 5 qps/clientIp and permit bursts of at most 5 tokens, and the limiter liveness time duration is 1 hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429) // handle exceed rate limit request
	}))

	//壓縮Response
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//寫Log

	r.LoadHTMLFiles("view/index.html")
	r.Static("/css", "view/css")
	r.Static("/js", "view/js")
	r.Static("/fonts", "view/fonts")
	r.StaticFile("/favicon.ico", "view/favicon.ico")
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

	r.Run(":9527") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
