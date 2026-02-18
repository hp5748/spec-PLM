import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import Layout from '@/components/layout/Layout.vue'

// 公开路由（无需登录）
export const publicRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' },
  },
  {
    path: '/403',
    name: '403',
    component: () => import('@/views/error/403.vue'),
    meta: { title: '无权限' },
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '页面不存在' },
  },
]

// 需要登录的路由
export const privateRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '首页', icon: 'HomeFilled', requiresAuth: true },
      },
    ],
  },
  {
    path: '/user',
    component: Layout,
    redirect: '/user/list',
    meta: { title: '用户管理', icon: 'User' },
    children: [
      {
        path: 'list',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: { title: '用户列表', permissions: ['user:view'] },
      },
      {
        path: 'organization',
        name: 'Organization',
        component: () => import('@/views/user/OrganizationList.vue'),
        meta: { title: '组织管理', permissions: ['user:view'] },
      },
      {
        path: 'department',
        name: 'Department',
        component: () => import('@/views/user/DepartmentList.vue'),
        meta: { title: '部门管理', permissions: ['user:view'] },
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/user/RoleList.vue'),
        meta: { title: '角色管理', permissions: ['role:view'] },
      },
    ],
  },
  {
    path: '/admin',
    component: Layout,
    redirect: '/admin/attribute-schema',
    meta: { title: '系统管理', icon: 'Setting', permissions: ['admin:all'] },
    children: [
      {
        path: 'attribute-schema',
        name: 'AttributeSchema',
        component: () => import('@/views/admin/AttributeSchemaList.vue'),
        meta: { title: '属性配置', permissions: ['admin:all'] },
      },
    ],
  },
  {
    path: '/material',
    component: Layout,
    redirect: '/material/list',
    meta: { title: '物料管理', icon: 'Box' },
    children: [
      {
        path: 'list',
        name: 'MaterialList',
        component: () => import('@/views/material/MaterialList.vue'),
        meta: { title: '物料列表', permissions: ['material:view'] },
      },
    ],
  },
  {
    path: '/document',
    component: Layout,
    redirect: '/document/list',
    meta: { title: '文档管理', icon: 'Document' },
    children: [
      {
        path: 'list',
        name: 'DocumentList',
        component: () => import('@/views/document/DocumentList.vue'),
        meta: { title: '文档列表', permissions: ['document:view'] },
      },
    ],
  },
  {
    path: '/bom',
    component: Layout,
    redirect: '/bom/list',
    meta: { title: 'BOM管理', icon: 'Share' },
    children: [
      {
        path: 'list',
        name: 'BOMList',
        component: () => import('@/views/bom/BOMList.vue'),
        meta: { title: 'BOM列表', permissions: ['bom:view'] },
      },
      {
        path: 'edit/:id?',
        name: 'BOMEdit',
        component: () => import('@/views/bom/BOMEdit.vue'),
        meta: { title: 'BOM编辑', permissions: ['bom:edit'] },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: [...publicRoutes, ...privateRoutes],
})

export default router
