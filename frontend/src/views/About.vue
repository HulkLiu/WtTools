<template>
  <n-space vertical :size="12">
    <n-space>
      <n-button @click="downloadCsv">
        导出 CSV（原始数据）
      </n-button>
      <n-button @click="exportSorterAndFilterCsv">
        导出 CSV（展示的数据）
      </n-button>
    </n-space>
    <n-data-table
        ref="tableRef"
        :columns="columns"
        :data="data"
        :pagination="pagination"
        :bordered="false"
    />
  </n-space>
</template>

<script>
import { defineComponent, ref } from "vue";


const columns = [
  {
    title: "Name",
    key: "name",
    sorter: "default"
  },
  {
    title: "Age",
    key: "age",
    sorter: (row1, row2) => row1.age - row2.age
  },
  {
    title: "Address",
    key: "address",
    filterOptions: [
      {
        label: "London",
        value: "London"
      },
      {
        label: "New York",
        value: "New York"
      }
    ],
    filter: (value, row) => {
      return !!~row.address.indexOf(value);
    }
  }
];

const data = [
  {
    key: 0,
    name: "John Brown",
    age: 18,
    address: "New York No. 1 Lake Park"
  },
  {
    key: 1,
    name: "Jim Green",
    age: 28,
    address: "London No. 1 Lake Park"
  },
  {
    key: 2,
    name: "Joe Black",
    age: 38,
    address: "Sidney No. 1 Lake Park"
  },
  {
    key: 3,
    name: "Jim Red",
    age: 48,
    address: "London No. 2 Lake Park"
  }
];

export default defineComponent({
  setup() {
    const tableRef = ref();

    const downloadCsv = () => tableRef.value?.downloadCsv({ fileName: "data-table" });

    const exportSorterAndFilterCsv = () => tableRef.value?.downloadCsv({
      fileName: "sorter-filter",
      keepOriginalData: false
    });

    return {
      data,
      tableRef,
      downloadCsv,
      exportSorterAndFilterCsv,
      columns,
      pagination: false
    };
  }
});
</script>