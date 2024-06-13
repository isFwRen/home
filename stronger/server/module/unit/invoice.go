package unitFunc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

type Data struct {
	Del     string `json:"del" form:"del"`
	IsRed   bool   `json:"isRed" form:"isRed"`
	IsPrint bool   `json:"isPrint" form:"isPrint"`
}
type Resp struct {
	Success     bool   `json:"success" form:"success"`
	Code        int    `json:"code" form:"code"`
	Data        Data   `json:"data" form:"data"`
	Message     string `json:"message" form:"message"`
	Description string `json:"description" form:"description"`
}

func Invoice(path string, bodyData map[string]interface{}) (error, Resp) {
	client := &http.Client{
		Timeout: 100 * time.Second,
	}
	url := "https://open.cs.zbj.com"
	// path := "/v2/invoice/query"
	// path := "/v2/eInvoice/query"
	// bodyData := make(map[string]interface{})
	// bodyData["fpdm"] = "65060222"
	// bodyData["fphm"] = "0003485353"
	// bodyData["kprq"] = "20230415"
	// bodyData["checkCode"] = "7436ad"
	// bodyData["money"] = "8345.36"
	// bodyData["noTaxAmount"] = "26263.00"

	// fmt.Println("--------path----------", path)

	// fmt.Println("--------bodyData----------", bodyData)

	byteBody, err := json.Marshal(bodyData)
	if err != nil {
		// fmt.Println("--------byteBody----------", err.Error())
		return err, Resp{}
	}

	req, err := http.NewRequest("POST", url+path, bytes.NewReader(byteBody))
	if err != nil {
		// fmt.Println("--------NewRequest----------", err.Error())
		return err, Resp{}
	}
	key := "360BBC83DE344DF499B8C05E411E9F1C"
	appSecret := "655886c62458220ebc8a3408d815c38d"
	id, _ := uuid.NewV4()
	timestamp := time.Now().Unix()
	// fmt.Println("----------------------------------", id.String(), strconv.FormatInt(timestamp, 10))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("X-CS-Authorization", "HMAC-SHA256")
	req.Header.Set("X-CS-Key", key)
	req.Header.Set("X-CS-Nonce", id.String())
	req.Header.Set("X-CS-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-CS-Version", "v2")

	signature := GetSignature(key, id.String(), strconv.FormatInt(timestamp, 10), appSecret)
	// fmt.Println("--------signature----------", signature)
	req.Header.Set("X-CS-Signature", signature)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("--------client.Do----------", err.Error())
		fmt.Println("--------body----------", resp)
		return err, Resp{}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("--------body111----------", string(body))
	if err != nil {
		fmt.Println("--------ioutil.ReadAll----------", err.Error())
		return err, Resp{}
	}

	// fmt.Println(string(body))

	respData := GetResValue(string(body), "")
	if !respData.Success {
		mes := "查验失败：" + strconv.Itoa(respData.Code) + "," + respData.Message + "," + respData.Description
		return errors.New(mes), respData
	}
	return nil, respData
}

func GetSignature(key, nonce, timestamp, appSecret string) string {
	value := "POST|X-CS-Authorization=HMAC-SHA256|X-CS-Key=" + key + "|X-CS-Nonce=" + nonce + "|X-CS-Timestamp=" + timestamp + "|X-CS-Version=v2"
	// fmt.Println("--------value----------", value)
	// cryto
	return HmacSha256(appSecret, value)
}

func HmacSha256(appSecret string, data string) string {
	mac := hmac.New(sha256.New, []byte(appSecret))
	_, _ = mac.Write([]byte(data))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func GetResValue(data, typeNum string) Resp {
	var resp Resp
	json.Unmarshal([]byte(data), &resp)
	// fmt.Println("--------value----------", resp.Data.Del, resp.Data.IsRed, resp.Data.IsPrint)
	return resp
}
