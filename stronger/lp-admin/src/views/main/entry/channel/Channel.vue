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
							<z-btn color="primary" depressed rounded smaller @click="entryItem(row)"
								>进入</z-btn
							>
						</template>
					</vxe-column>

					<vxe-column
						v-else
						:field="item.value"
						:title="item.text"
						:key="item.value"
						:width="item.width"
					>
					</vxe-column>
				</template>
			</vxe-table>
		</div>

		<task-dialog ref="channel"></task-dialog>
	</div>
</template>

<script>
import { sessionStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";

const isIntranet = lpTools.isIntranet();

export default {
	name: "Entry",
	mixins: [TableMixins],

	data() {
		return {
			cells,
			dispatchList: "GET_CHANNEL_LIST"
		};
	},

	beforeDestroy() {
		this.$store.commit("SET_PROJECT_INFO");
	},

	methods: {
		// 当前项目
		entryItem(row) {
			const { rowInfo } = sessionStorage.get("task");
			const { innerIp, inAppPort, outIp, outAppPort } = rowInfo;
			const baseURL = `https://${
				isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`
			}/api/`;

			this.$store.commit("SET_PROJECT_INFO", { code: row.proCode });
			this.$store.commit("UPDATE_CHANNEL", { rowInfo: row, baseURL });

			this.getTaskConfig();
		},

		// 理赔录入配置
		async getTaskConfig() {
			const result = await this.$store.dispatch("GET_LP_TASK_CONFIG");

			if (result.code === 200) {
				this.$store.commit("UPDATE_CHANNEL", { config: result.data });
				this.$router.replace({ path: "/main/entry/channel", query: { op: -1 } });
			} else {
				this.toasted.warning(result.msg);
			}
		}
	},

	components: {
		"task-dialog": () => import("./taskDialog")
	}
};
</script>
