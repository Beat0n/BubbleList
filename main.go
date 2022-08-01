package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)
type Todo struct{
	ID int`json:"id"`
	Title string `json:"title" form:"title"` 
	Status bool `json:"status" form:"status"`
}

func initMySQL()(err error){
	dsn := "root:200214@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil{
		return
	}
	return DB.DB().Ping()
}

func main() {
	err := initMySQL()
	if err != nil{
		panic(err)
	}

	defer DB.Close()
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context){
		c.HTML(200,"index.html",nil)
	})

	v1Group := r.Group("v1")
	{
		//添加
		v1Group.POST("/todo", func(c *gin.Context){
			var todo Todo
			c.ShouldBind(&todo)
			if err := DB.Create(&todo).Error;err !=nil {
				c.JSON(200, gin.H{"error": err.Error()})
			}else{
				c.JSON(200, todo)
			}
		})
		//查看所有
		v1Group.GET("/todo", func(c *gin.Context){
			var todoList []Todo
			if err := DB.Find(&todoList).Error;err !=nil {
				c.JSON(200, gin.H{"error": err.Error()})
			}else {
				c.JSON(200, todoList)
			}
		})
		//查看某一个
		v1Group.GET("/todo/:id", func(c *gin.Context){
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "id不存在"})
				return
			}
			var todo Todo
			if err := DB.Where("id=?", id).First(&todo).Error;err !=nil {
				c.JSON(200, gin.H{"error": err.Error()})
			}else {
				c.JSON(200, todo)
			}
		})
		//修改某一个
		v1Group.PUT("/todo/:id", func(c *gin.Context){
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "id不存在"})
				return
			}
			var todo Todo
			if err := DB.Where("id=?", id).First(&todo).Error;err !=nil{
				c.JSON(200, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error;err != nil {
				c.JSON(200, gin.H{"error": err.Error()})
			}else {
				c.JSON(200, todo)
			}
		})
		//删除某一个
		v1Group.DELETE("/todo/:id", func(c *gin.Context){
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "id不存在"})
				return
			}
			if err := DB.Where("id=?", id).Delete(&Todo{}).Error;err !=nil{
				c.JSON(200, gin.H{"error": err.Error()})
				return
			}else {
				c.JSON(200, gin.H{id: "Deleted"})
			}
		})
		//删除所有
	}
	r.Run(":8080")
}