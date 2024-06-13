<template>
	<v-card>
		<v-toolbar class="z-toolbar" color="primary" dark>
			<z-btn icon dark @click="onBack">
				<v-icon>mdi-arrow-left</v-icon>
			</z-btn>

			<v-spacer></v-spacer>

			<v-toolbar-items v-if="imageList.length > 1">
				<div class="align-center z-flex">
					<template v-for="item of cells.options">
						<lp-tooltip-btn
							bottom
							:btnClass="item.class"
							fab
							:icon="item.icon"
							outlined
							small
							:tip="item.tip"
							:key="item.value"
							@click="onChange(item)"
						>
						</lp-tooltip-btn>
					</template>
				</div>
			</v-toolbar-items>

			<v-spacer></v-spacer>

			<z-btn icon dark @click="onClose">
				<v-icon>mdi-close</v-icon>
			</z-btn>
		</v-toolbar>

		<v-card-text class="pt-16 case-image-card__text" v-if="imageList.length">
			<v-row class="pt-4">
				<v-col :cols="3">
					<template v-for="(item, index) of imageList">
						<z-btn
							:key="item.ID"
							class="mr-2 mb-2"
							color="primary"
							outlined
							small
							@click="showIndexPicture(index)"
						>
							{{ item.name }}
						</z-btn>
					</template>
				</v-col>

				<v-divider vertical></v-divider>

				<v-col class="z-flex justify-center" :cols="9">
					<div class="img-wrapper" :style="myStyle">
						<!-- <z-image class="elevation-1 preview-img" :src="`${baseURLApi}${imageList[imageIndex].picture}`" /> -->
						<lp-images
							class="elevation-1 preview-img"
							:src="`${imageList[imageIndex].base64}`"
						/>
					</div>
				</v-col>
			</v-row>
		</v-card-text>
	</v-card>
</template>

<script>
import { mapGetters } from "vuex";
import { tools as lpTools } from "@/libs/util";
import CaseMixins from "../../CaseMixins";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "ViewImage",
	mixins: [CaseMixins],

	data() {
		return {
			formId: "ViewImage",
			cells,
			rotate: 0,
			imageList: [],
			baseURLApi,
			imageIndex: 0,
			deg: 0,
			scale: 1,
			myStyle: "transform: rotate( 0 deg) scale(1);"
		};
	},

	created() {
		this.getCaseFieldImage();
	},

	// mounted() {
	//   document.addEventListener('keydown', this.watchKeydown)
	// },

	// beforeDestroy() {
	//   document.removeEventListener('keydown', this.watchKeydown)
	// },

	methods: {
		// 查看图片
		async getCaseFieldImage() {
			this.imageList = [];

			const { params, query } = this.$route;
			const body = {
				proCode: this.cases.caseInfo.proCode || window.sessionStorage.getItem("proCode"),
				blockId: params.blockId,
				fieldId: query.fieldId
			};

			const result = await this.$store.dispatch("GET_CASE_FIELD_IMAGE", body);

			let reg = new RegExp("files/files/", "g");
			let convert = new RegExp("/convert_", "g");

			if (result.code === 200) {
				result.data.forEach(async ele => {
					ele.picture = ele.picture.replace(reg, "files/");
					ele.picture = ele.picture.replace(convert, "/");
					ele.base64 = "";
					const newBase64 = await lpTools.getTokenImg(`${baseURLApi}${ele.picture}`);
					if (newBase64) {
						lpTools.getBase64(newBase64).then(base64String => {
							ele.base64 = base64String;
						});
					}
				});

				this.imageList = result.data;
			} else {
				this.toasted.error(result.msg);
			}
		},

		onChange({ value }) {
			switch (value) {
				case "switch":
					this.$router.push({ path: "/main/PM/case/image-thumb" });
					break;
				case "prev":
					this.imageIndex =
						this.imageIndex > 0 ? this.imageIndex - 1 : this.imageList.length - 1;
					break;
				case "left":
					this.deg = (this.deg + 90) % 360;
					break;
				case "right":
					this.deg = (this.deg - 90) % 360;
					break;
				default:
					this.imageIndex =
						this.imageIndex < this.imageList.length - 1 ? this.imageIndex + 1 : 0;
					break;
			}
		},

		watchKeydown(e) {
			var key = window.event ? e.keyCode : e.which;
			const keys = [112, 113, 76, 82];
			if (keys.indexOf(key) != -1) {
				window.event.preventDefault();
			}
			if (key === 112) {
				this.imageIndex =
					this.imageIndex < this.imageList.length - 1 ? this.imageIndex + 1 : 0;
			} else if (key === 113) {
				this.imageIndex =
					this.imageIndex > 0 ? this.imageIndex - 1 : this.imageList.length - 1;
			} else if (key === 76 && e.ctrlKey) {
				this.deg = (this.deg + 90) % 360;
			} else if (key === 82 && e.ctrlKey) {
				this.deg = (this.deg - 90) % 360;
			}
		},

		showIndexPicture(index) {
			this.imageIndex = index - 1;
		}
	},

	watch: {
		deg(val) {
			this.myStyle = "transform: rotate(" + this.deg + "deg) scale(" + this.scale + ");";
		}
	},

	computed: {
		...mapGetters(["cases"])
	},
	components: {
		"lp-images": () => import("@/components/lp-images")
	}
};
</script>

<style lang="scss">
.case-image-card__text {
	height: calc(100vh - 64px);
	overflow: hidden;
}

.img-wrapper {
	width: 100%;
	height: calc(100vh - 88px);
	text-align: center;

	.preview-img {
		max-width: 70%;
		max-height: 100%;
		vertical-align: middle;
	}
}
</style>
