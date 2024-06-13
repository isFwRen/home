import { tools } from "vue-rocket";
import { op0ThumbIdPrefix } from "../libs/constants";
import moment from 'moment'

export default {
	data() {
		return {
			ctVisible: true
		};
	},

	mounted() {
		window.addEventListener("keydown", this.fuckOp0Shortcut);
	},

	beforeDestroy() {
		window.removeEventListener("keydown", this.fuckOp0Shortcut);
	},

	methods: {
		async fuckOp0Shortcut(event) {
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
					ctrlKey && this.$refs.drawImage.eventMoveRight();
					break;

				// 上
				case 38:
					ctrlKey && this.$refs.drawImage.eventMoveBottom();
					break;

				// 右
				case 39:
					ctrlKey && this.$refs.drawImage.eventMoveLeft();
					break;

				// 下
				case 40:
					ctrlKey && this.$refs.drawImage.eventMoveTop();
					break;

				// 临时保存(ctrl + s)
				case 83:
					if (ctrlKey) {
						event.preventDefault();
						const valid = this.$refs.op0Form.validate();
						valid && this.$refs.drawImage.eventDone();
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
					if (this.op != 'op0') {
						this.getTask({ status: "modify", prevNums: ++this.prevNums });
					}
					break;

				// 临时保存(F4)
				case 115:
					event.preventDefault();
					const flag = await this.svDisableFieldsFirst()
					if (flag == false) break
					this.$refs.drawImage.eventDone();
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

			if (tools.isYummy(pictures)) {
				const lastIndex = this.bill.pictures?.length - 1;
				let id = null;

				switch (keyCode) {
					// 向前一页
					case 112:
						let prevIndex = this.newfieldsObjectArray[this.imgindex - 1] ?? this.thumbIndex - 1;
						if (this.imgindex == 0 && this.newfieldsObjectArray[0]) {
							this.imgindex = this.bill.pictures.length;
							prevIndex = this.newfieldsObjectArray[this.imgindex - 1];
						}
						if (prevIndex < 0) {
							prevIndex = lastIndex;
						}

						id = `${op0ThumbIdPrefix}${prevIndex}`;
						this.onSelectThumb({ id, thumbIndex: prevIndex, imgindex: this.imgindex - 1 });
						break;

					// 向后一页
					case 113:
						let nextIndex = this.newfieldsObjectArray[this.imgindex + 1] ?? this.thumbIndex + 1;
						if (this.imgindex + 1 == this.bill.pictures.length && this.newfieldsObjectArray[0]) {
							this.imgindex = -1;
							nextIndex = this.newfieldsObjectArray[0];
						}

						if (nextIndex > lastIndex) {
							nextIndex = 0;
						}

						id = `${op0ThumbIdPrefix}${nextIndex}`;
						this.onSelectThumb({ id, thumbIndex: nextIndex, imgindex: this.imgindex + 1 });
						break;

					// 跳回第一页
					case 33:
						// id = `${op0ThumbIdPrefix}${0}`;
						// this.onSelectThumb({ id, thumbIndex: this.newfieldsObjectArray[0] || 0 });
						this.autofocusToTopField()
						break;

					// 跳回最后一页
					case 34:
						// id = `${op0ThumbIdPrefix}${lastIndex}`;
						// this.onSelectThumb({
						// 	id,
						// 	thumbIndex: this.newfieldsObjectArray[this.newfieldsObjectArray.length] || lastIndex
						// });
						this.autofocusToBottomField()
						break;
				}
			}
		},

		// 图片缩放
		imageScales(keyCode) {
			switch (keyCode) {
				// i(缩小)
				case 73:
					this.$refs.drawImage.eventZoomOut();
					this.$store.commit('UPDATE_KEY', '按键:缩小' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
					break;

				// o(放大)
				case 79:
					this.$refs.drawImage.eventZoomIn();
					this.$store.commit('UPDATE_KEY', '按键:放大' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
					break;

				// z(还原)
				case 90:
					this.$refs.drawImage.eventZoomOrigin();
					this.$store.commit('UPDATE_KEY', '按键:还原' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页')
					break;
			}
		},
	}
};
