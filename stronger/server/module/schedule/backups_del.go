package schedule

import (
	"fmt"
	"io/ioutil"
	"server/global"
	"server/module/download/project"
	"server/utils"
	"strings"
	"time"
)

func BackupsAndDel() {
	path := global.GConfig.LocalUpload.FilePath
	// path := "D:/workspace/test"
	fmt.Println("-----path-----", path)
	err, items := ReadDir(path, "/", []string{})
	if err != nil {
		fmt.Println("-----err-----", err)
	}
	fmt.Println("-----items-----", items)
	for _, item := range items {
		fpath := path + "/." + item
		cmd := `rsync -avz --exclude={*.tif,*.png,*.jpg,*.pdf,*.dat,*.xml,*.txt} -R -r ` + fpath + ` root@192.168.202.14:/home/41`
		fmt.Println("-----正在备份-----", cmd)
		err, _, _ := project.ShellOut(cmd)
		if err != nil {
			fmt.Println("------备份失败--------:", cmd, err)
			continue
		}
		fpath = path + item
		cmd = `rm -rf ` + fpath
		err, _, _ = project.ShellOut(cmd)
		if err != nil {
			fmt.Println("-----删除失败--------:", cmd, err)
			continue
		}
	}

}

func ReadDir(path, item string, days []string) (error, []string) {
	// mark := "/"
	a := time.Now().AddDate(0, -8, 0)
	fmt.Println("-----a-----", a.Format("2006/01-02"))

	fpath := path + item
	fileInfoList, err := ioutil.ReadDir(fpath)

	if err != nil {
		fmt.Println("---err---", err)
		return err, days
	}

	// fmt.Println(len(fileInfoList))

	for i := range fileInfoList {
		name := fileInfoList[i].Name()
		// fmt.Println("=============name:", item+name) //打印当前文件或目录下的文件或目录名
		if strings.Index(item, "download") == -1 {
			if name == "B0108" {
				// if utils.RegIsMatch(`^B0\d{3}$`, name) {
				// files = append(files, item+fileInfoList[i].Name()+mark+"download"+mark)
				_, days = ReadDir(path, item+fileInfoList[i].Name()+"/download/", days)
			}
		} else {
			file := item + fileInfoList[i].Name() + "/"
			// fmt.Println("-------------file--------:", file)
			// fmt.Println("-------------222--------:", utils.RegIsMatch(`\/\d{4}\/\d{2}-\d{2}\/`, file))
			if utils.RegIsMatch(`\/\d{4}\/\d{2}-\d{2}\/`, file) {
				day := utils.RegMatchAll(file, `\d{4}\/\d{2}-\d{2}`)[0]
				b, _ := time.Parse("2006/01-02", day)
				if b.Before(a) {
					days = append(days, file)
				}
				// fmt.Println("-------------days--------:", days)
			} else if utils.RegIsMatch(`\/\d{4}\/\d{2}\/\d{2}\/`, file) {
				day := utils.RegMatchAll(file, `\d{4}\/\d{2}\/\d{2}`)[0]
				b, _ := time.Parse("2006/01/02", day)
				if b.Before(a) {
					days = append(days, file)
				}
			} else {
				_, days = ReadDir(path, file, days)
			}
			// files = append(files, file)
		}

	}
	// }

	return err, days
}

// func main() {
// 	BackupsAndDel()
// }
