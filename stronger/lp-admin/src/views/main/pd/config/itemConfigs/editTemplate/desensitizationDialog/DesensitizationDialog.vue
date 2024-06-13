<template>
	<lp-dialog ref="dialog" :title="title" width="960" fullscreen @dialog="handleDialog">
		<div slot="main" style="height: 100%">
			<div class="z-flex align-end pt-4 pb-8 lp-filters">
				<v-row class="z-flex align-end">
					<v-col
						v-for="(item, index) in cells.fields"
						:key="`entry_filters_${index}`"
						:cols="2"
					>
						<template v-if="item.inputType === 'input'">
							<z-text-field
								:formId="formId"
								:formKey="item.formKey"
								disabled
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:suffix="item.suffix"
								:defaultValue="shotArea[item.formKey]"
							>
							</z-text-field>
						</template>

						<template v-else>
							<z-select
								:formId="formId"
								:formKey="item.formKey"
								:hideDetails="item.hideDetails"
								:hint="item.hint"
								:label="item.label"
								:options="imageOptions"
								:suffix="item.suffix"
								:defaultValue="detail.picPage"
								@change="switchImg"
							></z-select>
						</template>
					</v-col>

					<div class="mb-3">
						<z-btn class="px-3" color="primary" @click="onSave"> 保存 </z-btn>
					</div>
				</v-row>
			</div>

			<div style="height: calc(100vh - 180px)">
				<lp-drawing-board
					v-if="dialog"
					id="op0DrawImage"
					ref="drawImage"
					imageExtension="image/png"
					:min-zoom-out="0.5"
					shot
					:shotArea="shotArea"
					:src="url"
					@clear="handleClear"
					@cut="handleCoordinate"
					@imageLoad="handleImageLoad"
				></lp-drawing-board>
			</div>
		</div>
	</lp-dialog>
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import { tools } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import cells from "./cells";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "DesensitizationDialog",
	mixins: [DialogMixins],

	props: {
		deviceType: {
			type: Boolean,
			default: false
		},

		imageOptions: {
			type: Array,
			default: () => []
		},

		imageList: {
			type: Array,
			default: () => []
		}
	},

	data() {
		return {
			formId: "desensitizationDialog",
			cells,
			url: "",
			imageWidth: 0,
			imageHeight: 0,
			shotArea: {
				x: 0,
				y: 0,
				width: 0,
				height: 0
			},
			coor: [],

			wCoordinate: [],

			index: -1,
			picPageIndex: 0,

			dialog: false
		};
	},

	computed: {
		...mapGetters(["config"])
	},

	watch: {
		"detail.picPage": {
			handler(index) {
				if (tools.isYummy(index)) {
					this.index = index;
					this.wCoordinate = this.detail.wCoordinate;
					let reg = new RegExp("files//", "g");
					if (!this.imageList[index]) {
						this.toasted.warning("图片路径有误");
					}
					this.url = (baseURLApi + this.imageList[index]).replace(reg, "files/");
				}
			},
			immediate: true
		},

		dialog: {
			handler(dialog) {
				this.dialog = dialog;
				this.$emit("dialog", dialog);
			}
		}
	},

	methods: {
		switchImg(value) {
			this.url = baseURLApi + this.imageList[value];
			this.setDefaultCutArea(value, this.detail.picPage);
		},

		async onSave() {
			const form = {
				blockId: this.detail.ID,
				coordinate: this.coor,
				coordinateType: this.detail.coordinateType,
				picPage: this.forms[this.formId].picPage
			};

			const result = await this.$store.dispatch("UPDATE_CONFIG_TEMP_CHUNK_SCREENSHOT", form);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.$emit("updated");
			}
		},

		// 图片加载完成
		handleImageLoad({ status, width, height }) {
			if (status === 1) {
				this.imageWidth = width;
				this.imageHeight = height;
				this.setDefaultCutArea(this.index, this.index);
			}
		},

		// 清空画布
		handleClear() {
			this.shotArea = {
				x: 0,
				y: 0,
				width: 0,
				height: 0
			};

			this.coor = [];
		},

		// 截图
		handleCoordinate({ x, y, width, height }) {
			this.shotArea = {
				x,
				y,
				width,
				height
			};

			const percentStartX = (x / this.imageWidth) * 100;
			const percentStartY = (y / this.imageHeight) * 100;
			const percentEndX = ((x + width) / this.imageWidth) * 100;
			const percentEndY = ((y + height) / this.imageHeight) * 100;

			this.coor = [
				String(percentStartX),
				String(percentStartY),
				String(percentEndX),
				String(percentEndY)
			];
		},

		// 默认截图
		setDefaultCutArea(index, picPage) {
			if (index === picPage) {
				if (tools.isYummy(this.wCoordinate)) {
					const [startX, startY, endX, endY] = this.wCoordinate;

					const rectX = this.imageWidth * (startX / 100);
					const rectY = this.imageHeight * (startY / 100);
					const rectW = (endX / 100) * this.imageWidth - rectX;
					const rectH = (endY / 100) * this.imageHeight - rectY;

					if (rectW && rectH) {
						this.shotArea = {
							x: rectX,
							y: rectY,
							width: rectW,
							height: rectH
						};
					}
					return;
				}

				this.shotArea = {
					x: 0,
					y: 0,
					width: 0,
					height: 0
				};

				this.coor = [];
			}
		}
	}
};
</script>
