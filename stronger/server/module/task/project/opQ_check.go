package project

import (
	"fmt"
	"server/module/load/model"
	"strings"
)

func OpQCheck(block model.ProjectBlock, fields []model.ProjectField, op string) []model.ProjectField {
	// op := block.Stage
	// mOp := strings.Replace(op, "o", "O", -1)
	fmt.Println("----------FieldsInputCheck----------", len(fields))
	for ff, field := range fields {
		// rfield := reflect.ValueOf(&field).Elem()

		op1Value := field.Op1Value
		op2Value := field.Op2Value
		if strings.Index(op1Value, "?") != -1 || strings.Index(op1Value, "？") != -1 || strings.Index(op2Value, "?") != -1 || strings.Index(op2Value, "？") != -1 {
			// rfield.FieldByName(mOp + "Input").Set(reflect.ValueOf("yes"))
			fields[ff].OpqInput = "yes"
		}
		if field.OpqInput == "" {
			// rfield.FieldByName(mOp + "Input").Set(reflect.ValueOf("no"))
			fields[ff].OpqInput = "no"
		} else if block.Op1Code != "0" && block.Op2Code != "0" && field.Op2Input != "no" && field.Op1Input != "no" && op1Value != op2Value {
			// rfield.FieldByName(mOp + "Input").Set(reflect.ValueOf("yes"))
			fields[ff].OpqInput = "yes"
		}
		// fmt.Println("----------rfield----------", rfield.FieldByName(mOp+"Input"), field.OpqInput, fields[ff].OpqInput)
	}
	return fields
}
