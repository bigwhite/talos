/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// nolint: dupl,golint
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/talos-systems/talos/cmd/osctl/pkg/client"
	"github.com/talos-systems/talos/cmd/osctl/pkg/helpers"
	"github.com/talos-systems/talos/internal/pkg/constants"
)

// kubeconfigCmd represents the kubeconfig command
var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Download the admin.conf from the node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := client.NewDefaultClientCredentials(talosconfig)
		if err != nil {
			helpers.Fatalf("error getting client credentials: %s", err)
		}
		c, err := client.NewClient(constants.OsdPort, creds)
		if err != nil {
			helpers.Fatalf("error constructing client: %s", err)
		}
		if err := c.Kubeconfig(); err != nil {
			helpers.Fatalf("error fetching kubeconfig: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(kubeconfigCmd)
}
