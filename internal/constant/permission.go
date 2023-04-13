package constant

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionRead   ActionType = "READ"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"

	PermissionFullAccess = string("FULL_ACCESS") // Full access to all object

	PermissionProductAll         = string("PRODUCT_ALL")          // Full access to product that they own
	PermissionProductCreate      = string("PRODUCT_CREATE")       // Only access to create product
	PermissionProductUpdate      = string("PRODUCT_UPDATE")       // Only access to update product
	PermissionProductDelete      = string("PRODUCT_DELETE")       // Only access to delete product
	PermissionProductRead        = string("PRODUCT_READ")         // Only access to read product
	PermissionProductReadOther   = string("PRODUCT_READ_OTHER")   // Only access to read other user product
	PermissionProductReadDeleted = string("PRODUCT_READ_DELETED") // only access to read deleted product
	PermissionProductModifyOther = string("PRODUCT_MODIFY_OTHER") // Only access to update/delete other user product
)

var (
	SeedPermissions = []string{
		PermissionFullAccess,
		PermissionProductAll,
		PermissionProductRead,
		PermissionProductReadOther,
		PermissionProductReadDeleted,
		PermissionProductCreate,
		PermissionProductUpdate,
		PermissionProductDelete,
		PermissionProductModifyOther,
	}

	SeedGroupPermissios = map[string][]string{
		"DEFAULT": {
			PermissionProductCreate,
			PermissionProductRead,
			PermissionProductUpdate,
			PermissionProductDelete,
		},
		"SUPER_USER": {
			PermissionFullAccess,
			PermissionProductAll,
			PermissionProductRead,
			PermissionProductReadOther,
			PermissionProductReadDeleted,
			PermissionProductCreate,
			PermissionProductUpdate,
			PermissionProductDelete,
			PermissionProductModifyOther,
		},
	}
)
