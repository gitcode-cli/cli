// Package schema exposes machine-readable command metadata.
package schema

import (
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

type schemaFlag struct {
	Name      string   `json:"name"`
	Shorthand string   `json:"shorthand,omitempty"`
	Usage     string   `json:"usage"`
	Type      string   `json:"type"`
	Default   string   `json:"default,omitempty"`
	Required  bool     `json:"required"`
	Enum      []string `json:"enum,omitempty"`
}

type schemaCommand struct {
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	Use         string          `json:"use"`
	Short       string          `json:"short,omitempty"`
	Long        string          `json:"long,omitempty"`
	Example     string          `json:"example,omitempty"`
	Hidden      bool            `json:"hidden"`
	Flags       []schemaFlag    `json:"flags,omitempty"`
	Subcommands []schemaCommand `json:"subcommands,omitempty"`
}

type Options struct {
	Root *cobra.Command
	Path string
}

// NewCmdSchema creates the schema command.
func NewCmdSchema(root *cobra.Command) *cobra.Command {
	opts := &Options{Root: root}

	cmd := &cobra.Command{
		Use:   "schema [command-path]",
		Short: "Print machine-readable command metadata",
		Long: heredoc.Doc(`
			Print machine-readable metadata for the command tree or a specific command.
		`),
		Example: heredoc.Doc(`
			# Print the entire command tree
			$ gc schema

			# Print schema for a specific command
			$ gc schema "issue view"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Path = strings.Join(args, " ")
			}

			target := root
			if opts.Path != "" {
				found, _, err := root.Find(strings.Fields(opts.Path))
				if err != nil {
					return cmdutil.NewUsageError(err.Error())
				}
				target = found
			}

			return cmdutil.WriteJSON(cmd.OutOrStdout(), buildSchema(target))
		},
	}

	return cmd
}

func buildSchema(cmd *cobra.Command) schemaCommand {
	entry := schemaCommand{
		Name:    cmd.Name(),
		Path:    cmd.CommandPath(),
		Use:     cmd.Use,
		Short:   cmd.Short,
		Long:    cmd.Long,
		Example: cmd.Example,
		Hidden:  cmd.Hidden,
		Flags:   buildFlags(cmd),
	}

	for _, child := range cmd.Commands() {
		if child.Name() == "help" {
			continue
		}
		entry.Subcommands = append(entry.Subcommands, buildSchema(child))
	}

	return entry
}

func buildFlags(cmd *cobra.Command) []schemaFlag {
	var flags []schemaFlag
	cmd.NonInheritedFlags().VisitAll(func(flag *pflag.Flag) {
		flags = append(flags, schemaFlag{
			Name:      flag.Name,
			Shorthand: flag.Shorthand,
			Usage:     flag.Usage,
			Type:      flag.Value.Type(),
			Default:   flag.DefValue,
			Required:  flag.Annotations != nil && len(flag.Annotations[cobra.BashCompOneRequiredFlag]) > 0,
		})
	})
	return flags
}
