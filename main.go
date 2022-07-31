package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
)

var (
	wf      *aw.Workflow
	query   string
	doCheck bool

	repo          = "drselump14/alfred-tunnelblick-manager"
	iconAvailable = &aw.Icon{Value: "update-available.png"}

	availableCommands = []string{"connect", "disconnect", "disconnect-all", "check"}
)

const tunnelblickctlBin = "./bin/tunnelblickctl"
const updateJobName = "checkForUpdate"

func init() {
	flag.BoolVar(&doCheck, "check", false, "check for a new version")

	wf = aw.New(update.GitHub(repo))
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
	flag.Parse()

	query = flag.Arg(0)

	if doCheck {
		wf.Configure(aw.TextErrors(true))
		log.Println("Checking for updates...")
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
		return
	}
	// ----------------------------------------------------------------
	// Script Filter
	// ----------------------------------------------------------------

	// Call self with "check" command if an update is due and a check
	// job isn't already running.
	if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
		log.Println("Running update check in background...")

		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(updateJobName, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	// Only show update status if query is empty.
	if query == "" {

		if wf.UpdateAvailable() {
			// Turn off UIDs to force this item to the top.
			// If UIDs are enabled, Alfred will apply its "knowledge"
			// to order the results based on your past usage.
			wf.Configure(aw.SuppressUIDs(true))

			// Notify user of update. As this item is invalid (Valid(false)),
			// actioning it expands the query to the Autocomplete value.
			// "workflow:update" triggers the updater Magic Action that
			// is automatically registered when you configure Workflow with
			// an Updater.
			//
			// If executed, the Magic Action downloads the latest version
			// of the workflow and asks Alfred to install it.
			wf.NewItem("Update available!").
				Subtitle("↩ to install").
				Autocomplete("workflow:update").
				Valid(true).
				Icon(iconAvailable)

		}

		for _, name := range availableCommands {
			wf.NewItem(name).
				Arg(name).
				UID(name).
				Valid(true)
		}
	} else if query == "list" {
		list()
	} else if query == "connect" {
		connect(args[1])
	} else if query == "disconnect" {
		disconnect(args[1])
	} else if query == "disconnect-all" {
		disconnectAll()
	}

	wf.WarnEmpty("No matching items", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
