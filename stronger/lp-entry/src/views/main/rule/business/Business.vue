<template>
	<div class="lp-business">
		<div>
			<div class="z-flex align-center">
				<z-btn class="mr-3" color="primary" small unlocked @click="selectAllCode">
					{{ projectFlag.length === codeHeaders.length ? "全不选" : "全选" }}
				</z-btn>
				<z-checkboxs ref="proCode" :formId="searchFormId" formKey="projectArr" :options="codeHeaders"></z-checkboxs>
			</div>

			<div class="z-flex align-center">
				<z-select
					:formId="searchFormId"
					formKey="ruleArr"
					clearable
					hideDetails
					label="规则类型"
					:options="cells.rule"
				></z-select>
				<v-btn color="primary" @click="onSearch">
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</v-btn>
			</div>
		</div>

		<div class="file">
			<div class="dad" v-for="item in dataSource" :key="item.model.ID" @click="viewDetail(item)">
				<div :class="['f_img', 'f_img_hover', item.class]">
					<div class="son">
						<p>{{ item.proCode }}</p>
						<p>
							{{ item.ruleName }}
							<span class="done_study" v-show="item.isRequired === 1 && item.isLearned === 1"> (已学)</span>
						</p>
					</div>
				</div>

				<div class="sons">
					<span class="must" v-show="item.isRequired === 1">必学 </span>
					<p>更新时间: {{ item.CreatedAt }}</p>
				</div>
			</div>
		</div>
		<z-pagination :options="pageSizes" @page="handlePage" :total="pagination.total"></z-pagination>
	</div>
</template>

<script>
import moment from "moment";
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";
import { sessionStorage, localStorage, tools } from "vue-rocket";
const { baseURLApi } = lpTools.baseURL();
import { driver } from "driver.js";
import "driver.js/dist/driver.css";
import _ from "lodash";

export default {
	name: "LPBusiness",
	mixins: [TableMixins],

	data() {
		return {
			pdfUrl: "",
			formId: "TeachingBusinessRules",
			cells,
			dispatchList: "RULE_GET_BUSINESS_LIST",
			projectFlag: false,
			projectArr: [],
			ruleFlag: false,
			ruleArr: [],
			codeHeaders: [],
			dataSource: [],
			steps: [],
			broadListen: null
		};
	},

	created() {
		const auths = localStorage.get("auth");
		this.codeHeaders = auths.proItems;
		this.broadListen = new BroadcastChannel("updateLearn");
		this.broadListen.addEventListener("message", event => {
			const id = event.data.id;
			const findItem = this.dataSource.find(item => item.model.ID === id);
			findItem.isLearned = 1;
		});
	},
	methods: {
		async viewDetail(row) {
			const filePath = baseURLApi + row.docsPath;
			const res = await this.$store.dispatch("GET_PDF_FILE", filePath);
			sessionStorage.set("pdfUrl", URL.createObjectURL(res));
			sessionStorage.set("post_data", {
				filishId: row.model.ID,
				projectCode: row.proCode,
				trainType: 1,
				isRequired: row.isRequired
			});

			window.open(`${location.origin}/normal/pdf-file`, "_blank", "toolbar=yes, scrollbars=yes, resizable=yes");
		},
		selectAllCode() {
			this.$refs.proCode.onSelectAll();
		},
		ruleClick() {
			if (this.ruleFlag) this.ruleArr = cells.rule.map(el => el.value);
			else this.ruleArr = [];
		},
		setDriver() {
			const driverObj = driver({
				showProgress: true,
				allowClose: false,
				steps: this.steps,
				doneBtnText: "完成",
				closeBtnText: "关闭",
				nextBtnText: "下一步",
				prevBtnText: "上一步"
			});
			driverObj.drive();
		}
	},
	destroy() {
		this.broadListen.close();
	},
	computed: {
		...mapGetters(["auth"])
	},

	watch: {
		desserts(val) {
			const sortArr = _.cloneDeep(val);

			// 先排序，必填排前面
			sortArr.sort((a, b) => {
				return b.isRequired - a.isRequired;
			});

			this.steps.length = 0;

			sortArr.forEach((item, index) => {
				item.CreatedAt = moment(item.UpdatedAt).format("YYYY-MM-DD");

				if (item.isRequired === 1) {
					item.class = `popover-node${index + 1}`;
					this.steps.push({
						element: `.${item.class}`,
						popover: {
							title: `${item.ruleName}必学`,
							description: "请点击进行学习",
							side: "bottom",
							align: "start"
						}
					});
				} else {
					item.class = "";
				}
			});
			this.dataSource = sortArr;
			if (this.steps.length > 0) {
				this.steps.push({
					element: `.popover-template`,
					popover: {
						title: `报销单模板`,
						description: "请点击进行学习",
						side: "bottom",
						align: "start"
					}
				});
				const timer = setTimeout(() => {
					this.setDriver();
					clearInterval(timer);
				}, 250);
			}
		}
	}
};
</script>

<style lang="scss" scoped>
.file {
	display: flex;
	justify-content: flex-start;
	flex-wrap: wrap;
	margin-top: 20px;
	gap: 20px;
	.dad {
		width: 160px;
		height: 200px;
		position: relative;
		.f_img_hover {
			transition: 0.2s;
			background: url("../../../../assets/images/files/file.png");
			cursor: pointer;
			color: #000;

			&:hover {
				background: url("../../../../assets/images/files/fileHover.png") !important;
				cursor: pointer;
				color: #1296db;
			}
		}

		.f_img {
			margin: 0 auto;
			width: 160px;
			height: 150px;
			display: flex;
			justify-content: center;
			align-items: center;
			img {
				width: 100%;
				height: 100%;
			}
		}

		.son {
			p {
				margin-bottom: 5px;
				text-align: center;
			}
		}

		.sons {
			position: relative;
			padding: 5px 0;
			p {
				text-align: center;
			}

			.must {
				position: absolute;
				display: inline-block;
				background-color: #27b148;
				width: 40px;
				height: 20px;
				text-align: center;
				color: white;
				top: -142px;
				left: 29px;
			}
		}
	}
}
</style>
