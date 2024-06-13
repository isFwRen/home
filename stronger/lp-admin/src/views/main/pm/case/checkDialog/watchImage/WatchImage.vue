<template>
	<div id="watchImage" class="watch-image">
		<ZImages ref="image" col-align="start" :min-zoom-out="0.5" :proportion="proportion" :src="src"></ZImages>
		<!-- v-if="isLoop" -->
		<div @click="show" class="rect" v-show="isLoop">框图 [YES / NO]</div>
		<div class="no">控制上下移动 (Shift + ↑ / Shift + ↓)</div>
		<canvas id="drawer" v-show="isLoop"></canvas>
	</div>
</template>

<script>
import { tools, sessionStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import { Drawer, DrawHelper, Polygon, DragableAndScalableRect } from "./rect";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "WatchImage",

	props: {
		bill: {
			type: Object,
			default: () => ({})
		},

		op: {
			validator(value) {
				return ["op0", "op1", "op2", "opq"].indexOf(value) !== -1;
			},
			required: false
		},

		src: {
			type: String,
			required: true
		},
		isLoop: {
			type: Boolean,
			default: false
		}
	},

	data() {
		return {
			width: 0,
			proportion: 0.75,
			react: "",
			content: "",
			drawer: "",
			rect1: ""
		};
	},

	watch: {
		$route: {
			handler({ query }) {
				switch (query.proCode) {
					case "B0108":
						this.proportion = 1;
						break;

					default:
						this.proportion = 0.75;
						break;
				}
			},
			immediate: true
		},

		"bill.pictures": {
			handler(pictures) {
				if (this.op === "opq") {
					const thumbs = [];

					if (tools.isYummy(pictures)) {
						pictures.map(picture => {
							thumbs.push({
								thumbPath: `${baseURLApi}files/${this.bill.downloadPath}A${picture}`,
								path: `${baseURLApi}files/${this.bill.downloadPath}${picture}`
							});
						});
					}

					sessionStorage.set("thumbs", thumbs);
				}
			},
			immediate: true
		}
	},

	mounted() {
		const watchImage = document.getElementById("watchImage");
		this.width = watchImage.offsetWidth;
		window.addEventListener("keydown", this.fuckOp1Op2OpqShortcut);
		this.react = document.querySelector("#drawer");
		this.react.style.right = 10000 + "px";
		document.onkeydown = function (event) {
			var t = document.querySelector("#drawer").offsetTop;
			var e = event || window.event || arguments.callee.caller.arguments[0];
			if (e.shiftKey && e.keyCode == 38) {
				// 按shift
				t -= 5;
			}
			if (e.shiftKey && e.keyCode == 40) {
				// 按shift
				t += 5;
			}
			document.querySelector("#drawer").style.top = t + "px";
		};
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckOp1Op2OpqShortcut);
	},

	methods: {
		fuckOp1Op2OpqShortcut(event) {
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
			}
		},

		navToViewImages() {
			window.open(
				`${location.origin}/normal/view-images`,
				"_blank",
				"width=" +
					(window.screen.availWidth - 10) +
					",height=" +
					(window.screen.availHeight - 30) +
					",top=0,left=0,toolbar=no,menubar=no,scrollbars=no, resizable=no,location=no, status=no"
			);
		},

		show() {
			this.content = document.querySelector(".rect");
			if (this.react.style.right === "10000px") {
				this.drawer = new Drawer("#drawer");
				this.rect1 = new DragableAndScalableRect({
					x: 200, // 中心点 x 坐标
					y: 200, // 中心点 y 坐标
					width: 200,
					height: 200,
					minWidth: 20,
					minHeight: 20,
					cornerWidth: 20
				});
				this.drawer.addPolygon(this.rect1);
				console.log(11111);
				this.react.style.right = 0;
				this.react.style.top = 0;
				this.content.style.color = "#348940";
			} else {
				console.log(222222);
				this.drawer = "";
				this.rect1 = "";
				this.react.style.right = 10000 + "px";
				this.content.style.color = "white";
			}
		}
	},

	components: {
		ZImages: () => import("./ZImage/ZImage.vue")
	}
};
</script>

<style scoped lang="scss">
.watch-image {
	width: 100%;
	height: 100%;
	position: relative;

	.rect {
		color: white;
		position: absolute;
		top: 0;
		right: 140px;
		font-weight: 600;
		width: 120px;
		text-align: center;
		font-size: 13px;
		line-height: 36px;
		cursor: pointer;
		z-index: 100000;
	}
	.no {
		width: 260px;
		height: 28px;
		font-size: 15px;
		line-height: 25px;
		text-align: center;
		top: 40px;
		right: 50px;
		border-radius: 3px;
		position: absolute;
		background-color: rgba(0, 0, 0, 0.5);
		color: white;
		display: none;
	}

	.rect:hover + .no {
		display: block;
	}
	#drawer {
		width: 100%;
		height: 100%;
		position: absolute;
		right: 10000px;
		top: 0px;
	}
}
</style>