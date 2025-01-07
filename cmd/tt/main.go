/*
Copyright © 2024 Cody Ernesti
*/
package main

import (
	"github.com/soarinferret/ticktask/cmd"
	"github.com/soarinferret/ticktask/internal/config"
)

func main() {

	config.LoadConfig()

	cmd.Execute()

}
