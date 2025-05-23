package model

var roleTypes = []string{
	"internal", "external", "clinic",
}

func GetRoleTypes() []string {
	return roleTypes
}
