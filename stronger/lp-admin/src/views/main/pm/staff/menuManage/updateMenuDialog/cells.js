const MENU_TYPE = 0;
const BUTTON_TYPE = 1;

const menuTypeOptions = [
	{
		label: "菜单",
		value: MENU_TYPE,
		icon: "mdi-menu"
	},

	{
		label: "按钮",
		value: BUTTON_TYPE,
		icon: "mdi-gesture-tap-button"
	}
];

const fields = [
	{
		inputType: "text",
		formKey: "parentName",
		cols: 12,
		colsClass: "pb-0",
		disabled: true,
		label: "上级",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true
	},

	{
		inputType: "text",
		formKey: "title",
		cols: 12,
		colsClass: "py-0",
		label: "标题",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "标题不能为空." }]
	},

	{
		inputType: "text",
		formKey: "action",
		cols: 12,
		colsClass: "py-0",
		disabled: true,
		label: "请求类型",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: false,
		validation: [{ rule: "required", message: "请求类型不能为空." }]
	},

	{
		inputType: "text",
		formKey: "api",
		cols: 12,
		colsClass: "py-0",
		disabled: true,
		label: "接口",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: false,
		validation: [{ rule: "required", message: "接口不能为空." }]
	},

	{
		inputType: "text",
		formKey: "apiId",
		cols: 12,
		colsClass: "py-0",
		disabled: true,
		label: "接口ID",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: false,
		validation: [{ rule: "required", message: "接口ID不能为空." }]
	},

	{
		inputType: "text",
		formKey: "component",
		cols: 12,
		colsClass: "py-0",
		label: "组件",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "组件不能为空." }]
	},

	{
		inputType: "text",
		formKey: "icon",
		cols: 12,
		colsClass: "py-0",
		label: "图标",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "图标不能为空." }]
	},

	{
		inputType: "text",
		formKey: "name",
		cols: 12,
		colsClass: "py-0",
		label: "名字",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "名字不能为空." }]
	},

	{
		inputType: "text",
		formKey: "path",
		cols: 12,
		colsClass: "py-0",
		label: "路径",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "路径不能为空." }]
	},

	{
		inputType: "text",
		formKey: "sort",
		cols: 12,
		colsClass: "py-0",
		label: "排序",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "排序不能为空." }]
	},

	{
		inputType: "switch",
		formKey: "isEnable",
		cols: 12,
		colsClass: "py-0",
		label: "是否启用",
		options: [],
		prependOuter: "*",
		prependOuterClass: "error--text",
		show: true,
		validation: [{ rule: "required", message: "请选择是否启用." }],
		defaultValue: true
	}

	// {
	//   inputType: 'switch',
	//   formKey: 'isFrame',
	//   cols: 12,
	//   colsClass: 'py-0',
	//   label: '是否弹窗',
	//   options: [],
	//   prependOuter: '*',
	//   prependOuterClass: 'error--text',
	//   show: true,
	//   validation: [
	//     { rule: 'required', message: '请选择是否弹窗.' }
	//   ],
	//   defaultValue: false
	// },
];

export { MENU_TYPE, BUTTON_TYPE, menuTypeOptions, fields };

export default {
	MENU_TYPE,
	BUTTON_TYPE,
	menuTypeOptions,
	fields
};
