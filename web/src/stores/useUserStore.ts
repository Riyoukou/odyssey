import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'
import { dynamicRoutes } from "@/router/routes"
import http from "@/api"


export default defineStore('userStore', () => {
  const userInfo = useStorage('userInfo', getInitUserInfo(), sessionStorage)

  const token = computed(() => {
    return userInfo.value.token
  })

  // 身份权限
  const role = computed(() => {
    return userInfo.value.role
  })

  // 菜单
  const menuData = computed(() => {
    return getMenuData(dynamicRoutes)
  })

  function getInitUserInfo() {
    return {
      name: "",
      id: "",
      token: "",
      role: "",
    }
  }

  function loginApp(data: any) {
    return http({
      url: import.meta.env.VITE_APP_BASE_URL + `/user/login`,
      method: 'post',
      data: data
    }).then((res: any) => {
      userInfo.value = res.data
      return res
    })
  }

  function registerApp(data: any) {
    return http({
      url: import.meta.env.VITE_APP_BASE_URL + `/user/register`,
      method: 'post',
      data: data
    }).then((res: any) => {
      return res
    })
  }

  function logoutApp() {
    return new Promise((reslove => {
      userInfo.value = getInitUserInfo()
      reslove(true)
    }))
  }

  interface MenuItem {
    children: MenuItem[];
    title: string;
    icon?: string;
    index: string;
    hidden: boolean;
  }

  interface RouteItem {
    path: string;
    meta: {
      hidden?: boolean;
      roles?: string[];
      title: string;
      icon?: string;
    };
    children?: RouteItem[];
  }

  function getMenuData(list: RouteItem[]): MenuItem[] {
    return list.map(item => {
      const isShow = item.meta.hidden !== true
      const hasRole = item.meta.roles ? item.meta.roles.includes(role.value) : true
      const menuItem: MenuItem = {
        children: [],
        title: item.meta.title,
        icon: item.meta.icon,
        index: item.path,
        hidden: !(isShow && hasRole)
      }
      if (item.children && item.children.length > 0) {
        menuItem.children = getMenuData(item.children)
      }
      return menuItem
    })
  }

  return {
    // 不能修改
    userInfo: userInfo,
    menuData,
    token,
    role,
    loginApp,
    registerApp,
    logoutApp
  }
})
