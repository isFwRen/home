<template>
	<div class="text-center">
		<v-dialog v-model="dialog" width="420">
			<v-card>
				<v-card-title class="text-h5 grey lighten-2">异常提示</v-card-title>

				<v-card-text class="mt-4">
					案件存在异常数据，请查看导出校验进行修改后再进行回传
				</v-card-text>

				<v-divider></v-divider>

				<v-card-actions>
					<v-spacer></v-spacer>
					<v-btn color="primary" text @click="dialog = false">确定 </v-btn>
				</v-card-actions>
			</v-card>
		</v-dialog>
	</div>
</template>
<script>
export default {
	data() {
		return {
			dialog: false
		};
	},
	methods: {
		open() {
			this.dialog = true;
		}
	}
};
</script>
<style scoped lang="scss"></style>
