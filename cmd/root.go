package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Version is overwritten by ldflags at build time
	Version = "0.1.0"
)

var banner = `
  ██████╗  ██████╗ ███████╗ ██████╗  █████╗  ███████╗
 ██╔════╝ ██╔═══██╗██╔════╝██╔════╝ ██╔══██╗ ██╔════╝
 ██║  ███╗██║   ██║███████╗██║      ███████║ █████╗  
 ██║   ██║██║   ██║╚════██║██║      ██╔══██║ ██╔══╝  
 ╚██████╔╝╚██████╔╝███████║╚██████╗ ██║  ██║ ██║     
  ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═╝     
`

var rootCmd = &cobra.Command{
	Use:   "goscaf",
	Short: "goscaf - enterprise-grade Go project scaffolder",
	Long: color.HiCyanString(banner) + "\n" +
		color.HiWhiteString("  goscaf scaffolds production-quality Go project boilerplate.\n") +
		color.HiBlackString("  Think create-react-app, but for Go services.\n"),
	Version: Version,
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
