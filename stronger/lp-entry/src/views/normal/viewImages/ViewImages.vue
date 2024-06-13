<template>
	<div class="view-images">
		<v-card slot="card">
			<v-toolbar class="z-toolbar" color="primary" dark>
				<v-spacer></v-spacer>

				<v-toolbar-items>
					<div class="align-center z-flex">
						<template v-for="item in cells.toolbarItems">
							<lp-tooltip-btn bottom :btnClass="item.class" fab :icon="item.icon" outlined small :tip="item.tip"
								:key="item.value" @click="switchAction(item)">
							</lp-tooltip-btn>
						</template>
					</div>
				</v-toolbar-items>

				<v-spacer></v-spacer>
			</v-toolbar>

			<v-card-text class="pt-16 case-image-card__text" v-if="images.length">
				<v-row class="thumbs-images-box pt-4">
					<v-col class="thumbs-box" cols="3">
						<v-row dense>
							<v-col v-for="(image, index) in images" :key="`view_images_${index}`" cols="4">
								<v-badge color="primary" :content="index + 1" offset-x="-6" offset-y="32">
								</v-badge>
								<v-lazy v-model="isActive" :options="{
									threshold: 0.5
								}" min-height="200" transition="fade-transition">
									<div :class="[
										'thumb-box',
										activedIndex === index ? 'actived' : ''
									]" @click="switchThumb(image, index)">
										<v-img max-height="160" :src="image.path"></v-img>
									</div>
								</v-lazy>
							</v-col>
						</v-row>
					</v-col>

					<v-divider vertical></v-divider>

					<v-col id="imagesBox" class="images-box">
						<ZImages ref="image" col-align="start" :min-zoom-out="0.5" :proportion="0.75" :src="selectedImage"></ZImages>
						<p>第 {{ activedIndex + 1 }} 张图</p>
					</v-col>
				</v-row>
			</v-card-text>
		</v-card>
	</div>
</template>

<script>
import { sessionStorage, tools } from "vue-rocket";
import * as cells from "./cells";
import axios from "axios";
import { localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
export default {
	name: "LPViewImages",

	data() {
		return {
			cells,
			activedIndex: 0,
			images: [],
			// 图片请求头
			instance: "",
			selectedImage: "",
			// 缩略图路径
			readerPath: [],
			isActive: false
		};
	},

	created() {
		const token = localStorage.get("token");
		const user = localStorage.get("user");

		this.instance = axios.create({
			headers: {
				"x-token": token,
				"x-user-id": user.id
			}
		});
	},

	mounted() {
		this.setViewImageWidth();
		window.addEventListener("keydown", this.fuckShortcut);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckShortcut);
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
					this.$refs.image.eventRotateLeft();
					break;

				case "right":
					this.$refs.image.eventRotateRight();
					break;
			}
		},

		// 上一张
		prevAction() {
			this.activedIndex -= 1;

			if (this.activedIndex <= 0) {
				this.activedIndex = 0;
			}

			this.selectedImage = this.images[this.activedIndex].path;
		},

		// 下一张
		nextAction() {
			console.log(this.images, "ss");
			const lastIndex = this.images.length - 1;
			this.activedIndex += 1;

			if (this.activedIndex >= lastIndex) {
				this.activedIndex = lastIndex;
			}

			this.selectedImage = this.images[this.activedIndex].path;
		},

		// 切换缩略图
		switchThumb(image, index) {
			this.activedIndex = index;
			this.selectedImage = image.path;
		},

		fuckShortcut(event) {
			const { ctrlKey, keyCode } = event || window.event;

			switch (keyCode) {
				// 左上右下
				case 37:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventMoveRight();
					}
					break;

				case 38:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventMoveBottom();
					}
					break;

				case 39:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventMoveLeft();
					}
					break;

				case 40:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventMoveTop();
					}
					break;

				// 图片缩放
				case 73:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventZoomOut();
					}
					break;

				case 79:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventZoomIn();
					}
					break;

				case 90:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventZoomOrigin();
					}
					break;

				// 图片旋转
				case 76:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventRotateLeft();
					}
					break;
				case 82:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.image.eventRotateRight();
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
		setViewImageWidth() {
			this.images = tools.deepClone(sessionStorage.get("thumbs"));
			let reg = new RegExp("/files/files/", "g");
			let convert = new RegExp("/convert_", "g");

			this.images.forEach(async (image, index) => {
				image.path = image.path.replace(reg, "/files/");
				image.path = image.path.replace(convert, "/");

				let item = await this.transform(image.path);

				this.getReader(item).then(res => {
					this.images[index].path = res;
				});
			});
			// console.log(this.images);
			// this.selectedImage = this.images[this.activedIndex].path;
			setTimeout(() => {
				this.selectedImage = this.images[this.activedIndex].path;
			}, 200);
		},

		// 图片转Base64格式
		async transform(el) {
			let code = "";
			const secret = localStorage.get("secret") || "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			let res = await this.instance.get(el, {
				responseType: "blob",
				headers: {
					"x-code": String(code)
				}
			});
			return res.data;
		},

		// 读取图片文件
		getReader(blob) {
			return new Promise((resolve, reject) => {
				const reader = new FileReader();
				reader.onloadend = () => {
					const base64String = reader.result;
					resolve(base64String);
				};
				reader.onerror = reject;
				reader.readAsDataURL(blob);
			});
		}
	},

	components: {
		ZImages: () =>
			import("../../main/entry/channel/taskDialog/components/watchImage/ZImage/ZImage.vue")
	}
};
</script>

<style scoped lang="scss">
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

		.thumb-box {
			max-height: 160px;
			border: 2px solid transparent;
			overflow: hidden;

			&.actived {
				border: 2px solid #1976d2;
			}
		}
	}
}
</style>
