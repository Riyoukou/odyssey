<template>
  <div class="login-page h-full flex-center bg-gray-100 dark:bg-[var(--el-bg-color-page)]">
    <div
      class="border border-gray-200 w-96 p-10 shadow-blue-200 shadow-md bg-white rounded-md dark:bg-[var(--el-bg-color)] dark:shadow-blue-700 dark:border-blue-900"
    >
      <div class="text-center text-2xl text-p mb-2">odyssey</div>
      <div class="text-gray-400 text-sm text-center">Please input name and password to login</div>
      <ElForm ref="elFormRef" :model="form.data" :rules="form.rules" class="mt-8" label-position="top">
        <ElFormItem prop="name">
          <ElInput v-model="form.data.name" placeholder="Please input name" @keydown.enter="submit"></ElInput>
        </ElFormItem>
        <ElFormItem prop="password">
          <ElInput
            v-model="form.data.password"
            placeholder="Please input password"
            type="password"
            @keydown.enter="submit"
          ></ElInput>
        </ElFormItem>
      </ElForm>
      <div class="flex justify-center items-center mt-6">
        <ElButton :loading="form.loading" class="w-full" type="primary" @click="submit">Login</ElButton>
        <ElButton :loading="form.loading" class="w-full" type="primary" @click="pushRegister">Register</ElButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import useUserStore from '@/stores/useUserStore'
import { homeRoute } from '@/router/routes'
import { useRouter } from 'vue-router'
import type { FormInstance } from 'element-plus'

const form = reactive({
  loading: false,
  data: {
    name: '',
    password: ''
  },
  rules: {
    name: { required: true, message: 'Please input Activity name', trigger: 'blur' },
    password: { required: true, message: 'Please input Activity password', trigger: 'blur' }
  }
})
const elFormRef = ref<FormInstance>()
const userStore = useUserStore()
const router = useRouter()
async function submit() {
  await elFormRef.value?.validate()
  form.loading = true
  userStore
    .loginApp(form.data)
    .then(() => {
      router.push({ path: homeRoute.path })
    })
    .catch(() => {
      form.loading = false
    })
}

function pushRegister() {
  router.push({ name: 'register' })
}
</script>
