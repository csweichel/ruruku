package types

type Permission string

const (
	PermissionNone              Permission = "none"
	PermissionUserAdd           Permission = "user.add"
	PermissionUserDelete        Permission = "user.delete"
	PermissionUserGrant         Permission = "user.grant"
	PermissionUserChpwd         Permission = "user.chpwd"
	PermissionSessionStart      Permission = "session.start"
	PermissionSessionClose      Permission = "session.close"
	PermissionSessionView       Permission = "session.view"
	PermissionSessionContribute Permission = "session.contribute"
)
