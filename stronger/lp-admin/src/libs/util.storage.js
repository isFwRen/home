const storage = {
	set(key, value) {
		localStorage.setItem(key, JSON.stringify(value));
	},

	get(key) {
		return JSON.parse(localStorage.getItem(key));
	},

	del(key) {
		localStorage.removeItem(key);
	}
};

export { storage };
export default storage;
