package nddscalegen

import (
	"gihub.com/yndd/ndd-scale-test/pkg/generator"
	"github.com/spf13/cobra"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	offset       int
	count        int
	outputDir    string
	templateFile string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:          "generate",
	Short:        "generate ndd scale templates",
	Aliases:      []string{"gen"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		zlog := zap.New(zap.UseDevMode(true), zap.JSONEncoder())
		log := logging.NewLogrLogger(zlog.WithName("nddscalegen"))
		log.Debug("generate nddscalegen templates ...")

		g, err := generator.NewGenerator(
			generator.WithIndexes(offset, count),
			generator.WithOutputDir(outputDir),
			generator.WithTemplate(templateFile),
		)
		if err != nil {
			return err
		}

		if err := g.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().IntVarP(&offset, "offset", "o", 100, "The offset of the index")
	generateCmd.Flags().IntVarP(&count, "count", "c", 10, "The number of templates generated")
	generateCmd.Flags().StringVarP(&templateFile, "templateFile", "t", "./templates/irb.tmpl", "the template used for the scale test generation")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "", "out/", "The directory that the Go package should be written to.")
}
