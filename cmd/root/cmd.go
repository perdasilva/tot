package root

import (
	"github.com/spf13/cobra"
	"tot/cmd/bugzilla"
	"tot/cmd/gh"
	"tot/cmd/jira"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tot",
		Short: "Tot is a command-line cli used to help automate some of the OLM team's workflows",
		Long: `For some, "tot" is short for "(T)ot has n(o)t been named ye(t)"; or for others: "(T)he (O)LM (t)eam."
It is meant to support the OLM team by automating some workflows, e.g. generating reports (sprint, bug, epic, etc.),
for querying the data sources we use more quickly, and for automating any workflow we can image. It is built as a 
command-line client to be as flexible as possible and adapt to our team's changing needs. It could also grow 
into something our customers can use to interface with us: e.g. submit bug reports, etc.`,
		Args: cobra.NoArgs,
	}
	rootCmd.AddCommand(bugzilla.NewCmd(), jira.NewCmd(), gh.NewCmd())
	return rootCmd
}
