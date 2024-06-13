import { mapState } from 'vuex'
import { tools } from 'vue-rocket'

const defaultParams = {
	pageSize: 10,
	pageIndex: 1
}

const defaultSabayon = {
	data: {},
	total: 0
}

export default {
	data() {
		return {
			expanded: false,
			searchId: 'Search',
			searchFormId: '',

			tableBorder: true,
			tableMaxHeight: 650,
			tableSize: 'mini',
			tableStripe: true,

			loading: true,
			loadingText: 'Loading... Please wait',

			isDeleteMore: false,
			params: defaultParams,
			effectParams: {},
			desserts: [],
			selected: [],
			ids: [],
			// detail: {},
			detailInfo: {},

			pagination: {
				total: 0
			},
			pageSizes: [
				{ label: '10条/页', value: 10 },
				{ label: '20条/页', value: 20 },
				{ label: '50条/页', value: 50 },
				{ label: '100条/页', value: 100 },
				{ label: '200条/页', value: 200 },
				{ label: '500条/页', value: 500 }
			],
			code: sessionStorage.getItem('proCode'),

			sabayon: defaultSabayon
		}
	},

	created() {
		this._stickFormId()

		if (!this.manual) {
			this.getList()
		}
	},

	methods: {
		/**
		 * @description 检索
		 */
		onSearch() {
			this.params = {
				...this.params,
				...this.forms[this.searchFormId],
				...defaultParams
			}
			this.getList()
		},

		/**
		 * @description 分页
		 * @param {{ 
		 * 	pageNum: Number, 
		 * 	pageSize: Number 
		 * }} 页码、条数
		 */
		handlePage(page) {
			this.params = {
				...this.params,
				pageSize: page.pageSize,
				pageIndex: page.pageNum
			}
			this.getList()
		},

		/**
		 * @description 获取列表
		 * @return {Object}
		 */
		async getList() {
			if (this.dispatchList) {
				const params = {
					...this.effectParams,
					...this.params
				}

				const result = await this.$store.dispatch(this.dispatchList, params)
				const { list, total } = result.data

				if (result.code === 200) {
					if (typeof list === 'object') {
						if (list instanceof Array) {
							this.desserts = list
						} else {
							this.desserts = []
						}
						this.pagination.total = total
					}
					else {
						this.desserts = result.data
						this.pagination.total = this.desserts.length
					}
				} else {
					this.toasted.error(result.msg)

					this.desserts = []
					this.pagination.total = 0
				}

				this.sabayon = result
			}

			this.loading = false

			return this.sabayon
		},

		/**
		 * @description 新增/编辑
		 * @param {Object} form
		 * @param {String} dispatchForm
		 * @param {String} ref
		 */
		async updateListItem(form, dispatchForm, ref = 'dynamic', list = 'getList') {

			const result = await this.$store.dispatch(dispatchForm, form)

			this.toasted.dynamic(result.msg, result.code)

			if (result.code === 200) {
				this[list]()
				this.$refs[ref].close()
			}

			return result
		},

		/**
		 * @description 获取详情
		 * @param {Object} row
		 */
		getDetail(row) {
			this.detailInfo = row
			console.log(this.detailInfo)
		},

		/**
		 * @description 修改单个
		 * @param {Object} form
		 */
		async modifyCell(form) {
			form.status = 1
			const result = await this.$store.dispatch(this.dispatchCellForm, form)
			this.toasted.dynamic(result.msg, result.code)
			if (result.code === 200) {
				this.getList()
			}
			return result
		},

		/**
		 * @description 批量删除
		 * @return {any}
		 */
		onDeleteMore() {
			this.$modal({
				visible: true,
				title: '批量删除提示',
				content: '请确认是否要批量删除？',
				confirm: () => {
					this.delRows()
				}
			})
		},

		/**
		 * @description 删除一个
		 * @return {any}
		 */
		deleteItem() {
			this.$modal({
				visible: true,
				title: '删除提示',
				content: '请确认是否要删除？',
				confirm: () => {
					this.delRows()
				}
			})
		},

		/**
		 * @description 删除一个/多个
		 * @return {any}
		 */
		async delRows() {
			if (this.dispatchDelete) {

				if (!this.isDeleteMore) {
					this.ids = [this.detailInfo.ID]
				}

				const length = this.ids.length

				const result = await this.$store.dispatch(this.dispatchDelete, this.ids)

				this.toasted.dynamic(result.msg, result.code)

				if (result.code === 200) {
					const oldTotal = this.pagination.total
					const newTotal = oldTotal - length
					let { pageSize, pageIndex } = this.params

					if (newTotal <= pageSize) {
						pageIndex -= 1
						if (pageIndex < 1) {
							pageIndex = 1
						}
					}

					this.params.pageIndex = pageIndex

					this.getList()
				}

				return result
			}
		},

		/**
		 * @description 全选/全不选
		 * @param {any} event
		 */
		handleSelectAll(event) {
			const { records } = event
			this.selected = records
		},

		/**
		 * @description 选中/不选中
		 * @param {any} event
		 */
		handleSelectChange(event) {
			const { records } = event
			this.selected = records
		},

		/**
		 * @description 每页递增序号
		 * @param {any} event
		 */
		increaseSeq({ seq }) {
			const { pageIndex, pageSize } = this.params
			return ((pageIndex - 1) * pageSize + seq)
		},

		/**
		 * @description 根据id显示文本
		 * @param {Array} value
		 * @param {String} target
		 */
		showText(value, target) {
			const result = tools.find(value, target)
			return result ? result.label : undefined
		},

		/**
		 * @description 更多/收起
		 */
		onExpand() {
			this.expanded = !this.expanded
		},

		/**
		 * @description searchFormId
		 */
		_stickFormId() {
			this.searchFormId = this.formId + this.searchId
		}
	},

	computed: {
		...mapState(['forms'])
	},

	watch: {
		selected: {
			handler() {
				this.ids = []
				for (let item of this.selected) {
					this.ids.push(item.ID || item.id)
				}
			}
		},

		ids: {
			handler() {
				if (this.ids.length) {
					this.isDeleteMore = true
				} else {
					this.isDeleteMore = false
				}
			}
		}
	}
}