const EntryRoutes = {
	path: "entry",
	name: "Entry",
	meta: {
		key: "entry",
		title: "录入管理"
	},
	component: () => import("@/views/main/entry"),
	children: [
		{
			path: "/main/entry",
			redirect: "channel"
		},

		{
			path: "channel",
			name: "Channel",
			meta: {
				key: "channel",
				pKey: "entry",
				realm: "channel",
				title: "录入通道"
			},
			component: () => import("@/views/main/entry/channel"),
			children: [
				{
					path: "op0",
					name: "Op0",
					meta: {
						key: "op0",
						pKey: "entry",
						realm: "channel",
						title: "初审",
						path: "op0"
					},
					component: () => import("@/views/main/entry/channel/taskDialog/op0")
				},

				{
					path: "op1",
					name: "Op1",
					meta: {
						key: "op1",
						pKey: "entry",
						realm: "channel",
						title: "一码",
						path: "op1"
					},
					component: () => import("@/views/main/entry/channel/taskDialog/op1")
				},

				{
					path: "op2",
					name: "Op2",
					meta: {
						key: "op2",
						pKey: "entry",
						realm: "channel",
						title: "二码",
						path: "op2"
					},
					component: () => import("@/views/main/entry/channel/taskDialog/op2")
				},

				{
					path: "opq",
					name: "Opq",
					meta: {
						key: "opq",
						pKey: "entry",
						realm: "channel",
						title: "问题件",
						path: "opq"
					},
					component: () => import("@/views/main/entry/channel/taskDialog/opq")
				}
			]
		}
	]
};

export default EntryRoutes;
