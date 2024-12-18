package listeners

type ListenerType = string

type unknown = map[string]any

const (
	HTTP    ListenerType = "http"
	TCP     ListenerType = "tcp"
	FOREIGN ListenerType = "foreign"
	SMB     ListenerType = "smb"
)

func AllListenerTypes() []ListenerType {
	return []ListenerType{
		HTTP, TCP, FOREIGN, SMB,
	}
}
