// Package view implements the user view command.
package view

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// ViewOptions configures the user view command.
type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Username string
	JSON     bool
}

// NewCmdView creates the user view command.
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view [<username>]",
		Short: "View a user profile",
		Long: heredoc.Doc(`
			View a GitCode user profile. With no arguments, shows the
			current authenticated user. Pass a username to view a
			specific user.
		`),
		Example: heredoc.Doc(`
			# View current user
			$ gc user view

			# View a specific user
			$ gc user view <username>

			# JSON output
			$ gc user view --json
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Username = strings.TrimSpace(args[0])
			}
			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func viewRun(opts *ViewOptions) error {
	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	var user *api.User
	if opts.Username != "" {
		user, err = api.GetUser(client, opts.Username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
	} else {
		user, err = api.CurrentUser(client)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, user)
	}

	printProfile(opts, user)
	return nil
}

func printProfile(opts *ViewOptions, u *api.User) {
	out := opts.IO.Out
	cs := opts.IO.ColorScheme()

	fmt.Fprintf(out, "%s\n", cs.Bold(u.Login))
	if u.Name != "" {
		fmt.Fprintf(out, "  Name:      %s\n", u.Name)
	}
	if u.Type != "" {
		fmt.Fprintf(out, "  Type:      %s\n", u.Type)
	}
	if u.Email != "" {
		fmt.Fprintf(out, "  Email:     %s\n", u.Email)
	}
	if u.Bio != "" {
		fmt.Fprintf(out, "  Bio:       %s\n", u.Bio)
	}
	if u.Company != "" {
		fmt.Fprintf(out, "  Company:   %s\n", u.Company)
	}
	if u.Blog != "" {
		fmt.Fprintf(out, "  Blog:      %s\n", u.Blog)
	}
	if u.HTMLURL != "" {
		fmt.Fprintf(out, "  Profile:   %s\n", u.HTMLURL)
	}
	if u.AvatarURL != "" {
		fmt.Fprintf(out, "  Avatar:    %s\n", u.AvatarURL)
	}
	if u.Followers > 0 || u.Following > 0 {
		fmt.Fprintf(out, "  Followers: %d\n", u.Followers)
		fmt.Fprintf(out, "  Following: %d\n", u.Following)
	}
}
