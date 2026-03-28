import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import ServiceView from '../views/ServiceView.vue';
import CultureView from '../views/CultureView.vue';
import ContactView from '../views/ContactView.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/service', name: 'service', component: ServiceView },
    { path: '/culture', name: 'culture', component: CultureView },
    { path: '/contact', name: 'contact', component: ContactView },
  ],
});

export default router;
