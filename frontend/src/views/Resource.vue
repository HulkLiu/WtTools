<template>
  <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
    <el-header class="header-area" style="height: 89px;">
      <n-button @click="dialogVisibleAdd"> 新增任务 </n-button>
      <n-button @click="downloadCsv"> naiveUI 导出 CSV（All）</n-button>
      <n-popover trigger="click" >
        <template #trigger>
          <n-button @click="checkUrl">检测 URL 可用</n-button>
        </template>
        <span>
          <div  v-for ="(item, index)  in checkMsg">
              {{item}} => {{index}}
          </div>

        </span>
      </n-popover>
    </el-header>
    <n-message-provider>

      <el-main class="main-area" style="--wails-draggable:no-drag">

        <n-data-table
            ref="tableRef"
            :columns="columns"
            :data="listData"
            :pagination="pagination"
            :bordered="false"

        />

        <n-modal v-model:show="dialogVisible" style="width: 1500px">
          <n-card
              style="width: 67%"
              title="新增/修改"
              :bordered="false"
              size="huge"
              role="dialog"
              aria-modal="true"

              :style="{maxWidth: '840px'}"
              :label-width="90"
          >

            <template #header-extra></template>

            <n-form  :model="form"  label-placement="left">
              <n-form-item v-show="false" path="id" label="id">
                <n-input :disabled="true" v-model:value="form.Id"  />
              </n-form-item>

              <n-form-item path="name" label="资源名称">
                <n-input v-model:value="form.Name" placeholder="资源名称" />
              </n-form-item>
              <n-form-item path="Description" label="URL">
                <n-input v-model:value="form.URL" placeholder="URL" />
              </n-form-item>
              <n-form-item label="状态" path="completed">
                <n-select
                    v-model:value="form.Status"
                    placeholder="Select"
                    :options="generalOptions"
                />
              </n-form-item>
              <n-form-item path="createdAt" label="创建时间">
                <n-input :disabled="true" v-model:value="form.CreatedAt" placeholder="创建时间" />
              </n-form-item>
              <n-form-item path="updatedAt" label="更新时间">
                <n-input :disabled="true" v-model:value="form.UpdatedAt" placeholder="更新时间" />
              </n-form-item>

              <n-form-item>
                <n-button v-if="form.Id === '' " type="primary" @click="onAdd"
                          :disabled="form.Name === '' ||
                          form.URL === '' ||
                          form.Status === ''"
                          :loading="loading"
                >创建</n-button>
                <n-button v-else type="primary" @click="onAdd">保存</n-button>
                <n-button @click="dialogVisible = false">取消</n-button>
              </n-form-item>

            </n-form>

            <template #footer>

            </template>
          </n-card>
        </n-modal>
      </el-main>
    </n-message-provider>
  </el-container>

</template>

<script>
import {h, defineComponent, onMounted, ref, reactive} from "vue";
import { NButton, useMessage } from "naive-ui";
import {AddSource,GetSourceList,EditSource,CheckSourceList} from "../../wailsjs/go/app/App.js";
import moment from 'moment';
import {ElNotification} from "element-plus";



export default defineComponent({

  setup() {
    const listData = ref()
    const tableRef = ref();
    const downloadCsv = () => tableRef.value?.downloadCsv({ fileName: "ResourceData" })

    const message = useMessage()
    const form = ref()

    const List = () =>{
      GetSourceList().then(res => {
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

    const Add = () =>{
      if (modelRef.value.CreatedAt==="" ){
        modelRef.value.CreatedAt = new Date().toISOString();
      }
      if ( modelRef.value.UpdatedAt===""){
        modelRef.value.UpdatedAt = new Date().toISOString();
      }
      console.log(modelRef.value.Id)
      if (modelRef.value.ID === null) {
        modelRef.value.Id = 0
        AddSource(JSON.stringify(modelRef.value) ).then(res => {
          // console.log(res)
          if (res.code === 0) {
            message.success("添加成功")
            dialogVisible.value = false
            List()

          }else{
            message.warning("添加失败")
            message.error(res.msg)
          }

        })
      }else{
        edit()
      }


    }

    const showEdit = (row) =>{
      dialogVisible.value = true
      modelRef.value = row
      console.log(modelRef.value )
    }

    const openUrl = (url) => {
      window.open(url, '_blank');

    }
    const edit = (row) =>{
      EditSource(JSON.stringify(modelRef.value) ).then(res => {
        // console.log(res)
        if (res.code === 0) {
          message.success("编辑成功")
          dialogVisible.value = false

        }else{
          message.warning("编辑失败")
          message.error(res.msg)
        }

        List()
      })


    }

    onMounted(() => {
      List()
    })
    const dialogVisible = ref(false)
    const modelRef = ref({
      ID:null ,
      Name:"Learn Go ",
      URL:"Study the basics of Go" ,
      Status:"使用中" ,
      CreatedAt:"2023-12-28 17:36:47.7258618 +0800 CST m=+16.554803701 ",
      UpdatedAt:"2023-12-28 17:36:47.7258618 +0800 CST m=+16.554803701",
    })
    const checkMsg = ref()
    const checkUrl = () => {
      CheckSourceList().then(res => {
        console.log(res)

        if (res.code === 0){
          if (res.data === null){
            checkMsg.value = res.data
          }else{
            checkMsg.value = res.data
          }
        }
      })

    }
    return {
      checkMsg,
      checkUrl,
      tableRef,
      downloadCsv,
      dialogVisible,
      listData,
      form:modelRef,
      pagination: { pageSize: 30 },
      List,
      onAdd(){
        Add()
      } ,

      dialogVisibleAdd(){
        modelRef.value = ({
          ID:null ,
          Name:"" ,
          URL:"" ,
          Status:"",
          CreatedAt:"",
          UpdatedAt:"",
        })
        dialogVisible.value = true
      },
      generalOptions: ["使用中","停用"].map(
          (v) => ({
            label: v,
            value: v
          })
      ),
      columns : [
        {
          type: "selection",
        },
        {
          title: "编号",
          key: "ID"
        },
        {
          title: "资源名称",
          key: "Name"
        },
        {
          title: "URL",
          key: "URL",
          render(row) {
            return h(
                NButton,
                {
                  strong: true,
                  tertiary: true,
                  size: "small",
                  onClick: () => openUrl(row.URL)
                },
                { default: () => row.URL }
            );
          }
        },
        {
          title: "状态",
          key: "Status",
          defaultFilterOptionValues: [],
          filterOptions: [
            {
              label: "使用中",
              value:"使用中",

            },
            {
              label: "停用",
              value: "停用"
            }
          ],
          filter(value, row) {
            return !!~row.Status.indexOf(value);
          },
          render (row) {
            return h(
                NButton,
                {
                  size: "small",
                  onClick: () => showEdit(row)
                },
                { default: () => row.Status}
            );
          }
        },

        // {
        //   title: "创建时间",
        //   key: "CreatedAt",
        //
        // },
        // {
        //   title: "修改时间",
        //   key: "UpdatedAt"
        // },
        // {
        //   title: "操作",
        //   key: "Title",
        //   render (row) {
        //     return h(
        //         NButton,
        //         {
        //           size: "small",
        //           onClick: () => alert(row)
        //         },
        //         { default: () => "删除"}
        //     );
        //   }
        // },
      ],


    };
  }
});
</script>