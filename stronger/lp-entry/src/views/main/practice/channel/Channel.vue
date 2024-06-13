<template>
	<div class="practice">
		<vxe-table :data="desserts" :border="tableBorder" :max-height="tableMaxHeight" :size="tableSize" :stripe="tableStripe">
			<template v-for="item in cells.headers">
				<vxe-column
					v-if="item.value === 'options'"
					:field="item.value"
					:title="item.text"
					:key="item.value"
					:width="item.width"
				>
					<template #default="{ row }">
						<v-btn color="primary" depressed rounded @click="enterProject(row)">开始练习</v-btn>
					</template>
				</vxe-column>

				<vxe-column v-else :field="item.value" :title="item.text" :key="item.value" :width="item.width"> </vxe-column>
			</template>
		</vxe-table>

		<task-dialog ref="channel"></task-dialog>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";
import { request } from "@/api/service";
import { sessionStorage, localStorage } from "vue-rocket";

const isIntranet = lpTools.isIntranet();

export default {
	name: "PracticeChannel",
	mixins: [TableMixins],

	data() {
		return {
			cells,
			desserts: []
		};
	},

	async created() {
		this.desserts = [];
		let result = await this.$store.dispatch("PRACTICE_CODE_LIST");
		this.desserts = result.data.filter(el => el.proCode == sessionStorage.get("proCode"));
		// sessionStorage.set("charP", this.desserts[0].character);
		// sessionStorage.set("rateP", this.desserts[0].accuracyRate);
		const proList = await this.$store.dispatch("INPUT_GET_TRANSIT_LIST");
		let row = proList.data.list.filter(el => el.proCode == sessionStorage.get("proCode"));
		const { innerIp, inAppPort, outIp, outAppPort } = row[0];
		const origin = `https://${isIntranet ? `${innerIp}:${inAppPort}` : `${outIp}:${outAppPort}`}`;
		const baseURL = `${origin}/api/`;
		this.$store.commit("UPDATE_PRACTICE", { display: true, topInfo: "true", rowInfo: row, baseURL });
	},

	methods: {
		async enterProject(row) {
			console.log(row, "rr");

			const result = await this.$store.dispatch("INPUT_GET_TASK_CONFIG");

			const code = sessionStorage.get("proCode");
			let mb001 = {};
			result.data.mb001.forEach(el => {
				mb001[el.code] = el.fields.length + 3;
			});
			localStorage.set("mb001", mb001);
			if (!code) return;

			if (result.code === 200) {
				this.$store.commit("UPDATE_PRACTICE", { config: result.data });
				this.$router.replace({ name: "PracticeChannel", query: { proCode: code, op: "opp" } });
			} else {
				this.toasted.warning(result.msg);
			}

			return;
		}
	},

	components: {
		"task-dialog": () => import("../../entry/channel/taskDialog")
	}
};
</script>

<style lang="scss" scoped>
.v-btn {
	height: 25px !important;
}
</style>