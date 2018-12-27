package types

type Permission string

const (
	PermissionNone              Permission = "none"
	PermissionUserAdd           Permission = "user.add"
	PermissionUserDelete        Permission = "user.delete"
	PermissionUserGrant         Permission = "user.grant"
	PermissionUserChpwd         Permission = "user.chpwd"
	PermissionUserList          Permission = "user.list"
	PermissionSessionStart      Permission = "session.start"
	PermissionSessionModify     Permission = "session.modify"
	PermissionSessionClose      Permission = "session.close"
	PermissionSessionView       Permission = "session.view"
	PermissionSessionContribute Permission = "session.contribute"
)

var AllPermissions = []Permission{
	PermissionNone,
	PermissionUserAdd,
	PermissionUserDelete,
	PermissionUserGrant,
	PermissionUserChpwd,
	PermissionUserList,
	PermissionSessionStart,
	PermissionSessionModify,
	PermissionSessionClose,
	PermissionSessionView,
	PermissionSessionContribute,
}
