// Package config implements user-visible CLI configuration commands.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

const defaultHost = "gitcode.com"

type configResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewCmdConfig creates the config command.
func NewCmdConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <command>",
		Short: "Read and write GitCode CLI configuration",
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSet(f))
	return cmd
}

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Read a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToLower(strings.TrimSpace(args[0]))
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			value, err := cfg.Get(defaultHost, key)
			if err != nil {
				return err
			}
			if key == "update.mode" && value == "" {
				value = "notify"
			}
			result := configResult{Key: key, Value: value}
			if jsonOutput {
				return cmdutil.WriteJSON(cmd.OutOrStdout(), result)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), value)
			return err
		},
	}
	cmdutil.AddJSONFlag(cmd, &jsonOutput)
	return cmd
}

func newCmdSet(f *cmdutil.Factory) *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Store a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.ToLower(strings.TrimSpace(args[0]))
			value := strings.TrimSpace(args[1])
			if key == "update.mode" {
				value = strings.ToLower(value)
			}
			if err := validateValue(key, value); err != nil {
				return cmdutil.NewUsageError(err.Error())
			}
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			if err := cfg.Set(defaultHost, key, value); err != nil {
				return err
			}
			result := configResult{Key: key, Value: value}
			if jsonOutput {
				return cmdutil.WriteJSON(cmd.OutOrStdout(), result)
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Set %s=%s\n", key, value)
			return err
		},
	}
	cmdutil.AddJSONFlag(cmd, &jsonOutput)
	return cmd
}

func validateValue(key, value string) error {
	if key != "update.mode" {
		return nil
	}
	switch value {
	case "auto", "notify", "off":
		return nil
	default:
		return errors.New("update.mode must be auto, notify, or off")
	}
}
