package permission

// Operation WebDAV 操作类型
type Operation string

const (
	OperationRead   Operation = "READ"
	OperationWrite  Operation = "WRITE"
	OperationCreate Operation = "CREATE"
	OperationDelete Operation = "DELETE"
)

// MapHTTPMethodToOperation 映射 HTTP 方法到操作
func MapHTTPMethodToOperation(method string) Operation {
	switch method {
	case "GET", "HEAD", "OPTIONS", "PROPFIND":
		return OperationRead
	case "PUT", "PATCH", "PROPPATCH":
		return OperationWrite
	case "POST", "MKCOL":
		return OperationCreate
	case "DELETE":
		return OperationDelete
	default:
		return OperationRead
	}
}

// MapOperationToPermission 映射操作到权限字符串
func MapOperationToPermission(op Operation) string {
	switch op {
	case OperationRead:
		return "R"
	case OperationWrite:
		return "U"
	case OperationCreate:
		return "C"
	case OperationDelete:
		return "D"
	default:
		return "R"
	}
}
