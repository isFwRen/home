<template>
	<div class="waitMakeData">
		<div class="menu">
			<z-btn-toggle
				class="pa-0 ma-0"
				formId="sexuals"
				formKey="sexual"
				color="primary"
				:group="true"
				mandatory
				:options="cells.menuList"
				@click="clickSelect"
			></z-btn-toggle>
		</div>
		<HomeMinTitle :text="dataType == 'today' ? '今日数据' : '昨日数据'">
			<a :href="dataType == 'today' ? '../PM/task' : '../PM/case'">
				<div class="more" @click="showMore('today')">
					<span>详情>></span>
				</div>
			</a></HomeMinTitle
		>

		<div class="mb-5"></div>
		<DataTable :DataType="dataType" />
	</div>
</template>

<script>
import DataTable from "./dataTable/dataTable.vue";
import cells from "./cells";
import HomeMinTitle from "../../../../../../../lp-entry/src/views/main/home/HomeMinTitle.vue";
export default {
	data() {
		return {
			cells,
			dataType: "today"
		};
	},
	methods: {
		showMore(e) {
			console.log(e);
		},
		clickSelect(e) {
			this.dataType = e;
		}
	},
	components: { DataTable, HomeMinTitle }
};
</script>

<style lang="scss" scoped>
.waitMakeData {
	padding: 0 16px;
	position: relative;
	::v-deep .minTitle {
		font-size: 1em;
	}

	.menu {
		position: absolute;
		top: -62px;
		right: 0;
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
}
</style>
