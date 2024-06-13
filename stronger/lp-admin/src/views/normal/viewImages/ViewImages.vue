<template>
	<div class="view-images">
		<v-card slot="card">
			<v-toolbar class="z-toolbar" color="primary" dark>
				<v-spacer></v-spacer>

				<v-toolbar-items>
					<div class="align-center z-flex">
						<template v-for="item in cells.toolbarItems">
							<lp-tooltip-btn
								bottom
								:btnClass="item.class"
								fab
								:icon="item.icon"
								outlined
								small
								:tip="item.tip"
								:key="item.value"
								@click="switchAction(item)"
							>
							</lp-tooltip-btn>
						</template>
					</div>
				</v-toolbar-items>

				<v-spacer></v-spacer>
			</v-toolbar>

			<v-card-text class="pt-16 case-image-card__text" v-if="images.length">
				<v-row class="thumbs-images-box pt-4">
					<v-col class="thumbs-box" cols="4">
						<v-row dense>
							<v-col
								v-for="(image, index) in images"
								:key="`view_images_${index}`"
								cols="4"
							>
								<v-badge
									color="primary"
									:content="index + 1"
									offset-x="-6"
									offset-y="32"
								>
								</v-badge>
								<v-lazy
									v-model="isActive"
									:options="{
										threshold: 0.5
									}"
									transition="fade-transition"
								>
									<!--7月14-->
									<div
										:class="[
											'thumb-box',
											activedIndex === index ? 'actived' : ''
										]"
										@click="switchThumb(image, index)"
									>
										<v-img height="160" :src="image.newThumbPath"></v-img>
										<span class="ocr" v-if="imagesType.length != 0">{{
											imagesType[index]
										}}</span>
									</div>
								</v-lazy>
							</v-col>
						</v-row>
					</v-col>

					<v-divider vertical></v-divider>

					<v-col id="imagesBox" class="images-box">
						<!-- <z-image v-if="width" ref="zImage" :src="selectedImage" :width="width"></z-image> -->

						<lp-images ref="zImage" :src="selectedImage" :width="width"></lp-images>
						<p>第 {{ activedIndex + 1 }} 张图</p>
					</v-col>
				</v-row>
			</v-card-text>
		</v-card>
	</div>
</template>

<script>
import { sessionStorage, localStorage, tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import * as cells from "./cells";
import _ from "lodash";

const { baseURL } = lpTools.baseURL();

const mapRotate = new Map([
	[76, "left"],
	[82, "right"]
]);

const mapScroll = new Map([
	[37, "left"],
	[38, "up"],
	[39, "right"],
	[40, "down"]
]);

const mapZoom = new Map([
	[73, "grow"],
	[79, "shrink"],
	[90, "origin"]
]);

export default {
	name: "LPViewImages",

	data() {
		return {
			cells,
			activedIndex: 0,
			images: [],
			selectedImage: "",
			rotateDeg: 0,
			imgStyle: {},
			isActive: false,
			width: 0,
			project: localStorage.get("project").code,
			imagesType: sessionStorage.get("imagesType") || []
		};
	},
	created() {
		this.setViewImageWidth();
	},
	mounted() {
		window.addEventListener("keydown", this.fuckShortcut);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckShortcut);
		this.images = [];
	},

	methods: {
		// 切换动作
		switchAction({ value }) {
			switch (value) {
				case "prev":
					this.prevAction();
					break;

				case "next":
					this.nextAction();
					break;

				case "left":
					this.$refs.zImage.eventRotateLeft();
					break;

				case "right":
					this.$refs.zImage.eventRotateRight();
					break;
			}
		},

		// 上一张
		async prevAction() {
			this.activedIndex -= 1;

			if (this.activedIndex <= 0) {
				this.activedIndex = 0;
			}

			this.imgStyle = {
				transform: "rotate(0deg)"
			};

			this.selectedImage = this.images[this.activedIndex].newThumbPath;
		},

		// 下一张
		async nextAction() {
			const lastIndex = this.images.length - 1;
			this.activedIndex += 1;

			if (this.activedIndex >= lastIndex) {
				this.activedIndex = lastIndex;
			}

			this.imgStyle = {
				transform: "rotate(0deg)"
			};

			this.selectedImage = this.images[this.activedIndex].newThumbPath;
		},

		// 切换缩略图
		async switchThumb(image, index) {
			this.activedIndex = index;
			this.selectedImage = image.newThumbPath;
		},

		fuckShortcut(event) {
			const { ctrlKey, keyCode } = event || window.event;

			switch (keyCode) {
				// 左上右下
				case 37:
				case 38:
				case 39:
				case 40:
					if (ctrlKey) {
						event.preventDefault();
						const direction = mapScroll.get(keyCode);
						this.$refs.zImage.scroll(direction);
					}
					break;

				// 图片缩放
				case 73:
				case 79:
				case 90:
					if (ctrlKey) {
						event.preventDefault();
						const zoom = mapZoom.get(keyCode);
						this.$refs.zImage.zoom(zoom);
					}
					break;

				// 图片旋转
				case 76:
				case 82:
					if (ctrlKey) {
						event.preventDefault();
						const direction = mapRotate.get(keyCode);
						this.$refs.zImage.rotate(direction);
					}
					break;

				// 查看图片(F12)
				case 123:
					if (this.op === "opq") {
						event.preventDefault();
						this.navToViewImages();
					}
					break;

				// 上一张
				case 112:
					event.preventDefault();
					this.prevAction();
					break;

				// 下一张
				case 113:
					event.preventDefault();
					this.nextAction();
					break;
			}
		},
		async setViewImageWidth() {
			this.images = _.cloneDeep(sessionStorage.get("thumbs") || []);
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

			this.$nextTick(() => {
				const imagesBox = document.getElementById("imagesBox");
				this.width = imagesBox.offsetWidth - 24;
			});
		}
	},
	components: {
		"lp-images": () => import("@/components/lp-images")
	}
};
</script>

<style lang="scss">
.case-image-card__text {
	height: 100vh;
	overflow: hidden;

	.thumbs-images-box {
		height: 100%;

		.v-badge,
		.v-badge__badge {
			z-index: 1;
		}
		.thumbs-box,
		.images-box {
			height: 100%;
			overflow: auto;
			position: relative;

			p {
				position: absolute;
				bottom: 0;
				text-align: center;
				width: 90%;
				font-weight: 600;
				z-index: 100;
			}
		}

		.z-image {
			& > .view {
				background-color: rgba(0, 0, 0, 0.1) !important;
			}
		}

		.thumb-box {
			height: 210px;
			max-height: 210px;
			border: 2px solid transparent;
			overflow: hidden;
			&.actived {
				border: 2px solid #1976d2;
			}
		}
	}
}
.ocr {
	display: block;
	width: 100%;
	min-width: 95px;
	text-align: center !important;
	font-size: 15px;
	font-weight: bolder;
	word-wrap: break-word;
}
</style>
