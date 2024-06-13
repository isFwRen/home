<template>
	<div class="lp-video">
		<div>
			<div class="z-flex align-center">
				<z-btn class="mr-3" color="primary" small unlocked @click="selectAllCode">
					{{ projectArr.length === codeHeaders.length ? "全不选" : "全选" }}
				</z-btn>
				<z-checkboxs ref="proCodeRef" :formId="searchFormId" formKey="proCode" :options="codeHeaders"></z-checkboxs>
				<v-btn color="primary" @click="onSearch">
					<v-icon class="text-h6">mdi-magnify</v-icon>
					查询
				</v-btn>
			</div>

			<!-- <div class="z-flex align-center">
				<z-btn class="mr-3" color="primary" small unlocked @click="selectAllRule">
					{{ ruleArr.length === codeHeaders.length ? "全不选" : "全选" }}
				</z-btn>
				<z-checkboxs ref="rule" :formId="searchFormId" formKey="ruleArr" :options="cells.rule"></z-checkboxs>
			</div> -->
		</div>

		<div class="file">
			<div v-for="item in dataSource" :key="item.model.ID">
				<div :class="['dad', item.class]" @click="viewDetail(item)" v-if="item.video && item.video.length > 0">
					<div class="video_w">
						<video :ref="`video_${item.model.ID}`" controls @ended="onEnded(item)">
							<source type="video/mp4" />
							您的浏览器不支持 HTML5 video 标签。
						</video>
					</div>
					<div class="son">
						<p>
							<span class="must" v-show="item.isRequired === 1">必学</span>
						</p>
						<p>
							<span>{{ item.sysBlockName }}</span>
							<span class="done_study" v-show="item.isRequired === 1 && item.isLearned === 1"> (已学)</span>
						</p>
					</div>
				</div>
			</div>
		</div>
		<z-pagination :options="pageSizes" @page="handlePage" :total="pagination.total"></z-pagination>
		<view-dialog ref="viewDialog" :proCode="proCode" :video="video"></view-dialog>

		<done-dialog ref="doneDialog" />
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import TableMixins from "@/mixins/TableMixins";
import cells from "./cells";
import { sessionStorage, localStorage, tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import _ from "lodash";
const { baseURLApi } = lpTools.baseURL();
import { driver } from "driver.js";
import "driver.js/dist/driver.css";

export default {
	name: "LPVideo",
	mixins: [TableMixins],
	tools,
	data() {
		return {
			formId: "LPVideo",
			cells,
			dispatchList: "RULE_VIDEO_GET_LIST",
			manual: true,
			proCode: "",
			codeHeaders: [],
			video: "",
			projectArr: [],
			ruleArr: [],
			dataSource: [],
			steps: []
		};
	},
	created() {
		const auths = localStorage.get("auth");
		this.codeHeaders = auths.proItems;
	},
	watch: {
		desserts: {
			handler(val) {
				const sortArr = _.cloneDeep(val);
				// 先排序，必填排前面
				sortArr.sort((a, b) => {
					return b.isRequired - a.isRequired;
				});

				this.steps.length = 0;

				sortArr.forEach((item, index) => {
					if (item.isRequired === 1) {
						item.class = `popover-node${index + 1}`;
						this.steps.push({
							element: `.${item.class}`,
							popover: {
								title: `${item.sysBlockName}必学`,
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
					const timer = setTimeout(() => {
						this.setDriver();
						clearInterval(timer);
					}, 250);
				}

				this.dataSource.forEach(async item => {
					if (!item.video || item.video?.length === 0) {
						item.videoSrc = "";
						return;
					}
					const path = item.video[0].path;
					const blob = await lpTools.getTokenImg(`${baseURLApi}${path}`);
					if (blob) {
						lpTools.getBase64(blob).then(base64String => {
							item.videoSrc = base64String;
							this.$refs[`video_${item.model.ID}`][0].src = base64String;
						});
					} else {
						item.videoSrc = "";
					}
				});
			}
		}
	},
	methods: {
		async onSearch() {
			this.params = {
				...this.params,
				...this.page
			};

			const total = await this.getList();
			if (typeof total !== "number") {
				return;
			}
			const index = this.params.pageIndex - 1;
			if (index * this.params.pageSize > total) {
				this.params.pageIndex = 1;
			}
		},
		handlePage(page) {
			this.params = {
				...this.params,
				pageSize: page.pageSize,
				pageIndex: page.pageNum
			};
			this.getList();
		},
		async getList() {
			if (this.dispatchList) {
				const params = {
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
				return total;
			}

			this.loading = false;

			return this.sabayon;
		},
		async viewDetail(row) {
			// this.$refs.viewDialog.onOpen({ status: -1 });
			// this.video = row.video ?? [];
		},
		selectAllCode() {
			this.$refs.proCodeRef.onSelectAll();
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
		},
		selectAllRule() {
			this.$refs.rule.onSelectAll();
		},
		async onEnded(row) {
			if (row.isRequired === 1) {
				const body = { finishId: row.model.ID, projectCode: this.forms[this.searchFormId].proCode[0], trainType: 3 };
				const res = await this.$store.dispatch("FILISH_FILE_READ", body);
				if (res.code === 200) {
					row.isLearned = 1;
				}
				this.toasted.dynamic(res.msg, res.code);
			}
			const filterRequired = this.dataSource.filter(item => item.isRequired === 1);
			const filterLearned = filterRequired.filter(item => item.isLearned === 1);
			if (filterLearned.length > 0) {
				this.$refs.doneDialog.open();
			}
			console.log(filterRequired, filterLearned, "filterLearned");
		}
	},

	computed: {
		...mapGetters(["auth"])
	},
	components: {
		"view-dialog": () => import("./viewDialog"),
		"done-dialog": () => import("./doneDialog/index.vue")
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
		width: 240px;
		height: 230px;
		overflow: hidden;
		.video_w {
			video {
				width: 100%;
				height: 180px;
				object-fit: contain;
			}
		}
		.son {
			top: 42px;
			left: 28px;

			p {
				margin-bottom: 5px;
				text-align: center;
			}

			.must {
				display: inline-block;
				background-color: #27b148;
				width: 40px;
				height: 20px;
				margin-right: 5px;
				border-radius: 5px;
				color: white;
			}
		}
	}

	.dad:hover {
		cursor: pointer;
	}
}
</style>
