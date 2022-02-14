package bugzilla

import (
	"fmt"
	"github.com/spf13/cobra"
	"tot/pkg/client"
)

func NewCmd() *cobra.Command {
	bzCmd := &cobra.Command{
		Use:   "bz",
		Short: "Bugzilla automations",
		RunE:  runE,
	}
	return bzCmd
}

func runE(command *cobra.Command, args []string) error {
	bzClient, err := client.NewTotClientFactory().NewBugzillaClient()

	if err != nil {
		return err
	}

	bugs, err := bzClient.QuickSearch(map[string]string{
		"bug_status":        "CLOSED",
		"component":         "OLM",
		"product":           "OpenShift Container Platform",
		"rh_sub_components": "OLM",
	}, 10000)

	//bugs, err := bzClient.SearchBugs(map[string]string{
	//	"bug_status":        "CLOSED",
	//	"component":         "OLM",
	//	"product":           "OpenShift Container Platform",
	//	"rh_sub_components": "OLM",
	//	"limit":             "100",
	//})

	if err != nil {
		return err
	}

	for _, bug := range bugs {
		fmt.Printf("%d\t%s\t%s\n", bug.ID, bug.Status, bug.Summary)
	}
	fmt.Printf("Found %d bugs\n", len(bugs))

	return nil
}
