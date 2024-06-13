<template>
	<v-dialog v-model="dialog" width="500" persistent>
		<v-card>
			<v-card-title class="blue darken-1"> 消息通知 </v-card-title>

			<v-card-text class="pt-3">
				<div class="message" v-for="(item, index) in data" :key="index">{{ item }}</div>
			</v-card-text>

			<v-divider></v-divider>

			<v-card-actions>
				<v-spacer></v-spacer>
				<v-btn color="primary" text @click="onOk"> 确定 </v-btn>
			</v-card-actions>
		</v-card>
	</v-dialog>
</template>

<script>
export default {
	data() {
		return { dialog: false };
	},
	props: {
		data: {
			type: Array
		}
	},
	methods: {
		openDialog() {
			this.dialog = true;
		},
		onOk() {
			this.$EventBus.$emit("clearSocketMessage");
			this.dialog = false;
		}
	}
};
</script>

<style lang="scss" scoped>
.message {
	line-height: 26px;
	color: #000;
	font-size: 18px;
}
</style>
