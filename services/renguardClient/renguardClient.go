package renguardClient

type RenguardClient interface {
	ComplainDelayedAddressSubmission([32]byte) error
	ComplainDelayedRequestorInitiation([32]byte) error
	ComplainWrongRequestorInitiation([32]byte) error
	ComplainDelayedResponderInitiation([32]byte) error
	ComplainWrongResponderInitiation([32]byte) error
	ComplainDelayedRequestorRedemption([32]byte) error
}
