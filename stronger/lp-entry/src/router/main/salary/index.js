const SalaryRoutes = {
  path: 'salary',
  name: 'Salary',
  meta: {
    key: 'salary',
    realm: 'salary',
    title: '我的工资',
    isRequired: true
  },
  component: () => import('@/views/main/salary')
}

export default SalaryRoutes