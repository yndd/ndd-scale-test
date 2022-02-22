package nddscalegen

import (
	"bytes"

	log "github.com/sirupsen/logrus"

	"gihub.com/yndd/ndd-scale-test/pkg/generator"
	"gihub.com/yndd/ndd-scale-test/pkg/output"
	"github.com/spf13/cobra"
)

var (
	offset       int
	count        int
	outputDir    string
	templateFile string
	kubeconfig   string
)

// fileCommand represents the generate command
var fileCommand = &cobra.Command{
	Use:   "file",
	Short: "generate ndd scale yaml files from template",
	//Aliases:      []string{"gen"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		g, err := generator.NewGenerator(
			generator.WithIndexes(offset, count),
			generator.WithTemplate(templateFile),
		)
		if err != nil {
			log.Fatal(err)
		}

		var data []*bytes.Buffer
		if data, err = g.Generate(); err != nil {
			log.Fatal(err)
		}

		outputPlugin := output.NewFileOutput(outputDir)

		for k, v := range data {
			err = outputPlugin.Commit(v, &output.OutputPluginInfo{Index: k})
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileCommand)
	rootCmd.AddCommand(k8sCommand)
	k8sCommand.AddCommand(k8sAddCommand)
	k8sCommand.AddCommand(k8sDeleteCommand)
	k8sCommand.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "kubeconfig file")

	rootCmd.PersistentFlags().IntVarP(&offset, "offset", "o", 0, "The offset of the index")
	rootCmd.PersistentFlags().IntVarP(&count, "count", "c", 1, "The number of templates generated")
	rootCmd.PersistentFlags().StringVarP(&templateFile, "templateFile", "t", "./templates/irb.tmpl", "the template used for the scale test generation")
	fileCommand.Flags().StringVarP(&outputDir, "output-dir", "", "out/", "The directory that the Go package should be written to.")
}

// k8sCommand represents the generate command
var k8sCommand = &cobra.Command{
	Use:          "k8s",
	Short:        "generate CustomResources from template and commit them to the api-server",
	Aliases:      []string{"k"},
	SilenceUsage: false,
}

// k8sCommand represents the generate command
var k8sDeleteCommand = &cobra.Command{
	Use:          "delete",
	Short:        "generate CustomResources from template and use the names to delete them from the api-server",
	Aliases:      []string{"d"},
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {

		g, err := generator.NewGenerator(
			generator.WithIndexes(offset, count),
			generator.WithTemplate(templateFile),
		)
		if err != nil {
			log.Fatal(err)
		}

		var data []*bytes.Buffer
		if data, err = g.Generate(); err != nil {
			log.Fatal(err)
		}

		outputPlugin := output.NewK8sOutput(kubeconfig)

		for k, v := range data {
			err = outputPlugin.Delete(v, &output.OutputPluginInfo{Index: k})
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	},
}

// k8sCommand represents the generate command
var k8sAddCommand = &cobra.Command{
	Use:          "commit",
	Short:        "generate CustomResources from template and commit them to the api-server",
	Aliases:      []string{"c"},
	SilenceUsage: false,
	RunE: func(cmd *cobra.Command, args []string) error {

		g, err := generator.NewGenerator(
			generator.WithIndexes(offset, count),
			generator.WithTemplate(templateFile),
		)
		if err != nil {
			log.Fatal(err)
		}

		var data []*bytes.Buffer
		if data, err = g.Generate(); err != nil {
			log.Fatal(err)
		}

		outputPlugin := output.NewK8sOutput(kubeconfig)

		for k, v := range data {
			err = outputPlugin.Commit(v, &output.OutputPluginInfo{Index: k})
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	},
}
