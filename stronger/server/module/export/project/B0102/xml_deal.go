package B0102

import (
	"server/global"
)

func XmlDeal(o interface{}, xmlValue string) (err error, newXmlValue string) {
	global.GLog.Info("------------------B0102:::XmlDeal-----------------------")
	// return nil, xmlValue
	// obj := o.(FormatObj)
	//fields := obj.Fields

	newXmlValue = xmlValue
	return err, newXmlValue
}
