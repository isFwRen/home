const TrainRoutes = {
  path: 'train',
  name: 'Train',
  meta: {
    key: 'train',
    realm: 'train',
    title: '培训流程',
    isRequired: true
  },
  component: () => import('@/views/main/train')
}

export default TrainRoutes