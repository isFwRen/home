const fields = [
	{
		formKey: "upMonth",
		inputType: "date",
		label: "月份",
		prependOuter: "*",
		validation: [{ rule: "required", message: "月份不能为空" }],
		defaultValue: undefined
	},

	{
		formKey: "payroll",
		inputType: "fileInput",
		label: "工资表",
		prependOuter: "*",
		validation: [{ rule: "required", message: "工资表不能为空" }],
		defaultValue: undefined
	}
];

export default {
	fields
};
