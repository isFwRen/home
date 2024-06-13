<template>
	<div class="RealTime">
		<div :class="projectBoxClass">
			<div class="head"><span>钉钉群名：</span></div>
			<div class="selectBox">
				<z-btn-toggle
					:formId="projectFormId"
					formKey="sexual"
					class="mt-n4"
					color="primary"
					mandatory
					:options="groupList"
					@change="onChange"
				></z-btn-toggle>
			</div>
			<div class="opition">
				<z-btn class="mr-4 mb-2" small color="primary" @click="isDownLoad = !isDownLoad">{{
					isDownLoad ? "收起" : "更多"
				}}</z-btn>
			</div>
		</div>
		<div class="beChooseCount">
			<span
				>[单量通知] 目前项目的男150、报销单50;二马单量120、报销单34
				请安排时间上线处理。</span
			>
		</div>

		<div class="table-box">
			<div class="table-button">
				<z-btn color="primary" small @click="AddRegularFir">
					<v-icon>mdi-plus</v-icon>
					新增
				</z-btn>
				<z-btn class="pl-3" small color="primary" outlined @click="isSaveFr()">
					<v-icon left v-if="FirTableDisabled"> mdi-pencil </v-icon>
					{{ FirTableDisabled ? "编辑" : "保存" }}
				</z-btn>
				<z-btn class="pl-3" color="warning" small @click="RE(1)">
					<v-icon class="text-h6">mdi-reload</v-icon>
					重置
				</z-btn>
			</div>
			<div class="table">
				<vxe-table :data="cells.FirdefaultVlue" :border="tableBorder" :size="tableSize">
					<template v-for="(item, i) in cells.headersOne">
						<vxe-column
							v-if="item.value == 'name'"
							:field="item.value"
							:title="item.text"
							:width="item.width"
							:key="i"
						>
							<template #default="{ row }">
								{{
									item["outPut"] ? item.outPut(row[item.value]) : row[item.value]
								}}
							</template>
						</vxe-column>
						<vxe-column
							v-else
							:field="item.value"
							:title="item.text"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-date-picker
									v-if="row.type != 'interval'"
									formId="fir"
									:formKey="row.key + item.value + row.line"
									label=""
									format="24hr"
									:immediate="false"
									mode="time"
									prepend-icon="mdi-alarm"
									time-use-seconds
									time-format="24hr"
									:defaultValue="
										item['outPut']
											? item.outPut(row[item.value])
											: row[item.value]
									"
									:disabled="FirTableDisabled"
								></z-date-picker>
								<z-text-field
									v-else
									formId="fir"
									:formKey="row.key + item.value + row.line"
									label=""
									:validation="[
										{
											regex: /^[1-9]([0-9]{0,2})$/,
											message: '必须是大于0的数字且不能大于3位'
										}
									]"
									:defaultValue="
										item['outPut']
											? item.outPut(row[item.value])
											: row[item.value]
									"
									:disabled="FirTableDisabled"
								>
								</z-text-field>
							</template>
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
		<div class="table-box">
			<div class="table-button">
				<z-btn color="primary" small @click="AddRegularSec">
					<v-icon>mdi-plus</v-icon>
					新增
				</z-btn>
				<z-btn class="pl-3" small color="primary" outlined @click="isSave()">
					<v-icon left v-if="SecTableDisabled"> mdi-pencil </v-icon>
					{{ SecTableDisabled ? "编辑" : "保存" }}
				</z-btn>
				<z-btn class="pl-3" color="warning" small @click="RE(2)">
					<v-icon class="text-h6">mdi-reload</v-icon>
					重置
				</z-btn>
			</div>
			<div class="table">
				<vxe-table :data="cells.SecdefaultVlue" :border="tableBorder" :size="tableSize">
					<template v-for="(item, i) in cells.headersTwo">
						<vxe-column
							v-if="item.value == 'name'"
							:field="item.value"
							:title="item.text"
							:width="item.width"
							:key="i"
						>
							<template #default="{ row }">
								{{
									item["outPut"] ? item.outPut(row[item.value]) : row[item.value]
								}}
							</template>
						</vxe-column>

						<vxe-column
							v-else
							:field="item.value"
							:title="item.text"
							:width="item.width"
						>
							<template #default="{ row }">
								<z-date-picker
									formId="sec"
									:formKey="row.key + item.value"
									label=""
									format="24hr"
									:immediate="false"
									mode="time"
									prepend-icon="mdi-alarm"
									time-use-seconds
									time-format="24hr"
									:defaultValue="
										item['outPut']
											? item.outPut(row[item.value])
											: row[item.value]
									"
									:disabled="SecTableDisabled"
								></z-date-picker>
							</template>
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
		<AddRegularFir ref="AddRegularFir" :id="groupID" @submitted="submitted" />
		<AddRegularSec
			ref="AddRegularSec"
			:id="groupID"
			:proCode="proCode"
			@submitted="submitted"
		/>
	</div>
</template>
<script>
import cells from "./cells";
import TableMixins from "@/mixins/TableMixins";
import AddRegularFir from "./Dlog/addRegularFir.vue";
import AddRegularSec from "./Dlog/addRegularSec.vue";
import { mapGetters } from "vuex";

export default {
	mixins: [TableMixins],
	data() {
		return {
			isDownLoad: false,
			projectBeChoosed: [],
			projectFormId: "projectForm",
			cells,
			FirTableDisabled: true,
			SecTableDisabled: true,
			FirdefaultVlue: cells.FirdefaultVlue,
			SecdefaultVlue: cells.SecdefaultVlue,
			groupList: [],
			groupObj: {},
			groupID: "",
			result: {},
			proCode: ""
		};
	},
	created() {
		this.getGroupList();
	},
	methods: {
		async isSaveFr() {
			if (!this.FirTableDisabled) {
				const form = this.forms["fir"];
				let oneObj = {};
				const objKeyArr = ["startTime", "endTime", "interval"];
				let oneArr = [];
				for (let key in form) {
					let keyArr = key.split("");
					if (!oneObj[keyArr[0]]) {
						oneObj[keyArr[0]] = {};
					}
					if (!oneObj[keyArr[0]][keyArr[keyArr.length - 2]]) {
						oneObj[keyArr[0]][keyArr[keyArr.length - 2]] = {};
					}
					oneObj[keyArr[0]][keyArr[keyArr.length - 2]][
						objKeyArr[keyArr[keyArr.length - 1]]
					] = form[key];
				}
				for (let k in oneObj) {
					for (let k2 in oneObj[k])
						oneArr.push({
							...this.result.one[k][k2],
							...oneObj[k][k2],
							dayOfWeek: +k2,
							block: +k,
							groupID: this.groupID,
							interval: +oneObj[k][k2].interval,
							proCode: this.groupObj[this.groupID].proCode
						});
				}
				let ones = oneArr.filter(e => {
					return e.startTime && e.endTime && e.interval;
				});
				const result = await this.$store.dispatch("EDIT_NOTIC_NEW_REGULAR_BYTIME", {
					ones,
					type: 1
				});
				this.toasted.dynamic(result.msg, result.code);

				this.getByproList(this.groupID);
			}
			this.FirTableDisabled = !this.FirTableDisabled;
		},

		async isSave() {
			if (!this.SecTableDisabled) {
				const form = this.forms["sec"];
				console.log(form);
				let twos = [];
				for (let key in form) {
					let strArr = key.split("");
					twos.push({
						...this.result.two[strArr[0]][strArr[strArr.length - 1]],
						block: +strArr[0],
						dayOfWeek: +strArr[strArr.length - 1],
						sendTime: form[key],
						groupID: this.groupID,
						proCode: this.groupObj[this.groupID].proCode
					});
				}
				twos = twos.filter(e => {
					return e.sendTime;
				});
				const result = await this.$store.dispatch("EDIT_NOTIC_NEW_REGULAR_BYTIME", {
					twos,
					type: 2
				});
				this.toasted.dynamic(result.msg, result.code);
				this.getByproList(this.groupID);
			}
			this.SecTableDisabled = !this.SecTableDisabled;
		},
		async RE(type) {
			const result = await this.$store.dispatch("RESET_NOTIC_REGULAR", {
				type: type,
				groupId: this.groupID
			});
			this.toasted.dynamic(result.msg, result.code);
			if (result.code == 200) {
				this.getByproList(this.groupID);
			}
		},
		async getByproList(proCode) {
			const result = await this.$store.dispatch("GET_NOTIC_BY_PRO_LIST", {
				proCode: proCode
			});
			if (result.code == 200) {
				let one = {},
					two = {};
				for (let key in result.data.one) {
					let blockArr = result.data.one[key];
					if (!one[key]) {
						one[key] = {};
					}
					for (let key2 in blockArr) {
						one[key][blockArr[key2].dayOfWeek] = blockArr[key2];
					}
				}
				for (let key in result.data.two) {
					let blockArr = result.data.two[key];
					if (!two[key]) {
						two[key] = {};
					}
					for (let key2 in blockArr) {
						two[key][blockArr[key2].dayOfWeek] = blockArr[key2];
					}
				}
				this.result = { one, two };
				cells.FirdefaultVlue = [];
				cells.SecdefaultVlue = [];
				for (let key in result.data.one) {
					cells.FirdefaultVlue = [
						...cells.FirdefaultVlue,
						...cells.ONEToDataOne(result.data.one[key], key)
					];
				}
				for (let key in result.data.two) {
					cells.SecdefaultVlue.push(cells.TWOToDataTwo(result.data.two[key], key));
				}
			}
		},
		onChange(e) {
			this.groupID = e;
			this.getByproList(e);
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
						this.groupObj[e.ID] = e;
						return { label: e.name, value: e.ID };
					});
				}
			}
		},
		AddRegularFir() {
			this.$refs.AddRegularFir.onOpen();
		},
		AddRegularSec() {
			this.$refs.AddRegularSec.onOpen();
		},
		submitted() {
			this.getByproList(this.groupID);
		}
	},
	computed: {
		projectBoxClass: function () {
			return "projectBox " + (this.isDownLoad ? "download" : "upHidde");
		},
		...mapGetters(["auth"])
	},
	components: { AddRegularFir, AddRegularSec }
};
</script>
<style lang="scss" scoped>
.projectBox {
	display: grid;
	grid-template-columns: 100px 8fr 70px;
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
		::v-deep button:not(:first-child) {
			margin-left: 10px;
			border-left-width: 1px !important;
		}
	}
}
::v-deep .v-text-field,
::v-deep .v-input__slot {
	padding: 0;
	margin: 0;
}

//  ::v-deep .v-text-field>.v-input__control>.v-input__slot:before{
//   border-style: none;
//  }
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
	color: #333;
	font-weight: 400;
	font-size: 0.9em;
	margin-top: 10px;
	span {
		display: inline-block;
		margin-right: 1.3em;
	}
	.table-button {
		display: flex;
		justify-content: right;
		margin-bottom: 10px;
	}
}
</style>
