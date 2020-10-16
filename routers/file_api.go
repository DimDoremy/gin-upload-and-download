package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Download struct {
	Name string `json:"name"`
	Unix int64  `json:"unix"`
}

var dir string

func SetDir(tmp string) {
	dir = tmp
}

func UploadApi(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, "上传出错")
	}

	unix := time.Now().UnixNano()
	dst := strconv.FormatInt(unix, 10) + file.Filename[0:2] + Ext(file.Filename)
	err = c.SaveUploadedFile(file, dir+"/"+dst)
	if err != nil {
		//log.Println(err)
		c.String(http.StatusInternalServerError, "上传文件出错")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": dst,
		"unix": unix,
	})
}

func DownloadApi(c *gin.Context) {
	var json Download
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体内容解析失败"})
		return
	}
	c.FileAttachment("./"+dir+"/"+json.Name, json.Name)
}
