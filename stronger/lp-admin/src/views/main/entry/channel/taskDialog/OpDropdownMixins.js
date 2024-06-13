import _ from "lodash";

export default {
	methods: {
		// 下拉选中回车
		onDropdownEnterField({ value, index }, field) {
			console.log(value);
			console.log(this.op);

			if (index > -1) {
				let val = value;

				field[`${this.op}Value`] = val;
				field.resultValue = val;
				field.autofocus = true;

				if (this.op === "op0") {
					_.set(this.memoFields, `${field.uniqueId}.value`, val);
					this.fieldsObject = { ...this.fieldsObject };
					console.log(this.fieldsObject);
				} else {
					this.fieldsList = [...this.fieldsList];
				}
			}
		},

		// 下拉按上键
		onDropdownUpField() {}
	}
};
