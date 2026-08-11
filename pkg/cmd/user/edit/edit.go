// Package edit implements the user edit command.
package edit

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"gitcode.com/gitcode-cli/cli/api"
	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// EditOptions configures the user edit command.
type EditOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	Name     string
	Bio      string
	Email    string
	Company  string
	Location string
	Website  string
	JSON     bool
}

// NewCmdEdit creates the user edit command.
func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit your user profile",
		Long: heredoc.Doc(`
			Update the current authenticated user's profile information.
			Only provided fields are updated.
		`),
		Example: heredoc.Doc(`
			# Update name
			$ gc user edit --name "New Name"

			# Update bio and company
			$ gc user edit --bio "New bio" --company "Acme Inc"
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return editRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Name, "name", "", "Display name")
	cmd.Flags().StringVar(&opts.Bio, "bio", "", "Bio / description")
	cmd.Flags().StringVar(&opts.Email, "email", "", "Public email")
	cmd.Flags().StringVar(&opts.Company, "company", "", "Company")
	cmd.Flags().StringVar(&opts.Location, "location", "", "Location")
	cmd.Flags().StringVar(&opts.Website, "website", "", "Website URL")
	cmdutil.AddJSONFlag(cmd, &opts.JSON)

	return cmd
}

func editRun(opts *EditOptions) error {
	req := &api.UpdateUserRequest{
		Nickname:    opts.Name,
		Description: opts.Bio,
		Email:       opts.Email,
		Company:     opts.Company,
		Location:    opts.Location,
		Website:     opts.Website,
	}

	if req.Nickname == "" && req.Description == "" && req.Email == "" &&
		req.Company == "" && req.Location == "" && req.Website == "" {
		return cmdutil.NewUsageError("at least one field must be provided to update")
	}

	if req.Description != "" {
		if err := cmdutil.ScanContentForSecrets(req.Description); err != nil {
			return err
		}
	}

	client, err := cmdutil.AuthenticatedClientFromFactory(opts.HttpClient)
	if err != nil {
		return err
	}

	user, err := api.UpdateUser(client, req)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if opts.JSON {
		return cmdutil.WriteJSON(opts.IO.Out, user)
	}

	fmt.Fprintf(opts.IO.ErrOut, "Updated profile for %s.\n", user.Login)
	return nil
}
