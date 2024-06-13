<template>
	<main-layout>
		<router-view></router-view>
	</main-layout>
</template>

<script>
import NotificationMixins from "./NotificationMixins";

export default {
	name: "LPMain",
	mixins: [NotificationMixins],
	data() {
		return {
			messages: []
		};
	},
	provide() {
		return {
			socketMessages: this.messages
		};
	}
};
</script>
