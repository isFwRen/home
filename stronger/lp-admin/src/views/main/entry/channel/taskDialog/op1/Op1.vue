<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="9" class="full image-contain">
				<watch-image :src="`${fileUrl}${bill.downloadPath}${block.picture}`"></watch-image>
			</v-col>

			<v-col class="full fields-wrapper" cols="3" id="fieldsWrapper">
				<h3 class="op-title">{{ block.name }}</h3>

				<v-form
					ref="op1Form"
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
									:bill="bill"
									:block="block"
									:disabled="field.disabled"
									:field="field"
									:fieldsIndex="focusFieldsIndex"
									:fieldsList="fieldsList"
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
									:defaultValue="field.op1Value"
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
								></op-text-field>
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
	name: "Op1",
	mixins: [
		OpMixins,
		OpDropdownMixins,
		Op1Op2OpqMixins,
		OpSpecificValidationsMixins,
		ScrollUpDnMixins
	],

	data() {
		return {
			formId: "Op1",
			op: "op1"
		};
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
