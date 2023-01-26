package repo

type PermissionStorageI interface {
	CheckPermission(userType, resource, action string) (bool, error)
}
