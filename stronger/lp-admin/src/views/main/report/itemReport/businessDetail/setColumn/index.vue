<template>
	<v-dialog v-model="dialog" width="1200">
		<v-card>
			<div class="x_container">
				<div class="lists">
					<div class="ther grey lighten-2">
						<div class="ther_left">
							<v-checkbox
								v-model="isAll1"
								label="列表1"
								hide-details
								@change="changeList1"
							></v-checkbox>
						</div>
						<div class="ther_right">{{ selected1Count }}/{{ check1Pool.length }}</div>
					</div>
					<div class="checks_list">
						<div class="check">
							<v-text-field
								placeholder="请输入关键字"
								v-model="originValue"
								@input="onInput1"
							></v-text-field>
						</div>
						<div class="scroll_wrap">
							<div class="check" v-for="item in check1Pool" :key="item.value">
								<v-checkbox
									v-model="forms1[item.value]"
									:label="item.label"
									hide-details
								></v-checkbox>
							</div>
						</div>
					</div>
				</div>
				<div class="operation">
					<div class="pb-2 z-flex justify-center">
						<v-btn class="mx-2" fab small @click="toRight">
							<v-icon dark> mdi-arrow-right </v-icon>
						</v-btn>
					</div>
					<div class="pt-2 z-flex justify-center">
						<v-btn class="mx-2" fab small @click="toLeft">
							<v-icon dark> mdi-arrow-left </v-icon>
						</v-btn>
					</div>
				</div>
				<div class="lists">
					<div class="ther grey lighten-2">
						<div class="ther_left">
							<v-checkbox
								v-model="isAll2"
								label="列表2"
								hide-details
								@change="changeList2"
							></v-checkbox>
						</div>
						<div class="ther_right">{{ selected2Count }}/{{ check2Pool.length }}</div>
					</div>
					<div class="checks_list">
						<div class="check">
							<v-text-field
								placeholder="请输入关键字"
								v-model="targetValue"
								@input="onInput2"
							></v-text-field>
						</div>
						<div class="scroll_wrap">
							<div class="check" v-for="item in check2Pool" :key="item.value">
								<v-checkbox
									v-model="forms2[item.value]"
									:label="item.label"
									hide-details
								></v-checkbox>
							</div>
						</div>
					</div>
				</div>
			</div>

			<v-card-actions>
				<v-spacer></v-spacer>
				<v-btn color="primary" text @click="dialog = false">取消</v-btn>
				<v-btn color="primary" @click="onOk">确定</v-btn>
			</v-card-actions>
		</v-card>
	</v-dialog>
</template>

<script>
import _ from "lodash";
export default {
	data() {
		return {
			dispatch: "REPORT_SETTING",
			dialog: false,
			isAll1: false,
			isAll2: false,
			forms1: {},
			forms2: {},
			originValue: "",
			targetValue: "",
			check1Pool: [],
			check2Pool: [],
			checkLists1: [],
			checkLists2: []
		};
	},
	props: {
		procode: {
			type: String
		}
	},
	created() {
		this.getReports();
	},
	computed: {
		selected1Count() {
			return Object.values(this.forms1).filter(bool => bool).length;
		},
		selected2Count() {
			return Object.values(this.forms2).filter(bool => bool).length;
		}
	},
	methods: {
		async getColumn() {
			const result = await this.$store.dispatch("REPORT_SETTING_CELL", {
				projectCode: this.procode
			});
			if (result.code === 200) {
				this.checkLists2 = result.data;

				const val = this.checkLists1.map(item => {
					const value = item.value;
					const findItem = this.checkLists2.find(node => node.value === value);
					if (!findItem) {
						return item;
					}
				});

				this.checkLists1 = val.filter(item => item);

				this.initPool();
			}
		},
		async getReports() {
			const result = await this.$store.dispatch(this.dispatch, {});
			if (result.code === 200) {
				this.checkLists1 = result.data;
				this.initPool();
			}
		},
		initPool() {
			this.check1Pool = _.cloneDeep(this.checkLists1);
			this.check2Pool = _.cloneDeep(this.checkLists2);
		},
		onInput1(val) {
			this.check1Pool = [];
			this.checkLists1.forEach(item => {
				if (item.label.includes(val)) {
					this.check1Pool.push(item);
				}
			});
		},
		onInput2(val) {
			this.check2Pool = [];
			this.checkLists2.forEach(item => {
				if (item.label.includes(val)) {
					this.check2Pool.push(item);
				}
			});
		},
		openDialog() {
			this.dialog = true;
			this.getColumn();
		},
		toRight() {
			const keys = Object.keys(this.forms1).forEach(key => {
				if (this.forms1[key]) {
					const index = this.checkLists1.findIndex(item => item.value === key);
					const popArr = this.checkLists1.splice(index, 1);
					this.checkLists2.push(...popArr);
				}
			});
			this.initPool();
			this.forms1 = {};
			this.isAll1 = false;
		},
		toLeft() {
			Object.keys(this.forms2).forEach(key => {
				if (this.forms2[key]) {
					const index = this.checkLists2.findIndex(item => item.value === key);
					const popArr = this.checkLists2.splice(index, 1);
					this.checkLists1.push(...popArr);
				}
			});
			this.initPool();
			this.forms2 = {};
			this.isAll2 = false;
		},
		changeList1(value) {
			if (value) {
				this.checkLists1.forEach(item => {
					this.forms1[item.value] = true;
				});
			} else {
				this.checkLists1.forEach(item => {
					this.forms1[item.value] = false;
				});
			}
		},
		changeList2(value) {
			if (value) {
				this.checkLists2.forEach(item => {
					this.forms2[item.value] = true;
				});
			} else {
				this.checkLists2.forEach(item => {
					this.forms2[item.value] = false;
				});
			}
		},
		async onOk() {
			if (this.procode === "") {
				this.toasted.warning("请先选择设置报表的项目");
			} else {
				if (this.check2Pool.length === 0) {
					this.toasted.warning("请选择设置报表的表头");
					return;
				}
				const tagsList = this.check2Pool.map(item => item.value);
				const result = await this.$store.dispatch("REPORT_SETTING_SET", {
					projectCode: this.procode,
					tagsList
				});
				this.toasted.dynamic(result.msg, result.code);
				if (result.code === 200) {
					this.getColumn();
					this.$emit("updateColumn");
					this.dialog = false;
				}
			}
		}
	},
	watch: {
		forms1: {
			handler(val) {
				this.isAll1 =
					this.selected1Count === this.check1Pool.length
						? this.selected1Count === 0
							? false
							: true
						: false;
			},
			deep: true
		},
		forms2: {
			handler(val) {
				this.isAll2 =
					this.selected2Count === this.check2Pool.length
						? this.selected2Count === 0
							? false
							: true
						: false;
			},
			deep: true
		}
	}
};
</script>

<style lang="scss" scoped>
.x_container {
	padding: 0 10px 0 0;
	display: flex;
	align-items: center;
	.operation {
		width: 10%;
	}
	.lists {
		width: 45%;
		border: 1px solid #eee;
		.ther {
			padding: 5px 10px;
			display: flex;
			align-items: center;
			justify-content: space-between;
			.ther_right {
				position: relative;
				top: 10px;
			}
		}
		.scroll_wrap {
			height: 500px;
			overflow-y: scroll;
		}
	}
	.checks_list {
		padding: 0 10px;
	}
}
</style>
