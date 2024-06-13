<template>
	<div id="watchImage" class="watch-image">
		<z-image ref="zImage" :src="src" :width="width"></z-image>
	</div>
</template>

<script>
import { tools, sessionStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

const { baseURLApi } = lpTools.baseURL();

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
		}
	},

	data() {
		return {
			width: 0
		};
	},

	mounted() {
		const watchImage = document.getElementById("watchImage");
		this.width = watchImage.offsetWidth;

		window.addEventListener("keydown", this.fuckOp1Op2OpqShortcut);
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
			}
		},

		navToViewImages() {
			window.open(
				`${location.origin}/normal/view-images`,
				"_blank",
				"toolbar=yes, scrollbars=yes, resizable=yes"
			);
		}
	},

	watch: {
		"bill.pictures": {
			handler(pictures) {
				if (this.op === "opq") {
					const thumbs = [];

					if (tools.isYummy(pictures)) {
						pictures.map(picture => {
							thumbs.push({
								thumbPath: `${baseURLApi}files/${this.bill.downloadPath}${picture}`,
								path: `${baseURLApi}files/${this.bill.downloadPath}${picture}`
							});
						});
					}

					sessionStorage.set("thumbs", thumbs);
				}
			},
			immediate: true
		}
	}
};
</script>

<style scoped lang="scss">
.watch-image {
	width: 100%;
	height: 100%;
}
</style>
