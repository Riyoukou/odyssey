<template>
  <div>
    <ElButton type="primary" @click="loginAdmin">切换到管理员</ElButton>
    <ElButton type="primary" @click="loginUser">切换到用户</ElButton>
  </div>
</template>

<script setup lang="ts">
import useUserStore from '@/stores/useUserStore'
import { ElMessage } from 'element-plus'
const userStore = useUserStore()

function loginAdmin(): void {
  userStore
    .loginApp({
      name: 'admin',
      password: '123'
    } as any)
    .then((): void => {
      ElMessage.success('切换管理员角色成功')
    })
}

function loginUser(): void {
  userStore
    .loginApp({
      name: 'user',
      password: '123'
    } as any)
    .then((): void => {
      ElMessage.success('切换用户角色成功')
    })
}
</script>

<style scoped></style>
