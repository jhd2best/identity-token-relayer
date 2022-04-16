package lib

import (
	"fmt"
	"go.uber.org/zap"
	"identity-token-relayer/log"
	"testing"
)

func TestNewPagerIncident(t *testing.T) {
	// send PagerDuty incident
	pagerSummary := fmt.Sprintf("sync transaction failed\n\ntxHash:%s\ncontract:%s\ntokenId:%d", "0x2a96e24dad86cc67917cb8e0730054384c0e8f813ea6c7b1c1a8fac282ba8a98", "0xeeb0a631E178ec91Cb735032F50F11Fbe617876E", 3222)
	pagerErr := NewPagerIncident("failed transaction 0x2a96e24dad86cc67917cb8e0730054384c0e8f813ea6c7b1c1a8fac282ba8a98", "trigger", "sync transaction failed", pagerSummary, "critical")
	if pagerErr != nil {
		log.GetLogger().Warn("send PagerDuty incident failed.", zap.String("error", pagerErr.Error()))
	}
}
