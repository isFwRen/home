<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="9" class="full image-contain">
				<!-- <watch-image
					v-if="bill.proCode == 'B0108'"
					ref="watchImage"
					:op="op"
					:bill="bill"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image>
				<watch-image
					v-else-if="bill.proCode == 'B0114'"
					ref="watchImage"
					:op="op"
					:bill="bill"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image>
				<watch-image
					v-else-if="bill.proCode == 'B0118'"
					ref="watchImage"
					:op="op"
					:bill="bill"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image>
				<watch-image
					v-else
					ref="watchImage"
					:op="op"
					:bill="bill"
					:src="`${fileUrl.replace('files/', '')}${bill.downloadPath}${block.picture}`"
					:isLoop="block.isLoop"
				></watch-image> -->

				<!-- <watch-image
					ref="watchImage"
					:op="op"
					:bill="bill"
					:src="imgSrc(bill.proCode, fileUrl, bill.downloadPath, block.picture, block)"
					:isLoop="block.isLoop"
				></watch-image> -->

				<watch-image ref="watchImage" :op="op" :bill="bill" :src="imgTrail" :isLoop="block.isLoop"></watch-image>
			</v-col>

			<v-col id="opqFieldsNode" class="full fields-wrapper" cols="3" ref="col">
				<div class="loading t-loading" data-value="tLoading"></div>

				<!-- fields标题 BEGIN -->
				<op-fields-title
					:blockCode="block.code"
					:isLoop="block.isLoop"
					:proCode="bill.proCode"
					:title="block.name"
				></op-fields-title>
				<!-- fields标题 END -->

				<v-form ref="opqForm" v-model="valid" lazy-validation @submit="preventForm">
					<template v-if="fieldsList.length">
						<div v-for="(fields, fieldsIndex) in fieldsList" :key="`fields_${fieldsIndex}`">
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field
									v-if="showRange(field, fieldIndex, fieldsIndex)"
									:key="field.uniqueKey"
									:autofocus="field.autofocus"
									:bill="bill"
									:block="block"
									class="mb-n4"
									:clearNVs="clearNVs"
									:disabled="field.disabled"
									:field="field"
									:fieldsIndex="fieldsIndex"
									:fieldsList="fieldsList"
									:firstDiffIndex="firstDiffIndex"
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
									:defaultValue="field.opqValue"
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
									@focusClear="clearCount"
									@input="onInputField($event, field, fieldsIndex, fieldIndex)"
									@keyup="onOpqKeydownField($event, field)"
									@keydown.38="
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
								>
									<div :key="field.uniqueKey" class="diff">
										<p
											class="mb-0 op2"
											v-html="lpTools.compareString(field.op2Value, field.op1Value).targetHtml"
										></p>
										<p
											class="mb-0 op1"
											v-html="lpTools.compareString(field.op1Value, field.op2Value).targetHtml"
										></p>
									</div>
								</op-text-field>
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
import { tools as lpTools } from "@/libs/util";
import OpMixins from "../mixins/OpMixins";
import OpDropdownMixins from "../mixins/OpDropdownMixins";
import Op1Op2OpqMixins from "../mixins/Op1Op2OpqMixins";
import OpSpecificValidationsMixins from "../mixins/OpSpecificValidationsMixins";
import ScrollUpDnMixins from "../mixins/ScrollUpDnMixins";
import axios from "axios";
import { localStorage } from "vue-rocket";

// 获取字符串问号的下标
const questionMarkIndexs = (value, symbol = "?") => {
	return value
		.split("")
		.map((v, i) => {
			if (v === symbol) return i;
		})
		.filter(f => f || f === 0);
};

export default {
	name: "Opq",
	mixins: [OpMixins, OpDropdownMixins, Op1Op2OpqMixins, OpSpecificValidationsMixins, ScrollUpDnMixins],

	data() {
		return {
			formId: "Opq",
			op: "opq",
			lpTools,
			firstDiffIndex: -1,
			indexsIndex: 0,
			fieldIndex: 0,
			// 问题件请求
			instance: "",
			// 问题件拼接路径
			imgPath: "",
			// 问题件展示路径
			imgTrail: ""
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

	methods: {
		onOpqKeydownField({ keyCode }, field) {
			if (keyCode === 13) {
				return;
			}

			this.indexsIndex = -1;
			this.opqNextDiffIndex(field);
		},

		// 找到一码二码第一个不同的值，并在 field 选中
		opqMarkFieldFirstDiffIndex(field, symbol = "?") {
			if (this.op !== "opq") {
				return;
			}

			if (field.opqValue) {
				this.firstDiffIndex = field.opqValue.indexOf(symbol);
			} else {
				this.firstDiffIndex = -1;
			}
		},

		// 选中下一个
		opqNextDiffIndex(field, symbol = "?") {
			if (this.op !== "opq") {
				return;
			}
			const opqValue = field.opqValue;

			const indexs = questionMarkIndexs(opqValue);

			++this.indexsIndex;

			const lastIndex = indexs.length - 1;

			if (this.indexsIndex === lastIndex) {
				this.indexsIndex = lastIndex;
			}

			this.firstDiffIndex = indexs[this.indexsIndex];
		},

		clearCount(count) {
			this.indexsIndex = count;
			this.firstDiffIndex = count;
		},

		// 设置默认值
		opqSetFieldEffectKeyValue(field) {
			if (this.op !== "opq") {
				return;
			}

			// 前端自行根据一码二码的值动态设置问题件 opqValue 的值(当第一次领取时)
			if (field.disabled == false && field.show === true) {
				const initTime = this.block.opqSubmitAt?.slice(0, 10);

				if (initTime === "0001-01-01") {
					field.opqValue = lpTools.compareString(field.op2Value, field.op1Value).diffValue;
					field.resultValue = field.opqValue;
				}
			}
		},

		navToViewImages() {
			this.$refs.watchImage.navToViewImages();
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

		// 图片路径
		imgSrc(proCode, fileUrl, downloadPath, picture, block) {
			console.log(block);
			if (proCode == "B0108" || proCode == "B0114" || proCode == "B0118" || proCode == "B0116") {
				return `${fileUrl}${downloadPath}${picture}`;
			} else {
				return `${fileUrl.replace("files/", "")}${downloadPath}${picture}`;
			}
		}
	},

	watch: {
		focusFieldIndex: {
			handler(index, oldValue) {
				if (oldValue !== undefined && index !== oldValue) {
					this.indexsIndex = -1;
					this.firstDiffIndex = -1;
					console.log("Watch-----");
				}
			},
			immediate: true
		},

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
</style>