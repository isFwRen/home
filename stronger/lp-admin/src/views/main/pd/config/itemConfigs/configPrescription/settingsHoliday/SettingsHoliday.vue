<template>
	<div class="settings-holiday">
		<lp-calendar
			clearSelectedItems
			disabledGray
			:defaultValue="selectedItems"
			@select="getCalendar"
			@change:date="changeCalendar"
		></lp-calendar>
	</div>
</template>

<script>
import TableMixins from "@/mixins/TableMixins";

const year = new Date().getFullYear();
const month = new Date().getMonth();

export default {
	name: "SettingsHoliday",
	mixins: [TableMixins],

	data() {
		return {
			currentDate: {},
			id: null,
			selectedItems: [],
			items: []
		};
	},

	methods: {
		/**
		 * @description 当前日历
		 * @param {object} value
		 */
		async getCalendar(value) {
			const selectedItems = [];

			if (!value) return;

			this.currentDate = value;

			for (let day = 1; day <= value.days; day++) {
				const weekend = this.getDayOfWeek(value.year, value.month - 1, day);

				if (weekend === 0 || weekend === 6) {
					selectedItems.push({
						year: value.year,
						month: value.month,
						day: day,
						selected: true
					});
				}
			}

			this.selectedItems = selectedItems;

			const params = {
				startDate: "" + value.year + value.month,
				endDate: "" + value.year + value.month
			};

			const result = await this.$store.dispatch("GET_CONFIG_PRESCRIPTION_HOLIDAY", params);

			if (result.code === 200) {
				const item = result.data.list[0] || {};
				this.items = [];

				this.id = item.ID;

				if (this.id) {
					for (let key in item) {
						if (key.includes("day")) {
							this.items.push({
								year: value.year,
								month: value.month,
								day: +key.substr(3),
								selected: item[key]
							});
						}
					}
				} else {
					for (let day = 1; day <= value.days; day++) {
						const weekend = this.getDayOfWeek(value.year, value.month - 1, day);
						const selected = weekend === 0 || weekend === 6 ? true : false;

						this.items.push({
							year: value.year,
							month: value.month,
							day: day,
							selected
						});
					}
				}
			}

			if (!this.items.length) return;

			this.selectedItems = [];

			for (let item of this.items) {
				if (item.selected) {
					this.selectedItems.push({ ...item });
				}
			}

			console.log(this.selectedItems);
		},

		async changeCalendar(item, items) {
			const { year, month, day, days, selected } = item;

			for (let record of this.items) {
				if (day === record.day) {
					if (selected) {
						record.selected = true;
					} else {
						record.selected = false;
					}
				}
			}

			const holiday = {
				id: this.id,
				date: "" + year + month
			};

			for (let record of this.items) {
				holiday[`day${record.day}`] = record.selected;
			}

			const result = await this.$store.dispatch(
				"UPDATE_CONFIG_PRESCRIPTION_HOLIDAY",
				holiday
			);
			this.toasted.dynamic(result.msg, result.code);
		},

		/**
		 * @description 某年某月某日是星期几
		 * @param {number} 年 year
		 * @param {number} 月 month
		 * @param {number} 日 date = 1
		 */
		getDayOfWeek(year, month, date = 1) {
			return new Date(year, month, date).getDay();
		}
	}
};
</script>
