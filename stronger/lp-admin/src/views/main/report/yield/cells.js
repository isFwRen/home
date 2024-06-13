const prefixPath = "/main/report/yield/";

const tabsOptions = [
	{
		label: "人员产量统计",
		key: "staff-yield",
		path: `${prefixPath}staff-yield`
	},

	{
		label: "OCR产量统计",
		key: "ocr-yield",
		path: `${prefixPath}ocr-yield`
	}
];

export default {
	tabsOptions
};
