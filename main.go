package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	aw "github.com/deanishe/awgo"
)

var wf *aw.Workflow

var tunnelblickctlBin string

func init() {
	tunnelblickctlBin = "./bin/tunnelblickctl"
	wf = aw.New()
}

func list() {
	out, err := exec.Command(tunnelblickctlBin, "list").Output()
	if err != nil {
		log.Fatal(err)
	}

	listConfigs := strings.Trim(string(out), "\n")
	configs := strings.Split(listConfigs, "\n")

	for _, config := range configs {
		wf.NewItem(config).Valid(true).Arg(config)
	}
	wf.SendFeedback()
}

func connect(config string) {
	_, err := exec.Command(tunnelblickctlBin, "connect", config).Output()
	if err != nil {
		log.Fatal(err)
	}

	// golang string interpolation
	successMessage := fmt.Sprint("Connecting to ", config)
	log.Println(successMessage)
}

// disconnect vpn
func disconnect(config string) {
	_, err := exec.Command(tunnelblickctlBin, "disconnect", config).Output()
	if err != nil {
		log.Fatal(err)
	}

	// golang string interpolation
	successMessage := fmt.Sprint("Disconnected from ", config)
	log.Println(successMessage)
}

func disconnectAll() {
	_, err := exec.Command(tunnelblickctlBin, "disconnect", "--all").Output()
	if err != nil {
		log.Fatal(err)
	}

	// golang string interpolation
	successMessage := fmt.Sprint("Disconnected from all VPN")
	log.Println(successMessage)
}

func run() {
	// parse args
	args := wf.Args()
	log.Println(args)
	if len(args) == 0 {
		fmt.Println("No args is passed")
		return
	}

	if args[0] == "list" {
		list()
	} else if args[0] == "connect" {
		connect(args[1])
	} else if args[0] == "disconnect" {
		disconnect(args[1])
	} else if args[0] == "disconnect-all" {
		disconnectAll()
	} else {
		fmt.Println("Unknown command")
	}
}

func main() {
	wf.Run(run)
}
