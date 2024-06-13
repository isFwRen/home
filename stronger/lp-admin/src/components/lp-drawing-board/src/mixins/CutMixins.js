import { fabric } from "fabric";
import { cutRectStrokeWidth, cornerSize } from "../utils/constants";

export default {
	methods: {
		// 创建剪裁方框
		createCutRect(pointer) {
			if (JSON.stringify(this.downPoint) === JSON.stringify(pointer)) return;

			const activeObject = this.canvas.getActiveObject();

			if (activeObject) {
				return;
			}

			let top = Math.min(this.downPoint.y, pointer.y) - this.translate.y;
			let left = Math.min(this.downPoint.x, pointer.x) - this.translate.x;
			let width = Math.abs(this.downPoint.x - pointer.x);
			let height = Math.abs(this.downPoint.y - pointer.y);

			const cutRect = new fabric.Rect({
				top,
				left,
				width,
				height,
				fill: "rgba(0, 0, 0, .3)",
				stroke: "#4caf50",
				strokeWidth: cutRectStrokeWidth / this.imageScale,

				hasRotatingPoint: false,
				hasBorders: false,
				lockRotation: true,
				cornerColor: "#4caf50",
				cornerSize: cornerSize / this.imageScale,
				transparentCorners: false,
				strokeUniform: true
			});

			cutRect.type = "cut";
			cutRect.unique = Date.now();

			cutRect.on("selected", () => {
				this.activeIndex = cutRect.unique;
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
	}
};
