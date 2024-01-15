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
    <n-space>
      <n-input v-model:value="keyword" placeholder="请输入关键字" @keyup.enter="List" />
    </n-space>
    <n-data-table
        ref="tableRef"
        :columns="columns"
        :data="listData"
        :pagination="pagination"
        :bordered="false"
    />
  </n-space>
</template>

<script>
import { defineComponent, ref } from "vue";
import {SearchCourse} from "../../wailsjs/go/app/App.js";


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
    const listData = ref()

    const downloadCsv = () => tableRef.value?.downloadCsv({ fileName: "data-table" });

    const exportSorterAndFilterCsv = () => tableRef.value?.downloadCsv({
      fileName: "sorter-filter",
      keepOriginalData: false
    });
    const List = () =>{
      alert(1)
      SearchCourse(JSON.stringify(keyword.value)).then(res => {
        console.log(res)

        if (res.code === 0){
          if (res.data === null){
            listData.value = []
          }else{
            listData.value = res.data;

            // listData.value.CreatedAt = moment(listData.value.CreatedAt).format("YYYY-MM-DD HH:mm:ss")
            // listData.value.UpdatedAt = moment(listData.value.UpdatedAt).format("YYYY-MM-DD HH:mm:ss")

          }
          console.log(listData.value)
        }
      })
    }
    const keyword = ref()

    return {
      keyword,
      List,
      listData,
      tableRef,
      downloadCsv,
      exportSorterAndFilterCsv,
      columns,
      pagination: false
    };
  }
});
</script>