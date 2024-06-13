<template>
	<div class="dataTable">
		<div class="table">
			<vxe-table
				:data="desserts"
				:border="tableBorder"
				:size="tableSize"
				:sort-config="{ multiple: true, trigger: 'cell' }"
				@sort-change="handleSort"
			>
				<vxe-column type="seq" title="序号" width="60"></vxe-column>

				<template v-for="item in cells.headers">
					<vxe-column :field="item.value" :title="item.text" :width="item.width" sortable>
					</vxe-column>
				</template>
			</vxe-table>
		</div>
		<z-pagination
			:options="pageSizes"
			@page="handlePage"
			:total="pagination.total"
		></z-pagination>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import moment from "moment";
import cells from "./cells";
export default {
	mixins: [TableMixins],
	data() {
		return {
			cells,
			daySelcted: {
				yearstday: moment().subtract(1, "day").format("YYYY-MM-DD"),
				today: moment().format("YYYY-MM-DD")
			}
		};
	},
	props: {
		DataType: {
			type: String,
			default: "today"
		}
	},
	mounted() {
		this.search();
	},

	methods: {
		async search() {
			const result = await this.$store.dispatch("GET_HOME_WAIT_MAKE_DATA_LIST", {
				queryDay: this.queryDay
			});
			if (result.code == 200) {
				this.desserts = result.data.list;
			}
		}
	},
	computed: {
		queryDay() {
			return this.daySelcted[this.DataType];
		}
	},
	watch: {
		queryDay() {
			this.search();
		}
	}
};
</script>

<style lang="scss" scoped></style>
