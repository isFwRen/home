<template>
	<div class="loading-wrap" v-if="loading">
		<div class="progress">
			<v-progress-circular
				:rotate="360"
				:size="100"
				:width="8"
				color="primary"
				indeterminate
			></v-progress-circular>
		</div>
	</div>
</template>
<script>
export default {
	name: "LPLoading",
	data() {
		return {
			loading: false
		};
	},
	methods: {
		open() {
			this.loading = true;
		},
		close() {
			this.loading = false;
		}
	}
};
</script>
<style scoped lang="scss">
.loading-wrap {
	position: fixed;
	left: 0;
	bottom: 0;
	top: 0;
	right: 0;
	z-index: 120;
	background-color: rgba(0, 0, 0, 0.2);
}
.progress {
	width: 200px;
	height: 200px;
	position: absolute;
	left: 50%;
	top: 50%;
	transform: translateX(-50%);
	transform: translateY(-50%);
}
</style>
