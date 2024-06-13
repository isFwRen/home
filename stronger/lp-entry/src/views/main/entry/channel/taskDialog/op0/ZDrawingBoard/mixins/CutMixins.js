import { fabric } from "fabric";
import { cutRectStrokeWidth, cornerSize } from "../libs/constants";
import { RabbitLegacy } from "crypto-js";
import moment from 'moment'

export default {
	methods: {
		// 创建剪裁方框
		createCutRect(pointer) {
			if (JSON.stringify(this.downPoint) === JSON.stringify(pointer)) return;

			const activeObject = this.canvas.getActiveObject();

			if (activeObject) {
				return;
			}
			// console.log("创建剪裁方框");
			let top = Math.min(this.downPoint.y, pointer.y);
			let left = Math.min(this.downPoint.x, pointer.x);
			let width = Math.abs(this.downPoint.x - pointer.x);
			let height = Math.abs(this.downPoint.y - pointer.y);

			const cutRect = new fabric.Rect({
				top,
				left,
				width,
				height,
				fill: "rgba(0, 0, 0, 0)",
				stroke: "rgba(25, 118, 210)",
				strokeWidth: cutRectStrokeWidth / this.imageScale,
				// strokeDashArray: [20, 20], // 设置虚线样式
				// strokeWidth: 15,

				hasRotatingPoint: false,
				hasBorders: false,
				lockRotation: true,
				cornerColor: "#2080f0",
				// cornerStyle: "circle",
				cornerSize: cornerSize / this.imageScale,
				// cornerSize: 35,
				transparentCorners: false,
				strokeUniform: true
			});

			cutRect.type = "cut";
			cutRect.unique = Date.now();

			cutRect.on("selected", () => {
				// 修复Bug 使截图框唯一 只存在一个截图框
				// this.activeIndex = cutRect.unique;
				this.activeIndex = -1;
			});

			cutRect.bringToFront();
			cutRect.setControlsVisibility({ mtr: false });

			this.canvas.add(cutRect).setActiveObject(cutRect);

			this.ctxList.push(cutRect);

			this.cutArea = {
				x: this.downPoint.x,
				y: this.downPoint.y,
				width,
				height
			};

			this.$emit("cut", this.cutArea);

			this.downPoint = null;
		}
	},
	watch: {
		cutArea: {
			handler(newValue) {
				let x = newValue.x
				let y = newValue.y
				let width = newValue.width
				let height = newValue.height
				if (x) {
					let value = '按键:切图' + ',' + moment(new Date).format('YYYY/MM/DD HH:mm:ss') + ',' + '第' + this.$store.state['recordKey'].page + '页' + ',' + '起始坐标点:' + `(${x}, ${y})` + ',' + '长:' + width + '高:' + height
					this.$store.commit('UPDATE_KEY', value)
				}
			},
		}
	}
};
