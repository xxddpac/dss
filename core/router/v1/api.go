package v1

import "github.com/gin-gonic/gin"

func Register(v1 *gin.RouterGroup) {
	// this rule defines the scanning host and port
	rule := v1.Group("/rule")
	{
		rule.POST("", Rule.Post)          // add rule
		rule.GET("", Rule.Get)            // list rule
		rule.PUT("", Rule.Put)            // modify rule
		rule.DELETE("", Rule.Delete)      // delete rule
		rule.GET("/type/enum", Rule.Enum) // display rule type enum
	}
	// task scheduler
	task := v1.Group("/task")
	{
		task.POST("", Task.Post) // push task to redis
	}
	//scan result
	scan := v1.Group("/scan")
	{
		scan.GET("", Scan.Get)              // list scan result
		scan.GET("trend", Scan.Trend)       // last 7 days scan result trend
		scan.GET("remind", Scan.Remind)     // send notification if new port open  (schedule)
		scan.GET("location", Scan.Location) // location select from web front
		scan.DELETE("clear", Scan.Clear)    // save last 7 days result (schedule)
	}
}
