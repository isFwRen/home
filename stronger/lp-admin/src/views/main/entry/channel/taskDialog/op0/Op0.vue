<template>
	<div v-if="hasOp" class="op op1op2opq">
		<v-row class="full ma-0">
			<v-col cols="3" class="full">
				<v-row v-if="bill.downloadPath" dense>
					<v-col v-for="(img, index) in bill.pictures" :key="`billPho${index}`" cols="4">
						<div
							:class="['thumb-wrapper', thumbIndex === index ? 'actived' : '']"
							@click="onSelectThumb({ thumbIndex: index })"
						>
							<img
								style="max-height: 160px"
								:src="`${fileUrl}${bill.downloadPath}${img}`"
							/>
						</div>
					</v-col>
				</v-row>
			</v-col>

			<v-col ref="drawImageContain" cols="7" class="full image-contain">
				<canvas-toolbar v-if="ctVisible" @select="onSelect"></canvas-toolbar>

				<div class="draw-image-wrapper">
					<z-draw-image
						id="op0DrawImage"
						ref="drawImage"
						:fileName="drewImageName"
						:imageWidth="canvasWidth"
						:src="`${fileUrl}${bill.downloadPath}${modifyImage}`"
						@initialized="initializedImage"
						@save="drewImage"
					></z-draw-image>
				</div>
			</v-col>

			<v-col class="full fields-wrapper" cols="2" id="fieldsWrapper">
				<h3 class="op-title">{{ block.name }}</h3>

				<v-form
					ref="op0Form"
					v-model="valid"
					lazy-validation
					@submit="event => event.preventDefault()"
				>
					<template v-if="fieldsObject[thumbIndex]">
						<div
							v-for="(fields, fieldsIndex) in fieldsObject[thumbIndex].fieldsList"
							:key="`fields_${fieldsIndex}`"
						>
							<template v-for="(field, fieldIndex) in fields">
								<op-text-field
									v-if="field.show !== false"
									:key="field.uniqueKey"
									:accent="field.name === `模板类型字段` ? true : false"
									:autofocus="field.autofocus"
									:field="field"
									:fieldsList="fieldsObject[thumbIndex].fieldsList"
									:id="field.uniqueId"
									:items="field.items ? field.items : []"
									:label="field.name"
									:labelTip="field.code"
									:op="op"
									:svHints="specificProject.op1op2opq.hints"
									:includes="
										svConstantsDD[field.code]
											? svConstantsDD[field.code].desserts
											: []
									"
									:svValidations="specificProject.op0.rules"
									:validations="field.rules"
									:defaultValue="field.op0Value"
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
import { tools } from "vue-rocket";
import OpMixins from "../OpMixins";
import OpDropdownMixins from "../OpDropdownMixins";
import Op0Mixins from "./Op0Mixins";
import Op0CanvasMixins from "./Op0CanvasMixins";
import Op0ShortcutMixins from "./Op0ShortcutMixins";
import OpSpecificValidationsMixins from "../OpSpecificValidationsMixins";
import ScrollUpDnMixins from "../ScrollUpDnMixins";

export default {
	name: "Op0",
	mixins: [
		OpMixins,
		OpDropdownMixins,
		Op0Mixins,
		Op0CanvasMixins,
		Op0ShortcutMixins,
		OpSpecificValidationsMixins,
		ScrollUpDnMixins
	],

	data() {
		return {
			formId: "Op0",
			op: "op0"
		};
	},

	methods: {
		// 选择缩略图
		onSelectThumb({ thumbIndex }) {
			this.thumbIndex = thumbIndex;

			this.$refs.drawImage.clear();

			this.fieldsObject[this.thumbIndex].sessionStorage = false;

			const { initFieldsList } = this.fieldsObject[this.thumbIndex];

			// 未按F4
			for (let thumbIndex in this.fieldsObject) {
				const sessionStorage = this.fieldsObject[thumbIndex].sessionStorage;

				if (!sessionStorage) {
					// const fieldsList = this.fieldsObject[thumbIndex].fieldsList
					// const flatFieldsList = tools.flatArray(fieldsList)

					// for(let field of flatFieldsList) {
					// delete this.memoFields[field.uniqueId]
					// delete this.svMemoFields[field.code]
					// }

					this.fieldsObject[thumbIndex].fieldsList = tools.deepClone(initFieldsList);
				}
			}

			// 按F8后将，唯一一个[模板类型字段]被清除掉
			if (tools.isLousy(this.fieldsObject[this.thumbIndex]?.fieldsList)) {
				this.fieldsObject[this.thumbIndex].fieldsList = tools.deepClone(initFieldsList);
			}

			this.bill.thumbIndex = this.thumbIndex;

			if (tools.isLousy(this.fieldsObject)) {
				this.toasted.warning("字段为空！");
				return;
			}

			this.$emit("bill", this.bill);
		}
	},

	components: {
		"canvas-toolbar": () => import("../components/canvasToolbar"),
		"op-text-field": () => import("../components/opTextField")
	}
};
</script>

<style scoped lang="scss">
@import "../op1op2opq.scss";

/* 缩略图 */
.thumb-wrapper {
	max-height: 160px;
	border: 2px solid transparent;
	text-align: center;
	overflow: hidden;

	&.actived {
		border: 2px solid #1976d2;
	}
}

/* 画布 */
.draw-image-wrapper {
	height: 100%;
}
</style>
