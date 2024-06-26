<template>
	<div class="z-drawing-board">
		<top-bar
			:isCut="isCut"
			:isRect="isRect"
			:isText="isText"
			:src="src"
			@topBarEvent="topBarEvent"
		></top-bar>

		<div
			class="view"
			id="view"
			ref="view"
			:style="{
				'justify-content': rowAlign,
				'align-items': colAlign
			}"
		>
			<canvas class="canvas" id="canvas" ref="canvas">
				The browser does not support canvas
			</canvas>

			<span v-show="size" class="size">{{ computedSize }}</span>
		</div>

		<v-overlay :opacity="0.8" :value="overlay">
			<v-progress-circular indeterminate size="64"></v-progress-circular>
		</v-overlay>
	</div>
</template>

<script>
import { fabric } from "fabric";
import CanvasMixins from "./src/mixins/CanvasMixins";
import CutMixins from "./src/mixins/CutMixins";
import RectMixins from "./src/mixins/RectMixins";
import TextboxMixins from "./src/mixins/TextboxMixins";
import EventMixins from "./src/mixins/EventMixins";
import tools from "./src/utils/tools";
import { tools as lpTools } from "@/libs/util";
import containerEvent from "./src/utils/containerEvent";
import { TopBar } from "./src";
import { moveSpace, cutRectStrokeWidth } from "./src/utils/constants";
const debounce = (() => {
	let timer = null;

	return (fn, delay = 300) => {
		if (timer) {
			clearTimeout(timer);
		}

		timer = setTimeout(() => {
			fn();
		}, delay);
	};
})();

export default {
	name: "LPDrawingBoard",
	mixins: [CanvasMixins, TextboxMixins, CutMixins, RectMixins, EventMixins],

	props: {
		angle: {
			type: [Number, String],
			default: 5
		},

		// 图像垂直方向对齐方式
		colAlign: {
			validator(value) {
				return ["start", "center", "end"].includes(value);
			},
			default: "center"
		},

		// 图像默认坐标
		coord: {
			type: Object,
			required: false
		},

		// 图像默认方向
		direction: {
			validator(value) {
				return !!~["TOP", "RIGHT", "BOTTOM", "LEFT"].indexOf(value);
			},
			default: "TOP"
		},

		// 下载
		download: {
			type: Boolean,
			default: false
		},

		// 图片扩展名
		imageExtension: {
			type: String,
			default: "image/jpeg"
		},

		// 图片压缩质量
		imageCompress: {
			type: Number,
			default: 0.75
		},

		// 最大放大倍数
		maxZoomIn: {
			type: Number,
			required: false
		},

		// 最小缩小倍数
		minZoomOut: {
			type: Number,
			default: 1
		},

		// 图像名
		name: {
			type: String,
			default: "screenshot"
		},

		// 图像较短一边占视图的比例
		proportion: {
			type: Number,
			required: false
		},

		// 默认使用框图
		rect: {
			type: Boolean,
			default: false
		},

		// 图像水平方向对齐方式
		rowAlign: {
			validator(value) {
				return ["left", "center", "right"].includes(value);
			},
			default: "center"
		},

		// 默认使用截图
		shot: {
			type: Boolean,
			default: false
		},

		// 默认截图区域
		shotArea: {
			type: Object,
			required: false
		},

		size: {
			type: Number,
			required: false
		},

		// 图像源路径
		src: {
			type: String,
			required: true
		},

		// 默认使用文字输入框
		text: {
			type: Boolean,
			default: false
		},

		// 缩放
		zoom: {
			type: Number,
			default: 1
		}
	},

	data() {
		return {
			// view
			view: null,
			viewWidth: 0,
			viewHeight: 0,

			// image
			imageRealWidth: 0,
			imageRealHeight: 0,
			imageScale: 1,

			// container
			container: null,

			// canvas
			canvas: null,
			canvasWidth: 0,
			canvasHeight: 0,

			// 视网膜
			retinaWidth: 0,
			retinaHeight: 0,

			// 记录当前操作对象的状态
			isCut: false,
			isRect: false,
			isText: false,
			scale: 1,
			moveSpace,
			moveX: 0,
			moveY: 0,
			scaling: false,
			angling: false,
			isDragging: false,
			dragged: false,
			// 记录画布操作状态
			ctxList: [],
			activeIndex: -1,

			// 截图区域
			cutArea: {},

			// 鼠标按下的坐标
			downPoint: null,
			status: -1,
			translate: {
				x: 0,
				y: 0
			},

			// 记录旋转状态
			rotated: false,
			rotate: 0,
			directionCount: 0,

			initX: 0,
			initY: 0,

			// 记录画布的加载状态
			overlay: false
		};
	},

	computed: {
		params() {
			return {
				viewWidth: this.viewWidth,
				viewHeight: this.viewHeight,
				imageRealWidth: this.imageRealWidth,
				imageRealHeight: this.imageRealHeight,
				imageScale: this.imageScale,
				isRect: this.isRect,
				isText: this.isText,
				scale: this.scale,
				moveSpace: this.moveSpace,
				moveX: this.moveX,
				moveY: this.moveY
			};
		},

		computedSize() {
			if (this.size < 1024) {
				return `${this.size}KB`;
			}

			return `${(this.size / 1024).toFixed(2)}M`;
		}
	},

	watch: {
		src: {
			handler(src) {
				debounce(async () => {
					if (this.status === 0) return;

					this.status = 0;
					this.$emit("imageLoad", { status: this.status });
					this.$emit("load", { status: this.status });

					this.overlay = true;

					if (this.canvas) {
						this.canvas?.dispose();
						this.resetValues();
					}
					const newBase64 = await lpTools.getTokenImg(src);
					if (newBase64) {
						lpTools.getBase64(newBase64).then(base64String => {
							lpTools.loadImage(base64String, this.setImage);
						});
					}
					//tools.loadImage(src, this.setImage);
				});
			},
			immediate: true
		}
	},

	beforeDestroy() {
		document.onkeydown = null;
		document.onkeyup = null;

		if (this.view) {
			this.view.onmousedown = null;
			this.view.onmousemove = null;
			this.view.onmouseup = null;
		}
	},

	methods: {
		// 设置图片信息
		setImage(width, height) {
			// 图片不存在
			if (!width || !height) {
				this.status = -1;
				this.$emit("imageLoad", { status: this.status });
				this.$emit("load", { status: this.status });

				this.overlay = false;
				return;
			}

			this.getView();

			this.imageRealWidth = width;
			this.imageRealHeight = height;

			// 图片加载成功
			this.$emit("imageLoad", { status: 1, width, height });

			const scaleWidth = this.viewWidth / this.imageRealWidth;
			const scaleHeight = this.viewHeight / this.imageRealHeight;

			this.imageScale = Math.min(scaleWidth, scaleHeight);

			const imageScaleWidth = this.imageRealWidth * this.imageScale;
			const imageScalseHeight = this.imageRealHeight * this.imageScale;

			if (this.proportion && imageScaleWidth < imageScalseHeight) {
				this.retinaWidth = this.viewWidth * this.proportion;
				const magnification = this.retinaWidth / imageScaleWidth;
				this.retinaHeight = imageScalseHeight * magnification;
			} else {
				this.retinaWidth = imageScaleWidth;
				this.retinaHeight = imageScalseHeight;
			}

			this.createCanvas(this.src);
		},

		// 获取画布视图的宽高
		getView() {
			if (!this.view) {
				this.view = this.$refs.view;
				this.viewWidth = this.view.offsetWidth;
				this.viewHeight = this.view.offsetHeight;
			}

			this.moveCanvasWithMouse();
		},

		// 获取画布容器的实例
		getContainerInstance() {
			this.container = document.querySelector("#view .canvas-container");
		},

		/**
		 * @description 设置画布移动、旋转动画
		 * @param action zoom、move
		 */
		transformContainer(action) {
			if (action === "move") {
				this.$emit("move", { x: this.moveX, y: this.moveY });
			}

			if (this.scaling) {
				switch (this.colAlign) {
					case "start":
						this.container.style["transform-origin"] = "50% 0%";
						break;
				}
			} else {
				this.container.style["transform-origin"] = "50% 50% 0";
			}

			this.container.style.transform = `translate(${this.moveX}px, ${this.moveY}px) scale(${this.scale})`;
			this.container.style.transition = "transform .16s ease-out";
		},

		// 通过鼠标移动画布
		moveCanvasWithMouse() {
			document.onkeydown = ({ altKey }) => {
				if (altKey) {
					this.view.classList.add("drag");
				}
			};

			document.onkeyup = ({ altKey }) => {
				if (!altKey) {
					this.view.classList.remove("drag");
				}
			};

			let [downX, downY] = [void 0, void 0];

			// 按下鼠标
			this.view.onmousedown = downEvent => {
				const { altKey, x, y } = downEvent;

				if (altKey) {
					downX = x;
					downY = y;
				}

				// 移动鼠标
				this.view.onmousemove = moveEvent => {
					tools.throttle(() => {
						const { altKey, x, y } = moveEvent;

						if (altKey) {
							this.moveX = x - downX + this.initX;
							this.moveY = y - downY + this.initY;
							this.transformContainer("move");
						}
					});
				};
			};

			// 抬起鼠标
			this.view.onmouseup = () => {
				this.initX = this.moveX;
				this.initY = this.moveY;
				this.view.onmousemove = null;
			};
		},

		// 创建画布
		createCanvas(source) {
			const options = {
				enableRetinaScaling: true,
				width: this.imageRealWidth,
				height: this.imageRealHeight,
				backgroundColor: "rgb(255, 255, 255)",
				crossOrigin: "anonymous"
			};

			const dimensions = {
				width: this.retinaWidth + "px",
				height: this.retinaHeight + "px"
			};

			this.canvas = new fabric.Canvas("canvas", options);
			this.canvas.setDimensions(dimensions, {
				cssOnly: true
			});
			this.canvas.setCursor(this.pointerIcon);
			this.canvas.on("mouse:down", this.canvasMouseDown);
			this.canvas.on("mouse:up", this.canvasMouseUp);
			//this.canvas.on('mouse:move', this.canvasMove)
			this.canvas.on("mouse:wheel", this.canvasWheel);

			this.getContainerInstance();

			fabric.Image.fromURL(
				source,
				img => {
					img.type = "image";

					if (!this.canvas) return;

					this.canvas.add(img);

					this.canvas.item(0)["hasControls"] = false;
					this.canvas.item(0)["selectable"] = false;
					this.canvas.item(0)["evented"] = false;

					// 设置图像默认方向
					this.setDefaultDirection(() => {
						// 设置图像默认坐标
						this.setDefaultCoord();

						// 设置图像默认截图区域
						this.setDefaultCutArea();

						// 设置默认缩放
						if (this.zoom !== 1) {
							this.scale = this.zoom;
							this.transformContainer();
						}

						this.overlay = false;

						// 设置默认操作
						this.setDefaultOption();

						// 画布完成初始化
						this.status = 1;
						this.$emit("load", {
							status: this.status,
							width: this.imageRealWidth,
							height: this.imageRealHeight
						});
					});
				},
				{ crossOrigin: "anonymous" }
			);
		},

		// 旋转画布
		rotateCanvas(direction, func, manual) {
			this.setDirection(direction, manual);

			if (!direction) return;

			this.clearCutKlass();
			this.setContextsSelectable(false);

			const [imageRealWidth, imageRealHeight] = [this.imageRealWidth, this.imageRealHeight];
			const [retinaWidth, retinaHeight] = [this.retinaWidth, this.retinaHeight];

			const image = new Image();
			image.setAttribute("crossOrigin", "anonymous");

			image.onload = () => {
				// 旋转后交换画布的宽高
				{
					this.imageRealWidth = imageRealHeight;
					this.imageRealHeight = imageRealWidth;

					this.retinaWidth = retinaHeight;
					this.retinaHeight = retinaWidth;
				}

				// 生成旋转后的图片
				const canvas = document.createElement("canvas");

				{
					canvas.width = this.imageRealWidth;
					canvas.height = this.imageRealHeight;

					const ctx = canvas.getContext("2d");

					switch (direction) {
						case "RIGHT":
							ctx.translate(this.imageRealWidth, 0);
							ctx.rotate((90 * Math.PI) / 180);
							break;

						case "LEFT":
							ctx.translate(0, this.imageRealHeight);
							ctx.rotate((-90 * Math.PI) / 180);
							break;
					}

					ctx.drawImage(
						image,
						0,
						0,
						image.width,
						image.height,
						0,
						0,
						image.width,
						image.height
					);
				}

				const dataURL = canvas.toDataURL();

				this.clearCanvas();

				// 旋转后重新设置画布的宽高，重新绘制图片
				{
					this.canvas.setWidth(this.imageRealWidth);
					this.canvas.setHeight(this.imageRealHeight);

					const dimensions = {
						width: this.retinaWidth + "px",
						height: this.retinaHeight + "px"
					};

					this.canvas.setDimensions(dimensions, {
						cssOnly: true
					});

					const imgCtx = this.canvas.getObjects()[0];
					imgCtx.setSrc(
						dataURL,
						() => {
							this.canvas.renderAll();

							func && func();
						},
						{ crossOrigin: "anonymous" }
					);
				}
			};

			image.src = this.canvas.toDataURL();
		},

		// 鼠标按下
		canvasMouseDown(event) {
			const { pointer, e } = event;
			// if (!this.isCut && !this.isRect && !this.isText) {
			//   this.canvasDragDown(e, pointer)
			// }
			this.downPoint = pointer;
			const { klass } = this.getSelectedKlass();

			if (klass) {
				this.activeIndex = klass.unique;
			}

			// 文字
			if (this.isText) {
				this.generateTextbox(event);
			}
		},

		// 鼠标拖动
		canvasMove(event) {
			const { e, pointer } = event;
			if (this.isDragging) {
				let vpt = this.canvas.viewportTransform;
				// vpt[4] += e.clientX - this.lastPosX
				// vpt[5] += e.clientY - this.lastPosY
				vpt[4] += pointer.x - this.lastPosX;
				vpt[5] += pointer.y - this.lastPosY;
				this.translate.x += pointer.x - this.lastPosX;
				this.translate.y += pointer.y - this.lastPosY;
				this.canvas.requestRenderAll();
				this.lastPosX = pointer.x;
				this.lastPosY = pointer.y;
				this.dragged = true;
			}
		},
		// 鼠标抬起
		canvasMouseUp(event) {
			// 切图
			if (this.isCut) {
				const { klass } = this.getSelectedKlass();
				const isBear = this.isBear(event.pointer);

				if (!klass && !isBear) {
					this.clearCutKlass();
				}

				this.generateCutRect(event);
				this.canvas.discardActiveObject();
				return;
			}

			if (!this.isCut && !this.isRect && !this.isText) {
				this.canvas.setViewportTransform(this.canvas.viewportTransform);
				this.isDragging = false;
				this.canvas.selection = true;
			}

			// // 画框
			if (this.isRect) {
				this.createRect(event);
			}
		},

		// 生成切图方框
		generateCutRect({ pointer }) {
			const existKlass = this.ctxList.find(k => k?.type === "cut");

			// 只允存在一个切图框
			if (existKlass && this.activeIndex === -1) {
				const isBear = this.isBear(pointer);
				const index = this.ctxList.findIndex(c => c?.type === "cut");

				// 鼠标按下抬起大于误差范围则删除切图框
				if (!isBear) {
					this.ctxList.splice(index, 1);
					this.canvas.remove(existKlass);
				}
			}

			// 创建切图区域
			if (!existKlass || this.activeIndex === -1) {
				this.createCutRect(pointer);
			}

			// 返回最终切图区域
			{
				const { klass } = this.getSelectedKlass();

				// 计算切图区域
				if (klass?.type === "cut") {
					const { left, top, width, height, scaleX, scaleY } = klass;

					const finalWidth = width * scaleX;
					const finalHeight = height * scaleY;

					this.cutArea = {
						x: left,
						y: top,
						width: finalWidth,
						height: finalHeight
					};

					this.$emit("cut", this.cutArea);
				}
			}
		},

		// 生成输入框
		generateTextbox({ pointer }) {
			if (this.activeIndex === -1) {
				this.createTextbox(pointer.x - this.translate.x, pointer.y - this.translate.y);
				return;
			}

			const activeObject = this.canvas.getActiveObject();

			if (!activeObject) {
				this.activeIndex = -1;
			}

			if (!activeObject) {
				this.clearEmptyTextbox();
			}
		},

		// 滚轮滚动
		canvasWheel(event) {
			this.scaling = true;
			this.angling = false;

			tools.throttle(() => {
				// 放大
				if (event.e.wheelDelta > 0) {
					this.scale = containerEvent.zoomIn(this.params);
					this.limitZoomIn();
				}
				// 缩小
				else {
					this.scale = containerEvent.zoomOut(this.params);
					this.limitZoomOut();
				}

				this.$emit("zoom", this.scale);

				this.transformContainer();
			}, 10);
			event.e.preventDefault();
			event.e.stopPropagation();
		},

		// 图片拖动
		canvasDragDown(e, pointer) {
			if (!this.isCut && !this.isRect && !this.isText) {
				this.canvas.selection = false;
				this.isDragging = true;
				this.lastPosX = pointer.x;
				this.lastPosY = pointer.y;
			}
		},

		// 获取当前选中对象
		getSelectedKlass() {
			let [index, unique] = [-1, -1];
			const klass = this.canvas.getActiveObject();

			if (klass) {
				unique = klass.unique;
				index = this.ctxList.findIndex(k => k.unique === unique);
			}

			return { index, klass, unique };
		},

		// 获取画布上的切图对象
		getCutKlass() {
			return this.ctxList.find(c => c?.type === "cut");
		},

		// 设置所有对象均不可操作/可操作(截图对象保留当前状态)
		setContextsSelectable(selectable = true) {
			this.ctxList.map((klass, index) => {
				if (klass?.type !== "cut" || this.canvas.item(index)) {
					this.canvas.item(index)["selectable"] = selectable;
					this.canvas.item(index)["evented"] = selectable;
				}
			});
		},

		// 设置画布方向
		setDirection(direction, manual) {
			this.rotated = true;
			this.cutArea = {};

			switch (direction) {
				case "RIGHT":
					++this.directionCount;

					if (this.directionCount > 3) {
						this.directionCount = 0;
					}
					break;

				case "LEFT":
					--this.directionCount;

					if (this.directionCount < -3) {
						this.directionCount = 0;
					}
					break;
			}

			this.$emit("direction", tools.setDirection(this.directionCount), manual);
		},

		// 清空画布上当前操作对象
		clearActivatedCtx() {
			const { index, klass } = this.getSelectedKlass();
			this.ctxList.splice(index, 1);
			this.canvas.remove(klass);
		},

		// 清空未输入文字的文本框
		clearEmptyTextbox() {
			const len = this.ctxList.length;

			for (let i = 0; i < len; i++) {
				const ctx = this.ctxList[i];

				if (ctx?.type === "textbox" && !ctx.text) {
					this.ctxList.splice(i, 1);
					this.canvas.remove(ctx);
					break;
				}
			}
		},

		// 清空画布上的切图操作痕迹
		clearCutKlass() {
			const klass = this.getCutKlass();

			if (klass) {
				const index = this.ctxList.findIndex(c => c?.type === "cut");
				this.ctxList.splice(index, 1);
				this.canvas.remove(klass);
			}

			return klass;
		},

		// 清空画布上的操作痕迹
		clearCanvas() {
			for (let ctx of this.ctxList) {
				this.canvas.remove(ctx);
			}

			this.isCut = this.isCut;
			this.isRect = this.isRect;
			this.isText = this.isText;
			this.ctxList = [];

			this.$emit("clear");
		},

		// 限制放大
		limitZoomIn() {
			if (this.maxZoomIn && this.scale >= this.maxZoomIn) {
				this.scale = this.maxZoomIn;
			}
		},

		// 限制缩小
		limitZoomOut() {
			if (this.scale <= this.minZoomOut) {
				this.scale = this.minZoomOut;
			}
		},

		// 允许误差
		isBear({ x, y }) {
			const diffX = x - this.downPoint.x;
			const diffY = y - this.downPoint.y;

			if (diffX > 1 && diffY > 1) {
				return false;
			}

			return true;
		},

		// 保存编辑后的图片
		save() {
			const cutCtx = this.getCutKlass();
			{
				const klass = this.getCutKlass();
				this.canvas.remove(klass);
			}
			const dataURL = this.canvas.toDataURL();

			let args = {
				imageExtension: this.imageExtension,
				imageCompress: this.imageCompress
			};

			if (cutCtx) {
				const realStrokeWidth = cutRectStrokeWidth / this.imageScale;

				const imageWidth = this.cutArea.width - realStrokeWidth * 2;
				const imageHeight = this.cutArea.height - realStrokeWidth * 2;

				args = {
					...args,
					imageWidth,
					imageHeight,
					imageExtension: this.imageExtension,
					imageCompress: this.imageCompress,
					sx: this.cutArea.x + realStrokeWidth,
					sy: this.cutArea.y + realStrokeWidth,
					sw: this.cutArea.width - realStrokeWidth,
					sh: this.cutArea.height - realStrokeWidth,
					dx: 0,
					dy: 0,
					dw: this.cutArea.width - realStrokeWidth,
					dh: this.cutArea.height - realStrokeWidth
				};
			}

			tools.generateImage(dataURL, args, async ({ base64 }) => {
				// 下载
				if (this.download) {
					tools.downloadImage(base64, this.name);
				}

				const file = tools.base64ToFile(base64, this.name);

				const modified = !!this.ctxList.length || this.rotated;

				let cutArea = {};

				if (cutCtx) {
					cutArea = {
						sx: args.sx,
						sy: args.sy,
						sw: args.sw,
						sh: args.sh
					};
				}

				this.$emit("done", {
					file,
					modified,
					width: this.imageRealWidth,
					height: this.imageRealHeight,
					...cutArea
				});
			});
		}
	},

	components: {
		TopBar
	}
};
</script>

<style scoped lang="scss">
$color: #4c4c4c;

.z-drawing-board {
	display: flex;
	flex-direction: column;
	position: relative;
	height: inherit;

	.view {
		flex-grow: 1;
		display: flex;
		/* justify-content: center;
        align-items: center; */
		position: relative;
		background-color: $color;
		border-left: 1px solid $color;
		border-right: 1px solid $color;
		border-bottom: 1px solid $color;
		//overflow: hidden;
		overflow: scroll;
		box-sizing: border-box;
		padding-top: 210px;
		.size {
			position: absolute;
			bottom: 0;
			left: 0;
			padding: 0 0 4px 8px;
			color: #fff;
			font-size: 13px;
		}

		&::after {
			content: " ";
			position: absolute;
			width: 100%;
			height: 100%;
			cursor: grab;
			z-index: -1;
		}

		&.drag::after {
			z-index: 1;
		}
	}
}
</style>
