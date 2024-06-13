<template>
	<div class="main-layout">
		<app-bar @drawer="$refs.drawer.onToggle()"></app-bar>

		<!-- 侧边栏 BEGIN -->
		<lp-drawer ref="drawer" app :menus="dynamicMenus">
			<div class="drawer-top" slot="top"></div>
		</lp-drawer>
		<!-- 侧边栏 END -->

		<v-main>
			<div class="pa-4 z-main">
				<slot v-if="!overlay"></slot>
			</div>
		</v-main>

		<lp-spinners :overlay="overlay">
			<div slot="tips" class="tips">
				<h3 class="warning--text fw-bold">数据同步中，请勿刷新页面，否则将导致数据同步失败!</h3>
				<p class="mt-1 text-center">
					<template v-if="syncPercent !== '100.00'">
						<span>常量表更新进度：</span>
						<span class="percent">{{ syncPercent }}%</span>
					</template>
				</p>
			</div>
		</lp-spinners>
	</div>
</template>

<script>
import { localStorage, tools } from "vue-rocket";
import { mapGetters } from "vuex";
import localForage from "localforage";
import { menus, menuAuth } from "./cells";
import { LP2, syncLP2 } from "@/api/syncPouchDB";
import { sessionStorage } from "vue-rocket";

export default {
	name: "MainLayout",

	data() {
		return {
			menus,
			dynamicMenus: [],
			responded: false,
			overlay: true,
			tips: null,
			proCode: "",
			syncPercent: 0
		};
	},

	created() {
		this.$store.commit("UPDATE_AUTH");
		this.proCode = sessionStorage.get("proCode");
		const user = localStorage.get("user");
		const isPT = /^p\w*/i.test(user.code);

		if (isPT) {
			this.dynamicMenus = menus
				.map(item => {
					if (menuAuth.includes(item.key)) {
						return item;
					}
				})
				.filter(item => item);
		} else {
			this.dynamicMenus = menus;
		}
	},

	computed: {
		...mapGetters(["auth", "constants"])
	},

	watch: {
		auth: {
			handler({ menus }) {
				if (tools.isYummy(menus)) {
					this.responded = true;
				}
			},
			immediate: true
		},

		constants: {
			handler({ total, count }) {
				this.syncPercent = ((count / total) * 100).toFixed(2);
			},
			deep: true,
			immediate: true
		},

		$route: {
			async handler(route) {
				const { name, query } = route;
				//const proCode = sessionStorage.get("proCode");

				// 从中转平台跳转
				if (query.token && query.userId) {
					this.responded = true;
				} else {
					if (+query?.sync === 1) {
						this.$nextTick(() => {
							this.overlay = true;
						});

						this.syncConstants();
					}
				}

				if (name === "Channel" && !tools.isYummy(query)) {
					const forage = await localForage.getItem(LP2);
					window["constantsDB"] = forage;
					console.log(window["constantsDB"], "constant");
				}

				this.overlay = false;
			},
			immediate: true
		}
	},

	methods: {
		// 将常量挂载到 window 对象
		async syncConstants() {
			// const proCode = sessionStorage.get("proCode");
			// if (proCode) {
			// 	this.overlay = false;
			// 	this.$router.replace({ name: "Channel" });
			// 	return;
			// }
			const result = await syncLP2();

			if (result.code === 200) {
				const forage = await localForage.getItem(LP2);
				const assignData = { ...forage, ...result.data };
				await localForage.setItem(LP2, assignData);

				window["constantsDB"] = assignData;
			}

			this.overlay = false;

			this.$router.replace({ name: "Channel" });
		}
	},

	components: {
		"app-bar": () => import("../appBar"),
		"lp-drawer": () => import("@/components/lp-drawer"),
		"lp-spinners": () => import("@/components/lp-spinners")
	}
};
</script>

<style scoped>
.drawer-top {
	height: 64px;
}

.z-main {
	max-height: calc(100vh - 64px);
	overflow-y: auto;
}

.percent {
	font-size: 20px;
}
</style>
