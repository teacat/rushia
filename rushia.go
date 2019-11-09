package rushia

// New 會建立一個新的 SQL 建置工具。
func New() Builder {
	return newBuilder()
}

func NewMigration() Migration {
	return newMigration()
}
