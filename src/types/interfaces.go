package types

type ILogger interface {
	LogEvent(message any)
	LogError(message any)
	LogDebug(message any)
	LogWarning(message any)
}

type IStorage interface {
	DatabaseHas(table, column string, value any) bool
}
