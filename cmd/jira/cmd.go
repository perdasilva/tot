package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/spf13/cobra"
	"tot/pkg/client"
)

func NewCmd() *cobra.Command {
	jiraCmd := &cobra.Command{
		Use:   "jira",
		Short: "Jira automations",
		RunE:  runE,
	}
	return jiraCmd
}

func runE(command *cobra.Command, args []string) error {
	jiraClient, err := client.NewTotClientFactory().NewJiraClient()
	if err != nil {
		return err
	}
	jql := "issuetype not in (Sub-task) AND project in (MKTPLC, OLM) OR project = OCPBUGSM AND component in (OLM) ORDER BY Rank ASC"
	opt := &jira.SearchOptions{
		MaxResults: 5,
		StartAt:    0,
	}
	issues, _, err := jiraClient.Issue.Search(jql, opt)

	if err != nil {
		return err
	}

	for _, issue := range issues {
		fmt.Printf("%s\t%s\t%s\t%s\n", issue.Key, issue.Fields.Status.Name, issue.Fields.Type.Name, issue.Fields.Summary)
	}

	return nil
}
