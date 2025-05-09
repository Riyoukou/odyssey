<template>
  <div class="flex flex-col overflow-auto">
    <!-- 搜索表单 -->
    <ElForm v-if="search.show" :model="search.model" ref="searchRef" label-width="80px">
      <ElFormItem label="名称" prop="name">
        <ElInput v-model="search.model.name" placeholder="请输入名称" />
      </ElFormItem>
      <div class="flex justify-center mb-3">
        <ElButton icon="RefreshRight" @click="search.reset">重置</ElButton>
        <ElButton type="primary" icon="Search" @click="search.search">查询</ElButton>
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
        <ElButton icon="Refresh" round @click="table.request"></ElButton>
        <ElButton icon="Search" round @click="search.show = !search.show"></ElButton>
      </div>
    </div>

    <!-- 列表 -->
    <el-table :data="table.filteredData" style="width: 100%" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="工具名称" sortable width="200" />
      <el-table-column prop="type" label="工具类型" sortable width="200" />
      <el-table-column prop="url" label="工具地址" sortable width="200" />
      <el-table-column prop="credential_type" label="凭证类型" width="150" />
      <el-table-column prop="credential_name" label="凭证名称" width="150" />
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
        :rules="rules"
      >
        <ElFormItem label="工具名称" prop="name">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.name" placeholder="请输入工具名称" />
        </ElFormItem>
        <ElFormItem label="工具类型" prop="type">
          <ElSelect v-if="editForm.state !== 'view'" v-model="editForm.model.type" placeholder="请选择工具类型">
            <ElOption label="Git" value="git" />
            <ElOption label="Harbor" value="harbor" />
            <ElOption label="Jenkins" value="jenkins" />
            <ElOption label="ArgoCD" value="argocd" />
            <ElOption label="其他" value="other" />
          </ElSelect>
          <ElInput disabled v-else v-model="editForm.model.type" />
        </ElFormItem>
        <ElFormItem label="工具地址" prop="url">
          <ElInput :disabled="editForm.state === 'view'" v-model="editForm.model.url" placeholder="请输入工具地址" />
        </ElFormItem>
        <ElFormItem label="凭证类型" prop="credential_type">
          <ElSelect v-if="editForm.state !== 'view'" v-model="editForm.model.credential_type" placeholder="请选择工具凭证类型">
            <ElOption label="无" value="none" />
            <ElOption label="Token" value="token" />
            <ElOption label="KubeConfig" value="kube_config" />
            <ElOption label="用户名密码" value="basic" />
          </ElSelect>
          <ElInput disabled v-else v-model="editForm.model.credential_type"/>
        </ElFormItem>
        <ElFormItem label="凭证名称" prop="credential_name">
          <ElSelect v-if="editForm.state !== 'view'" v-model="editForm.model.credential_name" placeholder="请选择工具凭证" @focus="table.selectConfig">
            <ElOption  v-for="item in table.credentialData" :key="item.name" :label="item.type+'/'+item.name" :value="item.name" />
          </ElSelect>
          <ElInput disabled v-else v-model="editForm.model.credential_name" />
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
const searchRef = ref<FormInstance | null>(null);
const search = reactive({
  show: false,
  model: {
    name: '',
  },
  search: () => {
    table.filteredData = table.data.filter(
      (data: any) =>
        data.name?.toLowerCase().includes(search.model.name.toLowerCase())
    );
  },
  reset: () => {
    searchRef.value?.resetFields();
  },
})

const rules = computed(() => {
  return {
    name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
    type: [{ required: true, message: '请选择工具类型', trigger: 'blur' }],
    url: [{ required: true, message: '请输入工具地址', trigger: 'blur' }],
    credential_type: [{ required: true, message: '请选择凭证类型', trigger: 'blur' }],
    credential_name: [{ required: true, message: '请选择凭证名称', trigger: 'blur' }],
    description: [],
  }
})

// 表单配置
const editForm = reactive({
  ref: null as FormInstance | null,
  show: false,
  title: '',
  state: '',
  model: {
    name: '',
    type: '',
    url: '',
    credential_type: '',
    credential_name: '',
    description: '',
  } ,
  toAdd: () => {
    editForm.ref?.resetFields()
    editForm.show = true
    editForm.title = '新增集群'
    editForm.state = 'add'
    editForm.model = {
      name: '',
      type: '',
      url: '',
      credential_type: '',
      credential_name: '',
      description: ''
    }
  },
  toEdit: (row: any) => {
    editForm.show = true
    editForm.title = '编辑集群'
    editForm.state = 'edit'
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
      table.create(editForm.model)
    })
  },
  editSubmit: () => {
    editForm.ref?.validate().then(() => {
      editForm.show = false
      table.edit(editForm.model)
    })
  }
})
// 表格配置
const table = reactive({
  loading: false,
  border: true,
  data: [] as any[],
  filteredData: [] as any[],
  credentialData:[] as any[],
  request: () => {
    table.loading = true
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/cicd_tool`).then((res: any) => {
      table.data = res.data
      table.filteredData = res.data
      table.loading = false
    })
  },
  create: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/create/cicd_tool`, form).then((res: any) => {
      table.loading = false
      ElMessage.success('新增成功')
      table.request()
    })
  },
  edit: (form: any) => {
    table.loading = true
    http.post(import.meta.env.VITE_APP_BASE_URL + `/cicd/update/cicd_tool`, form).then((res: any) => {
      table.loading = false
      ElMessage.success('编辑成功')
    })
  },
  delete: (form: any) => {
    table.loading = true
    http.delete(import.meta.env.VITE_APP_BASE_URL + `/cicd/delete/cicd_tool/${form.id}`).then((res: any) => {
      table.loading = false
      table.request()
      ElMessage.success('删除成功')
    })
  },
  selectConfig: () => {
    http.get(import.meta.env.VITE_APP_BASE_URL + `/cicd/fetch/credential`).then((res: any) => {
      table.credentialData = res.data.filter((item: any) => item.type === editForm.model.credential_type)
      table.loading = false
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
    table.delete(row)
  })
}

// 初始化加载数据
table.request()
</script>
