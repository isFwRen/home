<template>
	<div class="business">
		<div class="mb-8 lp-filters">
			<v-row class="z-flex align-end">
				<v-col :cols="2">
					<z-select
						:formId="searchFormId"
						formKey="proCode"
						hideDetails
						label="项目编码"
						:options="auth.proItems"
						defaultValue="B0108"
					></z-select>
				</v-col>

				<v-col cols="3">
					<z-date-picker
						:formId="searchFormId"
						formKey="time"
						hideDetails
						label="发送时间"
						range
						z-index="10"
						:defaultValue="today"
					></z-date-picker>
				</v-col>

				<v-col cols="2">
					<z-select
						:formId="searchFormId"
						formKey="msgType"
						clearable
						hideDetails
						label="通知类型"
						:options="typeItems"
					></z-select>
				</v-col>

				<div class="z-flex">
					<z-btn class="pb-3 px-3" color="primary" @click="onSearch">
						<v-icon size="20">mdi-magnify</v-icon>
						查询
					</z-btn>
				</div>
			</v-row>
		</div>

		<!--通知list-start-->
		<div
			class="list mb-3"
			v-for="(item, index) in desserts"
			:key="item.ID"
			@click="handleAlert(item, index)"
			:class="[businessNotice.pushId === item.pushId ? 'focus-list' : '']"
		>
			<div class="list-content">
				<div class="list-avatar">{{ item.BusinessPush.proCode }}</div>
				<div class="list-content-right ml-4">
					<h4 class="list-content-title pb-1">
						<v-badge color="red" dot v-if="!item.isRead">
							{{ businessMsgType[item.BusinessPush.type] }}通知
						</v-badge>

						<span v-else>{{ businessMsgType[item.BusinessPush.type] }}通知</span>
					</h4>
					<div class="list-content-desc">
						{{ item.BusinessPush.msg }}
					</div>
				</div>
			</div>
			<div class="list-time pt-1 pb-1">
				<span class="time">{{ item.BusinessPush.CreatedAt | formTime }}</span>
			</div>
		</div>
		<!--通知list-end-->

		<z-pagination
			:options="pageSizes"
			@page="handlePage"
			:total="pagination.total"
		></z-pagination>

		<notifyDetail ref="notifyRef" :row="currentRow" @emitSearch="onSearch" />
	</div>
</template>
<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import moment from "moment";
import _ from "lodash";
const today = [moment().format("YYYY-MM-DD"), moment().format("YYYY-MM-DD")];
export default {
	name: "Business",
	mixins: [TableMixins],
	data() {
		return {
			today,
			formId: "Business",
			searchFormId: "Business",
			fab: false,
			dispatchList: "GET_BUSINESS_LIST",
			currentRow: {},
			typeItems: [],
			businessMsgType: {},
			firstTime: true
		};
	},
	filters: {
		formTime(time) {
			if (!time) return "";
			return moment(time).format("YYYY-MM-DD HH:mm:ss");
		}
	},
	created() {
		this.getConstant();
	},
	watch: {
		"businessNotice.pushId": {
			handler(value) {
				this.onSearch()
			}
		}
	},
	computed: {
		...mapGetters(["auth", "project", "forms", "businessNotice"])
	},
	methods: {
		async onSearch() {
			this.params = {
				...this.params,
				...this.page
			};
			this.getList();
		},
		async getList() {
			if (this.dispatchList) {
				const params = {
					firstTime: this.firstTime,
					...this.effectParams,
					...this.params,
					...this.forms[this.searchFormId]
				};

				const result = await this.$store.dispatch(this.dispatchList, params);
				const { list, total } = result.data;

				if (result.code === 200) {
					if (typeof list === "object") {
						if (list instanceof Array) {
							this.desserts = list;
						} else {
							this.desserts = [];
						}
						this.pagination.total = total;
					} else {
						this.desserts = result.data;
						this.pagination.total = this.desserts.length;
					}
				} else {
					this.toasted.error(result.msg);

					this.desserts = [];
					this.pagination.total = 0;
				}

				this.sabayon = result;
				this.firstTime = false;
			}

			this.loading = false;

			return this.sabayon;
		},
		async getConstant() {
			const result = await this.$store.dispatch("GET_DISCONECT");
			const arr = [];
			if (result.code === 200) {
				this.businessMsgType = result.data.businessMsgType;
				Object.keys(this.businessMsgType).forEach(key => {
					arr.push({
						label: this.businessMsgType[key],
						value: key
					});
				});
				this.typeItems = arr;
			}
		},
		handleAlert(item, index) {
			if (index === 0) {
				this.borderIndex = -1;
			}
			this.currentRow = item;
			this.$refs.notifyRef.open();
			if (item.pushId === this.businessNotice.pushId) {
				// 清空pushId
				this.$store.commit("CLEAR_NOTIFICATION_PUSHID");
				this.$store.commit("BUSINESS_NOTIFICATION_UPDATE_ITEM", {
					count: 0
				});
			}
		}
	},
	components: {
		notifyDetail: () => import("./notifyDetail.vue")
	}
};
</script>
<style scoped lang="scss">
.list {
	border: 2px solid #eee;
	border-radius: 10px;
	padding: 15px 15px 0 15px;
	box-sizing: border-box;
	width: 100%;
	font-size: 14px;
	position: relative;
	.list-content {
		display: flex;
		align-items: center;
		.list-avatar {
			min-width: 60px;
			min-height: 60px;
			background-color: #4b90e1;
			border-radius: 50%;
			display: flex;
			justify-content: center;
			align-items: center;
			color: #fff;
			border: 2px solid #fff;
			box-sizing: border-box;
		}
		.list-content-right {
			width: calc(100% - 60px);
			.list-content-title {
				font-weight: 600;
			}
			.list-content-desc {
				width: 99%;
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}
		}
	}
	.list-time {
		display: flex;
		justify-content: end;
		color: #64748b;
	}
	.lighting-wrap {
		position: absolute;
		top: 0;
		right: 0;
		.lighting {
			width: 20px;
			height: 20px;
			background: hsl(212, 100%, 96.27%);
			display: flex;
			align-items: center;
			justify-content: center;
			cursor: pointer;
		}
	}
}

.speed-wrapper {
	display: flex;
	justify-content: end;
	padding-top: 30px;
	.speed-card {
		width: 400px;
		background-color: #f8f8f8;
		padding: 10px;
		display: flex;
		font-size: 14px;
		border-radius: 10px;
		box-shadow: 1px 1px 2px 2px #eee;
		.speed-avatar {
			min-width: 55px;
			min-height: 55px;
			background-color: #4b90e1;
			display: flex;
			justify-content: center;
			align-items: center;
			border-radius: 50%;
			color: #fff;
		}
		.speed-right {
			flex: 1;
			padding-left: 5px;
			.speed-right-notify {
				display: flex;
				justify-content: space-between;
				.time {
					font-size: 12px;
					color: gray;
				}
			}
			.speed-subtitle {
				border-bottom: 2px solid #eee;
				padding-bottom: 10px;
			}
			.speed-action {
				display: flex;
				width: 70%;
				justify-content: space-between;
				margin-top: 10px;
				.speed-action-download {
					display: flex;
					align-items: center;
					span {
						margin-left: 10px;
					}
				}
			}
		}
	}
}

.focus-list {
	border: 2px solid #3b82f6;
}
</style>
