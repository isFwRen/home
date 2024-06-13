<template>
  <div class="lp-salary">
    <div class="z-flex align-end mb-4 lp-filters">
      <v-row>
        <v-col cols="4">
          <z-date-picker
            :formId="searchFormId"
            formKey="date"
            label="日期"
            picker-type="month"
            z-index="10"
            range
            @input="onSearch"
          >
          </z-date-picker>
        </v-col>
      </v-row>

      <p class="mb-0">说明：当月工资为预估数据，最终工资于次月底统计的为准。</p>
    </div>

    <div class="table entry-table">
      <vxe-table :border="tableBorder" :data="desserts" :size="tableSize">
        <template v-for="item in cells.headers">
        <vxe-column 
          v-if="item.value === 'payDay'" 
          :field="item.value" 
          :title="item.text" 
          :key="item.value"
        >
          <template #default="{ row }">
            {{ row.payDay | dateFormat('YYYY-MM') }}
          </template>
        </vxe-column>

        <vxe-column 
          v-else 
          :field="item.value" 
          :title="item.text" 
          :key="item.value"
        > 
        </vxe-column>
        </template>
      </vxe-table>
    </div>
    <z-pagination
      :options="pageSizes"
      :total="pagination.total"
      @page="handlePage"
    ></z-pagination>
  </div>
</template>

<script>
  import TableMixins from '@/mixins/TableMixins'
  import cells from './cells'

  export default {
    name: 'Salary',
    mixins: [TableMixins],

    data() {
      return {
        cells,
        formId: 'salary',
        dispatchList: 'GET_SALARY_PT_LIST'
      }
    }
  }
</script>