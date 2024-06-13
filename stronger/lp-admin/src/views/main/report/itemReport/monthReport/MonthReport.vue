<template>
	<div>
		<DayReport :type="type" />
	</div>
</template>

<script>
import DayReport from "../dayReport/DayReport.vue";
export default {
	components: { DayReport },
	data() {
		return {
			type: 3
		};
	}
};
</script>
