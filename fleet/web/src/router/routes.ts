import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      {
        path: 'devices',  children: [
          { path: '', component: () => import('pages/DevicesPage.vue') },
          { path: ':id', component: () => import('pages/DevicePage.vue'), children: [
            { path: 'events', component: () => import('pages/DevicePage.vue') },
            { path: 'configuration', component: () => import('pages/DevicePage.vue') },
            { path: 'changes', component: () => import('pages/DevicePage.vue') },
            { path: 'stanza', component: () => import('pages/DevicePage.vue') },
          ]},
        ]
      },
      {
        path: 'admin',
        children: [
          { path: '', component: () => import('pages/AdminPage.vue') },
          { path: 'notifications', component: () => import('pages/NotificationsPage.vue') }
        ]
      },

    ],
  },


  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;
