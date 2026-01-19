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

# Fetch data from an external API
workspaces := res.body if {
	allow

	res := http.send({
		"method": "GET",
		"url": sprintf("http://localhost:10012/workspaces?user_id=%s", [input.user.id]),
		"force_json_decode": true,
	})

	res.status_code == 200
}

decision := {
	"allow": allow,
	"workspaces": workspaces,
}
