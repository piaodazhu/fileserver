package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/russross/blackfriday/v2"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

var configFile string

type ConfigModel struct {
	Ip       string `yaml:"ip"`
	Port     uint16 `yaml:"port"`
	RootPath string `yaml:"rootPath"`
	Password string `yaml:"password"`
	DocFile  string `yaml:"docFile"`
}

var config ConfigModel

func readInConfig(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

var uploadTaskMap sync.Map

type uploadTask struct {
	uploadpath    string
	localpath   string
	totalSize   int64
	writtenSize int64
	createAt    time.Time
}

func initRootPath() {
	if !strings.HasPrefix(config.RootPath, "/var/") {
		panic("rootPath must start with /var/")
	}
	os.MkdirAll(config.RootPath, os.ModePerm)
}

func main() {
	flag.StringVar(&configFile, "c", "fileserver.yaml", "config file")
	flag.Parse()

	readInConfig(configFile)
	initRootPath()

	doc, err := os.ReadFile(config.DocFile)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Static("/static", config.RootPath)
	engine.POST("/upload/*path", func(c *gin.Context) {
		if c.GetHeader("password") != config.Password {
			c.String(http.StatusUnauthorized, "[ERR] Unauthorized\n")
			return
		}

		file, _ := c.FormFile("file")
		if file == nil {
			c.String(http.StatusBadRequest, "[ERR] Invalid file\n")
			return
		}

		dst := config.RootPath + c.Param("path")
		if dst[len(dst)-1] == '/' {
			dst += file.Filename
		}

		c.SaveUploadedFile(file, dst)
		c.String(http.StatusOK, fmt.Sprintf("[OK] %s --> %s\n", file.Filename, c.Param("path")))
	})

	engine.POST("/rawupload/*path", func(c *gin.Context) {
		body := c.Request.Body
		defer body.Close()

		if c.GetHeader("password") != config.Password {
			c.String(http.StatusUnauthorized, "[ERR] Unauthorized\n")
			return
		}

		targetPath := c.Param("path")
		if targetPath[len(targetPath)-1] == '/' {
			c.String(http.StatusBadRequest, "[ERR] Must provide a complete file path, not only directory\n")
			return
		}

		localPath := config.RootPath + targetPath
		stat, err := os.Stat(targetPath)
		if err == nil && stat.IsDir() {
			c.String(http.StatusBadRequest, "[ERR] Provided path is an existing directory\n")
			return
		}

		task := uploadTask{
			uploadpath: targetPath,
			localpath: localPath,
			totalSize: c.Request.ContentLength,
			writtenSize: 0,
			createAt: time.Now(),
		}
		if v, exists := uploadTaskMap.LoadOrStore(targetPath, &task); exists {
			c.String(http.StatusBadRequest, fmt.Sprintf("[ERR] Provided path is being upload now, start at %s\n", v.(*uploadTask).createAt))
			return
		}
		defer uploadTaskMap.Delete(targetPath)
		
		os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
		targetFile, err := os.Create(localPath)
		if err != nil {
			c.String(http.StatusInternalServerError, "[ERR] Cannot create target file on server\n")
			return
		}
		defer targetFile.Close()

		buf := make([]byte, 65536)
		for {
			n, err := body.Read(buf)
			if err != nil {
				break
			}
			task.writtenSize += int64(n)

			written, err := targetFile.Write(buf[:n])
			if err != nil || written != n {
				c.String(http.StatusInternalServerError, fmt.Sprintf("[ERR] Write target file error: %s\n", err.Error()))
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("[OK] %s upload finish\n", c.Param("path")))
	})

	engine.GET("/progress/*path", func(c *gin.Context) {
		targetPath := c.Param("path")
		v, exists := uploadTaskMap.Load(targetPath)
		if !exists {
			c.Status(http.StatusNotFound)
			return
		}
		taskP := v.(*uploadTask)
		if taskP == nil {
			c.Status(http.StatusNotFound)
			return
		}
		task := *taskP
		
		cost := time.Since(task.createAt).Seconds()
		speed := float64(task.writtenSize) / (1024 * 1024) / cost
		estimatedCost := cost * (float64(task.totalSize)/float64(task.writtenSize))
		c.String(http.StatusOK, fmt.Sprintf("%.2f%% [%.1f s / %.1f s] [%d B / %d B] %.2f MB/s\n", float64(task.writtenSize)*100/float64(task.totalSize), cost, estimatedCost, task.writtenSize, task.totalSize, speed))
	})

	engine.GET("/list/*path", func(c *gin.Context) {
		path := c.Param("path")
		localPath := config.RootPath + path
		info, err := os.Stat(localPath)
		if err != nil {
			c.String(http.StatusNotFound, "[ERR] Not found\n")
			return
		}

		var html strings.Builder
		if !info.IsDir() {
			html.WriteString(fmt.Sprintf("<a href=\"/static%s\" style=\"color: black; text-decoration: underline;\" download>%s</a><br>", path, info.Name()))
			c.Data(http.StatusOK, "html", []byte(html.String()))
			return
		}

		entries, err := os.ReadDir(localPath)
		if err != nil {
			c.String(http.StatusNotFound, "[ERR] Not found\n")
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				html.WriteString(fmt.Sprintf("<a href=\"/static%s/%s\" style=\"color: black; text-decoration: underline;\" download>%s</a><br>", path, entry.Name(), entry.Name()))
			} else {
				html.WriteString(fmt.Sprintf("<a href=\"/list%s/%s\" style=\"color: blue; text-decoration: underline; font-weight: bold\">%s</a><br>", path, entry.Name(), entry.Name()))
			}
		}
		c.Data(http.StatusOK, "html", []byte(html.String()))
	})
	engine.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", blackfriday.Run(doc))
	})
	if err := engine.Run(fmt.Sprintf("%s:%d", config.Ip, config.Port)); err != nil {
		panic(err)
	}
}
