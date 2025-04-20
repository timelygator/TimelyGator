//go:build darwin

package lib

import "fmt"

const (
	StrategyJXA         = "jxa"
	StrategyAppleScript = "applescript"
)

func GetCurrentWindowMacOS(strategy string) (*WindowInfo, error) {
	switch strategy {
	case StrategyJXA:
		info, err := GetInfo()
		if err != nil {
			return nil, err
		}
		return &WindowInfo{
			App:   info["app"],
			Title: info["title"],
			URL:   StringPtr(info["url"]),
			Incognito: func() *bool {
				if info["incognito"] == "" {
					return nil
				}
				b := info["incognito"] == "true"
				return &b
			}(),
		}, nil
	case StrategyAppleScript:
		info, err := GetInfoAppleScript()
		if err != nil {
			return nil, err
		}
		return &WindowInfo{App: info["app"], Title: info["title"]}, nil
	default:
		return nil, FatalError{fmt.Sprintf("invalid strategy %q", strategy)}
	}
}
