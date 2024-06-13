<template>
	<lp-dialog ref="dialog" title="查看配置" :width="800" @dialog="handleDialog">
		<div slot="main">
			<pre v-highlightjs="xml">
        <code class="highlight_s"></code>
      </pre>
		</div>

		<div class="mt-n14 mr-2" slot="actions">
			<z-btn color="primary" @click="onClose">关闭</z-btn>
		</div>
	</lp-dialog>
</template>

<script>
import DialogMixins from "@/mixins/DialogMixins";

export default {
	name: "ViewConfigDialog",
	mixins: [DialogMixins],

	props: {
		xml: {
			type: String,
			default: ""
		}
	}
};
</script>
