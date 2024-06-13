package initialize

import (
	"fmt"
	_ "github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
	"go.uber.org/zap"
	"os"
	"server/global"
)

func InitCouchdb(){
	couchCig := global.GConfig.Couchdb
	fmt.Println("couchCig", couchCig)
	if couchCig.Url != "" {
		//在使用kivik连接couchdb的时候，除了引入kivik/v3，还要引入couchdb/v3,否则连不上驱动
		client, err := kivik.New("couch", couchCig.Url)
		if err != nil{
			global.GLog.Error("CouchDBClient初始化失败", zap.Any("err", err))
			os.Exit(0)
		}else {
			fmt.Println("CouchDB初始化成功")
			fmt.Println("CouchDB client", client)
			global.GCouchdbClient = client
		}
		//数据库名称
		//constDb := client.DB(context.Background(),"lp2")
		//if constDb != nil {
		//	fmt.Println("CouchDB初始化成功")
		//	global.GCouchdb = constDb
		//} else {
		//	global.GLog.Error("CouchDB初始化失败", zap.Any("err", err))
		//	os.Exit(0)
		//}
	}
}
