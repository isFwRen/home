const regex_phone = /^1(?:3\d|4[4-9]|5[0-35-9]|6[67]|7[013-8]|8\d|9\d)\d{8}$/;

const usernameField = {
	formKey: "account",
	label: "工号",
	validation: [{ rule: "required", message: "请输入工号!" }]
};

const phoneField = {
	formKey: "account",
	label: "手机号",
	rules: [
		value => !!value || "请输入手机号.",
		value => regex_phone.test(value) || "手机格式不正确"
	]
};

const passwordField = {
	formKey: "password",
	label: "密码",
	rules: [value => !!value || "请输入密码."]
};

const captchaField = {
	formKey: "captcha",
	label: "验证码",
	rules: []
};

export const authImages = {
	imgUrl1: require("@/assets/images/usage/loginGuide/authenGuide/authImg-1.png"),
	imgUrl2: require("@/assets/images/usage/loginGuide/authenGuide/authImg-2.png"),
	imgUrl3: require("@/assets/images/usage/loginGuide/authenGuide/authImg-3.png"),
	imgUrl4: require("@/assets/images/usage/loginGuide/authenGuide/authImg-4.png")
};

export { regex_phone };

export default {
	usernameField,
	phoneField,
	passwordField,
	captchaField,
	authImages
};
