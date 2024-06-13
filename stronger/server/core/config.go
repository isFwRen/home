package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"server/global"
	global2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
)

func InitConfig() {
	//pflag.StringP("configFile", "c", "", "choose config file.")
	//pflag.StringP("downloadPath", "p", "", "set download path.")
	//pflag.StringP("downloadBillEnd", "b", "", "set bill End.")
	//pflag.StringP("port", "P", "", "set port.")
	//pflag.StringP("projectCode", "n", "", "set Project Code.")
	//pflag.Parse()

	// 优先级: 命令行 > 环境变量 > 默认值
	v := viper.New()
	err := v.BindPFlags(pflag.CommandLine)
	if err != nil {
		return
	}
	v.SetEnvPrefix("gva")
	err1 := v.BindEnv("configFile")
	if err1 != nil {
		return
	} // GVA_CONFIGFILE

	//configFile := v.GetString("configFile")
	//if configFile == "" {
	//	configFile = defaultConfig
	//}

	var confFile = map[string]string{
		"develop": "config.yaml",
		"test":    "config_test.yaml",
		"prod":    "config_prod.yaml",
	}

	v.SetConfigFile(confFile[global2.GConfig.System.Env])
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GConfig); err != nil {
		fmt.Println(err)
	}
	global.GVp = v

	downloadPath := v.GetString("downloadPath")
	downloadBillEnd := v.GetString("downloadBillEnd")
	port := v.GetInt("port")
	projectCode := v.GetString("projectCode")
	global.GConfig.System.DownloadPath = downloadPath
	global.GConfig.System.DownloadBillEnd = downloadBillEnd
	global.GConfig.System.Port = port
	global.GConfig.System.ProCode = projectCode
	if port == 0 {
		port = global.GConfig.System.Port
		if port == 0 {
			port = 9999
		}
	}
	global.GConfig.System.Port = port

	if projectCode == "" {
		global.GConfig.System.ProCode = "common"
	}
	fmt.Println(global.GConfig.System)
	// 获取命令行参数 （可以直接用端口跑，后期接受参数是项目名称，根据项目名称去找配置判断使用哪一个端口）
	//if len(os.Args) == 1 {
	//	panic("请指定端口")
	//}

	//fmt.Println("命令行参数:")
	//for k, v := range os.Args {
	//	fmt.Printf("args[%v]=[%v]\n", k, v)
	//}
	//global.GConfig.System.Port = 9999
	//if len(os.Args) > 1 {
	//	global.GConfig.System.Port, _ = strconv.Atoi(os.Args[1])
	//}
	//global.GConfig.System.ProCode = "common"
	//if len(os.Args) > 2 {
	//	global.GConfig.System.ProCode = os.Args[2]
	//}
}
