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
import {defineComponent, h, onMounted, ref} from "vue";
import {SearchCourse} from "../../wailsjs/go/app/App.js";
import {NButton,NTag} from "naive-ui";


export default defineComponent({
  setup() {
    const tableRef = ref();
    const listData = ref()

    const downloadCsv = () => tableRef.value?.downloadCsv({ fileName: "data-table" });

    const exportSorterAndFilterCsv = () => tableRef.value?.downloadCsv({
      fileName: "sorter-filter",
      keepOriginalData: false
    });
    onMounted(() => {
      List()
    })
    const List = () =>{
      // alert(1)
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
    const openUrl = (url) => {
      window.open(url, '_blank');

    }

    return {
      keyword,
      List,
      listData,
      tableRef,
      downloadCsv,
      exportSorterAndFilterCsv,
      columns : [
        {
          title: "分类",
          key: "LanguageType",
          width: 100,
          ellipsis: true,
          render (row) {
            return h(
                NTag,
                {
                  size: "small",
                },
                { default: () => row.Payload.LanguageType}
            );
          }
        },
        {
          title: "最后更新",
          key: "LastTime",
          width: 100,
          ellipsis: true,
          render (row) {
            return h(
                NTag,
                {
                  size: "small",
                },
                { default: () => row.Payload.LastTime}
            );
          }
        },
        {
          title: "标题",
          key: "Title",
          width: 300,
          ellipsis: true,
          render (row) {
            return h(
                NButton,
                {
                  size: "small",
                  onClick: () => showDataModal(row)
                },
                { default: () => row.Payload.Title}
            );
          }
        },

        {
          title: "Url",
          key: "Url",
          width: 300,
          ellipsis: true,
          render(row) {
            return h(
                NButton,
                {
                  strong: true,
                  tertiary: true,
                  size: "small",
                  onClick: () => openUrl(row.Url)
                },
                { default: () => row.Url }
            );
          }
        },

        // {
        //   title: "操作",
        //   key: "Title",
        //   render (row) {
        //     return h(
        //         NButton,
        //         {
        //           size: "small",
        //           onClick: () => deleteRow(row)
        //         },
        //         { default: () => "删除"}
        //     );
        //   }
        // },
        // {
        //   title: '本地链接',
        //   key: 'ShortLocal',
        // }
      ],

      pagination: 20
    };
  }
});
</script>