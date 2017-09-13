export default {
  listSSHKeys ({userId}) {
    return this.http.get(`users/${userId}/ssh_keys`)
  },
  createSSHKey ({userId, name, publicKey}) {
    return this.http.post(`users/${userId}/ssh_keys/create`, {name, publicKey})
  },
  destroySSHKey ({id}) {
    return this.http.post(`ssh_keys/${id}/destroy`)
  }
}
