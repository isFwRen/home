<template>
	<v-card id="scroll-target" style="max-height: 100vh" class="overflow-y-auto">
		<LPLoading ref="lploading" />
		<v-card
			style="padding-bottom: 50px; box-shadow: 0 0 0 0 !important; position: relative"
			class="scrollContainer"
			v-scroll:#scroll-target="onScroll"
		>
			<v-toolbar class="z-toolbar" color="primary" dark>
				<z-btn dark icon @click="onBack">
					<v-icon>mdi-arrow-left</v-icon>
				</z-btn>

				<lp-tooltip-btn
					bottom
					btnIcon
					fab
					icon="mdi-swap-horizontal"
					small
					tip="切换至质检界面"
					@click="onSwitch"
				>
				</lp-tooltip-btn>

				<span class="pl-2">案件号 {{ cases.caseInfo.billNum }}</span>

				<span class="pl-2">机构号 {{ cases.caseInfo.billInfo.agency }} </span>

				<v-spacer></v-spacer>
				<lp-tooltip-btn
					bottom
					btnIcon
					fab
					icon="mdi-receipt-text-check-outline"
					small
					tip="发票查验"
					@click="getdata()"
				>
				</lp-tooltip-btn>
				<lp-tooltip-btn
					v-for="item of cells.icons"
					bottom
					btnIcon
					fab
					:icon="item.icon"
					:key="item.value"
					small
					:tip="item.tip"
					@click="onNavigate(item)"
				></lp-tooltip-btn>

				<z-btn dark icon @click="onClose">
					<v-icon>mdi-close</v-icon>
				</z-btn>
			</v-toolbar>

			<v-card-text class="py-8">
				<div id="contain" class="contain">
					<!-- 查询 BEGIN -->
					<div class="z-flex align-end mt-4 mb-8 rounded-md lp-filters searchFlex">
						<v-row class="z-flex align-end">
							<v-col
								v-for="(item, index) in cells.fields"
								:key="`entry_filters_${index}`"
								:cols="2"
							>
								<v-text-field
									:formId="searchFormId"
									:formKey="item.formKey"
									:hideDetails="item.hideDetails"
									:label="item.label"
									v-model="form[item.formKey]"
								>
								</v-text-field>
							</v-col>

							<z-btn class="pl-3 pb-3" color="primary" @click="onSearch">
								<v-icon class="text-h6">mdi-magnify</v-icon>
								查询
							</z-btn>
						</v-row>
					</div>
					<!-- 查询 END -->

					<v-expansion-panels v-model="panel" multiple>
						<v-expansion-panel
							ref="expantionRef"
							v-for="(item, index) in desserts"
							:key="`invoice_${index}`"
							@change="handlePanel(index)"
						>
							<template v-if="item.leftdesserts.length || item.rightdesserts.length">
								<v-expansion-panel-header color="#cfe5fb">
									{{ item.showIndex + "." + item.title }}
								</v-expansion-panel-header>

								<v-expansion-panel-content color="#cfe5fb">
									<v-row justify="center">
										<v-col cols="6">
											<div
												v-for="(fields, fieldsIndex) of item.leftdesserts"
												:key="`fields_${fieldsIndex}`"
											>
												<template v-for="(field, fieldIndex) of fields">
													<div
														v-if="
															field.visible != false &&
															field.resultInput !== 'no'
														"
														:key="`field_${fieldIndex}`"
														class="z-flex justify-center align-center mb-4 action"
													>
														<div class="mr-2 text-right fw-bold label">
															<span>
																{{
																	`${field.code}${field.name}`
																}}</span
															>
														</div>

														<div
															class="mt-n2 field"
															:class="`${field.ID}${field.code}`"
														>
															<v-text-field
																:ref="`field_${field.ID}`"
																:formId="formId"
																:formKey="`field_${field.ID}`"
																dense
																hide-details
																persistent-hint
																width="250"
																:value="field.resultValue"
																@keydown.enter="onFieldEnter"
																@change="
																	onFieldChange({
																		value: $event,
																		field
																	})
																"
															>
															</v-text-field>
															<z-list
																:dataSource="[
																	{ label: '23', value: '23' }
																]"
															></z-list>
														</div>

														<div
															class="px-2 py-1 ml-2 mt-n1 rounded-sm hint"
															@click="onHint(field)"
														>
															<span
																style="
																	max-width: 180px !important;
																	overflow: hidden;
																	display: inline-block;
																	white-space: nowrap;
																	text-overflow: ellipsis;
																	padding: 4px 4px 0 4px;
																"
																@mouseover="
																	openTooltip(
																		$event,
																		field.finalValue,
																		'left'
																	)
																"
																@mouseleave="closeTooltip"
															>
																{{ field.finalValue }}
															</span>

															<template
																v-if="
																	field.issues &&
																	field.issues.length
																"
															>
																<span
																	class="error--text"
																	v-for="(
																		issue, index
																	) of field.issues"
																	:key="index"
																>
																	{{
																		`[${issue.code}: ${issue.message}]`
																	}}<i
																		v-show="
																			index <
																			field.issues.length - 1
																		"
																		>、</i
																	>
																</span>
															</template>
														</div>

														<lp-tooltip-btn
															btnColor="primary"
															depressed
															fab
															icon="mdi-image-area"
															right
															small
															tip="查看图片"
															@click="onViewImage(field)"
														></lp-tooltip-btn>
													</div>
													<p
														:key="`field_${fieldIndex} + 1`"
														v-if="
															field.code == 'fc067' &&
															field.name == '社保自费' &&
															saleChannel == '1'
														"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														特约件， 乙类自付无需录入
													</p>
													<p
														:key="fieldIndex"
														v-if="
															field.code == 'fc035' &&
															proCode == 'B0108' &&
															agency.slice(0, 3) == '001'
														"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														无需录入总金额、统筹、自费、发票大项和清单
													</p>
													<p
														:key="fieldIndex"
														v-if="
															field.code == 'fc057' &&
															proCode == 'B0108' &&
															agency.slice(0, 3) == '001'
														"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														无需录入总金额、统筹、自费、发票大项和清单
													</p>
													<!-- <p
														:key="`field_${fieldIndex} + 2`"
														v-if="field.code == 'fc275' && flag"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														江苏徐州件
													</p> -->
												</template>
											</div>
										</v-col>

										<!--分组-->
										<v-col cols="6">
											<div
												v-for="(fields, fieldsIndex) of item.rightdesserts"
												:key="`fields_${fieldsIndex}`"
											>
												<template v-for="(field, fieldIndex) of fields">
													<div
														v-if="
															field.visible != false &&
															field.resultInput !== 'no'
														"
														:key="`field_${fieldIndex}`"
														class="z-flex justify-center align-center mb-4 action"
													>
														<div class="mr-2 text-right fw-bold label">
															<span>
																{{
																	`${field.code}${field.name}`
																}}</span
															>
														</div>

														<div
															class="mt-n2 field"
															:class="`${field.ID}${field.code}`"
														>
															<v-text-field
																:ref="`field_${field.ID}`"
																:formId="formId"
																:formKey="`field_${field.ID}`"
																dense
																hide-details
																persistent-hint
																width="250"
																:value="field.resultValue"
																@keydown.enter="onFieldEnter"
																@change="
																	onFieldChange({
																		value: $event,
																		field
																	})
																"
															>
															</v-text-field>
															<z-list
																:dataSource="[
																	{ label: '23', value: '23' }
																]"
															></z-list>
														</div>

														<div
															class="px-2 py-1 ml-2 mt-n1 rounded-sm hint"
															@click="onHint(field)"
														>
															<span
																style="
																	max-width: 180px !important;
																	overflow: hidden;
																	display: inline-block;
																	white-space: nowrap;
																	text-overflow: ellipsis;
																	padding: 4px 4px 0 4px;
																"
																@mouseover="
																	openTooltip(
																		$event,
																		field.finalValue,
																		'right'
																	)
																"
																@mouseleave="closeTooltip"
															>
																{{ field.finalValue }}
															</span>
															<template
																v-if="
																	field.issues &&
																	field.issues.length
																"
															>
																<span
																	class="error--text"
																	v-for="(
																		issue, index
																	) of field.issues"
																	:key="index"
																>
																	{{
																		`[${issue.code}: ${issue.message}]`
																	}}<i
																		v-show="
																			index <
																			field.issues.length - 1
																		"
																		>、</i
																	>
																</span>
															</template>
														</div>

														<lp-tooltip-btn
															btnColor="primary"
															depressed
															fab
															icon="mdi-image-area"
															right
															small
															tip="查看图片"
															@click="onViewImage(field)"
														></lp-tooltip-btn>
													</div>
													<p
														:key="`field_${fieldIndex} + 1`"
														v-if="
															field.code == 'fc067' &&
															field.name == '社保自费' &&
															saleChannel == '1'
														"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														特约件， 乙类自付无需录入
													</p>
													<p
														:key="fieldIndex"
														v-if="
															field.code == 'fc035' &&
															proCode == 'B0108' &&
															agency.slice(0, 3) == '001'
														"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														无需录入总金额、统筹、自费、发票大项和清单
													</p>
													<p
														:key="fieldIndex"
														v-if="
															field.code == 'fc057' &&
															proCode == 'B0108' &&
															agency.slice(0, 3) == '001'
														"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														无需录入总金额、统筹、自费、发票大项和清单
													</p>
													<!-- <p
														:key="`field_${fieldIndex} + 2`"
														v-if="field.code == 'fc275' && flag"
														style="
															color: red;
															text-align: center;
															margin-bottom: 0px !important;
														"
													>
														江苏徐州件
													</p> -->
												</template>
											</div>
										</v-col>
									</v-row>
								</v-expansion-panel-content>
							</template>
						</v-expansion-panel>
					</v-expansion-panels>
				</div>
			</v-card-text>

			<div class="z-card_fixed-actions">
				<v-card-actions class="z-flex justify-center px-4">
					<v-row>
						<v-col :cols="6">
							<z-pagination
								:pageNum="pageNum"
								:pageSize="1"
								:showQuickJumper="false"
								:showSizeChanger="false"
								:total="pagination.total"
								@page="handlePage"
							></z-pagination>
						</v-col>

						<v-col :cols="4">
							<z-switch
								v-if="desserts?.length"
								formId="all"
								formKey="all"
								class="mt-n2 ml-4"
								:label="selectedAll ? '全不选' : '全选'"
								hide-details
								:defaultValue="selectedAll"
								@change="onSelectAll"
							>
							</z-switch>
						</v-col>

						<v-col class="z-flex justify-end mt-1" :cols="2">
							<z-btn class="mr-3" color="normal" @click="onBack"> 返回 </z-btn>
							<v-btn color="primary" v-if="desserts?.length" @click="onSaveAsXML">
								保存为xml
							</v-btn>
						</v-col>
					</v-row>
				</v-card-actions>
			</div>

			<field-info-dialog ref="info" :info="hintInfo"></field-info-dialog>

			<change-log-dialog ref="log"></change-log-dialog>

			<export-validation ref="export" :wrongNote="wrongNote"></export-validation>
		</v-card>

		<searchResultData ref="searchDialog" :resultData="searchResultSet" @toLink="handleLink" />

		<content-com :data="contents" @toContents="handleContents" :activeIndex="contentsIndex" />
	</v-card>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import { tools, sessionStorage } from "vue-rocket";
import { tools as lpTools } from "@/libs/util";
import DialogMixins from "@/mixins/DialogMixins";
import TableMixins from "@/mixins/TableMixins";
import CaseMixins from "../CaseMixins";
import _ from "lodash";

import cells, { invoiceKeys, invoiceNamesMap } from "./cells";
import LPLoading from "../../../../../components/lp-loading";

const { baseURLApi } = lpTools.baseURL();

export default {
	name: "UpdateEntryData",
	mixins: [DialogMixins, TableMixins, CaseMixins],

	data() {
		return {
			isSearch: false,
			isContents: false,
			formId: "UpdateEntryData",
			cells,
			selectedAll: false,
			panel: [],
			loading: false,
			memoPanels: {},
			invoiceList: [],
			billInfo: {},
			desserts: [],
			hintInfo: {},
			fields: [],
			modifiedFields: [],
			isEqualSearch: false,
			pageNum: 1,
			sabayon: {},
			panelIndex: 0,
			wrongNote: "",
			isEnter: false,
			viewHeight: 0,
			distance: 0,
			enabled: true,
			inputDom: null,
			formKey: "",
			keyword: "",
			isKeywordSearch: false,
			form: {
				id: "",
				code: "",
				money: "",
				keyword: ""
			},
			lastScrollTop: 0,
			searchResultSet: [],
			noScrollExpandCard: [],
			saleChannel: "",
			agency: "",
			proCode: "",
			contents: [],
			contentsIndex: 0,
			objInvoice: [],
			flag: ""
		};
	},
	created() {
		this.getCaseEntryData();
		this.proCode = sessionStorage.get("CaseSearch").proCode;
	},
	mounted() {
		window.addEventListener("keydown", this.fuckEvents);
		// 获取视口高度
		this.viewHeight = document.documentElement.clientHeight;
		this.$nextTick(() => {
			this.scrollDom = document.querySelector(".v-dialog");
			this.scrollDom.style.overflowY = "hidden";
			this.scrollCard = document.querySelector(".overflow-y-auto");
		});
	},
	beforeDestroy() {
		this.resetData();
		window.removeEventListener("keydown", this.fuckEvents);
	},

	computed: {
		...mapState(["forms"]),
		...mapGetters(["cases"]),
		row() {
			return this.cases.rowInfo;
		}
	},
	watch: {
		panelIndex(val) {
			if (this.panel.length >= val) {
				return;
			}
			if (val < this.desserts.length) {
				this.panel = this.genePanel(val);
			} else {
				this.panel = this.genePanel(this.desserts.length - 1);
			}
		}
		// pageNum(val) {
		// 	if (this.billInfo.proCode == "B0118") {
		// 		for (let i = 0; i < this.objInvoice.length; i++) {
		// 			if (i + 1 == val) {
		// 				let arr1 = this.objInvoice[i].invoice[0].filter(
		// 					item => item.code == "fc180" && item.resultValue == "徐州"
		// 				);
		// 				let arr2 = this.objInvoice[i].baoXiaoDan[0].filter(
		// 					item => item.code == "fc275"
		// 				);
		// 				if (arr1.length != 0 && arr2.length != 0) {
		// 					this.flag = true;
		// 				} else {
		// 					this.flag = false;
		// 				}
		// 			}
		// 		}
		// 	}
		// }
		// objInvoice: {
		// 	handler(val) {
		// 		if (this.billInfo.proCode == 'B0118') {
		// 			for (let i = 0; i < val.length; i++) {
		// 				if (i == this.pageNum - 1) {
		// 					let arr1 = val[i].invoice[0].filter(item => item.code == 'fc180' && item.resultValue == '徐州')
		// 					let arr2 = val[i].baoXiaoDan[0].filter(item => item.code == 'fc275')
		// 					if (arr1.length != 0 && arr2.length != 0) {
		// 						this.flag = true
		// 					} else {
		// 						this.flag = false
		// 					}
		// 				}
		// 			}
		// 		}
		// 	},
		// 	deep: true,
		// 	immediate: true,
		// }
	},
	methods: {
		openTooltip(event, content, dir) {
			if (!content) {
				return;
			}

			let left = event.target.getBoundingClientRect().left;
			left =
				dir === "left"
					? left + event.target.getBoundingClientRect().width + 15
					: left - 450 - 15;
			let top = event.target.getBoundingClientRect().top;
			this.$tooltipshow({
				content,
				styles: {
					left,
					top,
					dir
				}
			});
		},
		closeTooltip() {
			this.$tooltiphide();
		},
		showTooltip(name) {
			this.$refs[name][0].showTooltip();
		},
		hideTooltip(name) {
			this.$refs[name][0].hideTooltip();
		},
		/**function 收集跳转页面信息
		 * @param {isSearch} is keyword search?
		 * @param {pathPanel} 手风琴展开页信息
		 * @param {noScrollExpandCard} 手动关闭手风琴索引，手动关闭后不再自动展开
		 * @param {contentsIndex} 当前目点击录索引
		 */
		handleLink(path) {
			this.noScrollExpandCard = [];
			this.pathPanel = path.panel;
			this.contentsIndex = 0;
			this.isSearch = true;
			//this.$refs.lploading.open();
			this.pageNum = path.index + 1;
			this.scrollToFunOne(path);
		},
		// 查询关键字后设置前端数据结构
		setSearchEndData(result, path) {
			const { invoice } = result.data;
			this.billInfo = invoice.bill;
			this.wrongNote = result.data.bill?.wrongNote;
			if (!this.isEnter) {
				// 恢复默认
				{
					this.invoiceList = [];

					this.pagination = {
						total: 0
					};

					this.desserts = [];

					if (!tools.isYummy(this.memoPanels))
						if (tools.isLousy(invoice)) {
							return;
						}
				}
			}
			// 同页跳转展开手风琴
			this.panel = this.genePanel(this.pathPanel);
			invoice.map((item, index) => {
				this.invoiceList[index] = [];

				for (let key in item) {
					if (invoiceKeys.includes(key)) {
						if (invoice[index][key] != null && invoice[index][key].length > 0) {
							if (invoice[index][key][0] != null) {
								invoice[index][key].map(arr => {
									const arrs = _.chunk(arr, Math.floor(arr.length / 2));
									let leftdesserts = [],
										rightdesserts = [];
									if (arrs.length > 2) {
										const temp = [...arrs[1], ...arrs[2]];
										leftdesserts = [arrs[0]];
										rightdesserts = [temp];
									} else {
										leftdesserts = [arrs[0]];
										rightdesserts = [arrs[1]];
									}

									this.invoiceList[index] = [
										...this.invoiceList[index],
										{
											key,
											title: invoiceNamesMap.get(key),
											desserts: [arr],
											leftdesserts,
											rightdesserts
										}
									];
								});
							} else {
								this.invoiceList.splice(index, 1);
							}
						}
					}
				}
			});

			// 记忆最开始的resultValue
			this.invoiceList?.map(invoice => {
				invoice?.map(item => {
					item.desserts?.map(fields => {
						fields?.map(field => {
							field.freezeResultValue = field.resultValue;
						});
					});
				});
			});

			const arrIndexs = [];
			// 过滤数组为空的元素
			this.invoiceList.forEach((item, index) => {
				if (item.length === 0) {
					arrIndexs.push(index);
				}
			});

			if (arrIndexs.length !== 0) {
				for (let i = arrIndexs.length; i >= 0; i--) {
					this.$delete(this.invoiceList, arrIndexs[i]);
				}
			}

			this.pagination = {
				total: this.invoiceList.length
			};

			// 设置渲染页面数据
			this.setFieldEffectKeyValue({
				desserts: this.invoiceList[this.pageNum - 1]
			});
		},
		/** 表单聚焦 function
		 * @param {inputDom} 聚焦表单
		 */
		scrollToFunOne(path) {
			this.form.keyword = "";

			// this.setSearchEndData(tools.deepClone(this.sabayon), path); 临时
			this.setSearchEndData(this.sabayon, path);

			this.$nextTick(() => {
				const timer = setTimeout(() => {
					this.inputDom = this.$refs[`field_${path.ID}`][0];
					this.inputDom.focus();
					//this.$refs.lploading.close();
					clearTimeout(timer);
				}, 1000);
			});
		},
		/** 搜索关键字 function
		 * @param {searchResultSet}搜索结果
		 * @_includes 方法模糊搜索
		 */
		searchKeyword() {
			this.searchResultSet = [];
			this.invoiceList.map((invoice, parentIndex) => {
				invoice.map((item, index) => {
					item.desserts.map((desserts, idx) => {
						for (let i = 0; i < desserts.length; i++) {
							if (
								desserts[i].resultInput !== "no" &&
								(desserts[i].name == this.form.keyword ||
									desserts[i].code == this.form.keyword ||
									_.includes(
										desserts[i].resultValue.toLowerCase(),
										this.form.keyword.toLowerCase()
									))
							) {
								this.searchResultSet.push({
									panel: index,
									ID: desserts[i].ID,
									index: parentIndex,
									title: item.title,
									code: desserts[i].code,
									name: desserts[i].name,
									resultValue: desserts[i].resultValue
								});
								continue;
							}
						}
					});
				});
			});
			this.$refs.searchDialog.open();
		},
		computeEleHeight(dessertIndex, fieldIndex, arr) {
			let scrollContainer = document.querySelector(".scrollContainer");
			let len = fieldIndex;
			let arrlen = 0;
			for (let i = 0; i <= dessertIndex; i++) {
				arrlen += this.desserts[i].desserts[0].length;
				len += this.desserts[i - 1] ? this.desserts[i - 1].desserts[0].length : len;
			}
			let proportion = len / arrlen;

			let Height = scrollContainer.scrollHeight * proportion - 110;

			this.scrollCard.scrollTo({
				top: Height,
				left: 0,
				behavior: "smooth"
			});
		},
		/**function 表单聚焦
		 * @param {isSearch} is keyword search?
		 * @param {pathPanel} 手风琴展开页信息
		 * @param {noScrollExpandCard} 手动关闭手风琴索引，手动关闭后不再自动展开
		 * @param {contentsIndex} 当前目点击录索引
		 */
		// scrollToFun(dessertIndex, fieldIndex, arr) {
		// 	this.form.keyword = "";
		// 	//this.setFrontEndData(tools.deepClone(this.sabayon)); 临时
		// 	this.setFrontEndData(this.sabayon);
		// 	//this.computeEleHeight(dessertIndex, fieldIndex, arr);
		// 	this.$nextTick(() => {
		// 		const timer = setTimeout(() => {
		// 			this.inputDom.focus();
		// 			this.$refs.lploading.close();
		// 			//this.$refs[this.formKey][0].focus()
		// 			clearTimeout(timer);
		// 		}, 1100);
		// 	});
		// },
		/**function 保存展开手风琴的索引
		 * * @param {directly} true代表展开点击目录对应的分块
		 */
		genePanel(len, directly = false) {
			let arr = [];
			if (directly) {
				this.$set(this.panel, len, len);
				this.$set(this.memoPanels, len, true);
			} else {
				for (let i = 0; i <= len; i++) {
					arr[i] = i;
					this.memoPanels[i] = true;
				}

				this.noScrollExpandCard.map(index => {
					this.memoPanels[index] = false;
					arr[index] = -1;
				});
			}
			return arr;
		},
		resetData() {
			this.noScrollExpandCard = [];
			this.contents = [];
			this.contentsIndex = 0;
			this.lastScrollTop = 0;
			this.selectedAll = false;
			this.panel = [];
			this.memoPanels = {};
			this.invoiceList = [];
			this.billInfo = {};
			this.desserts = [];
			this.hintInfo = {};
			this.fields = [];
			this.modifiedFields = [];
			this.pageNum = 1;
			this.sabayon = {};
			this.isEqualSearch = false;
			this.wrongNote = "";
			this.isEnter = false;
			this.isContents = false;
			this.form = {
				id: "",
				code: "",
				money: "",
				keyword: ""
			};
			this.saleChannel = "";
			this.agency = "";
		},
		/**function 点击目录对应分块
		 */
		scrollAnimation() {
			let top = Math.abs(
				this.$refs.expantionRef[this.contentsIndex]?.$el.getBoundingClientRect().top
			);

			if (top < 100) {
				this.isContents = false;
			} else {
				this.$refs.expantionRef[this.contentsIndex]?.$el.scrollIntoView({
					block: "start",
					behavior: "smooth"
				});
			}
		},
		// 目录跳转
		handleContents(index) {
			this.isContents = true;
			this.contentsIndex = index;
			this.genePanel(index, true);
			this.$nextTick(() => {
				const timer = setTimeout(() => {
					this.$refs.expantionRef[this.contentsIndex].$el.scrollIntoView({
						block: "start",
						behavior: "smooth"
					});
					clearTimeout(timer);
				}, 300);
			});
		},
		// 搜索
		onSearch() {
			const form = this.form;
			// 清除搜索框空格
			let bool = Object.values(form).every(item => {
				if (item) {
					item = item.replace(/ /g, "");
				}
				return item === undefined || item === "";
			});

			// const result = tools.deepClone(this.sabayon); 临时
			const result = this.sabayon;
			if (!bool && (form.keyword === "" || form.keyword === undefined)) {
				if (result.data.invoice?.length === 0) {
					return; // select null
				}
				// select.code ?= currentdata.code
				this.searchCurrentPage(result.data.invoice);
				// 如果相同把整个result设置过去
				if (this.isEqualSearch) {
					// this.setFrontEndData(tools.deepClone(this.sabayon)); 临时
					this.setFrontEndData(this.sabayon);
					return;
				}
			}

			if (form.keyword && form.id === "" && form.code === "" && form.money === "") {
				this.searchKeyword();
				return;
			}

			this.setFrontEndData(result);
		},

		// compute pageNum
		searchCurrentPage(arr) {
			//const { id, code, keyword, money } = this.forms[this.searchFormId]
			const { id, code, keyword, money } = this.form;
			let dessertsArr = this.sabayon.data.invoice;
			// 单条数据，多条数据待处理
			let obj = arr.length === 1 ? arr[0] : arr;
			if (id) {
				this.mathCompute(dessertsArr, obj, id, "id");
			}
			if (!id && code) {
				this.mathCompute(dessertsArr, obj, code, "code");
			}
			if (!id && !code && money) {
				this.mathCompute(dessertsArr, obj, money, "money");
			}
		},
		mathCompute(dessertsArr, obj, condition, attr) {
			let index = dessertsArr.findIndex(item => item[attr] == condition);

			if (index === -1) {
				return;
			}
			if (dessertsArr[index][attr] === obj[attr]) {
				this.isEqualSearch = true;
			}
			this.pageNum = index + 1;
		},
		// 查看图片
		async onViewImage(row) {
			let routeData = this.$router.resolve({
				path: `view-image/${row.blockID}`,
				query: { fieldId: row.ID }
			});

			window.sessionStorage.setItem("proCode", this.cases.caseInfo.proCode);
			window.open(routeData.href, "_blank");
		},
		// 分页
		handlePage(page) {
			// 清空用户手动关闭card
			this.noScrollExpandCard = [];
			this.pageNum = page.pageNum;
			this.desserts = [];
			this.memoPanels = {};
			this.panelIndex = 0;
			this.contentsIndex = 0;
			const timer = setTimeout(() => {
				this.setFieldEffectKeyValue({
					desserts: this.invoiceList[page.pageNum - 1]
				});

				if (this.isSearch) {
					// 不同页跳转展开card
					this.isSearch = false;
					this.panel = this.genePanel(this.pathPanel);
				} else {
					this.panel = [0];
					this.memoPanels[0] = true;
				}
				clearTimeout(timer);
			}, 150);
		},
		// 接口
		async getdata() {
			const data = {
				editType: 2,
				billId: this.cases.caseInfo.caseId,
				proCode: this.cases.caseInfo.proCode,
				fields: this.fields
			};
			const results = await this.$store.dispatch("GET_OR_UPDATE_CASE_ENTRY_DATAS", data);
			this.toasted.dynamic(results.msg, results.code);

			if (results.code == 200) {
				this.setFields(results);
				this.setFrontEndData(results);
				this.toasted.success("发票查验成功 ");
			} else {
				this.setFields(results);
				this.setFrontEndData(results);
			}

			this.sabayon = results;
			this.billInfo = results.data.bill;
		},

		// 查询录入数据
		async getCaseEntryData() {
			const data = {
				editType: 1,
				billId: this.cases.caseInfo.caseId,
				proCode: this.cases.caseInfo.proCode
			};

			const result = await this.$store.dispatch("GET_OR_UPDATE_CASE_ENTRY_DATA", data);
			this.toasted.dynamic(result.msg, result.code);

			if (result.code == 200) {
				this.setFields(result);
				this.setFrontEndData(result);
				this.agency = result.data.bill.agency;
				this.saleChannel = result.data.bill.saleChannel;
			} else {
				this.onClose();
			}
			this.sabayon = result;
			this.billInfo = result.data.bill;
			this.objInvoice = result.data.invoice;
		},

		// 记忆输入框修改的值
		onFieldChange({ value, field }) {
			const modifiedField = {
				...field,
				resultValue: value,
				isChange: true
			};

			// modifiedFields
			{
				const field = tools.find(this.modifiedFields, {
					ID: modifiedField.ID
				});
				const memoModifiedField = {
					fieldId: modifiedField.ID,
					code: modifiedField.code,
					name: modifiedField.name,
					beforeVal: modifiedField.freezeResultValue,
					endVal: modifiedField.resultValue
				};

				if (!field) {
					this.modifiedFields.push(memoModifiedField);
				} else {
					const modifiedFieldsIndex = this.modifiedFields
						.map(item => item.ID === field.ID)
						.indexOf(true);
					this.modifiedFields[modifiedFieldsIndex] = memoModifiedField;
				}
			}

			// fields
			{
				const field = tools.find(this.fields, {
					ID: modifiedField.ID
				});

				if (tools.isYummy(field)) {
					const fieldIndex = this.fields.map(item => item.ID === field.ID).indexOf(true);
					this.fields[fieldIndex] = modifiedField;
				}
			}
		},
		/**function 滚动条事件
		 *@param {expantionRef} 手风琴dom元素
		 *1、区分点击目录展开手风琴和滚轮事件自动展开手风琴
		 *2、滚轮事件，当next card距离视口的top < 视口高度 + 视口高度 /4 。即展开next card
		 */
		onScroll(e) {
			let scrollTop = e.target.scrollTop;

			if (this.isContents) {
				lpTools.requestAnimationFrameFun(this.scrollAnimation, 700);
			} else {
				if (this.lastScrollTop < scrollTop) {
					if (!this.$refs.expantionRef[this.panelIndex]) {
						return;
					}
					let top =
						this.$refs.expantionRef[this.panelIndex].$el.getBoundingClientRect().top;
					let elHeight = this.$refs.expantionRef[this.panelIndex].$el.clientHeight;

					if (elHeight <= this.viewHeight + this.viewHeight / 5) {
						this.panelIndex = this.panelIndex + 1;
					} else {
						if (top < this.viewHeight + this.viewHeight / 5) {
							this.panelIndex = this.panelIndex + 1;
						}
					}
				}
			}

			this.lastScrollTop = scrollTop;
		},

		// 回车提交输入框修改值
		async onFieldEnter() {
			//this.isEnter = true;
			const body = {
				editType: 2,
				billId: this.cases.caseInfo.caseId,
				fields: this.fields,
				modifiedFields: this.modifiedFields,
				proCode: this.cases.caseInfo.proCode
			};

			const result = await this.$store.dispatch("GET_OR_UPDATE_CASE_ENTRY_DATA", body);

			this.toasted.dynamic(result.msg, result.code);

			this.isEnter = true;
			if (result.code === 200) {
				this.setFields(result);
				this.setFrontEndData(result, "enter");
				this.wrongNote = result.data.bill?.wrongNote;
				this.billInfo = result.data.bill;
				this.sabayon = result;
				this.objInvoice = result.data.invoice;
				// console.log(this.objInvoice);
				//this.onSearch()
			}
		},

		// 获取字段、分块和字段配置信息
		onHint(item) {
			this.hintInfo = {
				proCode: this.cases.caseInfo.proCode,
				fieldId: item.ID
			};
			this.$refs.info.onOpen();
		},

		// 全选/不全选
		async onSelectAll(value) {
			if (tools.isLousy(this.desserts)) {
				return;
			}

			this.selectedAll = value;

			let blockIds = [];

			this.desserts.map(item => {
				item.desserts?.map(fields => {
					fields?.map(field => {
						blockIds = [...blockIds, field.blockID];
					});
				});
			});

			const data = {
				blockIds,
				isPractice: value,
				proCode: this.cases.caseInfo.proCode
			};

			const result = await this.$store.dispatch("UPDATE_CASE_BLOCK_PRACTICE", data);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code === 200) {
				this.desserts.map(item => {
					item.desserts?.map(fields => {
						fields.map(field => {
							field.isPractice = value;
						});
					});
				});

				this.desserts = [...this.desserts];
			} else {
				this.selectedAll = false;
			}
		},

		// 选择/不选择
		async onSingleSelect({ value, field, index, fieldsIndex, fieldIndex }) {
			this.desserts[index]["desserts"][fieldsIndex][fieldIndex].isPractice = value;

			const data = {
				blockIds: [field.blockID],
				isPractice: value,
				proCode: this.cases.caseInfo.proCode
			};

			const result = await this.$store.dispatch("UPDATE_CASE_BLOCK_PRACTICE", data);

			this.toasted.dynamic(result.msg, result.code);

			if (result.code !== 200) {
				this.desserts[index]["desserts"][fieldIndex].isPractice = false;
			}

			let flatern = [];

			this.desserts.map(item => {
				if (tools.isYummy(item.desserts)) {
					flatern = [...flatern, ...item.desserts];
				}
			});

			this.selectedAll = flatern.every(item => item.isPractice);
		},

		onSwitch() {
			this.$router.push({
				path: "/main/PM/case/view-result-data"
			});
		},

		onNavigate({ value }) {
			if (value === 1) {
				this.$router.push({
					path: "/main/PM/case/report-table"
				});
			} else if (value === 2) {
				let reg = new RegExp("/files/files/", "g");
				let convert = new RegExp("/convert_", "g");

				let prefixUrl = `${baseURLApi}files/${this.row.downloadPath}`;
				prefixUrl = prefixUrl.replace(reg, "/files/");
				prefixUrl = prefixUrl.replace(convert, "/");

				const thumbs = [];

				this.row.pictures.map(image => {
					thumbs.push({
						thumbPath: `${prefixUrl}A${image}`,
						path: `${prefixUrl}${image}`
					});
				});
				sessionStorage.set("thumbs", thumbs);

				window.open(
					`${location.origin}/normal/view-images`,
					"_blank",
					"toolbar=yes, scrollbars=yes, resizable=yes"
				);
			} else if (value === 3) {
				this.$refs.export.onOpen();
			} else {
				const data = {
					...this.cases.caseInfo
				};
				this.$refs.log.getLogs(data);
				this.$refs.log.onOpen();
			}
		},

		// 保存为xml
		async onSaveAsXML() {
			const body = {
				editType: 2,
				billId: this.cases.caseInfo.caseId,
				proCode: this.cases.caseInfo.proCode,
				fields: this.fields
			};

			if (this.row.stage === 5) {
				this.row.stage = 3;
				//this.$store.commit('UPDATE_CASE', { rowInfo: this.row })
				this.$EventBus.$emit("updateStage", this.row);
			}

			const result = await this.$store.dispatch("SAVE_CASE_ENTRY_DATA_AS_XML", body);

			if (result.code === 200) {
				this.toasted.dynamic("保存成功", result.code);
			} else {
				this.toasted.dynamic(result.msg, result.code);
			}

			this.wrongNote = result.data;

			console.log(this.sabayon.data, "XML-before");
			this.$set(this.sabayon.data.bill, "wrongNote", result.data);
			console.log(this.sabayon.data, "XML-after");
		},

		// 设置字段
		setFields(result) {
			this.fields = [];
			const { invoice } = result.data;

			invoice?.map(item => {
				for (let key in item) {
					if (invoiceKeys.includes(key)) {
						item[key]?.map(fields => {
							fields?.map(field => {
								//console.log(field.name, field)

								this.fields.push({
									...field,
									freezeResultValue: field.resultValue
								});
							});
						});
					}
				}
			});
		},

		// 设置前端数据结构
		setFrontEndData(result, status) {
			const { invoice } = result.data;

			this.billInfo = invoice.bill;
			this.wrongNote = result.data.bill?.wrongNote;

			if (!this.isEnter) {
				// 恢复默认
				{
					this.invoiceList = [];

					this.pagination = {
						total: 0
					};

					this.desserts = [];

					if (!tools.isYummy(this.memoPanels)) this.panel = [];

					if (tools.isLousy(invoice)) {
						return;
					}
				}
			}

			invoice.map((item, index) => {
				this.invoiceList[index] = [];

				for (let key in item) {
					if (invoiceKeys.includes(key)) {
						if (invoice[index][key] != null && invoice[index][key].length > 0) {
							if (invoice[index][key][0] != null) {
								invoice[index][key].map(arr => {
									const arrs = _.chunk(arr, Math.floor(arr.length / 2));
									let leftdesserts = [],
										rightdesserts = [];
									if (arrs.length > 2) {
										const temp = [...arrs[1], ...arrs[2]];
										leftdesserts = [arrs[0]];
										rightdesserts = [temp];
									} else {
										leftdesserts = [arrs[0]];
										rightdesserts = [arrs[1]];
									}

									this.invoiceList[index] = [
										...this.invoiceList[index],
										{
											key,
											title: invoiceNamesMap.get(key),
											desserts: [arr],
											leftdesserts,
											rightdesserts
										}
									];
								});
							} else {
								this.invoiceList.splice(index, 1);
							}
						}
					}
				}
			});

			// 记忆最开始的resultValue
			this.invoiceList?.map(invoice => {
				invoice?.map(item => {
					item.desserts?.map(fields => {
						fields?.map(field => {
							field.freezeResultValue = field.resultValue;
						});
					});
				});
			});

			const arrIndexs = [];
			// 过滤数组为空的元素
			this.invoiceList.forEach((item, index) => {
				if (item.length === 0) {
					arrIndexs.push(index);
				}
			});

			if (arrIndexs.length !== 0) {
				for (let i = arrIndexs.length; i >= 0; i--) {
					this.$delete(this.invoiceList, arrIndexs[i]);
				}
			}

			this.pagination = {
				total: this.invoiceList.length
			};

			this.setFieldEffectKeyValue({
				desserts: this.invoiceList[this.pageNum - 1]
			});

			// if (!tools.isYummy(this.memoPanels)) {
			//   for (let i = 0; i < this.desserts.length; i++) {
			//     this.panel = [...this.panel, i]
			//     this.memoPanels[i] = true
			//   }
			// }

			if (!this.isEnter) {
				this.panel = [0];
				this.memoPanels[0] = true;
				if (this.isKeywordSearch) {
					this.panel = this.genePanel(this.desserts.length);
				}
			}

			this.isKeywordSearch = false;
			this.isEnter = false;
		},

		setFieldEffectKeyValue({ desserts }) {
			this.contents = [];
			let countIndex = 0;
			desserts?.map((dessert, index) => {
				dessert.showIndex = index;
				this.contents.push({
					key: dessert.key,
					index,
					title: index + "." + dessert.title
				});

				dessert?.desserts?.map(fields => {
					fields?.map(field => {
						field.countIndex = countIndex;
						++countIndex;
					});
				});
			});
			this.desserts = desserts;
		},

		handlePanel(index) {
			this.memoPanels[index] = !this.memoPanels[index];
			this.panel = [];

			for (let key in this.memoPanels) {
				const value = this.memoPanels[key];
				if (value) {
					this.panel[key] = key;

					const idx = this.noScrollExpandCard.findIndex(val => val === key);

					if (idx !== -1) {
						this.noScrollExpandCard.splice(idx, 1);
					}
				} else {
					// 手风琴关闭后，不再自动展开
					this.noScrollExpandCard.push(key);
					this.noScrollExpandCard = Array.from(new Set(this.noScrollExpandCard));
					this.panel[key] = -1;
				}
			}
		},
		// 快捷键
		fuckEvents(event) {
			event = event || window.event;
			const { ctrlKey, keyCode } = event;

			switch (keyCode) {
				case 83:
					if (ctrlKey) {
						event.preventDefault();
						this.onSaveAsXML();
					}
					break;
			}
		}
	},

	components: {
		"field-info-dialog": () => import("./fieldInfoDialog"),
		"change-log-dialog": () => import("./changeLogDialog"),
		"export-validation": () => import("./exportValidation"),
		"content-com": () => import("./contents/contents.vue"),
		searchResultData: () => import("./searchResultData/searchResultData.vue"),
		LPLoading
	}
};
</script>

<style lang="scss">
/* .contain {
    position: relative;
    padding-top: 64px;
    height: calc(100vh - 132px);
  } */

.hint {
	min-height: 22px;
	color: rgba(0, 0, 0, 0.6);
	cursor: pointer;
}

.searchFlex {
	position: sticky;
	top: 57px;
	z-index: 100;
	background-color: #fff;
	box-shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.2);
	padding: 10px;
	box-sizing: border-box;
}

.z-card_fixed-actions {
	position: fixed;
	width: 100%;
	bottom: 0;
	background-color: #fff;
	border-top: 1px solid rgba(0, 0, 0, 0.12);
	z-index: 1;
}

.action {
	.field {
		position: relative;
	}

	.label {
		width: 300px;
	}

	.hint {
		width: 200px;
	}
}
</style>
