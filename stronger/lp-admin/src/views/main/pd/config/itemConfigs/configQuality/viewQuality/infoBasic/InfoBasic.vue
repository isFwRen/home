<template>
	<div class="info-basic">
		<v-expansion-panels v-model="panel" multiple>
			<v-expansion-panel>
				<v-expansion-panel-header>申请人信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>

			<v-expansion-panel>
				<v-expansion-panel-header>被保人信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>

			<v-expansion-panel>
				<v-expansion-panel-header>受托人信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>

			<v-expansion-panel>
				<v-expansion-panel-header>其他信息</v-expansion-panel-header>
				<v-expansion-panel-content> </v-expansion-panel-content>
			</v-expansion-panel>
		</v-expansion-panels>
	</div>
</template>

<script>
export default {
	name: "InfoBasic",

	data() {
		return {
			panel: [0, 1, 2, 3]
		};
	}
};
</script>
