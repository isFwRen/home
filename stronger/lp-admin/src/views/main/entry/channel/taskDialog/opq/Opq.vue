<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="9" class="full image-contain">
				<watch-image
					ref="watchImage"
					:op="op"
					:bill="bill"
					:src="`${fileUrl}${bill.downloadPath}${block.picture}`"
				></watch-image>
			</v-col>

			<v-col class="full fields-wrapper" cols="3" id="fieldsWrapper">
				<h3 class="op-title">{{ block.name }}</h3>

				<v-form
					ref="opqForm"
					v-model="valid"
					lazy-validation
					@submit="event => event.preventDefault()"
				>
					<template v-if="fieldsList.length">
						<div
							v-for="(fields, fieldsIndex) in fieldsList"
							:key="`fields_${fieldsIndex}`"
						>
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field
									v-if="focusFieldsIndex === fieldsIndex && field.show"
									:key="field.uniqueKey"
									:autofocus="field.autofocus"
									:disabled="field.disabled"
									:field="field"
									:fieldsIndex="focusFieldsIndex"
									:fieldsList="fieldsList"
									:firstDiffIndex="firstDiffIndex"
									:id="field.uniqueId"
									:items="
										svDropdownFields[field.uniqueId]
											? svDropdownFields[field.uniqueId].items
											: []
									"
									:label="`(${field.code}_${fieldsIndex}-${fieldIndex})${field.name}`"
									:labelTip="field.code"
									:op="op"
									:svHints="specificProject.op1op2opq.hints"
									:includes="
										svConstantsDD[field.code]
											? svConstantsDD[field.code].desserts
											: []
									"
									:svValidations="specificProject.op1op2opq.rules"
									:validations="field.rules"
									:defaultValue="field.opqValue"
									@enter="onEnterField($event, field, fieldsIndex, fieldIndex)"
									@focus="onFocusField($event, field, fieldsIndex, fieldIndex)"
									@input="onInputField($event, field, fieldsIndex, fieldIndex)"
									@keydown.38="onDnKey($event, field, fieldsIndex, fieldIndex)"
									@dropdownEnter="
										onDropdownEnterField($event, field, fieldsIndex, fieldIndex)
									"
									@dropdownUp="
										onDropdownUpField($event, field, fieldsIndex, fieldIndex)
									"
								>
									<div :key="field.uniqueKey" class="diff">
										<p
											class="mb-0 op2"
											v-html="
												_tools.compareString(field.op2Value, field.op1Value)
													.targetHtml
											"
										></p>
										<p
											class="mb-0 op1"
											v-html="
												_tools.compareString(field.op1Value, field.op2Value)
													.targetHtml
											"
										></p>
									</div>
								</op-text-field>
							</template>
						</div>
					</template>
				</v-form>
			</v-col>
		</v-row>
	</div>
</template>

<script>
import OpMixins from "../OpMixins";
import OpDropdownMixins from "../OpDropdownMixins";
import Op1Op2OpqMixins from "../Op1Op2OpqMixins";
import OpSpecificValidationsMixins from "../OpSpecificValidationsMixins";
import ScrollUpDnMixins from "../ScrollUpDnMixins";

export default {
	name: "Opq",
	mixins: [
		OpMixins,
		OpDropdownMixins,
		Op1Op2OpqMixins,
		OpSpecificValidationsMixins,
		ScrollUpDnMixins
	],

	data() {
		return {
			formId: "Opq",
			op: "opq",
			firstDiffIndex: -1
		};
	},

	methods: {
		// 找到一码二码第一个不同的值，并在 field 选中
		opqGetFieldFirstDiffIndex(field, symbol = "?") {
			if (field.opqValue) {
				this.firstDiffIndex = field.opqValue.indexOf(symbol);
			} else {
				this.firstDiffIndex = -1;
			}
		},

		// 设置默认值
		opqSetFieldEffectKeyValue(field) {
			// 前端自行根据一码二码的值动态设置问题件 opqValue 的值(当第一次领取时)
			if (field.disabled == false && field.show === true) {
				const initTime = this.block.opqSubmitAt?.slice(0, 10);

				if (initTime === "0001-01-01") {
					field.opqValue = this._tools.compareString(
						field.op2Value,
						field.op1Value
					).diffValue;
					field.resultValue = field.opqValue;
				}
			}
		},

		navToViewImages() {
			this.$refs.watchImage.navToViewImages();
		}
	},

	components: {
		"op-text-field": () => import("../components/opTextField"),
		"watch-image": () => import("../components/watchImage")
	}
};
</script>

<style scoped lang="scss">
@import "../op1op2opq.scss";
</style>
