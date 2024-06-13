<template>
	<div class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="9" class="full image-contain">
				<v-progress-linear v-model="progress" color="#3894ff" height="10" style="margin-bottom: 2px; border-radius: 5px">
					<template v-slot:default="{ value }">
						<strong>{{ Math.ceil(value) }}%</strong>
					</template>
				</v-progress-linear>
				<!-- <watch-image
					:src="`https://img1.baidu.com/it/u=3893389324,4043822814&fm=253&fmt=auto&app=120&f=JPEG?w=1280&h=800`"
				></watch-image> -->
				<watch-image :src="imgTrail"></watch-image>
			</v-col>

			<v-col id="op2FieldsNode" class="full fields-wrapper" cols="3">
				<v-form ref="opeForm" v-model="valid" lazy-validation @submit="preventForm">
					<template v-if="fieldsList.length">
						<div v-for="(fields, fieldsIndex) in fieldsList" :key="`fields_${fieldsIndex}`">
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field v-if="field.problemType == 1" :key="field.assessmentSingleId" :autofocus="field.autofocus"
									class="mb-n4" :field="field" :fieldsIndex="fieldsIndex" :fieldsList="fieldsList"
									:items="field.items || []" :label="`${field.fieldIndex + 1}-${field.question}`" :op="op"
									:defaultValue="field.optionList[0] || ''" :id="field.assessmentSingleId"
									@enter="onEnterField({ event: $event, field, fields, fieldsIndex, fieldIndex })"
									@focus="onFocusField({ event: $event, field, fields, fieldsIndex, fieldIndex })"
									@input="onInputField($event, field, fieldsIndex, fieldIndex)"
									@keyup.38="onDnKey({ event: $event, field, fields, fieldsIndex, fieldIndex })"
									@dropdown="onDropdownField($event, field, fieldsIndex, fieldIndex)"
									@dropdownUp="onDropdownUpField($event, field, fieldsIndex, fieldIndex)">
									<div class="message" v-if="errorIndex == fieldIndex">
										{{ "该题必填，尚未未作答" }}
									</div>
								</op-text-field>
								<field-radio v-if="field.problemType == 2" :key="field.assessmentSingleId" :field="field"
									:fieldsIndex="fieldsIndex" :fieldsList="fieldsList" :label="`${field.fieldIndex + 1}-${field.question}`"
									:op="op" :radioOptions="field.optionList" :id="field.assessmentSingleId"
									@change="radioChange($event, fieldsIndex, fieldIndex)">
									<div class="messages" v-if="errorIndex == fieldIndex">
										{{ "该题必填，尚未未作答" }}
									</div>
								</field-radio>
								<field-check-box v-if="field.problemType == 3" :key="field.assessmentSingleId" :field="field"
									:fieldsIndex="fieldsIndex" :fieldsList="fieldsList" :label="`${field.fieldIndex + 1}-${field.question}`"
									:op="op" :boxOptions="field.optionList" :id="field.assessmentSingleId"
									@change="boxChange($event, fieldsIndex, fieldIndex)">
									<div class="messages" v-if="errorIndex == fieldIndex">
										{{ "该题必填，尚未未作答" }}
									</div>
								</field-check-box>
							</template>
						</div>
					</template>
				</v-form>

				<div class="bottom">
					<v-btn depressed color="primary" v-if="page < pages - 1" @click="nextPage"> 下一页 </v-btn>
					<v-btn depressed color="primary" v-if="page == pages - 1" @click="fEightSubmit"> 提交 </v-btn>
				</div>
			</v-col>
		</v-row>

		<v-dialog v-model="dialog" transition="dialog-top-transition" max-width="600" persistent>
			<template>
				<v-card>
					<v-toolbar color="primary" dark style="font-size: 22px">考核结果</v-toolbar>
					<div class="p">
						<p v-if="examResult.isPass == true">
							您本次的考核分数为<span style="color: red">{{ examResult.point }}</span>，恭喜你已经达到考核要求，管理员近期会联系安排您上岗。
						</p>
						<p v-else>
							您本次的考核分数为<span style="color: red">{{ examResult.point
							}}</span>，还未达到考核要求，请继续进行考核，点击确定，关闭弹出框，重新点击考核领取新的考核题进行录入。
						</p>
					</div>
					<v-card-actions class="justify-end">
						<v-btn color="primary" @click="close(examResult.isPass)">确定</v-btn>
					</v-card-actions>
				</v-card>
			</template>
		</v-dialog>

		<!-- 字段规则 BEGIN -->
		<!-- <field-rules-dialog ref="fieldRules" :fieldName="fieldName" :proCode="bill.proCode"></field-rules-dialog> -->
		<!-- 字段规则 END -->
	</div>
</template>

<script>
import opMMixins from "../mixins/opMMixins";
// import OpDropdownMixins from "../mixins/OpDropdownMixins";
import OpeMixins from "./OpeMixins";
import axios from "axios";
import { localStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";

export default {
	name: "Ope",
	mixins: [OpeMixins, opMMixins],
	// mixins: [opMMixins, OpDropdownMixins, OpeMixins],

	data() {
		return {
			formId: "Ope",
			op: "ope",
			valid: false,
			dialog: false,
			progress: 0,
			// 图片拼接路径
			imgPath: "",
			// 图片展示路径
			imgTrail: ""
		};
	},

	methods: {
		// 图片转Base64格式
		async transform(el) {
			const token = localStorage.get("token");
			const user = localStorage.get("user");
			let code = "";
			const secret = localStorage.get("secret") || "";
			if (secret) {
				code = lpTools.GetCode(secret);
			}
			let res = await axios.get(el, {
				responseType: "blob",
				headers: {
					"x-code": String(code),
					"x-token": token,
					"x-user-id": user.id
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

		close(bol) {
			this.dialog = false;
			if (bol) {
				this.$router.push(`/main/exam/examChannel`);
			} else {
				this.$router.push(`/main/exam/examChannel?proCode=${this.$route.query.proCode}&op=${this.op}`);
			}
		}
	},

	watch: {
		fieldImg: {
			async handler(val) {
				this.imgPath = `${this.fileUrl}${val}`;
				console.log(this.imgPath);
				let item = await this.transform(this.imgPath);
				this.getReader(item).then(res => {
					this.imgTrail = res;
				});
			},
			immediate: true
		}
	},

	components: {
		"op-text-field": () => import("./opTextField"),
		"watch-image": () => import("../components/watchImage"),
		"field-radio": () => import("../components/fieldRadio"),
		"field-check-box": () => import("../components/fieldCheckBox")
	}
};
</script>

<style scoped lang="scss">
@import "../op.scss";

.bottom {
	display: flex;
	justify-content: center;
	margin-top: 30px;
}

.p {
	display: flex;
	justify-content: center;
	font-size: 25px;
	margin-top: 25px;

	p {
		width: 90%;
	}
}

.message {
	margin-bottom: 10px;
	color: red;
}

.messages {
	color: red;
}

::v-deep .v-messages {
	display: none;
}

::v-deep .v-text-field__details {
	display: none;
}
</style>