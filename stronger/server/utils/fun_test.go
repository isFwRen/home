/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/13 4:11 下午
 */

package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"server/global"
	"testing"
)

func TestName(t *testing.T) {
	//convey.Convey("test fun", t, func() {
	//	Fun("nobug")
	//})

	//sh := "curl 'ftp://192.168.202.3/' -u 'myftp:myftp' -s -x 'socks5://127.0.0.1:8090'"
	//arr := strings.Split(sh," ")
	//fmt.Println(arr)
	////cmd := exec.Command("curl", "ftp://192.168.202.3/", "-u", "myftp:myftp", "-s", "-x", "socks5://127.0.0.1:8090")
	//cmd := exec.Command("/bin/sh", sh)
	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr
	//err := cmd.Run()
	//fmt.Println(err)
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Println(outStr)
	//fmt.Println(errStr)
	//err, s, s2 := ShellOut("curl -T '/Users/mjl/project/stronger/trunk/server/files/B0118/upload_xml/2021/November/11/198273192123.xml' 'ftp://192.168.202.3/lipei2.0/198273192123.xml' -u 'myftp:myftp' -f -s -x socks5://127.0.0.1:8090")
	//fmt.Println(err)
	//fmt.Println("1111111111")
	//fmt.Println(s)
	//fmt.Println("000000000000")
	//fmt.Println(s2)
	//if err != nil {
	//	return
	//}
	//https://oapi.dingtalk.com/robot/send?access_token=19ea3b4bef45e4279b323262fac79f71a7070dcf770418cea5a052ee2ca8a47e
	robot := NewRobot("19ea3b4bef45e4279b323262fac79f71a7070dcf770418cea5a052ee2ca8a47e", "SEC43282bac1b60e8825baaaf9e72a0a7ce4985f727e503ba16adf9b9be28955db6")
	err := robot.SendTextMessage("皮皮虾", []string{"18319748981"}, false)
	//err := robot.SendLinkMessage("皮皮虾", "哈喽","http://113.106.108.93:5000/","https://s1.ax1x.com/2022/07/06/jUdjAI.jpg" )
	//err = robot.SendMarkdownMessage("皮皮虾", "导出：\n```bash\nmongodump -u test -p test --host 127.0.0.1 --port 37017 --out /home/mongodb_bak\n```\n导入：\n```bash\nmongorestore -u test -p test --host 127.0.0.1 --port 57017 --dir /home/mongodb_bak\n```\n启动：\n```bash\nmongod -f /etc/mongod5.conf --master --journal --fork --bind_ip 0.0.0.0\n```",[]string{},false )
	fmt.Println(err)
}

func uploadExcel(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		return
	}
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		global.GLog.Error("open excel error:" + err.Error())
		return
	}
	//fmt.Println(xlsx)
	//解析excel的数据
	rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))

	for _, row := range rows {
		fmt.Println(len(row))
		//fmt.Println(row)
	}

	p, b := xlsx.GetPicture("Sheet1", "A5")
	fmt.Println(p)
	fmt.Println(b)
	if err := ioutil.WriteFile(p, b, 0644); err != nil {
		fmt.Println(err)
	}
}
