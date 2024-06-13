<template>
	<v-dialog v-model="dialog" width="500">
		<v-card>
			<v-card-title class="text-h6 grey lighten-2"> 已完成提示 </v-card-title>

			<v-card-text class="pt-2"> 恭喜{{ user.name }}，您已完成文档学习，请前往上岗练习进行系统训练 </v-card-text>

			<v-divider></v-divider>

			<v-card-actions>
				<v-spacer></v-spacer>
				<v-btn color="primary" text @click="handleClick">上岗练习 </v-btn>
			</v-card-actions>
		</v-card>
	</v-dialog>
</template>

<script>
import { localStorage, sessionStorage } from "vue-rocket";
export default {
	data() {
		return {
			dialog: false,
			user: {}
		};
	},
	mounted() {
		this.user = localStorage.get("user");
	},
	methods: {
		open() {
			this.dialog = true;
		},
		handleClick() {
			this.$router.push("/main/practice/practiceChannel");
			this.dialog = false;
		}
	}
};
</script>

<style lang="scss" scoped></style>
