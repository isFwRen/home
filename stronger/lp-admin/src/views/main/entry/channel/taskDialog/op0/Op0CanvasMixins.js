let drawAction = "drawScreenshot";

export default {
	mounted() {
		window.addEventListener("keydown", this.canvasFuckEvents);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.canvasFuckEvents);
	},

	methods: {
		initializedImage() {
			this.$refs.drawImage[drawAction]();
		},

		// 选择绘图方式
		onSelect({ key }) {
			switch (key) {
				case "screenshot":
					drawAction = "drawScreenshot";
					this.$refs.drawImage.drawScreenshot();
					break;

				case "rectangle":
					drawAction = "drawRectangle";
					this.$refs.drawImage.drawRectangle();
					break;

				case "text":
					this.$refs.drawImage.drawText();
					break;

				case "rotateL":
					this.$refs.drawImage.rotateImage("left");
					break;

				case "rotateR":
					this.$refs.drawImage.rotateImage("right");
					break;

				case "clear":
					this.$refs.drawImage.clear();
					break;

				case "save":
					this.$refs.drawImage.onSave();
					break;
			}
		},

		// 通过按键选择绘图方式或者提交
		canvasFuckEvents(event) {
			event = event || window.event;
			const { altKey, ctrlKey, keyCode } = event;

			switch (keyCode) {
				// 切图(A)
				case 65:
					if (ctrlKey) {
						event.preventDefault();
						drawAction = "drawScreenshot";
						this.$refs.drawImage.drawScreenshot();
					}
					break;

				// 矩形(X)
				case 88:
					if (ctrlKey) {
						event.preventDefault();
						drawAction = "drawRectangle";
						this.$refs.drawImage.drawRectangle();
					}
					break;

				// 文字(E)
				case 69:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.drawImage.drawText();
					}
					break;

				// 左旋转(L)
				case 76:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.drawImage.rotateImage("left");
					}
					break;

				// 右旋转(R)
				case 82:
					if (ctrlKey) {
						event.preventDefault();
						this.$refs.drawImage.rotateImage("right");
					}
					break;

				// 清除(Z)
				case 90:
					if (ctrlKey && !altKey) {
						event.preventDefault();
						this.$refs.drawImage.clear();
					} else if (ctrlKey && altKey) {
						event.preventDefault();
						this.screenshotOrRectangle();
					}
					break;
			}
		},

		// 设置画布宽高
		setCanvasSize() {
			this.$nextTick(() => {
				const drawImageContain = this.$refs.drawImageContain;
				this.canvasWidth = drawImageContain.offsetWidth - 30;

				window.addEventListener("resize", function () {
					this.canvasWidth = drawImageContain.offsetWidth - 30;
				});
			});
		}
	},

	watch: {
		hasOp: {
			handler(hasOp) {
				if (hasOp) {
					this.setCanvasSize();
				}
			},
			immediate: true
		}
	}
};
