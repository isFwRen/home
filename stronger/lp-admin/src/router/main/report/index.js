const prefixPath = "/main/report/";

const ReportRoutes = {
	path: "report",
	name: "Report",
	meta: {
		key: "report",
		title: "报表管理"
	},
	component: () => import("@/views/main/report"),
	children: [
		{
			path: prefixPath,
			redirect: "yield"
		},

		{
			path: "yield",
			name: "yield",
			meta: {
				key: "yield",
				pKey: "report",
				realm: "yield",
				title: "产量统计"
			},
			component: () => import("@/views/main/report/yield"),

			children: [
				{
					path: `${prefixPath}yield`,
					redirect: "staff-yield"
				},

				{
					path: "staff-yield",
					name: "staffYield",
					meta: {
						key: "staff-yield",
						pKey: "report",
						path: "staff-yield",
						realm: "yield",
						title: "人工产量"
					},
					component: () => import("@/views/main/report/yield/staffYield")
				},

				{
					path: "ocr-yield",
					name: "ocrYield",
					meta: {
						key: "ocr-yield",
						pKey: "report",
						path: "ocr-yield",
						realm: "yield",
						title: "OCR产量"
					},
					component: () => import("@/views/main/report/yield/ocrYield")
				}
			]
		},

		{
			path: "error",
			name: "error",
			meta: {
				key: "error",
				pKey: "report",
				realm: "error",
				title: "错误查询"
			},
			component: () => import("@/views/main/report/error"),

			children: [
				{
					path: `${prefixPath}error`,
					redirect: "error-detail"
				},

				{
					path: "error-detail",
					name: "errorDetail",
					meta: {
						key: "error-detail",
						pKey: "report",
						path: "error-detail",
						realm: "error",
						title: "错误明细"
					},
					component: () => import("@/views/main/report/error/errorDetail")
				},

				{
					path: "error-analysis",
					name: "errorAnalysis",
					meta: {
						key: "error-analysis",
						pKey: "report",
						path: "error-analysis",
						realm: "error",
						title: "错误分析"
					},
					component: () => import("@/views/main/report/error/errorAnalysis")
				},
			]
		},

		{
			path: "salary",
			name: "salary",
			meta: {
				key: "salary",
				pKey: "report",
				realm: "salary",
				title: "工资数据"
			},
			component: () => import("@/views/main/report/salary"),

			children: [
				{
					path: `${prefixPath}salary`,
					redirect: "inside-salary"
				},

				{
					path: "inside-salary",
					name: "insideSalary",
					meta: {
						key: "inside-salary",
						pKey: "report",
						path: "inside-salary",
						realm: "salary",
						title: "内部工资"
					},
					component: () => import("@/views/main/report/salary/insideSalary")
				},

				{
					path: "pt-salary",
					name: "ptSalary",
					meta: {
						key: "pt-salary",
						pKey: "report",
						path: "pt-salary",
						realm: "salary",
						title: "PT工资"
					},
					component: () => import("@/views/main/report/salary/ptSalary")
				}
			]
		},

		{
			path: "itemReport",
			name: "itemReport",
			meta: {
				key: "itemReport",
				pKey: "report",
				realm: "itemReport",
				title: "项目报表"
			},
			component: () => import("@/views/main/report/itemReport"),

			children: [
				{
					path: `${prefixPath}itemReport`,
					redirect: "business-detail"
				},

				{
					path: "business-detail",
					name: "businessDetail",
					meta: {
						key: "business-detail",
						pKey: "report",
						path: "business-detail",
						realm: "itemReport",
						title: "业务明细"
					},
					component: () => import("@/views/main/report/itemReport/businessDetail")
				},

				{
					path: "character-statistics",
					name: "characterStatistics",
					meta: {
						key: "character-statistics",
						pKey: "report",
						path: "character-statistics",
						realm: "itemReport",
						title: "字符统计表"
					},
					component: () => import("@/views/main/report/itemReport/characterStatistics")
				},

				{
					path: "day-report",
					name: "dayReport",
					meta: {
						key: "day-report",
						pKey: "report",
						path: "day-report",
						realm: "itemReport",
						title: "日报表"
					},
					component: () => import("@/views/main/report/itemReport/dayReport")
				},

				{
					path: "week-report",
					name: "WeekReport",
					meta: {
						key: "week-report",
						pKey: "report",
						path: "week-report",
						realm: "itemReport",
						title: "周报表"
					},
					component: () => import("@/views/main/report/itemReport/weekReport")
				},

				{
					path: "month-report",
					name: "MonthReport",
					meta: {
						key: "month-report",
						pKey: "report",
						path: "month-report",
						realm: "itemReport",
						title: "月报表"
					},
					component: () => import("@/views/main/report/itemReport/monthReport")
				},

				{
					path: "year-report",
					name: "YearReport",
					meta: {
						key: "year-report",
						pKey: "report",
						path: "year-report",
						realm: "itemReport",
						title: "年报表"
					},
					component: () => import("@/views/main/report/itemReport/yearReport")
				},

				{
					path: "project-report",
					name: "ProjectReport",
					meta: {
						key: "project-report",
						pKey: "report",
						path: "project-report",
						realm: "itemReport",
						title: "项目报表"
					},
					component: () => import("@/views/main/report/itemReport/projectReport")
				},

				{
					path: "project-settlement",
					name: "ProjectSettlement",
					meta: {
						key: "project-settlement",
						pKey: "report",
						path: "project-settlement",
						realm: "itemReport",
						title: "项目结算"
					},
					component: () => import("@/views/main/report/itemReport/projectSettlement")
				}
			]
		},

		{
			path: "specialReport",
			name: "specialReport",
			meta: {
				key: "specialReport",
				pKey: "report",
				realm: "specialReport",
				title: "特殊报表"
			},
			component: () => import("@/views/main/report/specialReport"),
			children: [
				{
					path: `${prefixPath}specialReport`,
					redirect: "income-analysis"
				},

				{
					path: "income-analysis",
					name: "incomeAnalysis",
					meta: {
						key: "income-analysis",
						pKey: "report",
						path: "income-analysis",
						realm: "specialReport",
						title: "来量分析"
					},
					component: () => import("@/views/main/report/specialReport/incomeAnalysis")
				},

				{
					path: "return-analysis",
					name: "returnAnalysis",
					meta: {
						key: "return-analysis",
						pKey: "report",
						path: "return-analysis",
						realm: "specialReport",
						title: "回传分析"
					},
					component: () => import("@/views/main/report/specialReport/returnAnalysis")
				},

				{
					path: "directory-out",
					name: "directoryOut",
					meta: {
						key: "directory-out",
						pKey: "report",
						path: "directory-out",
						realm: "specialReport",
						title: "目录外数据"
					},
					component: () => import("@/views/main/report/specialReport/directoryOut")
				},

				{
					path: "abnormal-part",
					name: "abnormalPart",
					meta: {
						key: "abnormal-part",
						pKey: "report",
						path: "abnormal-part",
						realm: "specialReport",
						title: "异常件数据"
					},
					component: () => import("@/views/main/report/specialReport/abnormalPart")
				},

				{
					path: "institutional-extraction",
					name: "institutionalExtraction",
					meta: {
						key: "institutional-extraction",
						pKey: "report",
						path: "institutional-extraction",
						realm: "specialReport",
						title: "机构抽取"
					},
					component: () => import("@/views/main/report/specialReport/institutionalExtraction")
				},

				{
					path: "identify-statistics",
					name: "identify-statistics",
					meta: {
						key: "identify-statistics",
						pKey: "report",
						path: "identify-statistics",
						realm: "specialReport",
						title: "ocr识别统计"
					},
					component: () => import("@/views/main/report/specialReport/identifyStatistics")
				},
				{
					path: "destruction-report",
					name: "destruction-report",
					meta: {
						key: "destruction-report",
						pKey: "report",
						path: "destruction-report",
						realm: "specialReport",
						title: "销毁报告"
					},
					component: () => import("@/views/main/report/specialReport/destructionReport")
				},
				{
					path: "deduction-details",
					name: "deduction-details",
					meta: {
						key: "deduction-details",
						pKey: "report",
						path: "deduction-details",
						realm: "specialReport",
						title: "扣费明细"
					},
					component: () => import("@/views/main/report/specialReport/deductionDetails")
				},

			]
		}
	]
};

export default ReportRoutes;
