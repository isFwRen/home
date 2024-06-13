export default prefixPath => ({
	path: "teaching",
	name: "Teaching",
	meta: {
		key: "teaching",
		realm: "teaching",
		pKey: "PM",
		title: "教学管理"
	},
	component: () => import("@/views/main/pm/teaching"),

	children: [
		{
			path: `${prefixPath}teaching`,
			redirect: "business-rules"
		},

		{
			path: "business-rules",
			name: "BusinessRules",
			meta: {
				key: "business-rules",
				pKey: "PM",
				path: "business-rules",
				realm: "teaching",
				title: "业务规则"
			},
			component: () => import("@/views/main/pm/teaching/businessRules")
		},

		{
			path: "expense-template",
			name: "ExpenseTemplate",
			meta: {
				key: "expense-template",
				pKey: "PM",
				path: "expense-template",
				realm: "teaching",
				title: "报销单模板"
			},
			component: () => import("@/views/main/pm/teaching/expenseTemplate")
		},

		{
			path: "field-rules",
			name: "FieldRules",
			meta: {
				key: "field-rules",
				pKey: "PM",
				path: "field-rules",
				realm: "teaching",
				title: "字段规则"
			},
			component: () => import("@/views/main/pm/teaching/fieldRules")
		},

		{
			path: "teaching-video",
			name: "TeachingVideo",
			meta: {
				key: "teaching-video",
				pKey: "PM",
				path: "teaching-video",
				realm: "teaching",
				title: "教学视频"
			},
			component: () => import("@/views/main/pm/teaching/teachingVideo")
		}
	]
});
