<template>
	<div class="staff">
		<v-row>
			<v-col :cols="6">
				<vxe-table :data="list" :size="tableSize">
					<template v-for="item in cells.headers">
						<vxe-column
							:field="item.value"
							:title="item.text"
							:key="item.value"
							:width="item.width"
						>
						</vxe-column>
					</template>
				</vxe-table>
			</v-col>

			<v-col class="z-flex justify-center" :cols="6">
				<div v-if="list.length" ref="chart" class="main"></div>
			</v-col>
		</v-row>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";
import ChartMixins from "@/mixins/ChartMixins";
import cells from "./cells";

export default {
	name: "Staff",
	mixins: [TableMixins, ChartMixins],

	props: {
		list: {
			type: Array,
			requried: true
		}
	},

	data() {
		return {
			cells
		};
	},

	watch: {
		list: {
			handler(list) {
				const items = [];

				list.map(item => {
					items.push({
						value: item.nums,
						name: item.people
					});
				});

				this.$nextTick(() => {
					this.setPieChart("chart", {
						name: "按人员",
						items
					});
				});
			}
		}
	}
};
</script>

<style lang="scss">
.main {
	width: 500px;
	height: 500px;
}
</style>
