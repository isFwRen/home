<template>
	<v-app>
		<div class="lp-app">
			<router-view></router-view>
		</div>
	</v-app>
</template>

<script>
export default {
	name: "App"
};
</script>

<style lang="scss">
.vxe-table--body-wrapper {
	.vxe-cell {
		font-size: 14px;
	}
}

.lp-app {
	width: 100%;
	height: 100%;
}
</style>
