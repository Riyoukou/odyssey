<template>
  <div class="flex flex-col overflow-auto">
    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="form.toAdd">新增</ElButton>
        <el-popconfirm title="Are you sure to delete this?"  @confirm="deleteSelected">
          <template #reference>
            <ElButton icon="Delete">删除</ElButton>
          </template>
        </el-popconfirm>
      </div>
      <div>
        <ElButton icon="Refresh" round @click="table.request"></ElButton>
      </div>
    </div>
    <!-- 列表 -->
    <AgelTable class="flex-1" v-bind="table" @selection-change="handleSelectionChange"> </AgelTable>
    <!-- 弹窗表单 -->
    <ElDialog v-model="form.show" :title="form.title" width="800px" top="10vh">
      <ElForm :ref="(v) => (form.ref = v)" :model="form.model" label-width="80px">
        <AgelFormDesc 
          :items="form.items.filter(item => item.prop !== 'config' || form.state !== 'edit')" 
          :view-model="form.state === 'view'"
        ></AgelFormDesc>
      </ElForm>
      <template #footer>
        <ElButton v-if="form.state === 'add'" type="primary" @click="form.submit">提交</ElButton>
        <ElButton v-if="form.state === 'edit'" type="primary" @click="form.submitEdit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="jsx" setup>
import { reactive, ref, nextTick } from 'vue'
import http from '@/api'
import yaml from 'js-yaml'

const selectedRows = ref([]) // 用来保存选中的行

const form = reactive({
  test: false,
  ref: null,
  show: false,
  title: '',
  state: 'edit',
  model: { name: '', api_server: '', region: '', version: '', description: '', config: '' },
  items: [
    { label: '集群名称', prop: 'name', span: 3, required: true },
    { label: 'APIServer', span: 3, prop: 'api_server', required: true },
    { label: 'Region', span: 3, prop: 'region' },
    { label: '集群版本', span: 3, prop: 'version' },
    {
      label: '介绍',
      prop: 'description',
      span: 3,
      attrs: { type: 'textarea', rows: 5 }
    },
    {
      label: '集群凭证',
      prop: 'config',
      span: 3,
      attrs: { type: 'textarea', rows: 20 },
      required: true
    }
  ],
  toAdd: () => {
    form.show = true
    form.title = '新增集群'
    form.state = 'add'
    nextTick(() => {
      form.ref?.resetFields()
    })
  },
  toEdit: (row) => {
    form.show = true
    form.title = '编辑集群'
    form.state = 'edit'
    //form.config = atob(form.config)
    nextTick(() => {
      form.model = { ...row }
    })
  },
  submit: () => {
    form.ref?.validate().then(() => {
      form.show = false
      table.create(form.model)
    })
  },
  submitEdit: () => {
    form.ref?.validate().then(() => {
      form.show = false
      table.edit(form.model)
    })
  }
})

const table = reactive({
  loading: false,
  border: true,
  data: [],
  columns: [
    { label: '#', type: 'selection' },
    { label: '集群名称', prop: 'name', width: 200 },
    { label: 'APIServer', prop: 'api_server', width: 200 },
    { label: 'Region', prop: 'region', width: 150 },
    { label: '集群版本', prop: 'version', width: 150 },
    { label: '介绍', prop: 'description', showOverflowTooltip: true },
    {
      width: '200px',
      label: '操作',
      fixed: 'right',
      slot: (scope) => {
        return (
          <div>
            <el-button link type="primary" onClick={() => form.toEdit(scope.row)}>
              编辑
            </el-button>
          </div>
        )
      }
    }
  ],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res) => {
      table.data = res.result
      table.loading = false
    })
  },
  create: (form) => {
    form.config = yaml.dump(form.config)
      .replace(/^ *\|-\n/, '')    // 去掉 |-\n 后的两个空格
      .replace(/\\n {2}/g, '\n')    // 去掉 \n 后的两个空格
      .toString()
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/cluster`, form).then((res) => {
      table.loading = false
      table.request()
    })
  },  
  edit: (form) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/cluster`, form).then((res) => {
      table.loading = false
      table.request()
    })
  },
  delete: (form) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/cluster/${form.id}`).then((res) => {
      table.loading = false
      table.request()
    })
  }
})

// 获取选中的行
function handleSelectionChange(selected) {
  selectedRows.value = selected
}

// 删除选中的行
function deleteSelected() {
  if (selectedRows.value.length === 0) {
    return alert('请选择至少一行进行删除！')
  }

  selectedRows.value.forEach(row => {
    table.delete(row)  // 调用 delete 方法删除每一行
  })
}

table.request()
</script>
