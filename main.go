package main

import (
	"gin-servers/routers"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	m := melody.New()
	defer m.Close()
	r.LoadHTMLGlob("template/*")

	r.MaxMultipartMemory = 2 << 30

	dir, err := ioutil.TempDir("./", "tmp")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(dir string) {
		err := os.RemoveAll(dir)
		if err != nil {
			log.Fatalln(err)
		}
	}(dir)

	routers.SetDir(dir)

	api := r.Group("/api")
	{
		api.POST("/upload", routers.UploadApi)
		api.POST("/download", routers.DownloadApi)
		api.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		})
	}

	r.GET("/ws/:name", func(c *gin.Context) {
		_ = m.HandleRequest(c.Writer, c.Request)
	})

	// 如果用户ws连接成功
	m.HandleConnect(func(s *melody.Session) {
		_ = m.Broadcast([]byte(s.Request.RemoteAddr + "已连接"))
	})

	// 如果用户ws连接断开
	m.HandleDisconnect(func(s *melody.Session) {
		_ = m.Broadcast([]byte(s.Request.RemoteAddr + "断开连接"))
	})

	// 发送消息广播
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		_ = m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path && s != q
		})
	})

	srv := &http.Server{
		Addr:    ":23344",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	routers.GraceApi(srv)
}
