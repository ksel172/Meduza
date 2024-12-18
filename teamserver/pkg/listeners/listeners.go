package listeners

type ListenerType = string

type unknown = map[string]any

const (
	HTTP    ListenerType = "http" // HTTP is a constant for all HTTP listener types (e.g. HTTP/1 and etc.)
	TCP     ListenerType = "tcp"  // TCP is a constant for TCP bind & reverse listeners.
	FOREIGN ListenerType = "foreign"
	SMB     ListenerType = "smb" // SMB is a constant for SMB named pipe bind  & reverse listeners.
)

func AllListenerType() []ListenerType {
	return []ListenerType{
		HTTP, TCP, FOREIGN, SMB,
	}
}
