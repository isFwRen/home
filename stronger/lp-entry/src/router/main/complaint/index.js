const ComplaintRoutes = {
  path: 'complaint',
  name: 'Complaint',
  meta: {
    key: 'complaint',
    realm: 'complaint',
    title: '客户投诉',
    isRequired: true
  },
  component: () => import('@/views/main/complaint')
}

export default ComplaintRoutes