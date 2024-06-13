<template>
	<div class="settings-prescription">
		<v-expansion-panels v-model="panel" multiple>
			<settings-basic></settings-basic>

			<settings-download></settings-download>

			<settings-upload></settings-upload>
		</v-expansion-panels>
	</div>
</template>

<script>
export default {
	name: "SettingsPrescription",

	data() {
		return {
			panel: [0, 1, 2]
		};
	},

	components: {
		"settings-basic": () => import("./settingsBasic"),
		"settings-download": () => import("./settingsDownload"),
		"settings-upload": () => import("./settingsUpload")
	}
};
</script>
