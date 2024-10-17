/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/go-caja/pscan/scan"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List the hosts in hosts list",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}

		return listAction(os.Stdout, hostsFile, args)
	},
}

func listAction(out io.Writer, hostsFile string, arg []string) error {
	h := &scan.HostsList{}

	if err := h.Load(hostsFile); err != nil {
		return err
	}

	for _, host := range h.Hosts {
		if _, err := fmt.Fprintln(out, host); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
