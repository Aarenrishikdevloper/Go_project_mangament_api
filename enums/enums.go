package enums

type ProviderEnum string

const (
	Google   ProviderEnum = "GOOGLE"
	Github   ProviderEnum = "GITHUB"
	Facebook ProviderEnum = "FACEBOOK"
	Email    ProviderEnum = "EMAIL"
)

type RoleType string

const (
	Owner  RoleType = "OWNER"
	Admin  RoleType = "ADMIN"
	Member RoleType = "MEMBER"
)

type PermissionType string

const (
	CreateWorkspace         PermissionType = "CREATE_WORKSPACE"
	DeleteWorkspace         PermissionType = "DELETE_WORKSPACE"
	EditWorkspace           PermissionType = "EDIT_WORKSPACE"
	ManageWorkspaceSettings PermissionType = "MANAGE_WORKSPACE_SETTINGS"
	AddMember               PermissionType = "ADD_MEMBER"
	ChangeMemberRole        PermissionType = "CHANGE_MEMBER_ROLE"
	RemoveMember            PermissionType = "REMOVE_MEMBER"
	CreateProject           PermissionType = "CREATE_PROJECT"
	EditProject             PermissionType = "EDIT_PROJECT"
	DeleteProject           PermissionType = "DELETE_PROJECT"
	CreateTask              PermissionType = "CREATE_TASK"
	EditTask                PermissionType = "EDIT_TASK"
	DeleteTask              PermissionType = "DELETE_TASK"
	ViewOnly                PermissionType = "VIEW_ONLY"
)

type TaskSTatusEnum string

const (
	Backlog    TaskSTatusEnum = "BACKLOG"
	Todo       TaskSTatusEnum = "TODO"
	inProgress TaskSTatusEnum = "IN_PROGRESS"
	INReview   TaskSTatusEnum = "IN_REVIEW"
	Done       TaskSTatusEnum = "DONE"
)

type TaskPriorityEnum string

const (
	Low    TaskPriorityEnum = "LOW"
	Medium TaskPriorityEnum = "MEDIUM"
	High   TaskPriorityEnum = "HIGH"
)
