/**
 * @Author: xingqiyi
 * @Description:白名单路由不需要登录的 支持正则
 * @Date: 2022/1/13 9:39 上午
 */

package middleware

var WhiteList = []string{
	"/sys-base/sys-login/login",
	"/sys-base/sys-user-operation-pw/LogWhitelists",
	"/dinging/captcha",
	"socket.io",
	"global-export",
	"/favicon.ico",
	//"^/files/",
	"/sys-user/sys-user-operation-pw/changePassword",
	"/swagger/.*",
	"/pro-manager/bill-list/qqq",
	//"/task/releaseBill",
	//"/task/op",
	//"/task/submit",
	"/metrics",
	"/pig-pig-hungry/.*",
	"/sys-socket-io-notice/customer-notice",
	"/sys-socket-io-notice/business-push",
	"/sys-base/sys-login/get-user-qrCode",
	"/sys-base/login",
	"/sys-base/send-code",
	"/sys-base/forget-pwd",
	"/sys-base/change-pwd-common",

	// "/pro-config/sys-project/move-conf",
	// "/pro-config/sys-field-check/list",
	// "/pro-config/sys-field-check/edit",
	"/practice/.*",
	"/files/common/*",
	//"/pro-config/refresh-pro-conf",
}

var LogWhitelists = []string{
	"socket.io",
	"/favicon.ico",
	"^/files/",
	"/swagger/.*",
	"/metrics",
	"/pro-manager/bill-list/edit-bill-result-data",
}
