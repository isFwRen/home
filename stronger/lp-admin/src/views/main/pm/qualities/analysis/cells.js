import moment from "moment";

const typeOptions = [
	{
		value: "1",
		label: "按项目"
	},
	{
		value: "2",
		label: "按字段"
	},
	{
		value: "3",
		label: "按人员"
	}
];

const DEFAULT_DATE = (() => {
	const date = new Date();

	const [year, month] = [date.getFullYear(), date.getMonth() + 1];

	return [moment(`${year}-${month}-1`).format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];
})();

export default {
	typeOptions,
	DEFAULT_DATE
};
