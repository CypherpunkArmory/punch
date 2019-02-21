package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var githubAPIToken string

var githubRepo string

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update cli version",
	Long:  `update cli to latest release on github`,
	Run: func(cmd *cobra.Command, args []string) {
		confirmAndSelfUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func confirmAndSelfUpdate() {
	up, err := selfupdate.NewUpdater(selfupdate.Config{
		APIToken: githubAPIToken,
	})
	latest, found, err := up.DetectLatest(githubRepo)
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return
	}
	v := semver.MustParse(version)
	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return
	}

	fmt.Print("Do you want to update to: ", latest.Version, "? (y/n): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input != "y\n" && input != "n\n") {
		log.Println("Invalid input")
		return
	}
	if input == "n\n" {
		return
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path")
		return
	}
	if err := up.UpdateTo(latest, exe); err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
}
