package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"iris_study/iris_study/CmsProject/config"
	"iris_study/iris_study/CmsProject/controller"
	"iris_study/iris_study/CmsProject/datasource"
	"iris_study/iris_study/CmsProject/service"
	"time"
)

func main() {
	app := newApp()


	//应用App设置
	configation(app)

	//路由设置
	mvcHandle(app)
	//从配置文件中读取配置信息
	config := config.InitConfig()
	addr := ":" + config.Port
	app.Run(
		iris.Addr(addr),                               //在端口9000进行监听
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

//构建App
func newApp() *iris.Application {
	app := iris.New()

	//设置日志级别  开发阶段为debug
	app.Logger().SetLevel("debug")

	//注册静态资源
	//如果浏览器请求时使用的是"/manage/static"，则会替换成"./static"
	//app.StaticWeb("/static", "./static")
	//app.StaticWeb("/manage/static", "./static")

	//注册视图文件
	app.RegisterView(iris.HTML("./static", ".html"))
	app.Get("/", func(context context.Context) {
		context.View("index.html")
	})

	return app
}

/**
 * MVC 架构模式处理
 */
func mvcHandle(app *iris.Application) {

	//启用session
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookie",
		Expires: 24 * time.Hour,
	})

	//构建数据库操作引擎
	engine := datasource.NewMysqlEngine()

	//管理员模块功能
	//将数据库操作引擎放入service管理模块当中
	adminService := service.NewAdminService(engine)
	//解析/admin路由组下的请求
	admin := mvc.New(app.Party("/admin"))
	//
	admin.Register(
		adminService,
		sessManager.Start,
	)
	//使用Handle方法将AdminController传入进去设置好，当出现admin路由组请求时，使用AdminController进行解析和处理
	admin.Handle(new(controller.AdminController))

	//用户功能模块
	userService := service.NewUserService(engine)

	user := mvc.New(app.Party("/v1/users"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controller.UserController))

	//统计功能模块
	statisService := service.NewStatisService(engine)
	statis := mvc.New(app.Party("/statis/{model}/{date}/"))
	statis.Register(
		statisService,
		sessManager.Start,
	)
	statis.Handle(new(controller.StatisController))
}

/**
 * 项目设置
 */
func configation(app *iris.Application) {

	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error ",
			"data":   iris.Map{},
		})
	})
}

