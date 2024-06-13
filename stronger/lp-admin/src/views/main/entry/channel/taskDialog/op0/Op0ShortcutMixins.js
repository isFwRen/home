import { R } from "vue-rocket";
const SCALE = 0.5;

export default {
	data() {
		return {
			ctVisible: true,

			scale: 1
		};
	},

	mounted() {
		window.addEventListener("keydown", this.fuckOp0Shortcut);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckOp0Shortcut);
	},

	methods: {
		fuckOp0Shortcut(event) {
			event = event || window.event;
			const { ctrlKey, keyCode } = event;

			switch (keyCode) {
				// 跳回第一页(PgUp)
				case 33:
					event.preventDefault();
					this.computedThumbPage(33);
					break;

				// 跳回最后一页(PgDn)
				case 34:
					event.preventDefault();
					this.computedThumbPage(34);
					break;

				// 图片缩放
				case 73:
				case 79:
				case 90:
					if (ctrlKey) {
						event.preventDefault();
						this.imageScales(keyCode);
					}
					break;

				// 左
				case 37:
					ctrlKey && this.$refs.drawImage.scrollLeft();
					break;

				// 上
				case 38:
					ctrlKey && this.$refs.drawImage.scrollTop();
					break;

				// 右
				case 39:
					ctrlKey && this.$refs.drawImage.scrollRight();
					break;

				// 下
				case 40:
					ctrlKey && this.$refs.drawImage.scrollBottom();
					break;

				// 临时保存(ctrl + s)
				case 83:
					if (ctrlKey) {
						event.preventDefault();
						const valid = this.$refs.op0Form.validate();
						valid && this.$refs.drawImage.onSave();
					}
					break;

				// 向前一页(F1)
				case 112:
					event.preventDefault();
					this.computedThumbPage(112);
					break;

				// 向后一页(F2)
				case 113:
					event.preventDefault();
					this.computedThumbPage(113);
					break;

				// 返回修改(F3)
				case 114:
					event.preventDefault();
					this.getTask({ status: "modify", prevNums: ++this.prevNums });
					break;

				// 临时保存(F4)
				case 115:
					event.preventDefault();
					this.$refs.drawImage.onSave();
					break;

				// // 临时保存(F6)
				// case 117:
				//   event.preventDefault()
				//   this.setInitImage()
				//   break;

				// 隐藏、显示影像编辑功能(F7)
				case 118:
					event.preventDefault();
					this.ctVisible = !this.ctVisible;
					break;
			}
		},

		// 录入-F1、F2、PgUp、PgDn
		computedThumbPage(keyCode) {
			const pictures = this.bill.pictures;

			if (R.isYummy(pictures)) {
				const lastIndex = this.bill.pictures?.length - 1;

				switch (keyCode) {
					// 向前一页
					case 112:
						let prevIndex = this.thumbIndex - 1;
						if (prevIndex <= 0) {
							prevIndex = 0;
						}
						this.onSelectThumb({ thumbIndex: prevIndex });
						break;

					// 向后一页
					case 113:
						let nextIndex = this.thumbIndex + 1;
						if (nextIndex >= lastIndex) {
							nextIndex = lastIndex;
						}
						this.onSelectThumb({ thumbIndex: nextIndex });
						break;

					// 跳回第一页
					case 33:
						this.onSelectThumb({ thumbIndex: 0 });
						break;

					// 跳回最后一页
					case 34:
						this.onSelectThumb({ thumbIndex: lastIndex });
						break;
				}
			}
		},

		// 图片缩放
		imageScales(keyCode) {
			const img = document.getElementById(`${this.op}DrawImage`);

			switch (keyCode) {
				// i(放大)
				case 73:
					this.scale += SCALE;
					img.style.transform = `scale(${this.scale})`;
					break;

				// i(缩小)
				case 79:
					this.scale -= SCALE;
					img.style.transform = `scale(${this.scale})`;
					break;

				// z(还原)
				case 79:
					this.scale = 1;
					img.style.transform = `scale(${this.scale})`;
					break;
			}
		}
	}
};
