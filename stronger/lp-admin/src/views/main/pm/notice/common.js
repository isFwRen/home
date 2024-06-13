function getDate() {
	const date = new Date();
	this.getYear = (n = 0) => {
		let year = date.getFullYear() + n;
		return year;
	};
	this.getMonth = (n = 0) => {
		let month = date.getMonth() + n;
		return month < 9 ? "0" + (month + 1) : month + 1;
	};
	this.getDay = (n = 0) => {
		let day = date.getDay() + n;
		return day < 10 ? "0" + (day + 1) : day + 1;
	};
	this.getDate = ({ format, day = 0, year = 0, month = 0 }) => {
		const map = {
			mm: this.getMonth(month),
			yy: this.getYear(year).toString().slice(-2),
			dd: this.getDay(day),
			yyyy: this.getYear(year)
		};
		return format.replace(/mm|dd|yyyy|yy/gi, matched => map[matched]);
	};
}

export { getDate };
