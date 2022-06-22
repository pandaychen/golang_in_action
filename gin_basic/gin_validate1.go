package main

//gin binding rules

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpParam struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

func main() {
	r := gin.Default()

	//client: curl -H "Content-type: application/json" -X POST -d '{"name":"pandaychen","age":18,"email":"panda@123.com","password":"11111","re_password":"11111"}' http://127.0.0.1:8999/signup
	r.POST("/signup", func(c *gin.Context) {
		var u SignUpParam
		//if err := c.ShouldBindJSON(&u); err != nil {
		if err := c.ShouldBind(&u); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, "success")
	})

	_ = r.Run(":8999")
}
