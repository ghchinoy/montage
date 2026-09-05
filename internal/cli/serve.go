package cli

import (
	"github.com/ghchinoy/montage/internal/server"
	"github.com/spf13/cobra"
)

var portFlag int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Montage web application and API service",
	Long: `Starts the Montage HTTP web service, serving the Lit WebComponent user interface,
the /api/assess repository evaluation endpoint, and in-memory artifact zip exports.
Automatically respects the $PORT environment variable for Cloud Run deployments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resolvedPort := server.ResolvePort(portFlag)
		return server.StartServer(resolvedPort)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&portFlag, "port", "p", 8080, "Port to listen on (defaults to $PORT or 8080)")
}
