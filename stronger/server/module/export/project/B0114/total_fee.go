/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/15 10:50 上午
 */

package B0114

import (
	"fmt"
	"server/global"
	model4 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func TotalFee(formatQuality map[int][]model4.SysQuality, xml string) (err error, formatQuality2 map[int][]model4.SysQuality, obj model.BillAmount) {
	constMap := constDealSee("B0114")

	for ii, items := range formatQuality {
		for jj, item := range items {
			if item.FieldCode == "fc035" {
				for kk, _ := range item.XmlNodeVal {
					fc111 := getValue(items, kk, "fc111")
					fc112 := getValue(items, kk, "fc112")
					t1, err := time.Parse("2006/01/02", fc111)
					if err != nil {
						continue
					}
					t2, err := time.Parse("2006/01/02", fc112)
					if err != nil {
						continue
					}
					// 计算两个时间点的相差天数
					diff := t2.Sub(t1)
					fmt.Println("-----------diff--------------", fc111, fc112, diff.Hours())
					aa := diff.Hours() / 24
					// days := diff / (24 * time.Hour)
					// value, _ = strconv.ParseFloat(aa, 64)
					formatQuality[ii][jj].XmlNodeVal[kk] = fmt.Sprintf("%v", aa)
					if formatQuality[ii][jj].XmlNodeVal[kk] == "0" {
						formatQuality[ii][jj].XmlNodeVal[kk] = "1"
					}
				}

			}
			if item.FieldCode == "fc105" {
				for kk, XmlNodeVal := range item.XmlNodeVal {
					if XmlNodeVal == "1" {
						formatQuality[ii][jj].XmlNodeVal[kk] = "门诊"
					} else if XmlNodeVal == "2" {
						formatQuality[ii][jj].XmlNodeVal[kk] = "住院"
					}
				}
			}
			if item.FieldCode == "fc106" {
				for kk, XmlNodeVal := range item.XmlNodeVal {
					v, ok := constMap["yiYuanMingChengMap"][XmlNodeVal]
					if ok {
						formatQuality[ii][jj].XmlNodeVal[kk] = v
					}
				}
			}
		}
	}

	//账单金额总和thirdInSIAmount
	//扣费金额总和paymentFee
	//报销金额总和fee
	billAmount := utils.GetNodeData(xml, "billAmount")
	paymentFee := utils.GetNodeData(xml, "otherAbateAmount")
	selfPayAmount := utils.GetNodeData(xml, "selfPayAmount")
	paymentFee = append(paymentFee, selfPayAmount...)
	fee := utils.GetNodeData(xml, "sociSecuAbateAmount")
	AdjustmentAmount := utils.GetNodeData(xml, "ratifyAmount")

	obj = model.BillAmount{
		BillAmount:       totalStr(billAmount),
		FeeAmount:        totalStr(fee),
		DeductionAmount:  totalStr(paymentFee),
		AdjustmentAmount: totalStr(AdjustmentAmount),
	}

	return nil, formatQuality, obj
}

func getValue(items []model4.SysQuality, kk int, code string) string {
	for _, item := range items {
		if item.FieldCode == code && len(item.XmlNodeVal) >= kk {
			return item.XmlNodeVal[kk]
		}
	}
	return ""
}

func totalStr(str []string) decimal.Decimal {
	total := decimal.NewFromFloat(0)
	for _, s := range str {
		fromString, err := decimal.NewFromString(s)
		if err != nil {
			return total
		}
		total = decimal.Sum(total, fromString)
	}
	return total
}

func constDealSee(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"yiYuanMingChengMap", "B0114_华夏理赔_医院名称表", "1", "2"},
	}
	for _, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		tempNumMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])
				//else {
				tempMap[strings.TrimSpace(arr[k])] = arr[v]
				//}
			}
		}
		constObj[item[0]] = tempMap
		constObj[item[0]+"Num"] = tempNumMap
	}

	return constObj
}
