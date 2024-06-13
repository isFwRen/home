<template>
	<div class="channel">
		<div class="table entry-table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:max-height="tableMaxHeight"
				:size="tableSize"
				:stripe="tableStripe"
			>
				<template v-for="item in cells.headers">
					<vxe-column
						v-if="item.value === 'options'"
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
						<template #default="{ row }">
							<z-btn color="primary" depressed rounded smaller @click="enterProject(row)">进入</z-btn>
						</template>
					</vxe-column>

					<vxe-column v-else :field="item.value" :title="item.text" :key="item.value" :width="item.width"> </vxe-column>
				</template>
			</vxe-table>
		</div>

		<task-dialog ref="channel"></task-dialog>

		<lp-spinners :overlay="overlay"></lp-spinners>
	</div>
</template>

<script>
import { sessionStorage, localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import AuthMixins from "@/views/login/AuthMixins";

const isIntranet = lpTools.isIntranet();

export default {
	name: "Entry",
	mixins: [TableMixins, AuthMixins],

	data() {
		return {
			cells,
			dispatchList: "INPUT_GET_LIST",
			overlay: false
		};
	},
	beforeRouteEnter(to, from, next) {
		const { token, userId, secret, isApp, height, width } = to.query;
		if (token && userId) {
			localStorage.delete(["secret", "token", "user"]);
			console.log(token, localStorage, "localStorage");
			localStorage.set("tokens", token);
			localStorage.set("userId", userId);
			localStorage.set("secret", secret);
			localStorage.set("viewport", { height, width });
			sessionStorage.set("isApp", { isApp });
			next(vm => {
				if (token && userId) {
					vm.getRoleSysMenu();
					// 获取user信息
					vm.getUserInfo({ token, userId, secret });
				}
			});
		} else {
			next();
		}
	},
	methods: {
		async getUserInfo(data) {
			const form = {
				"x-token": data.token,
				"x-user-id": data.userId
			};
			const result = await this.$store.dispatch("GET_USER_INFO", form);
			if (result.code === 200) {
				localStorage.set({
					token: result.data.token,
					user: result.data.user
				});
			}
		},
		// 进入项目
		enterProject(row) {
			console.log(row, "rr");
			const proCode = sessionStorage.get("proCode");
			sessionStorage.set("projectCode", row.proCode);

			if (!proCode) {
				this.toasted.warning("请检查接口 api/sys-menu/role/get 是否返回参数 proCode!");
				return;
			}

			const { innerIp, inAppPort, outIp, outAppPort } = row;

			const origin = `https://${isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`}`;

			if (row.proCode === proCode) {
				this.overlay = true;

				const baseURL = `${origin}/api/`;

				this.$store.commit("SET_PROJECT_INFO", { code: row.proCode });
				this.$store.commit("UPDATE_CHANNEL", {
					rowInfo: row,
					baseURL,
					displayPrompt: true,
					displayTop: true,
					displayRight: false
				});

				this.getTaskConfig();

				return;
			}

			window.open(origin);
		},

		// 理赔录入配置
		async getTaskConfig() {
			const result = await this.$store.dispatch("INPUT_GET_TASK_CONFIG");

			const { code } = this.$store.getters.project;
			let mb001 = {};
			result.data.mb001.forEach(el => {
				mb001[el.code] = el.fields.length + 3;
			});
			localStorage.set("mb001", mb001);
			if (!code) return;

			if (result.code === 200) {
				this.$store.commit("UPDATE_CHANNEL", { config: result.data, display: true });
				this.$router.replace({ name: "Channel", query: { proCode: code, op: -1 } });
			} else {
				this.toasted.warning(result.msg);
			}

			this.overlay = false;
		}
	},

	components: {
		"task-dialog": () => import("./taskDialog"),
		"lp-spinners": () => import("@/components/lp-spinners")
	}
};
</script>
