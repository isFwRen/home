<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="9" class="full image-contain">
				<!-- <watch-image
					v-if="bill.proCode == 'B0108'"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image>
				<watch-image
					v-else-if="bill.proCode == 'B0114'"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image>
				<watch-image
					v-else-if="bill.proCode == 'B0118'"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image>
				<watch-image
					v-else
					:src="`${fileUrl.replace('files/', '')}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image> -->
				<watch-image :src="imgTrail" :isLoop="block.isLoop"></watch-image>
			</v-col>

			<v-col id="op2FieldsNode" class="full fields-wrapper" cols="3" ref="col">
				<div class="loading t-loading" data-value="tLoading"></div>

				<!-- fields标题 BEGIN -->
				<op-fields-title
					:blockCode="block.code"
					:isLoop="block.isLoop"
					:proCode="bill.proCode"
					:title="block.name"
				></op-fields-title>
				<!-- fields标题 END -->

				<!-- 循环分块记录上一组分块项目名称 BEGIN -->
				<div v-if="isLoop && focusFieldsIndex > 0">
					<ul class="pl-1">
						<li v-for="field in fieldsList[focusFieldsIndex - 1]" :key="field.uniqueKey">
							{{ computedPrevFieldsName(field) }}
						</li>
					</ul>
				</div>
				<!-- 循环分块记录上一组分块项目名称 END -->
				<!-- v-if="showRange(field, fieldIndex, fieldsIndex)" -->

				<v-form ref="op2Form" v-model="valid" lazy-validation @submit="preventForm">
					<template v-if="fieldsList.length">
						<div v-for="(fields, fieldsIndex) in fieldsList" :key="`fields_${fieldsIndex}`">
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field
									v-if="showRange(field, fieldIndex, fieldsIndex)"
									:key="field.uniqueKey"
									:lastDiffIndex="lastDiffIndex"
									:firstDiffIndex="firstDiffIndex"
									:autofocus="field.autofocus"
									:bill="bill"
									:block="block"
									class="mb-n4"
									:clearNVs="clearNVs"
									:disabled="field.disabled"
									:field="field"
									:fieldsIndex="fieldsIndex"
									:fieldsList="fieldsList"
									:hint="field.hint"
									:id="field.uniqueId"
									:includes="svDropdownFields[field.code] ? svDropdownFields[field.code].desserts : []"
									:items="field.items || []"
									:label="computedLabel({ field, fieldIndex, fieldsIndex })"
									:labelTip="field.code"
									:op="op"
									:ref="field.uniqueId"
									:svHints="computedSVHints"
									:svValidations="computedSVValidations"
									:validations="computedRules(field.rules)"
									:defaultValue="field.op2Value"
									@enter="
										onEnterField({
											event: $event,
											field,
											fields,
											fieldsIndex,
											fieldIndex
										})
									"
									@focus="
										onFocusField({
											event: $event,
											field,
											fields,
											fieldsIndex,
											fieldIndex
										})
									"
									@input="onInputField($event, field, fieldsIndex, fieldIndex)"
									@keyup.38="
										onDnKey({
											event: $event,
											field,
											fields,
											fieldsIndex,
											fieldIndex
										})
									"
									@dropdown="onDropdownField($event, field, fieldsIndex, fieldIndex)"
									@dropdownUp="onDropdownUpField($event, field, fieldsIndex, fieldIndex)"
									@ruleClick="handleFieldRulesDialog"
								></op-text-field>
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
import OpMixins from "../mixins/OpMixins";
import OpDropdownMixins from "../mixins/OpDropdownMixins";
import Op1Op2OpqMixins from "../mixins/Op1Op2OpqMixins";
import OpSpecificValidationsMixins from "../mixins/OpSpecificValidationsMixins";
import ScrollUpDnMixins from "../mixins/ScrollUpDnMixins";
import axios from "axios";
import { localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

// 获取字符串问号的下标
const questionMarkIndexs = (value, symbol = "?") => {
	return value
		.split("")
		.map((v, i) => {
			if (v === symbol && value[i - 1] != "?") return i;
		})
		.filter(f => f || f === 0);
};
export default {
	name: "Op2",
	mixins: [OpMixins, OpDropdownMixins, Op1Op2OpqMixins, OpSpecificValidationsMixins, ScrollUpDnMixins],

	data() {
		return {
			formId: "Op2",
			op: "op2",
			// 二码请求
			instance: "",
			// 二码拼接路径
			imgPath: "",
			// 图片展示路径
			imgTrail: "",
			// 108一二码选中问号
			firstDiffIndex: -1,
			indexsIndex: 0,
			lastDiffIndex: -1,
			fieldValue: ""
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
		"block.picture": {
			async handler(val) {
				if (
					this.bill.proCode == "B0108" ||
					this.bill.proCode == "B0114" ||
					this.bill.proCode == "B0118" ||
					this.bill.proCode == "B0116"
				) {
					this.imgPath = `${this.fileUrl}${this.bill.downloadPath}${val}`;
				} else {
					this.imgPath = `${this.fileUrl.replace("files/", "")}${this.bill.downloadPath}${val}`;
				}

				let item = await this.transform(this.imgPath);

				this.getReader(item).then(res => {
					this.imgTrail = res;
				});
			}
		},

		focusFieldIndex: {
			handler(index, oldValue) {
				if (oldValue !== undefined && index !== oldValue) {
					this.indexsIndex = 0;
					// this.firstDiffIndex = -1;
					// this.lastDiffIndex = -1;
					console.log("Watch-----");
				}
			},
			immediate: true
		}
	},

	methods: {
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

		// 找到一码二码第一个?，并在 field 选中
		opqMarkFieldFirstDiffIndex(field, symbol = "?") {
			if (field[`${this.op}Value`]) {
				this.fieldValue = field[`${this.op}Value`];
				this.firstDiffIndex = field[`${this.op}Value`].indexOf(symbol);
				let valueArr = [...field[`${this.op}Value`]];
				for (let i = this.firstDiffIndex; i < valueArr.length; i++) {
					if (valueArr[i] == "?" && i == valueArr.length - 1) {
						this.lastDiffIndex = i + 1;
						break;
					}
					if (valueArr[i] != "?") {
						this.lastDiffIndex = i;
						break;
					}
				}
			} else {
				this.firstDiffIndex = -1;
			}
		},
		// 选中下一个
		opqNextDiffIndex(field, symbol = "?") {
			const opqValue = field[`${this.op}Value`];

			const indexs = questionMarkIndexs(opqValue);
			console.log(indexs);
			if (this.fieldValue == field[`${this.op}Value`]) ++this.indexsIndex;
			else {
				this.fieldValue = field[`${this.op}Value`];
				this.indexsIndex = 0;
			}

			if (this.indexsIndex === indexs.length) {
				this.indexsIndex = 0;
			}
			if (!indexs[this.indexsIndex]) this.indexsIndex = 0;
			this.firstDiffIndex = indexs[this.indexsIndex];
			let valueArr = [...field[`${this.op}Value`]];
			for (let i = this.firstDiffIndex; i < valueArr.length; i++) {
				if (valueArr[i] == "?" && i == valueArr.length - 1) {
					this.lastDiffIndex = i + 1;
					break;
				}
				if (valueArr[i] != "?") {
					this.lastDiffIndex = i;
					break;
				}
			}
			const el = document.querySelector(`#input_${field.uniqueId}`);
			if (el.selectionStart == el.selectionEnd) {
				lpTools.setCursorPosition(el, this.firstDiffIndex, this.lastDiffIndex);
			}
			// console.log("this.firstDiffIndex-----回车------", this.firstDiffIndex);
			// console.log("this.lastDiffIndex------回车------", this.lastDiffIndex);
		},

		// 选中全部
		selectAll(field) {
			const el = document.querySelector(`#input_${field.uniqueId}`);
			if (el.selectionStart == el.selectionEnd) {
				lpTools.setCursorPosition(el, 0, field.resultValue.length);
			}
		}
	},

	components: {
		"op-text-field": () => import("../components/opTextField"),
		"op-fields-title": () => import("../components/opFieldsTitle"),
		"watch-image": () => import("../components/watchImage"),
		"field-rules-dialog": () => import("../components/fieldRulesDialog")
	}
};
</script>

<style scoped lang="scss">
@import "../op.scss";

::v-deep input::selection {
	background: #ffc0cb;
	/* 粉红色的背景 */
	color: red;
	/* 文字颜色 */
}
</style>
