<template>
	<div class="notify-detail">
		<v-dialog v-model="dialog" width="800" @click:outside="close">
			<v-card class="cards">
				<div class="card-icon">
					<v-btn class="mx-2" fab dark x-small color="black" @click="close">
						<v-icon dark> mdi-close</v-icon>
					</v-btn>
				</div>
				<v-card-title class="text-h6 lighten-2 font-weight-bold"> 业务通知 </v-card-title>

				<div class="proCode items">
					<div class="items-left">项目编码:</div>
					<div class="items-right">{{ BusinessPush.proCode }}</div>
				</div>
				<div class="notify-type items">
					<div class="items-left">通知类型:</div>
					<div class="items-right">{{ BusinessPush.title }}</div>
				</div>
				<div class="notify-time items">
					<div class="items-left">通知发送时间:</div>
					<div class="items-right">{{ BusinessPush.CreatedAt | formTime }}</div>
				</div>
				<div class="notify-content items">
					<div class="items-left">通知内容:</div>
					<div class="items-right">
						<v-textarea solo name="input-7-4" v-model="BusinessPush.msg"></v-textarea>
					</div>
				</div>
				<!-- <v-divider></v-divider> -->

				<v-card-actions>
					<v-spacer></v-spacer>
					<div class="z-flex justify-center" style="width: 100%">
						<v-btn color="primary" width="100px" @click="close"> 确定 </v-btn>
					</div>
				</v-card-actions>
			</v-card>
		</v-dialog>
	</div>
</template>
<script>
import moment from "moment";
export default {
	data() {
		return {
			dialog: false,
			content: "",
			BusinessPush: {},
			ID: ""
		};
	},
	props: {
		row: {
			type: Object,
			default: {}
		}
	},
	watch: {
		row(val) {
			this.BusinessPush = val.BusinessPush;
			this.ID = val.ID;
		},
		deep: true
	},
	filters: {
		formTime(time) {
			if (!time) return "";
			return moment(time).format("YYYY-MM-DD HH:mm:ss");
		}
	},
	methods: {
		async updateRead() {
			const body = {
				ids: [this.ID],
				proCode: this.BusinessPush.proCode
			};
			const result = await this.$store.dispatch("UPDATE_NOTIFICATION_READ", body);
			if (result.code === 200) {
				this.$emit("emitSearch");
			}
		},
		close() {
			this.updateRead();
			this.dialog = false;
		},
		open() {
			this.dialog = true;
		}
	}
};
</script>

<style lang="scss" scoped>
.cards {
	position: relative;
	padding-bottom: 20px;
	.card-icon {
		position: absolute;
		top: 15px;
		right: 5px;
	}
	.items {
		display: flex;
		justify-content: center;
		padding: 10px 0;
		.items-left {
			display: flex;
			justify-content: end;
			width: 150px;
		}
		.items-right {
			display: flex;
			width: 400px;
			margin-left: 30px;
		}
	}
}
</style>
