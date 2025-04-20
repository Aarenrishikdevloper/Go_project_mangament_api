package constants

import (
	"rishik.com/enums"
)

var RolePermissions = map[enums.RoleType][]enums.PermissionType{
	enums.Owner: {
		enums.CreateWorkspace,
		enums.EditWorkspace,
		enums.DeleteWorkspace,
		enums.ManageWorkspaceSettings,
		enums.AddMember,
		enums.ChangeMemberRole,
		enums.RemoveMember,
		enums.CreateProject,
		enums.EditProject,
		enums.DeleteProject,
		enums.CreateTask,
		enums.EditTask,
		enums.DeleteTask,
		enums.ViewOnly,
	},
	enums.Admin: {
		enums.AddMember,
		enums.CreateProject,
		enums.EditProject,
		enums.DeleteProject,
		enums.CreateTask,
		enums.EditTask,
		enums.DeleteTask,
		enums.ManageWorkspaceSettings,
		enums.ViewOnly,
	},
	enums.Member: {
		enums.ViewOnly,
		enums.CreateTask,
		enums.EditTask,
	},
}
