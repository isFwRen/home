<template>
	<div class="businessTrends">
		<div class="menu">
			<z-btn-toggle
				class="pa-0 ma-0"
				formId="sexuals"
				formKey="sexual"
				color="primary"
				:group="true"
				mandatory
				:options="cells.menuList"
				@click="changeType"
			></z-btn-toggle>
		</div>

		<div class="treeBox">
			<HomeMinTitle text=" ">
				<a href="../report/itemReport/day-report">
					<div class="more" @click="showMore('more')">
						<span>详情>></span>
					</div></a
				>
			</HomeMinTitle>
			<div class="mb-4"></div>
			<div class="statistics">
				<TreeView :inputData="treeData" />
			</div>
		</div>
		<div class="tableBox">
			<HomeMinTitle text="今日排行榜" />
			<div class="mb-4"></div>
			<div class="table">
				<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
					<template v-for="item in cells.headers">
						<vxe-column :field="item.value" :title="item.text" :width="item.width">
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
	</div>
</template>

<script>
import TreeView from "../../homepage/Sundries/TreeView.vue";
import HomeMinTitle from "../../homepage/HomeMinTitle.vue";
import cells from "./cells";
import TableMixins from "@/mixins/TableMixins";
export default {
	mixins: [TableMixins],
	components: { TreeView, HomeMinTitle },
	data() {
		return {
			cells,
			rankingType: 0,
			treeData: {
				X: [],
				Y: [],
				data: []
			},
			desserts: []
		};
	},
	methods: {
		async search() {
			const result = await this.$store.dispatch("GET_HOME_BUDINESS_RANKING_LIST", {
				rankingType: this.rankingType
			});
			if (result.code == 200) {
				this.treeData.X = [];
				this.treeData.data = [];
				result.data.list.sort(function (a, b) {
					return b.billCount - a.billCount;
				});
				this.desserts = result.data.list.map((e, i) => {
					this.treeData.X.push(e.proCode);
					this.treeData.data.push(e.billCount);
					return { ...e, rank: "TOP" + (i + 1) };
				});
				this.treeData.Y = this.getY(this.desserts[0].billCount);
			}
		},
		getY(bigger) {
			let max = 1,
				arr = [];
			while (max < bigger) {
				max *= 5;
			}
			for (let i = 0; i <= max; i += max / 5) {
				arr.push(i.toFixed(1));
			}
			return arr;
		},
		changeType(e) {
			this.rankingType = e;
			this.search();
		}
	},
	mounted() {
		this.search();
	}
};
</script>

<style lang="scss" scoped>
.businessTrends {
	padding: 0 16px;
	display: flex;
	justify-content: space-between;
	.statistics,
	.tableBox {
		display: inline-block;
	}
	.statistics {
		width: 60vw;
		min-width: 700px;
		height: 300px;
		border: 1px solid #999;
	}
	.tableBox {
		width: 35vw;
		min-width: 300px;
		margin-left: 20px;
	}
	.more {
		display: inline-block;
		position: absolute;
		color: rgb(91, 107, 115);
		font-size: 0.4em;
		right: 0;
		cursor: pointer;
		font-weight: 400;
	}
	.menu {
		position: absolute;
		top: 0;
		right: 0;
	}
	::v-deep .minTitle {
		font-size: 1em;
		color: red;
	}
}
</style>
