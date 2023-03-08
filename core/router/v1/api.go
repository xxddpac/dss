package v1

import "github.com/gin-gonic/gin"

func Register(v1 *gin.RouterGroup) {
	rule := v1.Group("/rule")
	{
		rule.POST("", Rule.Post)
		rule.GET("", Rule.Get)
		rule.PUT("", Rule.Put)
		rule.DELETE("", Rule.Delete)
		rule.GET("enum", Rule.Enum)
	}
	task := v1.Group("/task")
	{
		task.POST("", Task.Post)
		task.GET("", Task.Get)
		task.GET("query", Task.Query)
		task.GET("enum", Task.Enum)
	}
	scan := v1.Group("/scan")
	{
		scan.GET("", Scan.Get)
		scan.GET("trend", Scan.Trend)
		scan.GET("remind", Scan.Remind)
		scan.GET("location", Scan.Location)
		scan.DELETE("clear", Scan.Clear)
	}
	grpc := v1.Group("/grpc")
	{
		grpc.GET("client", Grpc.Get)
	}
}
