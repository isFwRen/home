const HomeRoutes = {
  path: 'home',
  name: 'Home',
  meta: {
    key: 'home',
    realm: 'home',
    title: '主页',
    isRequired: true
  },
  component: () => import('@/views/main/home')
}

export default HomeRoutes