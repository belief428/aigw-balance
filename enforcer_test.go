package aibalance

import "testing"

func TestNewEnforcer(t *testing.T) {
	NewEnforcer(WithPort(1111)).Enforcer()

	select {}
}
