/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/1 10:17 上午
 */

package export

import (
	"server/module/pro_manager/model"
	"testing"
)

func TestWrongSum(t *testing.T) {
	var a = model.ProCodeAndId{}
	err, _ := WrongSum(a)
	if err != nil {
		return
	}
}
