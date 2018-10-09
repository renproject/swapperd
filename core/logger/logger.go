package logger

type Logger interface {
	LogInfo(orderID [32]byte, msg string)
	LogDebug(orderID [32]byte, msg string)
	LogError(orderID [32]byte, msg string)
}
