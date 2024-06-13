<template>
	<lp-dialog ref="dialog" :title="title" width="500" @dialog="handleDialog">
		<div class="pt-6" slot="main">
			<div class="text">
				<label>名称:</label>
				{{ rowInfo.name }}
			</div>

			<div class="text">
				<label>修改时间:</label>
				{{ rowInfo.editTime }}
			</div>

			<div class="text">
				<label>修改人:</label>
				{{ rowInfo.editName }}
			</div>

			<div class="text">
				<label>类型:</label>
				{{ rowInfo.types }}文件
			</div>

			<div class="text">
				<label>大小:</label>
				{{ rowInfo.size }}KB
			</div>
		</div>

		<div class="z-flex" slot="actions">
			<z-btn class="mr-3" color="primary" outlined @click="onClose">关闭</z-btn>
		</div>
	</lp-dialog>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";

export default {
	name: "TeachingExpenseTemplateAttrDialog",
	mixins: [DialogMixins]
};
</script>

<style scoped lang="scss">
.text {
	label {
		display: inline-block;
		margin-bottom: 8px;
		width: 80px;
	}
}
</style>
