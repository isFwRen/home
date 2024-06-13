<template>
	<div class="lp-dialog">
		<v-dialog
			v-model="dialog"
			:close-delay="closeDelay"
			:content-class="contentClass"
			:fullscreen="fullscreen"
			:hide-overlay="hideOverlay"
			:min-width="minWidth"
			:max-width="maxWidth"
			no-click-animation
			:transition="transition"
			:width="width"
			:persistent="persistent"
		>
			<slot name="card" v-if="$slots.card"></slot>

			<v-card v-else>
				<slot name="toolbar">
					<v-toolbar
						class="z-toolbar"
						:color="toolbarColor"
						dark
						:rounded="false"
						:style="{ width: fullscreen ? '100%' : `${width}px` }"
					>
						<v-toolbar-title>{{ title }}</v-toolbar-title>
						<v-spacer></v-spacer>
						<!-- <span class="hover-bgcolor" icon dark @click="onClose">
              <v-icon>mdi-close</v-icon>
            </span> -->
						<v-btn icon dark @click="onClose">
							<v-icon>mdi-close</v-icon>
						</v-btn>
					</v-toolbar>
				</slot>

				<v-card-text
					:class="[{ 'pt-16': !cardTextClass }, cardTextClass]"
					:style="cardTextStyle"
				>
					<slot name="main"></slot>
				</v-card-text>

				<v-card-actions v-if="$slots.actions" class="z-flex justify-end">
					<slot name="actions"></slot>
				</v-card-actions>

				<div v-if="$slots.bottom" class="z-card-bottom">
					<slot name="bottom"></slot>
				</div>
			</v-card>
		</v-dialog>
	</div>
</template>

<script>
export default {
	name: "LPDialog",
	props: {
		cardTextClass: {
			type: String,
			default: ""
		},

		cardTextStyle: {
			type: Object,
			required: false
		},

		closeDelay: {
			type: [Number, String],
			default: 0
		},

		contentClass: {
			type: String,
			required: false
		},

		fullscreen: {
			type: Boolean,
			default: false
		},

		hideOverlay: {
			type: Boolean,
			default: false
		},

		maxWidth: {
			type: [Number, String],
			default: 1200
		},

		minWidth: {
			type: [Number, String],
			default: 300
		},

		persistent: {
			type: Boolean,
			default: false
		},

		title: {
			type: String,
			default: ""
		},

		toolbarColor: {
			type: String,
			default: "primary"
		},

		transition: {
			type: String,
			default: "fab-transition"
		},

		width: {
			type: [Number, String],
			default: 700
		}
	},

	data() {
		return {
			dialog: false
		};
	},
	methods: {
		onOpen() {
			this.dialog = true;
			this.$emit("open", this.dialog);
		},

		onClose() {
			this.dialog = false;
			this.$emit("close", this.dialog);
		},

		close() {
			this.dialog = false;
			this.$emit("close", this.dialog);
		},

		onToggle() {
			this.dialog = !this.dialog;
			this.$emit("toggle", this.dialog);
		}
	},

	watch: {
		dialog: {
			handler(dialog) {
				this.$emit("dialog", dialog);
			},
			immediate: true
		}
	}
};
</script>

<style lang="scss">
.z-toolbar {
	position: fixed;
	width: 100%;
	z-index: 10;
}

.z-card-bottom {
	position: absolute;
	bottom: 0;
	left: 0;

	width: 100%;
}

// .hover-bgcolor {
//     transition: all 0.2s;
//     border-radius: 50%;
//     cursor: pointer;
// }
</style>
