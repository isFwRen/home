const B0118 = {
	op0: {
		// 记录最后一次存储的合法field
		memoFields: ["fc054", "fc180", "fc181"],

		// fields 的值从 targets 里的值选择
		dropdownFields: [
			{
				targets: ["fc059"],
				fields: ["fc060", "fc061", "fc201"]
			}
		],

		// 校验规则
		rules: [
			// 6
			{
				fields: ["fc060", "fc061", "fc201"],
				validate6: function ({ value, items }) {
					const result = items.find(text => text.includes(value));

					if (result) {
						return true;
					} else {
						return "没有此发票，请核实！";
					}
				}
			},

			// 23
			{
				fields: ["fc054"],
				validate23: function ({ value, items }) {
					const result = items.find(text => text.includes(value));

					if (result) {
						return true;
					} else {
						return "录入内容不在代码表中，请录入“本代码表中不存在的其他医院”.";
					}
				}
			},

			// 33
			{
				fields: ["fc180"],
				validate33: function ({ value, items }) {
					// 录入内容为F时不做校验
					if (value === "F") {
						return true;
					}

					const result = items.find(text => text.includes(value));

					if (result) {
						return true;
					} else {
						return "区不在常量表中时，录入F.";
					}
				}
			}
		],

		// 字段已生成
		init: {
			methods: {
				// 14
				validate14: function ({ bill }) {
					const firstTwo = bill.billNum.slice(0, 2);

					if (firstTwo === "30") {
						alert("收据上有手术费的，需要切手术！");
					}
				},

				// 23
				validate23: function ({ op, flatFieldsList }) {
					const fields = ["fc054"];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_医院代码表",
								query: "名称"
							};
						}
					});
				},

				// 30
				validate30: function ({ op, flatFieldsList }) {
					const fields = ["fc152"];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_发票大项类型",
								query: "费用名称"
							};
						}
					});
				},

				// 33
				validate33: function ({ op, flatFieldsList }) {
					const fields = ["fc180"];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_地址库",
								query: "市中文名",
								targets: ["省中文名", "市中文名"]
							};
						}
					});
				},

				// 34
				validate34: function ({ op, flatFieldsList }) {
					const fields = ["fc181"];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_地址库",
								query: "省中文名"
							};
						}
					});
				}
			}
		},

		// 回车
		enter: {
			methods: {
				// 33 34
				validate33And34: function ({
					op,
					field,
					fieldsList,
					focusFieldsIndex,
					memoFields
				}) {
					if (field.code === "fc180") {
						const fields = fieldsList[focusFieldsIndex];

						const fc180Value = field.resultValue;
						const fc181Field = fields.find(field => field.code === "fc181");

						if (fc180Value.includes("-")) {
							const values = fc180Value.split("-");

							field[`${op}Value`] = values[1];
							field.resultValue = values[1];
							_.set(memoFields, `${field.uniqueId}.value`, values[1]);

							fc181Field[`${op}Value`] = values[0];
							fc181Field.resultValue = values[0];
							_.set(memoFields, `${fc181Field.uniqueId}.value`, values[0]);
						}
					}
				}
			}
		},

		// F8(提交前校验)
		beforeSubmit: {
			methods: {
				// 2
				validate2({ flatFieldsList }) {
					const fc059Field = flatFieldsList.find(field => field.code === "fc059");

					if (fc059Field) {
						const fc059Values = [];

						for (let field of flatFieldsList) {
							if (field.code === "fc089") {
								if (field.resultValue == 1) {
									return true;
								}
							}

							if (field.code === "fc059") {
								if (fc059Values.includes(field.resultValue)) {
									return {
										errorMessage: "发票属性不能重复."
									};
								} else {
									fc059Values.push(field.resultValue);
								}
							}
						}

						return true;
					}

					return true;
				},

				// 5
				validate5({ flatFieldsList }) {
					const fc060Field = flatFieldsList.find(field => field.code === "fc060");

					if (fc060Field) {
						const [fc059Values, fc060Values] = [[], []];

						flatFieldsList.map(field => {
							if (field.code === "fc059") {
								if (!fc059Values.includes(field.resultValue)) {
									fc059Values.push(field.resultValue);
								}
							}

							if (field.code === "fc060") {
								if (!fc060Values.includes(field.resultValue)) {
									fc060Values.push(field.resultValue);
								}
							}
						});

						for (let value of fc060Values) {
							if (!fc059Values.includes(value)) {
								return {
									errorMessage: `清单明细${value}没有匹配的发票，请检查！`
								};
							}
						}
					}

					return true;
				},

				// 6
				validate6({ flatFieldsList }) {
					const fc061Field = flatFieldsList.find(field => field.code === "fc061");

					if (fc061Field) {
						const [fc059Values, fc061Values] = [[], []];

						flatFieldsList.map(field => {
							if (field.code === "fc059") {
								if (!fc059Values.includes(field.resultValue)) {
									fc059Values.push(field.resultValue);
								}
							}

							if (field.code === "fc061") {
								if (!fc061Values.includes(field.resultValue)) {
									fc061Values.push(field.resultValue);
								}
							}
						});

						for (let value of fc061Values) {
							if (!fc059Values.includes(value)) {
								return {
									errorMessage: `报销单${value}没有匹配的发票，请检查！`
								};
							}
						}
					}

					return true;
				},

				// 7
				validate7({ flatFieldsList }) {
					const fc056Field = flatFieldsList.find(field => field.code === "fc056");

					if (fc056Field) {
						const fc056Values = [];

						for (let field of flatFieldsList) {
							if (field.code === "fc056") {
								fc056Values.push(+field.resultValue);
							}
						}

						if (fc056Values.includes(2)) {
							return true;
						} else {
							return {
								errorMessage: "没有录入发票，请确认."
							};
						}
					}

					return true;
				},

				// 8
				validate8({ flatFieldsList }) {
					const fc059Field = flatFieldsList.find(field => field.code === "fc059");
					const fc060Field = flatFieldsList.find(field => field.code === "fc060");

					if (fc059Field && fc060Field) {
						const [fc059Values, fc060Values, fc153Values] = [[], [], []];

						flatFieldsList.map(field => {
							if (field.code === "fc059") {
								!fc059Values.includes(field.resultValue) &&
									fc059Values.push(field.resultValue);
							}

							if (field.code === "fc060") {
								!fc060Values.includes(field.resultValue) &&
									fc060Values.push(field.resultValue);
							}

							if (field.code === "fc153") {
								!fc153Values.includes(field.resultValue) &&
									fc153Values.push(field.resultValue);
							}
						});

						const values = [...new Set([...fc059Values, ...fc060Values])];

						if (
							values.length === 1 &&
							fc153Values.length === 1 &&
							fc153Values[0] === "0"
						) {
							const [fc056Values, fc201Values] = [[], []];

							flatFieldsList.map(field => {
								if (field.code === "fc056") {
									!fc056Values.includes(field.resultValue) &&
										fc056Values.push(field.resultValue);
								}

								if (field.code === "fc201") {
									!fc201Values.includes(field.resultValue) &&
										fc201Values.push(field.resultValue);
								}
							});

							if (!fc056Values.includes("8")) {
								if (fc201Values.length === 1 && fc201Values[0] === values[0]) {
									return {
										errorMessage: `发票${values[0]}属性没有录入发票大项，请修改`
									};
								}
							} else {
								if (fc201Values.length === 1 && fc201Values[0] !== values[0]) {
									return {
										errorMessage: `发票${values[0]}属性没有对应发票大项内容，请修改`
									};
								}
							}
						}

						return true;
					}

					return true;
				},

				// 9
				validate9({ flatFieldsList }) {
					const fc056Field = flatFieldsList.find(field => field.code === "fc056");

					if (fc056Field) {
						const specialValues = ["2", "3", "8"];
						const fc056SpecialValues = [];
						const [fc059Values, fc060Values, fc201Values, fc153Values] = [
							[],
							[],
							[],
							[]
						];

						flatFieldsList.map(field => {
							if (field.code === "fc056") {
								if (
									!fc056SpecialValues.includes(field.resultValue) &&
									specialValues.includes(field.resultValue)
								) {
									fc056SpecialValues.push(field.resultValue);
								}
							} else if (field.code === "fc059") {
								!fc059Values.includes(field.resultValue) &&
									fc059Values.push(field.resultValue);
							} else if (field.code === "fc060") {
								!fc060Values.includes(field.resultValue) &&
									fc060Values.push(field.resultValue);
							} else if (field.code === "fc201") {
								!fc201Values.includes(field.resultValue) &&
									fc201Values.push(field.resultValue);
							} else if (field.code === "fc153") {
								!fc153Values.includes(field.resultValue) &&
									fc153Values.push(field.resultValue);
							}
						});

						for (let value of fc059Values) {
							if (fc060Values.includes(value) && fc201Values.includes(value)) {
								if (fc153Values.indexOf("0") === -1) {
									return {
										errorMessage: `发票${value}属性同时录入清单大项与发票大项，请修改`
									};
								}
							}
						}

						return true;
					}

					return true;
				},

				// 10
				validate10({ flatFieldsList }) {
					const fc201Field = flatFieldsList.find(field => field.code === "fc201");

					if (fc201Field) {
						const [fc059Values, fc201Values] = [[], []];

						flatFieldsList.map(field => {
							if (field.code === "fc059") {
								if (!fc059Values.includes(field.resultValue)) {
									fc059Values.push(field.resultValue);
								}
							}

							if (field.code === "fc201") {
								if (!fc201Values.includes(field.resultValue)) {
									fc201Values.push(field.resultValue);
								}
							}
						});

						for (let value of fc059Values) {
							if (!fc201Values.includes(value)) {
								return {
									errorMessage: `发票大项${value}没有匹配的发票，请检查！`
								};
							}
						}

						return true;
					}

					return true;
				},

				// 11
				validate11({ flatFieldsList }) {
					const fc056Field = flatFieldsList.find(field => field.code === "fc056");

					if (fc056Field) {
						const fc056Values = [];

						for (let field of flatFieldsList) {
							if (field.code === "fc056") {
								fc056Values.push(+field.resultValue);
							}
						}

						if (fc056Values.includes(5)) {
							return true;
						} else {
							return {
								errorMessage: "缺少诊断书，请确认."
							};
						}
					}

					return true;
				}
			}
		}
	},

	op1op2opq: {
		// 校验规则
		rules: [
			{
				fields: [
					"fc191",
					"fc192",
					"fc193",
					"fc194",
					"fc195",
					"fc196",
					"fc197",
					"fc198",
					"fc199",
					"fc200"
				],
				validateDot: function ({ value }) {
					if (!value) return true;

					if (!/^[A-Z]/.test(value)) {
						return "疾病诊断第一个必须为大写字母，请检查!";
					}

					if (/\./.test(value)) {
						const index = value.indexOf(".");
						const restValue = value.slice(index + 1);

						if (restValue.length > 1) {
							return "否则提示录入内容只能录入到“.”符号后的一个字符，请检查!";
						}
					}

					return true;
				}
			},

			{
				fields: ["fc005", "fc006", "fc007"],
				validateDate: function ({ value }) {
					if (!value) return true;

					if (/[A, \?]/.test(value)) {
						return true;
					}

					if (value.length !== 6) {
						return "日期格式错误! ";
					}

					return true;
				}
			},

			// 16
			{
				fields: [
					"fc053",
					"fc182",
					"fc183",
					"fc184",
					"fc185",
					"fc186",
					"fc187",
					"fc188",
					"fc189",
					"fc190"
				],
				validate16: function ({ value, items }) {
					// 录入内容为A、为空或包含?时不做校验
					if (value === "A" || !value || value.includes("?")) {
						return true;
					}

					const result = items.find(text => text === value);

					if (result) {
						return true;
					} else {
						return "疾病诊断录入错误，请根据下拉提示内容选录.";
					}
				}
			},

			// 23
			{
				fields: ["fc054"],
				validate23: function ({ value, items }) {
					if (!value) return true;

					const result = items.find(text => text === value);

					if (result) {
						return true;
					} else {
						return "录入内容不在代码表中，请录入“本代码表中不存在的其他医院”.";
					}
				}
			},

			// 26
			{
				fields: ["fc055"],
				validate26: function ({ value, field, items }) {
					// 录入内容为A、为空或包含?时不做校验
					if (value === "A" || !value || value.includes("?")) {
						return true;
					}

					const result = items.find(text => text === value);

					if (result) {
						return true;
					} else {
						return "手术录入错误，请根据下拉提示内容选录.";
					}
				}
			},

			// 27
			{
				fields: [
					"fc009",
					"fc011",
					"fc013",
					"fc015",
					"fc017",
					"fc019",
					"fc021",
					"fc023",
					"fc025",
					"fc027",
					"fc029",
					"fc031",
					"fc033",
					"fc035",
					"fc037",
					"fc039",
					"fc041",
					"fc043",
					"fc045",
					"fc047"
				],
				validate27: function ({ value, items }) {
					if (!value) return true;

					const result = items.find(text => text === value);

					if (result) {
						return true;
					} else {
						return "录入内容与代码表不一致，请选择相近的内容录入，如无相近则录入其他费.";
					}
				}
			},

			// 30
			{
				fields: ["fc152"],
				validate30: function ({ value, field, items }) {
					return true;
				}
			},

			// 37
			{
				fields: ["fc172", "fc173", "fc174", "fc175", "fc176", "fc177", "fc178", "fc179"],
				validate37: function ({ value, field, fieldsIndex, fieldsList }) {
					if (!value) return true;

					const mapCodesList = new Map([
						["fc172", ["fc154", "fc092"]],
						["fc173", ["fc155", "fc093"]],
						["fc174", ["fc156", "fc094"]],
						["fc175", ["fc157", "fc095"]],
						["fc176", ["fc158", "fc096"]],
						["fc177", ["fc159", "fc097"]],
						["fc178", ["fc160", "fc098"]],
						["fc179", ["fc161", "fc099"]]
					]);

					const codes = mapCodesList.get(field.code);

					if (codes) {
						const fields = fieldsList[fieldsIndex];
						const col1Field = fields.find(_field => _field.code === codes[0]);
						const col2Field = fields.find(_field => _field.code === codes[1]);

						if (+col1Field.resultValue === 4) {
							const col2FieldValue = col2Field.resultValue || 0;

							if (+value > col2FieldValue) {
								return "自付金额不能大于总金额.";
							}
						}
					}

					return true;
				}
			},

			// 38
			{
				fields: ["fc084", "fc085", "fc086", "fc087", "fc088", "fc089", "fc090", "fc091"],
				validate38: function ({ value, items }) {
					if (!value) {
						return true;
					}

					const result = items.find(text => text === value);

					if (result) {
						return true;
					} else {
						return "录入内容不在常量表中，请选择相近内容进行录入，如完全没有相近的则按单录入强过.";
					}
				}
			}
		],

		// 提示文本
		hints: [
			{
				fields: ["fc084", "fc067"],
				hintFc067: function ({ bill }) {
					const { billNum } = bill;

					if (/^64/.test(billNum)) {
						return "江苏徐州市的结算单模板，社保自费请按单录入“先行支付”+“自费”!";
					}

					return true;
				}
			},

			{
				fields: ["fc152"],
				hintFc152: function () {
					return '<p style="color: blue;">冻干人用狂犬病疫苗、破伤风疫苗、破伤风类毒素需录入西药费；工本费、病历费、卡费、复印费、陪护费、陪人费此类大项名称按单录入强过；照相费选择放射费、处置费选择治疗费、材料费选择卫材费!</p>';
				}
			}
		],

		// 字段已生成
		init: {
			methods: {
				disableFields: function ({ op, fieldsList, focusFieldsIndex }) {
					if (op === "op0") {
						return;
					}

					const codesList = [
						["fc057", "fc058"],
						["fc144", "fc145", "fc146"],
						["fc202", "fc203", "fc204", "fc211", "fc212", "fc209", "fc210"],
						["fc103", "fc150", "fc151", "fc170", "fc171"],
						["fc225", "fc226", "fc227"],
						["fc228", "fc229", "fc230"],
						["fc231", "fc232", "fc233"],
						[
							"fc104",
							"fc105",
							"fc106",
							"fc107",
							"fc108",
							"fc109",
							"fc110",
							"fc111",
							"fc112",
							"fc113",
							"fc114",
							"fc115",
							"fc116",
							"fc117",
							"fc118",
							"fc119",
							"fc120",
							"fc121",
							"fc122",
							"fc123",
							"fc124",
							"fc125",
							"fc126",
							"fc127",
							"fc128",
							"fc129",
							"fc130",
							"fc131",
							"fc132",
							"fc133",
							"fc134",
							"fc135",
							"fc136",
							"fc137",
							"fc138",
							"fc139",
							"fc140",
							"fc141",
							"fc142",
							"fc143"
						]
					];

					const flatCodesList = [];

					codesList.map(codes => {
						flatCodesList.push(...codes);
					});

					const fields = fieldsList[focusFieldsIndex];

					fields?.map(_field => {
						if (flatCodesList.includes(_field.code)) {
							_field.disabled = true;
						}
					});

					console.log(fieldsList);
				},

				// 16
				validate16: function ({ op, flatFieldsList }) {
					const fields = [
						"fc053",
						"fc182",
						"fc183",
						"fc184",
						"fc185",
						"fc186",
						"fc187",
						"fc188",
						"fc189",
						"fc190"
					];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_疾病诊断",
								query: "名称"
							};
						}
					});
				},

				// 26
				validate26: function ({ op, flatFieldsList }) {
					const fields = ["fc055"];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_手术编码",
								query: "名称"
							};
						}
					});
				},

				// 27
				validate27: function ({ op, flatFieldsList }) {
					const fields = [
						"fc009",
						"fc011",
						"fc013",
						"fc015",
						"fc017",
						"fc019",
						"fc021",
						"fc023",
						"fc025",
						"fc027",
						"fc029",
						"fc031",
						"fc033",
						"fc035",
						"fc037",
						"fc039",
						"fc041",
						"fc043",
						"fc045",
						"fc047"
					];

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							_field.table = {
								name: "B0118_中意理赔_发票大项类型",
								query: "费用名称"
							};
						}
					});
				},

				// 38
				validate38: function ({ op, flatFieldsList, codeValues }) {
					const fields = [
						"fc084",
						"fc085",
						"fc086",
						"fc087",
						"fc088",
						"fc089",
						"fc090",
						"fc091"
					];
					const constantsDB = window.constantsDB || {};
					const { fc180, fc181 } = codeValues || {};

					flatFieldsList.map(_field => {
						if (fields.includes(_field.code)) {
							const prefix = "B0118_中意理赔_省份-";

							if (constantsDB[`${prefix}${fc180}`]) {
								_field.table = {
									name: `${prefix}${fc180}`,
									query: "中文名称"
								};
							} else if (constantsDB[`${prefix}${fc181}`]) {
								_field.table = {
									name: `${prefix}${fc181}`,
									query: "中文名称"
								};
							} else {
								_field.table = {
									name: `${prefix}全国`,
									query: "项目名称"
								};
							}
						}
					});
				}
			}
		},

		// 回车
		enter: {
			methods: {
				// 17
				validate17({ field, fieldsList, focusFieldsIndex }) {
					const fields = fieldsList[focusFieldsIndex];

					if (field.code === "fc053") {
						const ifADisableCodes = [
							"fc053",
							"fc182",
							"fc183",
							"fc184",
							"fc185",
							"fc186",
							"fc187",
							"fc188",
							"fc189",
							"fc190"
						];
						const ifNotADisableCodes = [
							"fc191",
							"fc192",
							"fc193",
							"fc194",
							"fc195",
							"fc196",
							"fc197",
							"fc198",
							"fc199",
							"fc200"
						];

						if (field.resultValue === "A") {
							fields.map(_field => {
								if (ifADisableCodes.includes(_field.code)) {
									_field.disabled = true;
								}

								if (ifNotADisableCodes.includes(_field.code)) {
									_field.disabled = false;
								}
							});
						} else {
							fields.map(_field => {
								if (ifADisableCodes.includes(_field.code)) {
									_field.disabled = false;
								}

								if (ifNotADisableCodes.includes(_field.code)) {
									_field.disabled = true;
								}
							});
						}
					}
				},

				// 20
				validate20({ op, fieldsList, focusFieldsIndex }) {
					const fields = fieldsList[focusFieldsIndex];

					const fc003Field = fields.find(field => field.code === "fc003");
					const fc005Field = fields.find(field => field.code === "fc005");
					const fc006Field = fields.find(field => field.code === "fc006");
					const fc007Field = fields.find(field => field.code === "fc007");

					if (fc003Field?.resultValue == 1) {
						fc005Field.disabled = false;
						fc006Field.disabled = true;
						fc007Field.disabled = true;

						const fc005Value = fc005Field.resultValue;

						if (fc005Value) {
							fc006Field[`${op}Value`] = fc005Value;
							fc006Field.resultValue = fc005Value;

							fc007Field[`${op}Value`] = fc005Value;
							fc007Field.resultValue = fc005Value;
						} else {
							fc006Field[`${op}Value`] = "";
							fc006Field.resultValue = "";

							fc007Field[`${op}Value`] = "";
							fc007Field.resultValue = "";
						}
					} else if (fc003Field?.resultValue == 2) {
						fc005Field.disabled = true;
						fc006Field.disabled = false;
						fc007Field.disabled = false;

						fc005Field[`${op}Value`] = "";
						fc005Field.resultValue = "";
					}
				},

				// 22
				validate22({ field, fieldsList, focusFieldsIndex }) {
					if (field.code === "fc004") {
						const fields = fieldsList[focusFieldsIndex];
						const codes = [
							"fc063",
							"fc064",
							"fc065",
							"fc066",
							"fc067",
							"fc068",
							"fc069",
							"fc070",
							"fc071",
							"fc100",
							"fc072",
							"fc074",
							"fc075",
							"fc076",
							"fc077",
							"fc050"
						];

						fields.map(_field => {
							if (codes.includes(_field.code)) {
								_field.disabled = false;
							}
						});

						if (field.resultValue === "3" || field.resultValue === "4") {
							fields.map(_field => {
								if (codes.includes(_field.code)) {
									_field.disabled = true;
								}
							});
						}
					}
				},

				// 24 47
				validate24({ field, fieldsList, focusFieldsIndex }) {
					if (field.code === "fc101") {
						const fields = fieldsList[focusFieldsIndex];
						const fc079Field = fields.find(field => field.code === "fc079");
						const fc051Field = fields.find(field => field.code === "fc051");
						const fc102Field = fields.find(field => field.code === "fc102");
						const fc052Field = fields.find(field => field.code === "fc052");
						const fc067Field = fields.find(field => field.code === "fc067");

						fc079Field.disabled = false;
						fc051Field.disabled = false;
						fc102Field.disabled = false;
						fc052Field.disabled = false;
						fc067Field.disabled = false;

						if (field.resultValue === "1" || field.resultValue === "2") {
							fc051Field.disabled = true;
							fc102Field.disabled = true;
							fc079Field.disabled = true;
						} else {
							fc052Field.disabled = true;
							fc067Field.disabled = true;
						}
					}
				},

				// 36
				validate36({ field, fieldsList, focusFieldsIndex }) {
					const mapCodesList = new Map([
						["fc154", ["fc162", "fc172"]],
						["fc155", ["fc163", "fc173"]],
						["fc156", ["fc164", "fc174"]],
						["fc157", ["fc165", "fc175"]],
						["fc158", ["fc166", "fc176"]],
						["fc159", ["fc167", "fc177"]],
						["fc160", ["fc168", "fc178"]],
						["fc161", ["fc169", "fc179"]]
					]);

					const codes = mapCodesList.get(field.code);

					if (codes) {
						const fields = fieldsList[focusFieldsIndex];
						const rightFirField = fields.find(field => field.code === codes[1]);
						const rightSecField = fields.find(field => field.code === codes[0]);

						rightFirField.disabled = false;
						rightSecField.disabled = false;

						if (field.resultValue === "1") {
							rightFirField.disabled = true;
						} else if (field.resultValue === "4") {
							rightSecField.disabled = true;
						} else {
							rightFirField.disabled = true;
							rightSecField.disabled = true;
						}
					}
				},

				// 45
				validate45({ field, fieldsList, focusFieldsIndex }) {
					if (field.code === "fc205") {
						const fields = fieldsList[focusFieldsIndex];
						const fc206Field = fields.find(field => field.code === "fc206");
						const fc207Field = fields.find(field => field.code === "fc207");
						const fc208Field = fields.find(field => field.code === "fc208");

						fc206Field.disabled = false;
						fc207Field.disabled = false;
						fc208Field.disabled = false;

						if (field.resultValue === "A") {
							fc206Field.disabled = true;
							fc207Field.disabled = true;
						} else if (field.resultValue == 1 || field.resultValue == 2) {
							fc207Field.disabled = true;
							fc208Field.disabled = true;
						} else {
							fc206Field.disabled = true;
						}
					}
				}
			}
		},

		// 临时保存
		sessionSave: {
			methods: {
				// 18
				disable18({ fieldsList, focusFieldsIndex }) {
					const codesList = [
						[
							"fc053",
							"fc182",
							"fc183",
							"fc184",
							"fc185",
							"fc186",
							"fc187",
							"fc188",
							"fc189",
							"fc190"
						],
						[
							"fc191",
							"fc192",
							"fc193",
							"fc194",
							"fc195",
							"fc196",
							"fc197",
							"fc198",
							"fc199",
							"fc200"
						]
					];

					const fields = fieldsList[focusFieldsIndex];
					const focusField = fields.find(field => field.autofocus);

					let codes = [];

					for (let index in codesList) {
						if (codesList[index].includes(focusField.code)) {
							codes = codesList[index];
							break;
						}
					}

					const codeIndex = codes.indexOf(focusField.code);

					if (codeIndex > -1) {
						const restCodes = codes.slice(codeIndex + 1);
						const restFields = fields.slice(focusField.fieldIndex + 1);

						restFields?.map(restField => {
							if (restCodes.includes(restField.code)) {
								restField.disabled = true;
							}
						});
					}
				},

				// 31
				disable31({ fieldsList, focusFieldsIndex }) {
					const codesList = [
						["fc009", "fc010"],
						["fc011", "fc012"],
						["fc013", "fc014"],
						["fc015", "fc016"],
						["fc017", "fc018"],
						["fc019", "fc020"],
						["fc021", "fc022"],
						["fc023", "fc024"],
						["fc025", "fc026"],
						["fc027", "fc028"],
						["fc029", "fc030"],
						["fc031", "fc032"],
						["fc033", "fc034"],
						["fc035", "fc036"],
						["fc037", "fc038"],
						["fc039", "fc040"],
						["fc041", "fc042"],
						["fc043", "fc044"],
						["fc045", "fc046"],
						["fc047", "fc048"]
					];

					const col2Codes = [];

					codesList.map(codes => {
						col2Codes.push(codes[1]);
					});

					const fields = fieldsList[focusFieldsIndex];

					const focusField = fields.find(field => field.autofocus);
					const codeIndex = col2Codes.indexOf(focusField.code);

					if (codeIndex > -1) {
						const restCodes = [];
						let sliceIndex = -1;

						for (let codesIndex in codesList) {
							if (codesList[codesIndex].includes(focusField.code)) {
								sliceIndex = +codesIndex + 1;
								break;
							}
						}

						const restCodesList = codesList.slice(sliceIndex);

						restCodesList.map(codes => {
							restCodes.push(...codes);
						});

						const restFields = fields.slice(focusField.fieldIndex + 1);

						restFields?.map(restField => {
							if (restCodes.includes(restField.code)) {
								restField.disabled = true;
							}
						});
					}
				},

				// 41
				disable41({ fieldsList, focusFieldsIndex }) {
					const codesList = [
						["fc084", "fc154", "fc092", "fc162", "fc172"],
						["fc085", "fc155", "fc093", "fc163", "fc173"],
						["fc086", "fc156", "fc094", "fc164", "fc174"],
						["fc087", "fc157", "fc095", "fc165", "fc175"],
						["fc088", "fc158", "fc096", "fc166", "fc176"],
						["fc089", "fc159", "fc097", "fc167", "fc177"],
						["fc090", "fc160", "fc098", "fc168", "fc178"],
						["fc091", "fc161", "fc099", "fc169", "fc179"]
					];
					const fields = fieldsList[focusFieldsIndex];
					const focusField = fields.find(field => field.autofocus);
					let sliceIndex = -1;

					for (let codesIndex in codesList) {
						if (codesList[codesIndex].includes(focusField.code)) {
							sliceIndex = +codesIndex + 1;
							break;
						}
					}

					if (sliceIndex > -1) {
						const restCodesList = codesList.slice(sliceIndex);
						const restCodes = [];

						restCodesList.map(codes => {
							restCodes.push(...codes);
						});

						const restFields = fields.slice(focusField.fieldIndex + 1);

						restFields?.map(restField => {
							if (restCodes.includes(restField.code)) {
								restField.disabled = true;
							}
						});
					}
				}
			}
		},

		// 提交前
		beforeSubmit: {
			methods: {
				// 32
				validate32({ op, fieldsList }) {
					const fields = fieldsList[0];
					const codes = [
						"fc010",
						"fc012",
						"fc014",
						"fc016",
						"fc018",
						"fc020",
						"fc022",
						"fc024",
						"fc026",
						"fc028",
						"fc030",
						"fc032",
						"fc034",
						"fc036",
						"fc038",
						"fc040",
						"fc042",
						"fc044",
						"fc046",
						"fc048"
					];
					const fc008Field = fields.find(field => field.code === "fc008");
					let total = 0;

					if (!fc008Field) {
						return true;
					}

					for (let field of fields) {
						let value = field[`${op}Value`];

						if (/\?/.test(value)) {
							return true;
						}

						if (codes.includes(field.code)) {
							if (!value) {
								value = 0;
							}

							total += Number(value);
						}
					}

					const diff = fc008Field[`${op}Value`] - total;

					if (diff === 0) {
						return true;
					} else {
						return {
							errorMessage: `明细金额与总金额不一致，差额为${diff}，请确认并修改!`
						};
					}
				},

				// 39
				validate39({ block, fieldsList }) {
					if (block.code !== "bc002") {
						return true;
					}

					for (let fields of fieldsList) {
						for (let field of fields) {
							if (field.resultValue) {
								return true;
							}
						}
					}

					return {
						errorMessage:
							"清单不能空白提交，请检查，如清单内容无法录入则录入一组数据后按F8提交."
					};
				},

				// 40
				validate40({ fieldsList }) {
					const error = {
						errorMessage: "清单内容录入遗漏，请检查!"
					};

					for (let fields of fieldsList) {
						const codesList = [
							{ fc084: "", fc154: "", fc092: "", fc162: "", fc172: "" },
							{ fc085: "", fc155: "", fc093: "", fc163: "", fc173: "" },
							{ fc086: "", fc156: "", fc094: "", fc164: "", fc174: "" },
							{ fc087: "", fc157: "", fc095: "", fc165: "", fc175: "" },
							{ fc087: "", fc158: "", fc096: "", fc166: "", fc176: "" },
							{ fc088: "", fc159: "", fc097: "", fc167: "", fc177: "" },
							{ fc090: "", fc160: "", fc098: "", fc168: "", fc178: "" },
							{ fc091: "", fc161: "", fc099: "", fc169: "", fc179: "" }
						];

						for (let field of fields) {
							for (let codeItem of codesList) {
								if (codeItem.hasOwnProperty(field.code)) {
									codeItem[field.code] = field.resultValue;
								}
							}
						}

						for (let codeItem of codesList) {
							const keys = Object.keys(codeItem);
							const firKey = keys[0];
							const secKey = keys[1];
							const thiKey = keys[2];
							const fouKey = keys[3];
							const fifKey = keys[4];

							if (codeItem[secKey] == 1) {
								if (
									!codeItem[firKey] ||
									!codeItem[thiKey] ||
									!codeItem[fouKey] ||
									codeItem[fifKey]
								) {
									return error;
								}
							}

							if (codeItem[secKey] == 2 || codeItem[secKey] == 3) {
								if (
									!codeItem[firKey] ||
									!codeItem[thiKey] ||
									codeItem[fouKey] ||
									codeItem[fifKey]
								) {
									return error;
								}
							}

							if (codeItem[secKey] == 4) {
								if (
									!codeItem[firKey] ||
									!codeItem[thiKey] ||
									!codeItem[fifKey] ||
									codeItem[fouKey]
								) {
									return error;
								}
							}
						}
					}

					return true;
				}
			}
		}
	}
};

export default B0118;
