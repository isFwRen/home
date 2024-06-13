<template>
  <div class="lp-error">
    <div class="mb-8 lp-filters">
      <v-row class="z-flex align-end">
        <v-col :cols="2">
          <z-select
            :formId="searchFormId"
            formKey="proCode"
            clearable
            hideDetails
            label="项目"
            :options="auth.proItems"
          ></z-select>
        </v-col>

        <v-col :cols="2">
          <z-date-picker
            :formId="searchFormId"
            formKey="month"
            hideDetails
            label="月份"
            picker-type="month"
            :defaultValue="cells.DEFAULT_MONTH"
          ></z-date-picker>
        </v-col>

        <v-col :cols="2">
          <z-text-field
            :formId="searchFormId"
            formKey="billName"
            hideDetails
            label="案件号"
          >
          </z-text-field>
        </v-col>

        <v-col :cols="2">
          <z-text-field
            :formId="searchFormId"
            formKey="wrongFieldName"
            hideDetails
            label="错误字段"
          >
          </z-text-field>
        </v-col>

        <!-- 查询按钮 BEGIN -->
        <div class="z-flex">
          <z-btn class="pb-3" color="primary" @click="onSearch">
            <v-icon class="text-h6">mdi-magnify</v-icon>
            查询
          </z-btn>
        </div>
        <!-- 查询按钮 END -->
      </v-row>
    </div>

    <div class="table error-detail-table">
      <vxe-table
        :data="desserts"
        :border="tableBorder"
        :max-height="tableMaxHeight"
        :size="tableSize"
        :stripe="tableStripe"
        @checkbox-change="handleSelectChange"
      >
        <template v-for="item in cells.headers">
          <!-- 影像 BEGIN -->
          <vxe-column
            v-if="item.value === 'imagePath'"
            :field="item.value"
            :title="item.text"
            :key="item.value"
            :width="item.width"
          >
            <template #default="{ row, rowIndex }">
              <z-upload
                formId="paths"
                :formKey="`path${ rowIndex }`"
                show-only
                :defaultValue="setImages(row.imagePath)"
              ></z-upload>
            </template>
          </vxe-column>
          <!-- 影像 END -->

          <!-- 录入日期 BEGIN -->
          <vxe-column 
            v-else-if="item.value === 'entryDate'"
            :field="item.value" 
            :title="item.text"
            :key="item.value"
            :width="item.width"
          >
            <template #default="{ row }">
              {{ row[item.value] | dateFormat('YYYY-MM-DD') }}
            </template>
          </vxe-column>
          <!-- 录入日期 END -->

          <!-- 反馈日期 BEGIN -->
          <vxe-column 
            v-else-if="item.value === 'feedbackDate'"
            :field="item.value" 
            :title="item.text"
            :key="item.value"
            :width="item.width"
          >
            <template #default="{ row }">
              {{ row[item.value] | dateFormat('YYYY-MM-DD') }}
            </template>
          </vxe-column>
          <!-- 反馈日期 END -->

          <!-- 正确数据 BEGIN -->
          <vxe-column
            v-else-if="item.value === 'right'"
            :field="item.value"
            :title="item.text"
            :key="item.value"
            :width="item.width"
          >
            <template #default="{ row }">
              <span
                v-if="_tools.compareString(row.right, row.wrong, 'error--text').targetHtml"
                v-html="_tools.compareString(row.right, row.wrong, 'error--text').targetHtml"
              ></span>
              <span v-else>{{ row.right }}</span>
            </template>
          </vxe-column>
          <!-- 正确数据 END -->

          <!-- 错误数据 BEGIN -->
          <vxe-column
            v-else-if="item.value === 'wrong'"
            :field="item.value"
            :title="item.text"
            :key="item.value"
            :width="item.width"
          >
            <template #default="{ row }">
              <span 
                v-if="_tools.compareString(row.wrong, row.right, 'error--text').targetHtml"
                v-html="_tools.compareString(row.wrong, row.right, 'error--text').targetHtml"
              ></span>
              <span v-else>{{ row.wrong }}</span>
            </template>
          </vxe-column>
          <!-- 错误数据 END -->

          <vxe-column
            v-else
            :field="item.value"
            :title="item.text"
            :key="item.value"
            :width="item.width"
          ></vxe-column>
        </template>
      </vxe-table>
    </div>

    <z-pagination
      class="mt-4"
      :total="pagination.total"
      :options="pageSizes"
      @page="handlePage"
    ></z-pagination>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import TableMixins from '@/mixins/TableMixins'
  import { tools as lpTools } from '@/libs/util'
  import cells from './cells'

  const { baseURLApi } = lpTools.baseURL()

  export default {
    name: 'Complaint',
    mixins: [TableMixins],

    data() {
      return {
        formId: 'Complaint',
        cells,
        dispatchList: 'GET_COMPLAINT_LIST',
        manual: true
      }
    },

    computed: {
      ...mapGetters(['auth'])
    },

    methods: {
      setImages(paths) {
        return paths?.map(path => ({ url: `${ baseURLApi }${ path }` }))
      }
    }
  }
</script>
