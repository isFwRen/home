const headersOne = [
	{ text: "名称", value: "name" },
	{ text: "周一", value: "dayOfWeek1" },
	{ text: "周二", value: "dayOfWeek2" },
	{ text: "周三", value: "dayOfWeek3" },
	{ text: "周四", value: "dayOfWeek4" },
	{ text: "周五", value: "dayOfWeek5" },
	{ text: "周六", value: "dayOfWeek6" },
	{ text: "周日", value: "dayOfWeek0" }
];
const headerOneName = ["开始时间", "结束时间", "间隔时间"];

const headerOneTwo = ["固定时间"];
const headersTwo = [
	{ text: "名称", value: "name" },
	{ text: "周一", value: "dayOfWeek1" },
	{ text: "周二", value: "dayOfWeek2" },
	{ text: "周三", value: "dayOfWeek3" },
	{ text: "周四", value: "dayOfWeek4" },
	{ text: "周五", value: "dayOfWeek5" },
	{ text: "周六", value: "dayOfWeek6" },
	{ text: "周日", value: "dayOfWeek0" }
];

const FirdefaultVlue = [
	{ id: "1", name: "开始时间" },
	{ id: "2", name: "结束时间" },
	{ id: "3", name: "间隔时间" }
];

const SecdefaultVlue = [{ key: 1, groupId: "", name: "固定时间", type: "sendTime" }];

function ONEToDataOne(oneArr, key) {
	const groupId = oneArr[0].groupId;
	let FirdefaultVlue = [
		{ key, groupId: groupId, name: "开始时间", type: "startTime", line: 0 },
		{ key, groupId: groupId, name: "结束时间", type: "endTime", line: 1 },
		{ key, groupId: groupId, name: "间隔时间(min)", type: "interval", line: 2 }
	];
	oneArr.forEach(e => {
		for (let fir in FirdefaultVlue) {
			FirdefaultVlue[fir]["dayOfWeek" + e.dayOfWeek] = e[FirdefaultVlue[fir].type];
		}
	});
	return FirdefaultVlue;
}
function TWOToDataTwo(twoArr, key) {
	const groupId = twoArr[0].groupId;
	let SecdefaultVlue = { key, groupId: groupId, name: "固定时间", type: "sendTime" };
	twoArr.forEach(e => {
		SecdefaultVlue["dayOfWeek" + e.dayOfWeek] = e.sendTime;
	});
	return SecdefaultVlue;
}

export default {
	headersOne,
	headerOneName,
	headerOneTwo,
	headersTwo,
	FirdefaultVlue,
	SecdefaultVlue,
	ONEToDataOne,
	TWOToDataTwo
};
