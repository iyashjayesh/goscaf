package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var banner = `
 ██████╗  ██████╗ ███████╗████████╗ █████╗ ██████╗ ████████╗
██╔════╝ ██╔═══██╗██╔════╝╚══██╔══╝██╔══██╗██╔══██╗╚══██╔══╝
██║  ███╗██║   ██║███████╗   ██║   ███████║██████╔╝   ██║
██║   ██║██║   ██║╚════██║   ██║   ██╔══██║██╔══██╗   ██║
╚██████╔╝╚██████╔╝███████║   ██║   ██║  ██║██║  ██║   ██║
 ╚═════╝  ╚═════╝ ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝
`

var rootCmd = &cobra.Command{
	Use:   "gostart",
	Short: "gostart - enterprise-grade Go project scaffolder",
	Long: color.HiCyanString(banner) + "\n" +
		color.HiWhiteString("  gostart scaffolds production-quality Go project boilerplate.\n") +
		color.HiBlackString("  Think create-react-app, but for Go services.\n"),
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
