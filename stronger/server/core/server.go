package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"server/utils"
	//"github.com/fvbock/endless"
	//"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	socketio "server/initialize"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunServer(Router *gin.Engine, fun string) {
	//if global.GConfig.System.UseMultipoint {
	// 初始化redis服务
	//initialize.Redis()
	//}
	Router.Static("/static", "./resource/page")

	port := fmt.Sprintf(":%d", global.GConfig.System.Port)
	fmt.Println("port:", port)

	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GLog.Debug("server run success on ", zap.String("pong", port))

	utils.Fun(fun)
	fmt.Printf(`欢迎使用 %v
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:3000
`, global.GConfig.System.Name, port)

	//初始化socketIo
	socketio.InitSocketIO()
	Router.GET("/socket.io/*any", gin.WrapH(global.GSocketIo))
	Router.POST("/socket.io/*any", gin.WrapH(global.GSocketIo))

	go func() {
		if err = global.GSocketIo.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer global.GSocketIo.Close()
	s := initServer(port, Router)
	global.GLog.Error("asd", zap.Any("qwe", s.ListenAndServe()))
}
