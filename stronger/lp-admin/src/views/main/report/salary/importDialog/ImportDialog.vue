<template>
	<div class="import-dialog">
		<lp-dialog ref="dialog" title="导入" width="500" @dialog="handleDialog">
			<div class="pt-6" slot="main">
				<v-row>
					<v-col
						v-for="(item, index) in cells.fields"
						:key="`updateExport_${index}`"
						:cols="12"
					>
						<template v-if="item.inputType === 'date'">
							<!--
              <z-date-picker
                :formId="formId"
                :formKey="item.formKey"
                :hideDetails="item.hideDetails"
                :hint="item.hint"
                :label="item.label"
                :options="item.options"
                :suffix="item.suffix"
                :validation="item.validation"
                :defaultValue="item.defaultValue"
              >
                <span class="error--text" slot="prepend-outer">{{ item.prependOuter }}</span>
              </z-date-picker>
            -->
						</template>

						<template v-else-if="item.inputType === 'fileInput'">
							<z-file-input
								:formId="formId"
								:formKey="item.formKey"
								:label="item.label"
								action="http://113.106.108.93:13000/api/report-management/internal/salary/upload"
								hide-details
								placeholder="点击文件或将文件拖拽到这里"
								class="mt-n3 mr-3"
								:deleteIcon="false"
								parcel
								multiple
							>
								<span class="error--text" slot="prepend-outer">{{
									item.prependOuter
								}}</span>
							</z-file-input>
						</template>
					</v-col>
				</v-row>
			</div>

			<div class="z-flex mt-n4" slot="actions">
				<z-btn class="mr-3" color="normal" @click="onClose">取消</z-btn>

				<z-btn :formId="formId" btnType="validate" color="primary" @click="onConfirm"
					>确认</z-btn
				>
			</div>
		</lp-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import cells from "./cells";

export default {
	name: "ImportDialog",
	mixins: [DialogMixins],

	data() {
		return {
			formId: "importDialog",
			cells,
			action: "",
			dispatchUpload: "IMPORT_SALARY_INSIDE_UPLOAD"
		};
	},

	methods: {
		onConfirm() {
			const formData = new formData();

			const file = this.$store.dispatch("IMPORT_SALARY_INSIDE_UPLOAD");
		}
	},

	computed: {
		...mapGetters(["config"])
	}
};
</script>
