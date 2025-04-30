<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="search.ref" label-width="80px">
      <ElFormItem label="名称" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入名称" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="() => search.ref?.resetFields()">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="clusterTable.request">查询</ElButton>
      </div>
    </ElForm>

    <!-- 按钮条 -->
    <div class="table-bar flex justify-between items-center mb-3">
      <div>
        <ElButton icon="Plus" @click="editForm.toAdd">新增</ElButton>
        <el-popconfirm title="Are you sure to delete this?" @confirm="deleteSelected">
          <template #reference>
            <ElButton icon="Delete">删除</ElButton>
          </template>
        </el-popconfirm>
      </div>
      <div>
        <ElButton icon="Refresh" round @click="clusterTable.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table :data="clusterTable.data" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="集群名称" sortable width="200" />
      <el-table-column prop="api_server" label="APIServer" sortable width="200" />
      <el-table-column prop="region" label="Region" width="150" />
      <el-table-column prop="version" label="集群版本" width="150" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column label="操作" fixed="right" width="120">
        <template #default="{ row }">
          <el-button link type="primary" @click="editForm.toView(row)">查看</el-button>
          <el-divider direction="vertical" />
          <el-button link type="primary" @click="editForm.toEdit(row)">编辑</el-button>
        </template>
      </el-table-column>
    </el-table>

    <ElDialog v-model="editForm.show" :title="editForm.title" width="800px" top="10vh">
      <ElForm 
        :ref="(v: FormInstance | null) => (editForm.ref = v)" 
        :model="editForm.model" 
        label-width="80px"
        :rules="clusterRules"
      >
        <ElFormItem label="集群名称" prop="name">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入集群名称" />
        </ElFormItem>
        <ElFormItem label="APIServer" prop="api_server" >
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.api_server" placeholder="请输入APIServer" />
        </ElFormItem>
        <ElFormItem label="Region" prop="region">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.region" placeholder="请输入Region" />
        </ElFormItem>
        <ElFormItem label="集群版本" prop="version">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.version" placeholder="请输入集群版本" />
        </ElFormItem>
        <ElFormItem label="集群凭证" prop="config">
          <ElSelect v-if="editForm.state !== 'view'" v-model="editForm.model.config" placeholder="请选择集群凭证">
            <ElOption  v-for="item in clusterTable.credentialData" :key="item.name" :label="item.name" :value="item.name" />
          </ElSelect>
          <ElInput disabled v-else v-model="editForm.model.config" placeholder="请输入集群凭证" />
        </ElFormItem>
        <ElFormItem label="描述" prop="description">
          <ElInput
            v-model="editForm.model.description"
            type="textarea"
            :rows="5"
            placeholder="请输入描述"
            :disabled="editForm.state === 'view'"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton v-if="editForm.state === 'add'" type="primary" @click="editForm.submit">提交</ElButton>
        <ElButton v-else-if="editForm.state === 'edit'" type="primary" @click="editForm.editSubmit">提交</ElButton>
      </template>
    </ElDialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import http from '@/api'

// 搜索表单配置
const search = reactive({
  ref: null as FormInstance | null,
  show: false,
  model: {
    name: '',
    api_server: '',
    region: '',
    version: '',
    description: '',
    config: ''
  } 
})

const clusterRules = computed(() => {
  const baseRules = {
    name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
    api_server: [{ required: true, message: '请输入APIServer', trigger: 'blur' }],
    region: [],
    version: [],
    description: [],
  }

  // 仅在添加时校验 config（集群凭证）
  if (editForm.state === 'add') {
    return {
      ...baseRules,
      config: [{ required: true, message: '请输入集群凭证', trigger: 'blur' }]
    }
  }

  return baseRules
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: 'edit',
  model: {
    name: '',
    api_server: '',
    region: '',
    version: '',
    description: '',
    config: ''
  } ,
  toAdd: () => {
    editForm.show = true
    editForm.title = '新增集群'
    editForm.state = 'add'
    clusterTable.selectConfig()
    editForm.ref?.resetFields()
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑集群'
    editForm.state = 'edit'
    clusterTable.selectConfig()
    editForm.model = { ...row }
  },
  toView: (row: any) => {
    editForm.show = true
    editForm.title = '查看集群'
    editForm.state = 'view'
    editForm.model = { ...row }
  },
  submit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      clusterTable.create(editForm.model)
      editForm.ref?.resetFields()
    })
  },
  editSubmit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      clusterTable.edit(editForm.model)
      editForm.ref?.resetFields()
    })
  }
})
// 表格配置
const clusterTable = reactive({
  loading: false,
  border: true,
  data: [],
  credentialData:[],
  request: () => {
    clusterTable.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cluster`).then((res: any) => {
      clusterTable.data = res.data
      clusterTable.loading = false
    })
  },
  create: (form: any) => {
    clusterTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/cluster`, form).then((res: any) => {
      clusterTable.loading = false
      ElMessage.success('新增成功')
      clusterTable.request()
    })
  },
  edit: (form: any) => {
    clusterTable.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/cluster`, form).then((res: any) => {
      clusterTable.loading = false
      ElMessage.success('编辑成功')
      clusterTable.request()
    })
  },
  delete: (form: any) => {
    clusterTable.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/cluster/${form.id}`).then((res: any) => {
      clusterTable.loading = false
      clusterTable.request()
      if (res.status === 200){
        ElMessage.success('删除成功')
      }else{
        ElMessage.error('删除失败:'+ res.message)
      }
    })
  },
  selectConfig: () => {
    clusterTable.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/credential`).then((res: any) => {
      clusterTable.credentialData = res.data
      clusterTable.loading = false
    })
  }
})

const selectedRows = ref<any[]>([])

const handleSelectionChange = (selected: any) => {
  selectedRows.value = selected
}

const deleteSelected = () => {
  if (selectedRows.value.length === 0) {
    return alert('请选择至少一行进行删除！')
  }

  selectedRows.value.forEach(row => {
    clusterTable.delete(row)
  })
}

// 初始化加载数据
clusterTable.request()
</script>
