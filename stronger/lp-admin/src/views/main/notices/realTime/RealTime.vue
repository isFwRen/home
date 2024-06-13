<template>
	<div class="input">
		<div :class="projectBoxClass">
			<div class="head"><span>钉钉群名：</span></div>
			<div class="selectBox">
				<z-checkboxs
					:formId="projectFormId"
					formKey="groupId"
					@change="text"
					:options="groupList"
					:defaultValue="projectBeChoosed"
				>
				</z-checkboxs>
			</div>
			<div class="opition">
				<z-btn class="mr-4 mb-2" color="primary" small @click="chooseAll">全选</z-btn>
				<z-btn class="mr-4 mb-2" small @click="isDownLoad = !isDownLoad">{{
					isDownLoad ? "收起" : "更多"
				}}</z-btn>
			</div>
		</div>
		<div class="beChooseCount">
			<span></span> 所选钉钉群 {{ projectBeChoosed.length }}个:
			<span v-for="e in projectBeChoosed" class="beChoose">
				{{ groupObj[e] }}
			</span>
			<z-btn class="mr-4 mb-2 sendButton" color="primary" small @click="sendNoticeInput"
				>发送</z-btn
			>
		</div>

		<div class="text-box">
			<z-textarea
				:formId="projectFormId"
				formKey="msg"
				label="输入实时通知内容"
				outlined
				placeholder="请输入实时通知内容，长度控制在一百字以内."
				:defaultValue="msg"
			></z-textarea>
		</div>
		<div class="table-box">
			<span>录入实时记录</span>
			<div class="table">
				<vxe-table :data="desserts" :border="tableBorder" :size="tableSize">
					<template v-for="item in cells.headers">
						<vxe-column :field="item.value" :title="item.text" :width="item.width">
						</vxe-column>
					</template>
				</vxe-table>
			</div>
			<z-pagination
				:options="pageSizes"
				@page="handlePage"
				:total="pagination.total"
			></z-pagination>
		</div>
	</div>
</template>
<script>
import { rocket } from "vue-rocket";
import cells from "./cells";
import TableMixins from "@/mixins/TableMixins";
import { mapGetters } from "vuex";

export default {
	mixins: [TableMixins],
	data() {
		return {
			isDownLoad: false,
			projectBeChoosed: [],
			dispatchList: "GET_NOTICE_REALTIME_LIST",
			projectFormId: "projectForm",
			cells,
			msg: "",
			groupList: [],
			groupObj: {}
		};
	},
	created() {
		this.init();
	},
	methods: {
		init() {
			this.getGroupList();
		},
		async getGroupList(pageSize = 500) {
			const result = await this.$store.dispatch("GET_PT_MESSAGE_TABLE_LIST", {
				pageIndex: 1,
				pageSize: pageSize
			});
			if (result.code == 200) {
				let data = result.data;
				if (data.total > pageSize) {
					setTimeout(this.getGroupList(data.total), 100);
				} else {
					this.groupObj = {};
					this.groupList = data.list.map(e => {
						this.groupObj[e.ID] = e.name;
						return { label: e.name, value: e.ID };
					});
				}
			}
		},
		text(e) {
			this.projectBeChoosed = e;
		},
		chooseAll() {
			if (this.groupList.length !== this.projectBeChoosed.length) {
				for (let key in this.groupList) {
					this.projectBeChoosed.push(this.groupList[key].value);
				}
				return;
			}
			this.projectBeChoosed = [];
		},
		async sendNoticeInput() {
			const data = this.forms[this.projectFormId];

			const result = await this.$store.dispatch("SEND_NOTICE_REALTIME", data);
			this.toasted.dynamic(result.msg, result.code);
			if (result.code === 200) {
				this.onSearch();
				rocket.emit("ZHT_CLEAR_FORM", this.projectFormId);
			}
		}
	},
	computed: {
		projectBoxClass: function () {
			return "projectBox " + (this.isDownLoad ? "download" : "upHidde");
		},
		...mapGetters(["auth"])
	}
};
</script>
<style lang="scss" scoped>
.projectBox {
	display: grid;
	grid-template-columns: 100px 8fr 150px;
	border-bottom: 1px solid #ddd;

	.head,
	.opition {
		margin-top: 16px;
		padding-top: 8px;
	}
	.selectBox {
		display: flex;
		justify-content: left;
		flex-wrap: wrap;
	}
}
.upHidde {
	grid-template-rows: 70px;
	overflow: hidden;
}
.download {
	grid-template-rows: none;
	overflow: auto;
}
.text-box {
	margin-top: 1em;
}
.sendButton {
	position: absolute;
	right: 2em;
}
.table-box,
.beChooseCount {
	color: #aaa;
	font-size: 0.9em;
	margin-top: 10px;
	position: relative;
	span {
		display: inline-block;
		margin-right: 1.3em;
	}
}
</style>
