import LPModalComponent from "./LPModal";

export default {
	install(Vue) {
		let cancelFn = function () {};
		let confirmFn = function () {};

		const modalCtrl = Vue.extend(LPModalComponent);
		const instance = new modalCtrl();

		instance.$mount(document.createElement("div"));
		document.body.append(instance.$el);

		Vue.prototype.$modal = function (item) {
			instance.visible = item.visible;
			instance.title = item.title;
			instance.content = item.content;
			cancelFn = item.cancel;
			confirmFn = item.confirm;
		};

		Vue.prototype.$modal.cancel = function () {
			cancelFn && cancelFn.call(this);
		};

		Vue.prototype.$modal.confirm = function () {
			confirmFn && confirmFn.call(this);
		};
	}
};
