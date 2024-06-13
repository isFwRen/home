import HomeRoutes from './home'
import RuleRoutes from './rule'
import EntryRoutes from './entry'
import PracticeRoutes from './practice'
import ExamRoutes from './exam'
import ErrorRoutes from './error'
import ComplaintRoutes from './complaint'
import SalaryRoutes from './salary'
import YieldRoutes from './yield'
import TrainRoutes from './train'

const MainRoutes = {
  path: '/main',
  name: 'Main',
  meta: {
    key: 'main'
  },
  component: () => import('@/views/main'),

  children: [
    {
      path: '/main',
      redirect: 'home'
    },

    HomeRoutes,
    RuleRoutes,
    EntryRoutes,
    PracticeRoutes,
    ExamRoutes,
    ErrorRoutes,
    ComplaintRoutes,
    SalaryRoutes,
    YieldRoutes,
    TrainRoutes
  ]
}

export default MainRoutes