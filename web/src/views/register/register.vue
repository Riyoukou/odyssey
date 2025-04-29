<template>
  <div class="login-page h-full flex-center bg-gray-100 dark:bg-[var(--el-bg-color-page)]">
    <div
      class="border border-gray-200 w-96 p-10 shadow-blue-200 shadow-md bg-white rounded-md dark:bg-[var(--el-bg-color)] dark:shadow-blue-700 dark:border-blue-900"
    >
      <div class="text-center text-2xl text-p mb-2">Register</div>
      <div class="text-gray-400 text-sm text-center"> Please enter your name and password to register.</div>
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
        <ElFormItem prop="email">
          <ElInput v-model="form.data.email" placeholder="Please input email" @keydown.enter="submit"></ElInput>
        </ElFormItem>
        <ElFormItem prop="phone">
          <ElInput v-model="form.data.phone" placeholder="Please input phone" @keydown.enter="submit"></ElInput>
        </ElFormItem>
      </ElForm>
      <div class="flex justify-center items-center mt-6">
        <ElButton :loading="form.loading" class="w-full" type="primary" @click="submit">Login</ElButton>
      </div>
      <div class="flex justify-center items-center mt-4">
        <span class="text-blue-500 cursor-pointer" @click="goBack">
          返回
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import useUserStore from '@/stores/useUserStore'
import { useRouter } from 'vue-router'
import type { FormInstance } from 'element-plus'

const form = reactive({
  loading: false,
  data: {
    name: '',
    password: '',
    email: '',
    phone: ''
  },
  rules: {
    name: { required: true, message: 'Please input Activity name', trigger: 'blur' },
    password: { required: true, message: 'Please input Activity password', trigger: 'blur' },
    email: [
      { required: true, message: 'Please input email', trigger: 'blur' },
      { type: 'email', message: 'Please input valid email', trigger: ['blur'] }
    ],
    phone: [
      { required: true, message: 'Please input phone', trigger: 'blur' },
      {
        validator: (rule: any, value: any, callback: any) => {
          const phoneRegex = /^\d{11}$/
          if (!phoneRegex.test(value)) {
            callback(new Error('Please input valid 11-digit phone number'))
          } else {
            callback()
          }
        },
        trigger: ['blur']
      }  
    ]
  }
})

const elFormRef = ref<FormInstance>()
const userStore = useUserStore()
const router = useRouter()
async function submit() {
  await elFormRef.value?.validate()
  form.loading = true
  userStore
    .registerApp(form.data)
    .then(() => {
      router.push({ name: 'login' })
    })
    .catch(() => {
      form.loading = false
    })
}

function goBack() {
  router.push({ name: 'login' })
}
</script>
