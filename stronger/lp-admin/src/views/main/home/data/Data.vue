<template>
	<div class="data">
		<div class="memu">
			<z-btn-toggle
				class="pa-0 ma-0"
				formId="sexuals"
				formKey="sexual"
				color="primary"
				:group="true"
				mandatory
				:options="cells.menuList"
				@change="onChange($event)"
			></z-btn-toggle>
		</div>
		<WaitMakeData v-if="nowPage == 'waitMakeData'" />
		<BuinessTrends v-if="nowPage == 'businessTrends'" />
		<WordTrends v-if="nowPage == 'WordTrends'" />
	</div>
</template>

<script>
import cells from "./cells";
import WaitMakeData from "./waitMakeData/WaitMakeData.vue";
import BuinessTrends from "./businessTrends/BuinessTrends.vue";
import WordTrends from "./wordTrends/WordTrends.vue";
export default {
	data() {
		return {
			cells,
			nowPage: null
		};
	},
	methods: {
		onChange(e) {
			this.nowPage = e;
		}
	},
	components: { WaitMakeData, BuinessTrends, WordTrends }
};
</script>

<style scoped lang="scss">
.data {
	position: relative;
}
::v-deep .z-btn-toggle .z-input--btn-toggle-group {
	margin: 0;
	padding: 0;
}
.memu {
	padding: 0 13px;
}
</style>
