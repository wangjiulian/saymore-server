package routers

import (
	controllers2 "com.say.more.server/internal/app/controllers"
	"com.say.more.server/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	middleware.SetAuth(r)

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 版本分组
	v1 := r.Group("/api/v1")
	{

		base := v1.Group("/base")
		{
			base.POST("/upload-file", controllers2.UploadFile)
			base.GET("/dict", controllers2.Dict)
			base.GET("/banners", controllers2.Banners)
			base.POST("/feedback", controllers2.Feedback)
		}

		// 学生相关路由
		students := v1.Group("/student")
		{
			students.POST("/session-key", controllers2.StudentSessionKey) // 获取微信mini session key
			students.POST("/login", controllers2.StudentLogin)            // 登录注册
			students.PUT("/detail", controllers2.StudentEdit)             // 编辑
			students.GET("/detail", controllers2.StudentDetail)           // 详情
			students.GET("/trial", controllers2.StudentTrial)             // 试听机会
		}

		// 老师相关路由
		teachers := v1.Group("/teacher")
		{
			teachers.GET("/recommends", controllers2.TeacherRecommends)             // 推荐
			teachers.GET("/search", controllers2.TeacherSearch)                     // 搜索
			teachers.GET("/:id/detail", controllers2.TeacherDetail)                 // 详情
			teachers.GET("/:id/availabilities", controllers2.TeacherAvailabilities) // 详情
		}

		// 课程相关路由
		courses := v1.Group("/course")
		{
			courses.POST("/book-regular", controllers2.BookRegularCourse)                     // 预约正式课
			courses.POST("/book-trial", controllers2.BookTrialCourse)                         // 预约试听课
			courses.POST("/:course_id/cancel", controllers2.CancelCourse)                     // 取消课程
			courses.GET("/list", controllers2.Courses)                                        // 课程列表
			courses.GET("/:course_id/detail", controllers2.CourseDetail)                      // 课程详情
			courses.POST("/:course_id/evaluation", controllers2.AddCourseEvaluation)          // 添加课程评价
			courses.GET("/:course_id/evaluation/detail", controllers2.CourseEvaluationDetail) // 获取课程评价
			courses.GET("/evaluations", controllers2.CourseEvaluations)                       // 课程评价
		}

		// 学生课时包相关路由
		packages := v1.Group("/student-package")
		{
			packages.GET("/list", controllers2.StudentPackageList)                         // 课时包列表
			packages.GET("/:student_package_id/detail", controllers2.StudentPackageDetail) // 课时包明细
		}
	}

	return r
}
