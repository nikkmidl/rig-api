package github.access

deny if {
  input.user == input.blocked_users[_]
}
