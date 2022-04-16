package lib

import (
	"context"
	pd "github.com/PagerDuty/go-pagerduty"
	"identity-token-relayer/config"
)

func NewPagerIncident(incidentKey, action string, summary string, msg string, severity string) error {
	serviceKey := config.Get().Debug.PagerDutyKey

	// skip when not config key
	if serviceKey == "" {
		return nil
	}

	_, err := pd.ManageEventWithContext(context.Background(), pd.V2Event{
		RoutingKey: serviceKey,
		Action:     action,
		DedupKey:   incidentKey,
		Payload: &pd.V2Payload{
			Summary:  summary,
			Severity: severity,
			Details:  msg,
			Source:   "Prod",
		},
	})
	return err
}
