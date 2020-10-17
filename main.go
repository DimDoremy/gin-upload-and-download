package main

import (
	"gin-servers/routers"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
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

		api.GET("/send", routers.MessageSendApi)
		api.GET("/receive", routers.MessageReceiveApi)
	}

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
