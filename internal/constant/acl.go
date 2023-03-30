package constant

const (
	PermissionFullAccess = string("FULL_ACCESS") // Full access to all object

	PermissionProductAll         = string("PRODUCT_ALL")          // Full access to product that they own
	PermissionProductCreate      = string("PRODUCT_CREATE")       // Only access to create product
	PermissionProductUpdate      = string("PRODUCT_UPDATE")       // Only access to update product
	PermissionProductDelete      = string("PRODUCT_DELETE")       // Only access to delete product
	PermissionProductRead        = string("PRODUCT_READ")         // Only access to read product
	PermissionProductModifyOther = string("PRODUCT_MODIFY_OTHER") // Only allow access to modify other user product
)
