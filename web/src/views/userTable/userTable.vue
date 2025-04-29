<template>
  <div class="flex flex-col overflow-auto">
    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="createForm.toAdd">新增</ElButton>
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
    <ElDialog v-model="createForm.show" :title="createForm.title" width="800px" top="10vh">
      <ElForm :ref="(v) => (createForm.ref = v)" :model="createForm.model" label-width="80px">
        <AgelFormDesc 
          :items="createForm.items.filter(item => item.prop !== 'config' || form.state !== 'edit')" 
          :view-model="createForm.state === 'view'"
        ></AgelFormDesc>
      </ElForm>
      <template #footer>
        <ElButton v-if="createForm.state === 'add'" type="primary" @click="table.create(createForm.model)">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="jsx" setup>
import { reactive, ref, nextTick } from 'vue'
import http from '@/api'

const selectedRows = ref([]) // 用来保存选中的行

const createForm = reactive({
  test: false,
  ref: null,
  show: false,
  title: '',
  state: 'add',
  model: { name: '', api_server: '', region: '', version: '', description: '', config: '' },
  items: [
    { label: '用户名称', prop: 'name', width: 200 },
    { label: 'E-Mail', prop: 'email', width: 200 },
    { label: 'Phone', prop: 'phone', width: 150 },
    { label: '密码', prop: 'password', width: 150 },
  ],
  toAdd: () => {
    createForm.show = true
    createForm.title = '新增用户'
    createForm.state = 'add'
    nextTick(() => {
      createForm.ref?.resetFields()
    })
  }
})

const table = reactive({
  loading: false,
  border: true,
  data: [],
  columns: [
    { label: '#', type: 'selection' },
    { label: '用户名称', prop: 'name', width: 200 },
    { label: 'E-Mail', prop: 'email', width: 200 },
    { label: 'Phone', prop: 'phone', width: 150 },
    { label: '类型', prop: 'type', width: 150 },
    { label: '角色', prop: 'role', width: 150 },
    { label: '最后登陆', prop: 'last_login' },
    /*{
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
    }*/
  ],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/user/fetch/user`).then((res) => {
      table.data = res.result
      table.loading = false
    })
  },
  create: (form) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/user/register`, form).then((res) => {
      table.loading = false
      table.request()
    })
  },
  delete: (form) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/user/delete/user/${form.id}`).then((res) => {
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
