package listeners

import (
	"errors"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/pkg/listeners/http"
)

type ListenersService struct{}

func (s *ListenersService) CreateListener(listenerType string, config any) (ListenerRegistry, error) {
	parseConfig, err := ParseConfig(listenerType, config)
	if err != nil {
		return nil, err
	}
	switch listenerType {
	case "http", "https", "http3", "h2c":
		httpConfig, ok := parseConfig.(*http.Config)
		if !ok {
			return nil, errors.New("parsed config is not of type *http.Config")
		}
		return http.NewHttpListener(httpConfig.HostHeader, *httpConfig), nil
	default:
		return nil, fmt.Errorf("unsupported listener type: %s", listenerType)
	}
}
