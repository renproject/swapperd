package foundation

type Logger interface {
	LogInfo(SwapID, string)
	LogDebug(SwapID, string)
	LogError(SwapID, error)

	GlobalLogInfo(string)
	GlobalLogDebug(string)
	GlobalLogError(error)
}
