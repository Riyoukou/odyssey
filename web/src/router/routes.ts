import Home from '@/views/home/home.vue' // 避免生成 _plugin-vue_export-helper

export const homeRoute = {
  path: '/home',
  name: 'home',
  meta: {
    title: '首页',
    icon: 'HomeFilled',
    hidden: true
  },
  component: Home
}

export const dynamicRoutes = [
  homeRoute,
  {
    path: '/baseTable',
    name: 'baseTable',
    meta: {
      hidden: true,
      title: '基础列表',
      icon: 'Menu'
    },
    component: () => import('@/views/table/baseTable.vue')
  },
  {
    path: '/cicd',
    name: 'CICD',
    meta: {
      hidden: false,
      title: 'CICD',
      icon: 'Menu'
    },
    children: [
      {
        path: '/cicd/cicdOption',
        name: 'cicdOption',
        meta: {
          title: '资源设置',
          icon: 'Menu',
          roles: ['admin','user']
        },
        component: () => import('@/views/cicd/cicdOption/main.vue')
      },
      {
        path: '/cicd/serviceTable',
        name: 'serviceTable',
        meta: {
          title: '服务列表',
          icon: 'Menu',
          roles: ['admin','user']
        },
        component: () => import('@/views/cicd/serviceTable/serviceTable.vue'),
      },
      {
        path: '/cicd/buildRecord',
        name: 'buildRecord',
        meta: {
          title: '构建记录',
          icon: 'Menu',
          roles: ['admin','user']
        },
        component: () => import('@/views/cicd/buildRecord/buildRecordTable.vue'),
      },
      {
        path: '/cicd/deployRecord',
        name: 'deployRecord',
        meta: {
          title: '发布记录',
          icon: 'Menu',
          roles: ['admin','user']
        },
        component: () => import('@/views/cicd/deployRecord/deployRecordTable.vue'),
      },
      {
        path: '/cicd/cicd',
        name: 'cicd',
        meta: {
          hidden: true,
          title: 'CICD',
          icon: 'Menu',
          roles: ['admin','user']
        },
        component: () => import('@/views/cicd/cicdDetail.vue')
      },
    ]
  },
  {
    path: '/userTable',
    name: 'userTable',
    meta: {
      title: '用户列表',
      icon: 'Menu',
      roles: ['admin','user']
    },
    component: () => import('@/views/userTable/userTable.vue')
  },
  {
    path: '/cacheView',
    name: 'cacheView',
    meta: {
      hidden: true,
      title: '缓存页面',
      icon: 'Discount',
      keepAlive: true
    },
    component: () => import('@/views/cacheView/cacheView.vue')
  },
  {
    path: '/roles',
    name: 'roles',
    meta: {
      hidden: true,
      title: '权限相关',
      icon: 'UserFilled'
    },
    children: [
      {
        path: '/roles/switch',
        name: 'switch-roles',
        meta: {
          title: '权限切换',
          icon: 'UserFilled'
        },
        component: () => import('@/views/roles/switch-roles.vue')
      },
      {
        path: '/roles/admin',
        name: 'admin',
        meta: {
          title: '管理员角色可见',
          icon: 'UserFilled',
          roles: ['admin']
        },
        component: () => import('@/views/roles/admin.vue')
      },
      {
        path: '/roles/user',
        name: 'user',
        meta: {
          title: '用户角色可见',
          icon: 'UserFilled',
          roles: ['user']
        },
        component: () => import('@/views/roles/user.vue')
      }
    ]
  },
  {
    path: '/link',
    name: 'link',
    redirect: '/link/element-plus',
    meta: {
      title: '外部链接',
      icon: 'Connection'
    },
    children: [
      {
        path: '/link/vuejs',
        name: 'vuejs',
        meta: {
          title: 'vuejs',
          icon: 'ElementPlus',
          link: 'https://cn.vuejs.org/'
        },
        component: () => import('@/layouts/IframeLayout/IframeLayout.vue')
      },
      {
        path: '/link/element-plus',
        name: 'element-plus',
        meta: {
          title: 'element-plus',
          icon: 'ElementPlus',
          link: 'https://element-plus.org/zh-CN/'
        },
        component: () => import('@/layouts/IframeLayout/IframeLayout.vue')
      }
    ]
  },
  {
    path: '/menu',
    name: 'menu',
    redirect: '/menu/menu-1',
    meta: {
      hidden: true,
      title: '多级菜单',
      icon: 'Grid'
    },
    children: [
      {
        path: '/menu/menu-1',
        name: 'menu-1',
        meta: {
          title: '菜单1',
          icon: 'Menu'
        },
        component: () => import('@/views/menu/menu.vue')
      },
      {
        path: '/menu/menu-2',
        name: 'menu-2',
        meta: {
          title: '菜单2',
          icon: 'Menu'
        },
        component: () => import('@/views/menu/menu.vue')
      },
      {
        path: '/menu/menu-3',
        name: 'menu-3',
        redirect: '/menu/menu-3-1',
        meta: {
          title: '菜单3',
          icon: 'Menu'
        },
        children: [
          {
            path: '/menu/menu-3-1',
            name: 'menu-3-1',
            meta: {
              title: '菜单3-1',
              icon: 'Menu'
            },
            component: () => import('@/views/menu/menu.vue')
          },
          {
            path: '/menu/menu-3-2',
            name: 'menu-3-2',
            meta: {
              title: '菜单3-2',
              icon: 'Menu'
            },
            component: () => import('@/views/menu/menu.vue')
          }
        ]
      }
    ]
  }
]

export const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/layouts/ErroLayout/404.vue')
  },
  {
    path: '/401',
    name: '401',
    component: () => import('@/layouts/ErroLayout/401.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/login.vue')
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/register/register.vue')
  },
  {
    path: '/index',
    redirect: homeRoute.path,
    component: () => import('@/layouts/IndexLayout/dynamicLayout.vue'),
    children: dynamicRoutes
  }
]

export default routes
