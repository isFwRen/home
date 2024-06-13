const LoginRoutes = {
  path: '/login',
  name: 'Login',
  meta: {
    title: '登录',
    key: 'login',
    path: 'login'
  },
  component: () => import('@/views/login')
}

export default LoginRoutes