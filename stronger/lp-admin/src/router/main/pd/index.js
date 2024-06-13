const PDRoutes = {
	path: "PD",
	name: "PD",
	meta: {
		key: "PD",
		title: "项目开发"
	},
	component: () => import("@/views/main/pd"),
	children: [
		{
			path: "/main/PD",
			redirect: "config"
		},

		{
			path: "config",
			name: "Config",
			meta: {
				key: "config",
				pKey: "PD",
				realm: "config",
				title: "配置管理"
			},
			component: () => import("@/views/main/pd/config"),
			children: [
				{
					path: "template",
					name: "ConfigTemplate",
					meta: {
						key: "config-template",
						pKey: "PD",
						realm: "config",
						path: "template",
						title: "修改模板"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/editTemplate")
				},

				{
					path: "field",
					name: "ConfigField",
					meta: {
						key: "config-field",
						pKey: "PD",
						realm: "config",
						path: "field",
						title: "字段配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configField")
				},

				{
					path: "export",
					name: "ConfigExport",
					meta: {
						key: "export",
						pKey: "PD",
						realm: "config",
						path: "export",
						title: "导出配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configExport")
				},

				{
					path: "audit",
					name: "ConfigAudit",
					meta: {
						key: "config-audit",
						pKey: "PD",
						realm: "config",
						path: "audit",
						title: "审核配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configAudit")
				},

				{
					path: "quality",
					name: "ConfigQuality",
					meta: {
						key: "config-quality",
						pKey: "PD",
						realm: "config",
						path: "quality",
						title: "质检配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configQuality")
				},

				{
					path: "prescription",
					name: "ConfigPrescription",
					meta: {
						key: "config-prescription",
						pKey: "PD",
						realm: "config",
						path: "prescription",
						title: "时效配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configPrescription")
				},

				{
					path: "path",
					name: "ConfigPath",
					meta: {
						key: "config-path",
						pKey: "PD",
						realm: "config",
						path: "path",
						title: "路径配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configPath")
				},


				{
					path: "monitor",
					name: "ConfigMonitor",
					meta: {
						key: "config-monitor",
						pKey: "PD",
						realm: "config",
						path: "monitor",
						title: "监控配置"
					},
					component: () => import("@/views/main/pd/config/itemConfigs/configMonitor")
				}
			]
		},

		{
			path: "constant",
			name: "Constant",
			meta: {
				key: "constant",
				pKey: "PD",
				realm: "constant",
				title: "常量管理"
			},
			component: () => import("@/views/main/pd/constant/Constant.vue")
		},
		{
			path: "constantClient",
			name: "ConstantClient",
			meta: {
				key: "constantClient",
				pKey: "PD",
				realm: "constantClient",
				title: "常量管理"
			},
			component: () => import("@/views/main/pd/constant/ConstantClient.vue")
		}
	]
};

export default PDRoutes;
