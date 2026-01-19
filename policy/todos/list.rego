package todos

import rego.v1

# Fetch data from an external API
workspaces := res.body if {
	input.user.id == input.resource.owner_id
	input.resource.type == "todo"
	"user" in input.user.roles
	input.action == "list"

	res := http.send({
		"method": "GET",
		"url": sprintf("http://localhost:10012/workspaces?user_id=%s", [input.user.id]),
		"force_json_decode": true,
	})

	res.status_code == 200
}

decision := {
	"workspaces": workspaces,
}
