import caseList from "./case";
import staff from "./staff";
import quality from "./quality";
import teaching from "./teaching";

const prefixPath = "/main/PM/";

const CaseRoutes = caseList(prefixPath);
const StaffRoutes = staff(prefixPath);
const QualityRoutes = quality(prefixPath);
const TeachingRoutes = teaching(prefixPath);

const PMRoutes = {
	path: "PM",
	name: "PM",
	meta: {
		key: "PM",
		title: "项目管理"
	},
	component: () => import("@/views/main/pm"),

	children: [
		{
			path: prefixPath,
			redirect: "case"
		},

		CaseRoutes,

		{
			path: "task",
			name: "Task",
			meta: {
				key: "task",
				realm: "task",
				pKey: "PM",
				title: "任务管理"
			},
			component: () => import("@/views/main/pm/task")
		},

		{
			path: "prescription",
			name: "Prescription",
			meta: {
				key: "prescription",
				realm: "prescription",
				pKey: "PM",
				title: "时效管理"
			},
			component: () => import("@/views/main/pm/prescription")
		},

		{
			path: "qualities",
			name: "Qualities",
			meta: {
				key: "qualities",
				realm: "qualities",
				pKey: "PM",
				title: "质量管理"
			},
			component: () => import("@/views/main/pm/qualities"),

			children: [
				{
					path: `${prefixPath}qualities`,
					redirect: "manage"
				},

				{
					path: "manage",
					name: "Manage",
					meta: {
						key: "manage",
						pKey: "PM",
						path: "manage",
						realm: "qualities",
						title: "质量管理"
					},
					component: () => import("@/views/main/pm/qualities/manage")
				},

				{
					path: "analysis",
					name: "Analysis",
					meta: {
						key: "analysis",
						pKey: "PM",
						path: "analysis",
						realm: "qualities",
						title: "质量分析"
					},
					component: () => import("@/views/main/pm/qualities/analysis")
				}
			]
		},

		QualityRoutes,

		StaffRoutes,

		{
			path: "quality",
			name: "Quality",
			meta: {
				key: "quality",
				realm: "quality",
				pKey: "PM",
				title: "质检管理"
			},
			component: () => import("@/views/main/pm/quality"),

			children: [
				{
					path: `${prefixPath}quality`,
					redirect: "sampling-data"
				}
			]
		},

		TeachingRoutes,

		{
			path: "notice",
			name: "Notice",
			meta: {
				key: "notice",
				realm: "notice",
				pKey: "PM",
				title: "公告管理"
			},
			component: () => import("@/views/main/pm/notice")
		}
	]
};

export default PMRoutes;
