package unitFunc

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func OcrZxing(filePath string) (error, string) {
	// fi, err := os.Open("D:/qqq/A9.png")
	fi, err := os.Open(filePath)
	if err != nil {
		fmt.Println("-------Open----------", filePath, err.Error())
		return err, ""
	}
	defer fi.Close()
	// 解析二维码
	img, aa, err := image.Decode(fi)
	if err != nil {
		fmt.Println("-------Decode----------", filePath, aa, err.Error())
		return err, ""
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		fmt.Println("-------NewBinaryBitmapFromImage----------", filePath, err.Error())
		return err, ""
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if result == nil || err != nil {
		fmt.Println("-------Decode----------", filePath, err.Error())
		return err, ""
	}
	fmt.Println("-------Decode----------", filePath, result.String())
	return nil, result.String()
	// base64解码
	// infoByte, err := base64.StdEncoding.DecodeString(result.String())
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(string(infoByte))
}
