package mote

const (
	CodeSetup     = "SETUP"
	CodeProtocol  = "PROTOCOL"
	CodeTransport = "TRANSPORT"
	CodeTimeout   = "TIMEOUT"
)

type MoteError struct {
	Code    string
	Message string
	Err     error
}

func (e *MoteError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *MoteError) Unwrap() error { return e.Err }

// Is reports a match when target is a MoteError carrying the same code, so
// errors.Is works against the exported sentinels below.
func (e *MoteError) Is(target error) bool {
	other, ok := target.(*MoteError)
	return ok && other.Code == e.Code
}

var (
	ErrSetup     = &MoteError{Code: CodeSetup, Message: "setup error"}
	ErrProtocol  = &MoteError{Code: CodeProtocol, Message: "protocol error"}
	ErrTransport = &MoteError{Code: CodeTransport, Message: "transport error"}
	ErrTimeout   = &MoteError{Code: CodeTimeout, Message: "timeout error"}
)

func newSetup(message string, err error) *MoteError {
	return &MoteError{Code: CodeSetup, Message: message, Err: err}
}

func newProtocol(message string, err error) *MoteError {
	return &MoteError{Code: CodeProtocol, Message: message, Err: err}
}

func newTransport(message string, err error) *MoteError {
	return &MoteError{Code: CodeTransport, Message: message, Err: err}
}

func newTimeout(message string, err error) *MoteError {
	return &MoteError{Code: CodeTimeout, Message: message, Err: err}
}
