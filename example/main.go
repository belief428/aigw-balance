package main

import (
	aibalance "github.com/belief428/aigw-balance"
)

func main() {
	enforcer := aibalance.NewEnforcer(aibalance.WithPort(11118))

	if err := enforcer.Enforcer(); err != nil {
		return
	}
	select {}
}
