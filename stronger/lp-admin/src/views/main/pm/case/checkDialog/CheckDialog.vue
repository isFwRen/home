<template>
	<v-row justify="center">
		<v-dialog
			v-model="dialog"
			fullscreen
			hide-overlay
			transition="dialog-bottom-transition"
			no-click-animation
			persistent
		>
			<v-card class="card">
				<v-toolbar dark color="primary">
					<v-btn icon dark @click="dialog = false">
						<v-icon>mdi-close</v-icon>
					</v-btn>
					<v-toolbar-title>案件号：{{ info.billNum }}</v-toolbar-title>
					<v-spacer></v-spacer>
				</v-toolbar>
				<div class="content">
					<div :class="['top', showImg ? 'height_60' : 'height_auto']">
						<div :class="['left', showImg ? 'width_60' : 'width_100']">
							<div class="tab">
								<v-tabs v-model="tab">
									<v-tab>基础信息</v-tab>
									<v-tab>受益人信息</v-tab>
									<v-tab>账单信息</v-tab>
									<v-tab>出险信息</v-tab>
									<v-tab>理算信息</v-tab>
								</v-tabs>
							</div>
							<div class="forms">
								<basic-info :BasicInfo="basicInfo" v-if="tab === 0"></basic-info>
								<benefit-info
									:BenefitInfo="benefitInfo"
									v-else-if="tab === 1"
								></benefit-info>
								<bill-info :BillInfo="billInfo" v-else-if="tab === 2"></bill-info>
								<insurance-info
									:InsuranceInfo="insuranceInfo"
									v-else-if="tab === 3"
								></insurance-info>
								<count-info :CountInfo="countInfo" v-else></count-info>
							</div>
						</div>
						<div class="right" v-if="showImg">
							<div class="img">
								<watch-image :src="selectedImage"></watch-image>
							</div>
							<div class="page">
								<v-pagination
									v-model="imgPage"
									circle
									:length="imgPages"
									:total-visible="10"
								></v-pagination>
							</div>
						</div>
					</div>
					<div class="bottom" v-if="showImg">
						<h1 class="error">错误信息</h1>
						<vxe-table
							ref="xTable"
							:data="desserts"
							border
							stripe
							max-height="190"
							class="mytable-scrollbar"
							:edit-config="{ trigger: 'dblclick', mode: 'cell' }"
							@edit-closed="editClosedEvent"
						>
							<vxe-column type="seq" width="60" title="序号"></vxe-column>
							<template v-for="item in cells.headers">
								<!-- 时间 BEGIN -->
								<vxe-column
									v-if="item.value === 'CreatedAt'"
									:field="item.value"
									:title="item.text"
									:key="item.value"
									:width="item.width"
									:sortable="item.sortable"
									:fixed="item.fixed"
								>
									<template #default="{ row }">
										{{ row.CreatedAt | dateFormat("YYYY-MM-DD") }}
									</template>
								</vxe-column>
								<!-- 时间 END -->

								<vxe-column
									v-else
									:field="item.value"
									:fixed="item.fixed"
									:title="item.text"
									:key="item.value"
									:width="item.width"
									:sortable="item.sortable"
									:edit-render="{ autofocus: '.vxe-input--inner' }"
								>
									<template #edit="{ row }">
										<vxe-input
											v-model="row[item.value]"
											type="text"
										></vxe-input>
									</template>
								</vxe-column>
							</template>
						</vxe-table>
						<div style="margin-left: 30px; margin-top: 30px">
							<v-btn
								depressed
								color="primary"
								style="margin-right: 25px"
								@click="onSave"
							>
								保存数据
							</v-btn>
							<v-btn depressed> 返回 </v-btn>
						</div>
					</div>
				</div>
			</v-card>
		</v-dialog>
	</v-row>
</template>

<script>
import cells from "./cells";
import { sessionStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import _ from "lodash";

const tabMap = new Map([
	[0, "basicInfo"],
	[1, "benefitInfo"],
	[2, "billInfo"],
	[3, "insuranceInfo"],
	[4, "countInfo"]
]);

export default {
	props: {
		checkDetail: {
			type: Object,
			default: () => {}
		}
	},
	data() {
		return {
			dialog: false,
			notifications: false,
			sound: true,
			widgets: false,
			tab: "",
			info: {},
			cells,
			desserts: [
				{
					errorType: "",
					errorContent: ""
				},
				{
					errorType: "",
					errorContent: ""
				}
			],
			images: [],
			imgPage: 1,
			imgPages: 20,
			showImg: true,
			activeIndex: "",
			selectedImage: "",

			basicInfo: {},
			benefitInfo: {},
			billInfo: {},
			insuranceInfo: {},
			countInfo: {}
		};
	},
	created() {
		this.$EventBus.$on("submitData", async forms => {
			let content = sessionStorage.get("checkForm");
			let tabArr = ["basicInfo", "benefitInfo", "billInfo", "insuranceInfo", "countInfo"];
			for (let el of tabArr) {
				const result = await this.$store.dispatch("UPDATE_INSPECT_INFO", {
					data: JSON.stringify(content),
					id: this.checkDetail.ID,
					type: el
				});
				if (result.code == 200) this.toasted.dynamic(result.msg, result.code);
			}
		});
	},
	methods: {
		async onOpen() {
			this.dialog = true;
		},
		async setViewImageWidth() {
			this.images = _.cloneDeep(sessionStorage.get("thumbs") || []);
			this.imgPages = this.images.length;
			let reg = new RegExp("/files/files/", "g");
			let convert = new RegExp("/convert_", "g");

			this.images.forEach(async (image, index) => {
				image.path = image.path.replace(reg, "/files/");
				image.path = image.path.replace(convert, "/");
				const newBase64 = await lpTools.getTokenImg(image.path);
				if (newBase64) {
					lpTools.getBase64(newBase64).then(base64String => {
						this.$set(this.images, index, {
							thumbPath: image.thumbPath,
							newThumbPath: base64String,
							path: image.path
						});
					});
				}
			});

			try {
				const selectedImage = await lpTools.getTokenImg(this.images[0].path);
				if (selectedImage) {
					lpTools.getBase64(selectedImage).then(base64String => {
						this.selectedImage = base64String;
					});
				}
			} catch (e) {
				console.log(e, "eeee");
			}
		},
		editClosedEvent() {
			const $table = this.$refs.xTable;
			let content = sessionStorage.get("checkForm");
			content.basicInfo.errorInfo = $table.getData();
			sessionStorage.set("checkForm", content);
		},
		async onSave() {
			let content = sessionStorage.get("checkForm");
			let tabArr = ["basicInfo", "benefitInfo", "billInfo", "insuranceInfo", "countInfo"];
			for (let el of tabArr) {
				const result = await this.$store.dispatch("UPDATE_INSPECT_INFO", {
					data: JSON.stringify(content),
					id: this.checkDetail.ID,
					type: el
				});
				if (result.code == 200) this.toasted.dynamic(result.msg, result.code);
			}

			let type = tabMap.get(newValue);
			const result = await this.$store.dispatch("GET_INSPECT_INFO", {
				id: this.checkDetail.ID,
				type
			});
			if (result.code == 200) {
				this[type] = typeData[type];
			}
		}
	},

	watch: {
		checkDetail: function (newValue) {
			this.setViewImageWidth();
			this.info = newValue;
		},
		tab: {
			async handler(newValue) {
				if (newValue == 5) return;
				let type = tabMap.get(newValue);
				const result = await this.$store.dispatch("GET_INSPECT_INFO", {
					id: this.checkDetail.ID,
					type
				});
				if (result.code == 200) this.toasted.dynamic(result.msg, result.code);
				let typeData = {};
				if (result.data == "") {
					this.basicInfo = {};
					this.benefitInfo = {};
					this.billInfo = {};
					this.insuranceInfo = {};
					this.countInfo = {};
				} else {
					typeData = JSON.parse(result.data);
					sessionStorage.set("checkForm", typeData);
					this[type] = typeData[type];
					if (typeData.basicInfo.errorInfo && typeData.basicInfo.errorInfo.length != 0) {
						this.desserts = typeData.basicInfo.errorInfo;
					}
				}
				this.countInfo = typeData["countInfo"];
				this.insuranceInfo = typeData["insuranceInfo"];
				this.showImg = newValue === 4 ? false : true;
			}
		},
		dialog: function (newValue) {
			this.tab = 5;
		},
		imgPage: function (newValue) {
			this.activeIndex = newValue - 1;
			this.selectedImage = this.images[this.activeIndex].newThumbPath;
		}
	},
	destroy() {
		this.$EventBus.$off("submitData");
	},
	components: {
		"watch-image": () => import("./watchImage"),
		"basic-info": () => import("./basicInfo"),
		"benefit-info": () => import("./benefitInfo"),
		"bill-info": () => import("./billInfo"),
		"insurance-info": () => import("./insuranceInfo"),
		"count-info": () => import("./countInfo")
	}
};
</script>

<style lang="scss">
.v-toolbar__title {
	font-size: 18px !important;
}
.width_100 {
	width: 100% !important;
}
.width_60 {
	width: 60%;
}
.height_60 {
	height: 60%;
}
.height_auto {
	height: auto !important;
}

.card {
	width: 100%;
	.content {
		width: 100%;
		height: calc(100vh - 64px);
		// overflow-y: scroll;
		.top {
			display: flex;
			justify-content: space-between;
			background-color: #f5f6fb;
			height: 60%;
			.left {
				display: flex;
				flex-direction: column;
				justify-content: flex-start;
				background-color: white;
				width: 60%;
				.forms {
					background-color: #f5f6fb;
					height: 100%;
					overflow-y: scroll;
					padding-bottom: 20px;
				}
			}
			.right {
				width: 40%;
				.img {
					height: 90%;
				}
				.page {
					height: 10%;
					display: flex;
					justify-content: center;
				}
			}
		}
		.bottom {
			height: 40%;
			box-sizing: border-box;
			padding-top: 8px;
			.error {
				height: 40px;
				color: #007aff;
				font-weight: bolder;
				font-size: 15px;
				line-height: 40px;
				border-radius: 10px;
				padding-left: 25px;
				margin-top: 15px;
				background-color: #e9f6ff !important;
			}
		}
	}
}

/*滚动条整体部分*/
.mytable-scrollbar ::-webkit-scrollbar {
	width: 10px;
	height: 10px;
}
/*滚动条的轨道*/
.mytable-scrollbar ::-webkit-scrollbar-track {
	background-color: #ffffff;
}
/*滚动条里面的小方块，能向上向下移动*/
.mytable-scrollbar ::-webkit-scrollbar-thumb {
	background-color: #bfbfbf;
	border-radius: 5px;
	border: 1px solid #f1f1f1;
	box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
}
.mytable-scrollbar ::-webkit-scrollbar-thumb:hover {
	background-color: #a8a8a8;
	cursor: pointer;
}
.mytable-scrollbar ::-webkit-scrollbar-thumb:active {
	background-color: #787878;
}
/*边角，即两个滚动条的交汇处*/
.mytable-scrollbar ::-webkit-scrollbar-corner {
	background-color: #ffffff;
}

/* 修改滚动条的宽度和颜色 */
::-webkit-scrollbar {
	width: 10px;
	height: 10px;
	background-color: #f5f5f5;
}

/* 修改滚动条上下箭头的样式 */
/* ::-webkit-scrollbar-button {
  border-radius: 5px;
  width: 10px;
  height: 10px;
  background-color: #ccc;
} */

/* 修改滚动条轨道的样式 */
::-webkit-scrollbar-track {
	border-radius: 5px;
	background-color: #ccc;
}

/* 修改滚动条滑块的样式 */
::-webkit-scrollbar-thumb {
	background-color: #888888;
	border-radius: 5px;
}

/* 修改滚动条滑块的悬停样式 */
::-webkit-scrollbar-thumb:hover {
	background-color: #888;
}
</style>
