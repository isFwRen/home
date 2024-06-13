const icons = [
	{
		icon: "mdi-clipboard-list-outline",
		tip: "查看报表",
		value: 1
	},

	{
		icon: "mdi-view-grid-outline",
		tip: "查看影像",
		value: 2
	},

	{
		icon: "mdi-spellcheck",
		tip: "导出校验",
		value: 3
	},

	{
		icon: "mdi-table",
		tip: "查看日志",
		value: 4
	}
];

const fields = [
	{
		formKey: "id",
		inputType: "text",
		hideDetails: true,
		label: "发票属性"
	},

	{
		formKey: "code",
		inputType: "text",
		hideDetails: true,
		label: "账单号 "
	},

	{
		formKey: "money",
		inputType: "text",
		hideDetails: true,
		label: "账单金额"
	},

  {
    formKey: 'keyword',
    inputType: 'text',
    hideDetails: true,
    label: '录入关键字'
  }
]

const invoiceKeys = [
	"baoXiaoDan",
	"invoice",
	"invoiceDaXiang",
	"qingDan",
	"thirdBaoXiaoDan1",
	"thirdBaoXiaoDan2",
	"thirdBaoXiaoDan3",
	"operation",
	"hospitalizationDate",
	"hospitalizationFee",
	"zhenDuan"
];

const invoiceNamesMap = new Map([
	["zhenDuan", "诊断"],
	["baoXiaoDan", "报销单"],
	["invoice", "发票"],
	["invoiceDaXiang", "发票大项"],
	["qingDan", "清单"],
	["thirdBaoXiaoDan1", "第三方报销单1"],
	["thirdBaoXiaoDan2", "第三方报销单2"],
	["thirdBaoXiaoDan3", "第三方报销单3"],
	["operation", "手术"],
	["hospitalizationDate", "住院日期"],
	["hospitalizationFee", "住院诊查费"]
]);

export { invoiceKeys, invoiceNamesMap };

export default {
	icons,
	fields
};