package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Project struct {
	ID         int    `json:"id" gorm:"column:id;"`
	Title      string `json:"title" gorm:"column:title;"`
	Decription string `json:"decription" gorm:"column:decription;"` // Sửa tên trường
	Day        string `json:"day" gorm:"column:day;"`
	Role       string `json:"role" gorm:"column:role;"`
	Member     int    `json:"member" gorm:"column:member;"`
	Tech       string `json:"tech" gorm:"column:tech;"`
}

func (Project) TableName() string {
	return "project"
}

type Projectcreate struct {
	Title      string `json:"title" gorm:"column:title;"`
	Decription string `json:"decription" gorm:"column:decription;"` // Sửa tên trường
	Day        string `json:"day" gorm:"column:day;"`
	Role       string `json:"role" gorm:"column:role;"`
	Member     int    `json:"member" gorm:"column:member;"`
	Tech       string `json:"tech" gorm:"column:tech;"`
}

func (Projectcreate) TableName() string { return Project{}.TableName() }

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/cvminh?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)
	item := Project{
		ID:         1,
		Title:      "anh bốn",
		Decription: "anh năm", // Sửa tên trường
		Day:        "01 / 10 / 2023",
		Role:       "FRONT-END",
		Member:     5,
		Tech:       "PHP",
	}
	r := gin.Default()
	r.Use(cors.Default())
	v1 := r.Group(`/cv`)
	{
		item := v1.Group(`/project`)
		{
			item.POST("", Create_project(db))
			item.GET("", Get_all_projects(db))
			item.GET("/:id", Get_project_byId(db))
			item.PATCH("/:id", Update_project(db))
			item.DELETE("/:id", Delete_project(db))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": item,
		})
	})
	r.Run()
}

func Create_project(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data Projectcreate
		err := c.ShouldBind(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		Err := db.Create(&data).Error
		if Err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": Err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func Get_project_byId(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data Project
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Không có param",
			})
			return
		}
		data.ID = id
		Err := db.First(&data).Error

		if Err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": Err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}
func Update_project(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data Projectcreate
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Không có param",
			})
			return
		}
		Err := c.ShouldBind(&data)
		if Err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": Err.Error(),
			})
			return
		}

		Erro := db.Where("id=?", id).Updates(&data).Error

		if Erro != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": Erro,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Thành Công",
		})
	}
}

func Delete_project(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Không có param",
			})
			return
		}

		Erro := db.Table(Project{}.TableName()).Where("id=?", id).Delete(nil).Error
		// Erro := db.Table(TodoItem{}.TableName()).Where("id=?", id).Updates(map[string]interface{
		// 	"status":"Delete",
		// }).Error
		if Erro != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": Erro,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Thành Công",
		})
	}
}
func Get_all_projects(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		var result []Project
		Erro := db.Table(Project{}.TableName()).Find(&result).Error
		if Erro != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": Erro.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Thành Công",
			"data":    result,
		})
	}
}
