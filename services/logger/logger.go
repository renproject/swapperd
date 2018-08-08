package logger

type Logger interface {
	LogError([32]byte, string)
	LogInfo([32]byte, string)
	LogDebug([32]byte, string)
}
