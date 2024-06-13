<template>
	<lp-dialog ref="dialog" title="问题件配置" :width="960" @dialog="handleDialog">
		<div slot="main">
			<div class="py-4">
				<p style="font-size: 15px">字段：{{ Row.code + "_" + Row.name }}</p>
				<div class="z-row header" style="font-weight: 800; font-size: 16px">
					<div class="z-col-4">录入值</div>
					<div class="z-col-4">转换值</div>
					<div class="z-col-4">问题件编码</div>
					<div class="z-col-11">问题件描述</div>
					<div class="z-col-1 z-flex justify-end">
						<v-btn icon @click="handleAdd">
							<v-icon>mdi-plus-circle-outline</v-icon>
						</v-btn>
					</div>
				</div>

				<div
					v-for="(item, index) in this.desserts[0].list"
					:key="`opq_config_${index}`"
					class="z-row gutter-x16 mb-4 solid"
				>
					<div class="z-col-4">
						<v-text-field v-model="desserts[0].list[index].inputVal"></v-text-field>
					</div>

					<div class="z-col-4">
						<v-text-field v-model="desserts[0].list[index].changeVal"></v-text-field>
					</div>

					<div class="z-col-4">
						<v-text-field v-model="desserts[0].list[index].code"></v-text-field>
					</div>

					<div class="z-col-11">
						<v-text-field v-model="desserts[0].list[index].desc"></v-text-field>
					</div>

					<div class="z-col-1 z-flex justify-end align-center">
						<v-btn icon @click="handleDelete(index)">
							<v-icon>mdi-trash-can-outline</v-icon>
						</v-btn>
					</div>
				</div>
			</div>
		</div>

		<div class="mt-n6 mr-2" slot="actions">
			<z-btn class="mr-4" @click="onClose">取消</z-btn>

			<z-btn color="primary" @click="onConfirm(desserts[0])">确认</z-btn>
		</div>
	</lp-dialog>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";

export default {
	name: "OpqConfigDialog",
	mixins: [DialogMixins],
	props: ["Row", "issueList"],

	data() {
		return {
			desserts: [{ list: [] }]
		};
	},

	watch: {
		issueList(newvalue) {
			this.desserts[0].fId = this.Row.ID;
			this.desserts[0].list = [...newvalue.list];
		},
		immediate: true
	},
	methods: {
		handleAdd() {
			this.desserts[0].fId = this.Row.ID;
			this.desserts[0].list.push({
				createdAt: this.Row.CreatedAt,
				updatedAt: this.Row.UpdatedAt,
				fid: this.Row.ID,
				proId: this.Row.proId
			});
		},

		handleDelete(index) {
			this.desserts[0].list.splice(index, 1);
			console.log(this.desserts[0]);
		}
	}
};
</script>
