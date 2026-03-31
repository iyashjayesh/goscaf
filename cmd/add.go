package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/iyashjayesh/goscaf/internal/generator"
)

var addCmd = &cobra.Command{
	Use:   "add <service-name>",
	Short: "Add a new service scaffold (DDD layout)",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}

		gen, err := generator.NewServiceGenerator(wd)
		if err != nil {
			return fmt.Errorf("failed to create servie generator : %w", err)
		}

		info, err := gen.Run(args[0])
		if err != nil {
			return fmt.Errorf("failed to run service generator : %w", err)
		}

		fmt.Println()
		color.HiGreen("  ✔ Service scaffold created successfully!")
		color.HiWhite("    Service: %s", info.StructName)
		color.HiWhite("    Package: %s", info.PackageName)
		color.HiWhite("    Domain:  internal/%s/service.go", info.DirectoryName)
		color.HiWhite("    Handler: internal/handler/%s_handler.go", info.FileBaseName)
		fmt.Println()
		color.HiCyan("  Next steps:")
		color.HiWhite("    wire %s.NewService() into your application dependencies", info.PackageName)
		color.HiWhite("    register handler routes in internal/server/server.go")
		fmt.Println()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
