/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package cmd

import (
	"github.com/gzuidhof/starboard-cli/starboard/internal/nbserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the notebook files in your current folder",
	Long:  `Serve serves the notebook files in the current folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		nbserver.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().StringP("port", "p", "8585", "Port to serve files on")
	serveCmd.Flags().String("port_secondary", "15742", "Port used as secondary origin (for additional sandboxing)")

	serveCmd.Flags().StringP("folder", "f", ".", "Folder (or file) to serve, defaults to the current working directory")

	serveCmd.Flags().String("static_folder", "", "Override where static assets are served from, it uses the embedded assets if not set")
	serveCmd.Flags().String("templates_folder", "", "Override where templates are loaded from, it uses the embedded assets if not set")

	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("port_secondary", serveCmd.Flags().Lookup("port_secondary"))

	viper.BindPFlag("static_folder", serveCmd.Flags().Lookup("static_folder"))
	viper.BindPFlag("templates_folder", serveCmd.Flags().Lookup("templates_folder"))

	viper.BindPFlag("serve.folder", serveCmd.Flags().Lookup("folder"))
}
