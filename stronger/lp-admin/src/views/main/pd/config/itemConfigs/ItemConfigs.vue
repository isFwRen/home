<template>
	<div class="item-configs">
		<!-- <lp-dialog
      ref="dialog"
      :title="$route.meta.title"
      fullscreen
      persistent
      @close="onConfigClose"
    >
      <div class="pt-6" slot="main">
        <router-view></router-view>
      </div>
    </lp-dialog> -->

		<lp-message-dialog ref="dialog" :title="$route.meta.title" @close="onConfigClose">
			<div class="pt-6" slot="main">
				<router-view></router-view>
			</div>
		</lp-message-dialog>
	</div>
</template>

<script>
import { mapGetters } from "vuex";
import DialogMixins from "@/mixins/DialogMixins";
import ConfigMixins from "../ConfigMixins";

export default {
	name: "ItemConfig",
	mixins: [DialogMixins, ConfigMixins],

	methods: {
		onConfigClose() {
			this.rememberIds({ tempId: "" });
			this.$router.push("/main/PD/config");
			this.$emit("close");
		}
	},
	computed: {
		...mapGetters(["config"])
	},
	watch: {
		$route: {
			handler(route) {
				if (route.meta.path) {
					this.$nextTick(() => {
						this.$refs.dialog.onOpen();
					});
				}
			},
			immediate: true
		}
	}
};
</script>
