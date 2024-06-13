<template>
	<v-overlay
		:absolute="absolute"
		:color="color"
		:dark="dark"
		:light="light"
		:opacity="opacity"
		:value="value"
		:zIndex="zIndex"
	>
		<Spinners />

		<!-- <div class="z-flex justify-center">
      <scaling-squares-spinner
        :animation-duration="spinner.duration"
        :size="spinner.size"
        :color="spinner.color"
      >
      </scaling-squares-spinner>
    </div> -->

		<div class="mt-4">
			<slot name="tips"></slot>
		</div>
	</v-overlay>
</template>

<script>
// import { AtomSpinner } from 'epic-spinners'
// import { ScalingSquaresSpinner } from 'epic-spinners'
import Spinners from "./Spinners";

export default {
	name: "LPSpinners",

	props: {
		absolute: {
			type: Boolean,
			default: false
		},

		color: {
			type: String,
			default: "#212121"
		},

		dark: {
			type: Boolean,
			default: true
		},

		delay: {
			type: Number,
			default: 1000
		},

		light: {
			type: Boolean,
			default: false
		},

		opacity: {
			type: Number,
			default: 0.46
		},

		overlay: {
			type: Boolean,
			default: false
		},

		timeout: {
			type: Number,
			default: 10000
		},

		zIndex: {
			type: Number,
			default: 7
		}
	},

	data() {
		return {
			value: false,
			startTime: null,
			endTime: null,

			spinner: {
				duration: 1250,
				size: 80,
				color: "#fff"
			}
		};
	},

	methods: {
		interval() {
			if (this.overlay) {
				this.startTime = Date.now();
			} else {
				this.endTime = Date.now();
			}

			const timer = setInterval(() => {
				this.endTime = Date.now();
				const diff = this.endTime - this.startTime;

				if (diff < this.delay) {
					this.value = false;
				} else {
					this.value = this.overlay;

					if (!this.overlay) {
						clearInterval(timer);
					}
				}
			}, 200);
		}
	},

	watch: {
		overlay: {
			handler() {
				this.interval();
			},
			immediate: true
		},

		value: {
			handler(value) {
				if (value) {
					const timer = setTimeout(() => {
						this.value = false;
						clearTimeout(timer);
					}, this.timeout);
				}
			},
			immediate: true
		}
	},

	components: {
		// AtomSpinner,
		// ScalingSquaresSpinner,
		Spinners
	}
};
</script>
