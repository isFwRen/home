<template>
	<div class="new-dialog">
		<lp-dialog ref="dialog" :title="title" :width="1000" @dialog="handleDialog">
			<div slot="main">
				<v-row style="font-weight: 800; font-size: 16px; padding-top: 15px">
					<v-col v-for="item in headers" :key="item.value">
						{{ item.text }}
					</v-col>
					<v-col cols="1">
						<v-btn icon @click="handleAdd">
							<v-icon>mdi-plus-circle-outline</v-icon>
						</v-btn>
					</v-col>
				</v-row>

				<v-row
					align="center"
					class="solid rounded-br-0"
					v-for="(forms, index) in formsList"
					:key="index"
					style="height: 75px"
				>
					<v-col v-for="form in forms" :key="form.value">
						<v-text-field v-model.trim="form.val"></v-text-field>
					</v-col>
					<v-col cols="1">
						<v-btn icon @click="handleDel(index)">
							<v-icon>mdi-trash-can-outline</v-icon>
						</v-btn>
					</v-col>
				</v-row>
			</div>
			<div class="z-flex" slot="actions">
				<v-btn class="mr-3" color="normal" @click="onClose">取消</v-btn>
				<v-btn color="primary" @click="onConfirm">确认</v-btn>
			</div>
		</lp-dialog>
	</div>
</template>
<script>
import DialogMixins from "@/mixins/DialogMixins";
import _ from "lodash";
export default {
	name: "NewDialog",
	mixins: [DialogMixins],
	props: {
		headers: {
			type: Array,
			required: true
		}
	},
	data() {
		return {
			formsList: [],
			formId: "NewDialog"
		};
	},
	watch: {
		dialog(val) {
			if (val) {
				const arr = _.cloneDeep(this.headers);
				this.formsList.push(arr);
			} else {
				this.formsList = [];
			}
		}
	},
	methods: {
		onConfirm() {
			const arr = [];
			this.formsList.map(item => {
				const obj = {};
				for (let key of item) {
					obj[key["text"]] = key["val"];
				}
				arr.push(obj);
			});
			this.$emit("add", arr);
		},
		handleAdd() {
			this.formsList.push(_.cloneDeep(this.headers));
		},
		handleDel(index) {
			this.formsList.splice(index, 1);
			console.log(this.formsList, index);
		}
	}
};
</script>
