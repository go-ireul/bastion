export default {
  fetchUser: function ({id}) {
    return this.http.get(`users/${id}`)
  },
  updateUser: function ({id, nickname}) {
    return this.http.post(`users/${id}/update`, {nickname})
  }
}
