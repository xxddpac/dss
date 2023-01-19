package v1

import "github.com/gin-gonic/gin"

func Register(v1 *gin.RouterGroup) {
	// this rule defines the scanning host and port
	rule := v1.Group("/rule")
	{
		rule.POST("", Rule.Post)     // add rule
		rule.GET("")                 // list rule
		rule.PUT("")                 // modify rule
		rule.DELETE("", Rule.Delete) // delete rule
	}
	// task scheduler
	task := v1.Group("/task")
	{
		task.POST("", Task.Post) // push task to redis
	}
}
