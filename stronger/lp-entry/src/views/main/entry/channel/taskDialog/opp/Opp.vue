<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="9" class="full image-contain">
				<!-- <watch-image
					:src="`https://img1.baidu.com/it/u=3893389324,4043822814&fm=253&fmt=auto&app=120&f=JPEG?w=1280&h=800`"
				></watch-image> -->
				<watch-image :src="imgTrail"></watch-image>
			</v-col>

			<v-col id="oppFieldsNode" class="full fields-wrapper" cols="3">
				<div class="loading t-loading" data-value="tLoading" style="margin-bottom: 10px; margin-top: 40px"></div>

				<!-- fields标题 BEGIN -->
				<op-fields-title
					:blockCode="block.code"
					:isLoop="block.isLoop"
					:proCode="bill.proCode"
					:title="block.name"
					:isOpenSource="bill.video"
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

				<v-form ref="oppForm" v-model="valid" lazy-validation @submit="preventForm">
					<template v-if="fieldsList.length">
						<div v-for="(fields, fieldsIndex) in fieldsList" :key="`fields_${fieldsIndex}`">
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field
									v-show="field.show"
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
									:hint="field.hint"
									:id="field.uniqueId"
									:includes="computedIncludes(field)"
									:items="field.items || []"
									:label="computedLabel({ field, fieldIndex, fieldsIndex })"
									:labelTip="field.code"
									:op="op"
									:ref="field.uniqueId"
									:svHints="computedSVHints"
									:svValidations="computedSVValidations"
									:validations="computedRules(field.rules)"
									:defaultValue="field.op1Value"
									@enter="onEnterField({ event: $event, field, fields, fieldsIndex, fieldIndex })"
									@focus="onFocusField({ event: $event, field, fields, fieldsIndex, fieldIndex })"
									@input="onInputField($event, field, fieldsIndex, fieldIndex)"
									@keyup.38="onDnKey({ event: $event, field, fields, fieldsIndex, fieldIndex })"
									@dropdown="onDropdownField($event, field, fieldsIndex, fieldIndex)"
									@dropdownUp="onDropdownUpField($event, field, fieldsIndex, fieldIndex)"
									@ruleClick="handleFieldRulesDialog"
								>
									<div v-if="field.answer">
										<span style="color: black">正确答案：</span>
										<span style="color: red">{{ field.answerValue }}</span>
										<p style="color: black; margin-bottom: 0px">解析：</p>
										<div style="color: #3d90fe">{{ field.analysis }}</div>
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
		<!-- <field-rules-dialog ref="fieldRules" :fieldName="fieldName" :proCode="bill.proCode"></field-rules-dialog> -->
		<!-- 字段规则 END -->
	</div>
</template>

<script>
import OpMixins from "./OpMixins";
import OpDropdownMixins from "./OpDropdownMixins";
import Op1Op2OpqMixins from "./Op1Op2OpqMixins";
import OpSpecificValidationsMixins from "../mixins/OpSpecificValidationsMixins";
import ScrollUpDnMixins from "../mixins/ScrollUpDnMixins";
import axios from "axios";
import { localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

export default {
	name: "Opp",
	mixins: [OpMixins, OpDropdownMixins, Op1Op2OpqMixins, OpSpecificValidationsMixins, ScrollUpDnMixins],

	data() {
		return {
			formId: "Opp",
			op: "opp",
			// 一码请求
			instance: "",
			// 图片拼接路径
			imgPath: "",
			// 图片展示路径
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

	watch: {
		"block.picture": {
			async handler(val) {
				console.log(111111, this.fieldsList);
				if (this.bills.type == 1) this.imgPath = `${this.outUrl}${this.bills.downloadPath}${val}`;
				else this.imgPath = `${this.fileUrl}${val}`;
				let item = await this.transform(this.imgPath);
				this.getReader(item).then(res => {
					this.imgTrail = res;
				});
			}
		}
	},

	computed: {
		computedIncludes() {
			return field => {
				return this.svDropdownFields[field.code]?.desserts || [];
			};
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
		}
	},

	components: {
		"op-text-field": () => import("./opTextField"),
		"op-fields-title": () => import("./opFieldsTitle"),
		"watch-image": () => import("../components/watchImage")
		// "field-rules-dialog": () => import("../components/fieldRulesDialog")
	}
};
</script>

<style scoped lang="scss">
@import "../op.scss";
</style>
