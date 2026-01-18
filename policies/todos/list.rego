package todos

import rego.v1

default allow := false

# Admins can list all todos
allow if {
	"admin" in input.user.roles
}

# Users can list their own todos
allow if {
	input.user.id == input.resource.owner_id
	input.resource.type == "todo"
	"user" in input.user.roles
	input.action == "list"
}
