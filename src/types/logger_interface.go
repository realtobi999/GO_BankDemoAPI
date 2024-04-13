package types

type ILogger interface {
	LogEvent(message string)
	LogError(message string)
	LogDebug(message string)
	LogWarning(message string)
}