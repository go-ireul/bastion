<template>
  <b-navbar toggleable="md" type="dark" variant="info">
    <b-nav-toggle target="nav_collapse"></b-nav-toggle>
    <router-link class="navbar-brand" :to="{ name: 'dashboard' }">Bastion</router-link>
    <b-collapse is-nav id="nav_collapse">
      <b-nav is-nav-bar>
        <b-nav-item :to="{name:'dashboard'}">工作台</b-nav-item>
      </b-nav>

      <b-nav is-nav-bar class="ml-auto">
        <b-nav-item right v-if="currentUser.isAdmin" :to="{name:'admin'}">系统管理</b-nav-item>
        <b-nav-item-dropdown right>
          <template slot="button-content">
            <em>{{ currentUser.nickname }}</em>
          </template>
          <b-dropdown-item :to="{name:'settings'}">设置</b-dropdown-item>
          <b-dropdown-divider></b-dropdown-divider>
          <b-dropdown-item @click="signout">退出登录</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-nav>
    </b-collapse>
  </b-navbar>
</template>

<script>
export default {
  name: 'navbar',
  computed: {
    currentUser () {
      return this.$store.state.user.currentUser
    }
  },
  data () {
    return {}
  },
  methods: {
    signout () {
      this.$api.destroyToken({ id: 'current' }).then(() => this.doSignout(), () => this.doSignout())
    },
    doSignout () {
      this.$store.commit('setCurrentToken', null)
      this.$store.commit('setCurrentUser', null)
      this.$router.push({ name: 'index' })
    }
  }
}
</script>
