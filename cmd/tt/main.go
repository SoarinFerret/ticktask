/*
Copyright © 2024 Cody Ernesti
*/
package main

import (
	"github.com/soarinferret/ticktask/cmd"
	"github.com/soarinferret/ticktask/internal/config"
)

var Version string

func main() {

	config.LoadConfig()

	if Version != "" {
		cmd.Version = Version
	} else {
		cmd.Version = "dev"
	}

	cmd.Execute()
}
