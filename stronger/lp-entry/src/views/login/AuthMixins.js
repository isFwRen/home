import { mapGetters } from 'vuex'
import { sessionStorage } from 'vue-rocket'

const hasAuth = ({ hasOp0, hasOp1, hasOp2, hasOpq }) => {
  if(hasOp0 || hasOp1 || hasOp2 || hasOpq) {
    return true
  }

  return false
}

export default {
  methods: {
    async getRoleSysMenu() {
      const result = await this.$store.dispatch('GET_ROLE_SYS_MENU')

      if(result.code === 200) {
        const { Menus, Perm, proCode = '' } = result.data
        const [mapPro, proItems] = this.resolvePerm(Perm, proCode)

        sessionStorage.set('proCode', proCode)
        this.$store.commit('UPDATE_AUTH', {
          menus: Menus,
          perm: Perm,
          mapPro,
          proItems
        })
      }
    },

    resolvePerm(permissions, proCode) {
      const mapPro = {}
      const proItems = []

      permissions.map(permission => {
        if(permission.proCode === proCode && hasAuth(permission)) {
          mapPro[permission.proCode] = permission

          proItems.push({ 
            label: permission.proCode, 
            value: permission.proCode 
          })
        }
      })

      return [mapPro, proItems]
    }
  },

  computed: {
    ...mapGetters(['auth'])
  }
}