package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 给Context实例设置一个值
		c.Set("geektutu", "1111")
		// 请求前
		c.Next()
		// 请求后
		latency := time.Since(t)
		log.Print(latency)
	}
}

// Credit: Gin middleware examples (https://sosedoff.com/2014/12/21/gin-middleware.html）
// all response header add header "X-Request-Id"
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuidV4, _ := uuid.NewV4()
		ctx.Header("X-Request-Id", uuidV4.String())

		ctx.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(RequestIDMiddleware()) // as Global

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World")
	})

	// GET /user/karen
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello, %s", name)
	})

	// GET /users?name=xxx&role=xxx，role is optional
	r.GET("/users", func(c *gin.Context) {
		name := c.Query("name")
		role := c.DefaultQuery("role", "engineer")
		c.String(http.StatusOK, "%s is a %s", name, role)
	})

	// POST /form
	r.POST("/form", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "000000")

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
		})
	})

	// Query and post form
	r.POST("/posts", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		username := c.PostForm("username")
		password := c.DefaultPostForm("username", "000000")

		c.JSON(http.StatusOK, gin.H{
			"id":       id,
			"page":     page,
			"username": username,
			"password": password,
		})
	})

	// Map as querystring or postform parameters
	r.POST("/post", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		c.JSON(http.StatusOK, gin.H{
			"ids":   ids,
			"names": names,
		})
	})

	// redirect
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/index")
	})

	r.GET("/goindex", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		r.HandleContext(c)
	})

	// Grouping Routes
	defaultHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"path": c.FullPath(),
		})
	}

	// group: v1
	v1 := r.Group("/v1")
	{
		v1.GET("/posts", defaultHandler)
		v1.GET("/series", defaultHandler)
	}

	// group: v2
	v2 := r.Group("/v2")
	{
		v2.GET("/posts", defaultHandler)
		v2.GET("/series", defaultHandler)
	}

	// upload simple document
	r.POST("/upload1", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		// c.SaveUploadedFile(file, dst)
		c.String(http.StatusOK, "%s uploaded!", file.Filename)
	})

	// upload multiple documents
	r.POST("/upload2", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, "%d files uploaded!", len(files))
	})

	type student struct {
		Name string
		Age  int8
	}

	// Template
	r.LoadHTMLGlob("templates/*")

	stu1 := &student{Name: "Karen", Age: 25}
	stu2 := &student{Name: "Mickey", Age: 18}
	r.GET("/arr", func(c *gin.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gin.H{
			"title":  "Guest",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	// as Global
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Run(":9999") // listen and serve on 0.0.0.0:8080
}
