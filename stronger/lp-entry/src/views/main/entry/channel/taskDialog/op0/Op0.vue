<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="3" class="pt-2 px-2 pb-0 full" id="op0ThumbsNode">
				<v-row v-if="bill.downloadPath &&
		bill.imagesType != null &&
		bill.imagesType[0] != '' &&
		proCode != 'B0110' &&
		proCode != 'B0106' &&
		proCode != 'B0103'
		" dense>
					<v-col v-for="(img, index) in bill.pictures" :key="`billPho${index}`" cols="4"
						:class="[thumbIndex === newfieldsObjectArray[index] ? 'actived' : '']">
						<div :class="[index < redimgArray.length ? 'select-group' : '', 'thumb-wrapper']"
							:id="`${op0ThumbIdPrefix}${newfieldsObjectArray[index]}`" @click="
		onSelectThumb({
			id: `${op0ThumbIdPrefix}${newfieldsObjectArray[index]}`,
			thumbIndex: newfieldsObjectArray[index],
			imgindex: index
		})
		">
							<img :style="{ width: '100%', height: '156px' }" :src="readerPath[index]" />
							<p>{{ imgName[index] }}</p>
						</div>
					</v-col>
				</v-row>
				<v-row v-else dense>
					<v-col v-for="(img, index) in bill.pictures" :key="`billPho${index}`" cols="4"
						:class="[thumbIndex === index ? 'actived' : '']">
						<div :class="[
		/[\u4E00-\u9FFF\u3400-\u4DFF\uF900-\uFAFF]/.test(bill.pictures[index]) || bill.imagesType
			? 'thumb-wrapper'
			: 'thumb-wrappers',
		bill.imagesType && imgType(bill.imagesType[index]) ? 'select-group' : ''
	]" :id="`${op0ThumbIdPrefix}${index}`" @click="
		onSelectThumb({
			id: `${op0ThumbIdPrefix}${index}`,
			thumbIndex: index,
			imgindex: index
		})
		">
							<img :style="{ width: '100%', height: '156px' }" :src="readerPath[index]" />
							<p v-if="/[\u4E00-\u9FFF\u3400-\u4DFF\uF900-\uFAFF]/.test(bill.pictures[index])">
								{{ imgDesc(bill.pictures[index]) }}
							</p>
							<p v-if="bill.imagesType">{{ imgType(bill.imagesType[index]) }}</p>
						</div>
					</v-col>
				</v-row>
			</v-col>

			<v-col ref="drawImageContain" cols="7" class="pt-2 pb-0 pl-2 pr-0 full image-contain">
				<!-- src中间图片， 路径B0113拼接 -->
				<AZDrawingBoard id="op0DrawImage" ref="drawImage" col-align="start" :coord="returnCoord()"
					:direction="returnDriection()" imageExtension="image/png" :min-zoom-out="0.5" :name="drewImageName"
					:proportion="0.75" shot :shotArea="returnArea()" :size="imageSize" :src="canvasImg" :zoom="returnScale()"
					@cut="handleCut" @direction="handleDirection" @done="handleDrewSave" @load="handleLoad" @move="handleMove"
					@restore="handleInitImage" @zoom="handleZoom">
				</AZDrawingBoard>
			</v-col>

			<v-col id="op0FieldsNode" cols="2" class="pt-2 pb-0 full fields-wrapper">
				<div class="loading t-loading" data-value="tLoading"></div>

				<!-- fields标题 BEGIN -->
				<op-fields-title :blockCode="block.code" :isLoop="block.isLoop" :proCode="bill.proCode"
					:title="block.name"></op-fields-title>
				<!-- fields标题 END -->

				<v-form ref="op0Form" v-model="valid" lazy-validation @submit="preventForm">
					<template v-if="fieldsObject[thumbIndex]">
						<div v-for="(fields, fieldsIndex) in fieldsObject[thumbIndex].fieldsList" :key="`fields_${fieldsIndex}`"
							class="mb-n4">
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field v-show="field.show !== false" :key="field.uniqueKey"
									:accent="field.name === `模板类型字段` ? true : false" :autofocus="field.autofocus" :bill="bill"
									:block="block" class="mb-n4" :clearNVs="clearNVs" :disabled="field.disabled" :field="field"
									:fieldsIndex="fieldsIndex" :fieldsList="fieldsObject[thumbIndex].fieldsList"
									:fieldsObject="fieldsObject" :focusFieldsIndex="focusFieldsIndex" :hint="field.hint"
									:id="field.uniqueId" :includes="computedIncludes(field)" :items="field.items || []"
									:sameFieldValue="field.sameFieldValue || []"
									:label="computedLabel({ field, fieldIndex, fieldsIndex })" :labelTip="field.code" :op="op"
									:ref="field.uniqueId" :svHints="computedSVHints" :svValidations="computedSVValidations"
									:thumbIndex="thumbIndex" :validations="computedRules(field.rules)" :defaultValue="field.op0Value"
									@enter="onEnterField($event, field, fieldsIndex, fieldIndex)"
									@focus="onFocusField($event, field, fieldsIndex, fieldIndex)"
									@input="onInputField($event, field, fieldsIndex, fieldIndex)"
									@keyup.38="onDnKey($event, field, fieldsIndex, fieldIndex)"
									@dropdown="onDropdownField($event, field, fieldsIndex, fieldIndex)"
									@dropdownUp="onDropdownUpField($event, field, fieldsIndex, fieldIndex)"
									@ruleClick="handleFieldRulesDialog"></op-text-field>
							</template>
						</div>
					</template>
				</v-form>

				<div class="loading b-loading" data-value="bLoading"></div>
			</v-col>
		</v-row>

		<!-- 字段规则 BEGIN -->
		<field-rules-dialog ref="fieldRules" :fieldName="fieldName" :proCode="bill.proCode"></field-rules-dialog>
		<!-- 字段规则 END -->
	</div>
</template>

<script>
import _ from "lodash";
import axios from "axios";
import { localStorage, sessionStorage } from "vue-rocket";
import { tools } from "vue-rocket";
import OpMixins from "../mixins/OpMixins";
import OpDropdownMixins from "../mixins/OpDropdownMixins";
import OpSpecificValidationsMixins from "../mixins/OpSpecificValidationsMixins";
import ScrollUpDnMixins from "../mixins/ScrollUpDnMixins";
import { tools as lpTools } from "@/libs/util";
import Op0Mixins from "./mixins/Op0Mixins";
import Op0CanvasMixins from "./mixins/Op0CanvasMixins";
import Op0ShortcutMixins from "./mixins/Op0ShortcutMixins";
import Op0ThumbsMixins from "./mixins/Op0ThumbsMixins";

import { op0ThumbIdPrefix } from "./libs/constants";

const imgTypesMap = new Map([
	["fapiao", "发票"],
	["qingdan", "清单"],
	["zhenduan", "诊断书"],
	["baoxiaodan", "报销单"]
]);

export default {
	name: "Op0",
	mixins: [
		OpMixins,
		OpDropdownMixins,
		Op0Mixins,
		Op0CanvasMixins,
		Op0ShortcutMixins,
		Op0ThumbsMixins,
		OpSpecificValidationsMixins,
		ScrollUpDnMixins
	],

	data() {
		return {
			formId: "Op0",
			op: "op0",
			op0ThumbIdPrefix,
			imgindex: 0,
			// 缩略图拼接的路径
			paths: [],
			// 缩略图路径
			readerPath: [],
			// 图片请求头
			instance: "",
			// canvas大图路径
			canvasImg: "",
			showOp: sessionStorage.get("hasOp")
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

	watch: {
		bill: {
			handler(newVal) {
				// 图片路径
				if (newVal.hasOwnProperty("pictures")) {
					if (
						this.proCode == "B0108" ||
						this.proCode == "B0114" ||
						this.proCode == "B0118" ||
						this.proCode == "B0116"
					) {
						newVal.pictures.forEach(img => {
							this.paths.push(`${this.fileUrl}${newVal.downloadPath}A${img}`);
						});
					} else {
						newVal.pictures.forEach(img => {
							this.paths.push(
								`${this.fileUrl.replace("files/", "")}${newVal.downloadPath}${img.replace("/", "/A")}`
							);
						});
					}
				}

				this.paths.forEach(async (el, index) => {
					let item = await this.transform(el);

					this.getReader(item).then(res => {
						// this.readerPath[index] = res;
						this.$set(this.readerPath, index, res);
					});
				});
			},
			immediate: true
		},

		async modifyImage(newVal) {
			if (this.proCode == "B0108" || this.proCode == "B0114" || this.proCode == "B0118" || this.proCode == "B0116") {
				let path = `${this.fileUrl}${this.bill.downloadPath}${newVal}`;
				let item = await this.transform(path);

				this.getReader(item).then(res => {
					this.canvasImg = res;
				});
			} else {
				let path = `${this.fileUrl.replace("files/", "")}${this.bill.downloadPath}${newVal}`;

				let item = await this.transform(path);

				this.getReader(item).then(res => {
					this.canvasImg = res;
				});
			}
		}
	},

	computed: {
		computedIncludes() {
			return field => {
				return this.svDropdownFields[field.code]?.desserts || [];
			};
		},
		imgDesc() {
			return desc => {
				return desc.match(/[^\x00-\x80]/g)?.join("");
			};
		},
		// B0110 B0106 B0103
		imgType() {
			return type => {
				return imgTypesMap.get(type);
			};
		}
	},

	methods: {
		// 选择缩略图
		onSelectThumb({ id, thumbIndex, imgindex }) {
			if (this.imageLoading) {
				return;
			}
			this.$store.commit('UPDATE_PAGE', imgindex)
			this.imgindex = imgindex;
			this.prevFocusFieldsIndex = -1;
			this.thumbIndex = thumbIndex;

			this.scrollThumbsUpDn({ id, thumbIndex, imgindex });

			// 获取临时保存的字段
			this.getSessionSaveField();

			// 按F8后将，唯一一个[模板类型字段]被清除掉
			if (tools.isLousy(this.fieldsObject[this.thumbIndex]?.fieldsList)) {
				this.fieldsObject[this.thumbIndex].fieldsList = tools.deepClone(this.tempFields);
			}

			if (this.thumbIndex - 1 >= 0 && !this.fieldsObject[this.thumbIndex - 1].sessionStorage) {
				this.fieldsObject[this.thumbIndex - 1].fieldsList = [tools.deepClone(this.tempFields)];
			}

			this.bill.thumbIndex = this.thumbIndex;

			if (tools.isLousy(this.fieldsObject)) {
				this.toasted.warning("字段为空！");
				return;
			}

			this.$emit("bill", this.bill);
			// console.log("this.svMemoFieldValues", this.svMemoFieldValues);
			// console.log("this.sameFieldValue", this.sameFieldValue);

			// console.log("this.clearValues", this.clearValues);
			// console.log("this.clearFieldValues", this.clearFieldValues);

			for (let el in this.clearValues) {
				if (this.svMemoFieldValues.hasOwnProperty(el)) {
					let arr = this.svMemoFieldValues[el].values;
					let clearArr = this.clearValues[el].values;
					for (let item of clearArr) {
						let flag = arr.indexOf(item);
						arr.splice(flag, 1);
					}
					_.set(this.svMemoFieldValues, `${el}.values`, arr);
					// this.svMemoFieldValues[el] = arr;
				}
			}

			for (let el in this.clearFieldValues) {
				if (this.sameFieldValue.hasOwnProperty(el)) {
					let arr = this.sameFieldValue[el].values;
					let clearArr = this.clearFieldValues[el].values;
					for (let item of clearArr) {
						let flag = arr.indexOf(item);
						arr.splice(flag, 1);
					}
					_.set(this.sameFieldValue, `${el}.values`, arr);
					// this.sameFieldValue[el] = arr;
				}
			}

			this.clearValues = {};
			this.clearFieldValues = {};
		},

		// 默认截图区域
		returnArea() {
			const { op0Value } = this.getRangeField();
			return !!op0Value ? void 0 : this.memoMarks[this.thumbIndex]?.area;
		},

		// 默认坐标
		returnCoord() {
			const { op0Value } = this.getRangeField();
			return !!op0Value ? void 0 : this.memoMarks[this.thumbIndex]?.coord;
		},

		// 默认旋转方向
		returnDriection() {
			const { op0Value } = this.getRangeField();
			return !!op0Value ? "TOP" : this.memoMarks[this.thumbIndex]?.direction;
		},

		// 默认缩放比例
		returnScale() {
			const { op0Value } = this.getRangeField();
			return !!op0Value ? 0.9 : this.memoMarks[this.thumbIndex]?.scale;
		},

		// 获取显示范围字段
		getRangeField() {
			return this.fieldsObject[this.thumbIndex]?.fieldsList?.[this.focusFieldsIndex]?.[2] || {};
		},

		// 临时保存的字段
		returnMemoField(uniqueId) {
			if (this.memoFields[uniqueId]) {
				const { items, sessionStorage, value } = this.memoFields[uniqueId];
				console.log("临时保存的字段", value);

				if (sessionStorage) {
					return { items, value };
				}
			}

			return {};
		},

		// 完成加载
		handleLoad({ status }) {
			this.imageLoading = !status;
		},

		// 记录移动坐标
		handleMove(coord) {
			this.updateMarks({ coord });
		},

		// 记录最后一次切图
		handleCut(area) {
			this.updateMarks({ area });
		},

		// 记录最后一次旋转
		handleDirection(direction, manual) {
			if (manual) {
				this.updateMarks({ direction });
			}
		},

		// 记录最后一次缩放
		handleZoom(scale) {
			this.updateMarks({ scale });
		},

		// 更新操作记录
		updateMarks({ area, coord, direction, scale }) {
			this.memoMarks = { ...this.memoMarks };

			if (!this.memoMarks[this.thumbIndex]) {
				this.memoMarks[this.thumbIndex] = {};
			}

			if (area) {
				this.memoMarks[this.thumbIndex].area = area;
			}

			if (coord) {
				this.memoMarks[this.thumbIndex].coord = coord;
			}

			if (direction) {
				this.memoMarks[this.thumbIndex].direction = direction;
			}

			if (scale) {
				this.memoMarks[this.thumbIndex].scale = scale;
			}
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
		},

		// 选中全部
		selectAll(field) {
			const el = document.querySelector(`#input_${field.uniqueId}`);
			if (el.selectionStart == el.selectionEnd) {
				lpTools.setCursorPosition(el, 0, field.resultValue.length);
			}
		}

		// 图片路径
		// imgSrc(proCode, fileUrl, downloadPath, img) {
		// 	if (proCode == "B0108" || proCode == "B0114" || proCode == "B0118") {
		// 		return `${fileUrl}${downloadPath}A${img}`;
		// 	} else {
		// 		return `${fileUrl.replace("files/", "")}${downloadPath}${img.replace("/", "/A")}`;
		// 	}
		// }
	},

	components: {
		"op-fields-title": () => import("../components/opFieldsTitle"),
		"op-text-field": () => import("../components/opTextField"),
		"field-rules-dialog": () => import("../components/fieldRulesDialog"),
		AZDrawingBoard: () => import("../op0/ZDrawingBoard/ZDrawingBoard.vue")
	}
};
</script>

<style scoped lang="scss">
@import "../op.scss";

/* 缩略图 */
.thumb-wrapper {
	box-sizing: border-box;
	height: 185px;
	border: 2px solid transparent;
	text-align: center;
	overflow: hidden;

	p {
		height: 25px;
		font-weight: 700;
	}
}

.thumb-wrappers {
	box-sizing: border-box;
	height: 160px;
	border: 2px solid transparent;
	text-align: center;
	overflow: hidden;
}

.actived {
	border: 2px solid #1976d2;
	padding: 2px;
}

.select-group {
	border: 2px solid #f80303;
}

/* 画布 */
.draw-image-wrapper {
	height: 100%;
}
</style>
