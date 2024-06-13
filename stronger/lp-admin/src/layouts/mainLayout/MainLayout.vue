<template>
	<div class="main-layout">
		<app-bar @drawer="$refs.drawer.onToggle()"></app-bar>

		<!-- 侧边栏 BEGIN -->
		<lp-drawer ref="drawer" app :menus="menus">
			<div class="drawer-top" slot="top"></div>
		</lp-drawer>
		<!-- 侧边栏 END -->

		<v-main>
			<div class="pa-4 z-main">
				<slot v-if="responded"></slot>
			</div>
		</v-main>

		<lp-spinners :overlay="overlay">
			<div slot="tips" class="tips">
				<h3 class="warning--text fw-bold">
					数据同步中，请勿刷新页面，否则将导致数据同步失败!
				</h3>
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
import { mapGetters } from "vuex";
import localForage from "localforage";
import { tools } from "vue-rocket";
import menus from "./menus";
import { LP2, syncLP2 } from "@/api/syncPouchDB";

export default {
	name: "MainLayout",

	data() {
		return {
			menus,
			responded: false,
			overlay: true,
			tips: null,
			newImgUrl: "",
			syncPercent: 0
		};
	},

	created() {
		this.$store.commit("UPDATE_AUTH");
	},

	computed: {
		...mapGetters(["auth", "constants"])
	},

	watch: {
		auth: {
			handler({ menus }) {
				if (tools.isYummy(menus)) {
					this.setMenus(menus);
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

				// if(+query?.sync === 1) {
				//   this.$nextTick(() => {
				//     this.overlay = true
				//   })

				//   this.syncConstants()
				// }

				if (name === "Case" && !tools.isYummy(query)) {
					const forage = await localForage.getItem(LP2);

					window["constantsDB"] = forage;

					console.log(window["constantsDB"]);
				}

				this.overlay = false;
			},
			immediate: true
		}
	},

	methods: {
		// 将常量挂载到 window 对象
		async syncConstants() {
			const result = await syncLP2();

			if (result.code === 200) {
				const forage = await localForage.getItem(LP2);
				const assignData = { ...forage, ...result.data };
				await localForage.setItem(LP2, assignData);

				window["constantsDB"] = assignData;
			}

			this.overlay = false;

			this.$router.replace({ name: "Case" });
		},

		// 设置侧边栏
		setMenus(menus) {
			const flatAuthMenus = [];

			menus?.map(menu => {
				flatAuthMenus.push({ link: menu.path });

				if (tools.isYummy(menu.children)) {
					menu.children.map(child => {
						flatAuthMenus.push({ link: child.path });
					});
				}
			});

			this.menus.map(menu => {
				const record = tools.find(flatAuthMenus, { link: menu.link });

				if (tools.isYummy(record)) {
					menu.visible = true;
				} else {
					menu.visible = false;
				}

				menu.children?.map(item => {
					const record = tools.find(flatAuthMenus, { link: item.link });

					if (tools.isYummy(record)) {
						item.visible = true;
					} else {
						item.visible = false;
					}
				});
			});
		}
	},

	components: {
		"app-bar": () => import("../appBar"),
		"lp-drawer": () => import("@/components/lp-drawer"),
		"lp-spinners": () => import("@/components/lp-spinners")
	}
};
</script>

<style lang="scss" scoped>
.pop_avatar_box {
	border: 1px solid #0f172a1a;
	z-index: 999;
	width: 145px;
	height: 150px;
	position: absolute;
	bottom: -170px;
	margin-left: -50%;
	border-radius: 10px;
	color: #334155;
	font-weight: 600;
	font-size: 14px;
	background-color: #ffffff;
	box-shadow: 0 10px 15px -3px #0000001a, 0 4px 6px -4px #0000001a;
	.pop_avatar_title {
		padding: 10px 0 5px 0;
	}
}
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
