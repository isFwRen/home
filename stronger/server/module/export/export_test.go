/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/25 10:39 上午
 */

package export

import (
	"fmt"
	"testing"
	"time"
)

func TestValChange(t *testing.T) {
	//f := model.ProjectField{
	//	FinalValue: "123456",
	//}
	b := []string{"1", "2"}
	a := make([]string, 2)
	copy(a, b)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Printf("%p", a)
	fmt.Println()
	fmt.Printf("%p", b)
	fmt.Println(int(time.Now().Month()))

	//valChange(&f, "3=ppp;5=oo")
	//valInsert(&f, "3=ppp;5=oo")
	//fmt.Println(f)
}
