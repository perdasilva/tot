package gh

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"tot/pkg/client"
)

func NewCmd() *cobra.Command {
	ghCmd := &cobra.Command{
		Use:   "gh",
		Short: "GitHub automations",
		RunE:  runE,
	}
	return ghCmd
}

func runE(command *cobra.Command, args []string) error {
	gh, err := client.NewTotClientFactory().NewGitHubClient()

	if err != nil {
		return err
	}

	repos, _, err := gh.Repositories.List(context.TODO(), "", nil)

	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Printf("%d\t%s\t%s\n", repo.ID, *repo.Name, *repo.URL)
	}

	return nil
}
