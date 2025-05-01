package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blackmou5e/grafana-annotator/internal/config"
	"github.com/blackmou5e/grafana-annotator/internal/grafana"
	"github.com/blackmou5e/grafana-annotator/internal/service"
	"github.com/blackmou5e/grafana-annotator/internal/validation"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	annotationTags string
	annotationText string
	log            *logrus.Logger
	cfg            *config.Config
)

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

var rootCmd = &cobra.Command{
	Use:   "grafana-annotator",
	Short: "A tool to create annotations in Grafana dashboards",
	Long:  `grafana-annotator allows you to create annotations across multiple Grafana dashboards simultaneously`,
}

var annotateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create annotations in Grafana dashboards",
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg, err := config.SetDefaultConfig()
		if err != nil {
			log.Printf("Failed to set default config values: %v", err)
		}

		err = config.LoadConfigFromFile(cfg)
		if err != nil {
			log.Printf("Failed to load config: %v", err)
		}

		err = config.LoadConfigFromEnv(cfg)
		if err != nil {
			log.Printf("Failed to read config from env: %v", err)
		}

		if cfg.Debug {
			log.SetLevel(logrus.DebugLevel)
		}
	},
	Run: runAnnotate,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n", version, commit, date)
	},
}

func init() {
	annotateCmd.Flags().StringVarP(&annotationTags, "tags", "t", "", "Comma-separated list of tags for the annotation")
	annotateCmd.Flags().StringVarP(&annotationText, "message", "m", "", "Text message for the annotation")

	annotateCmd.MarkFlagRequired("tags")
	annotateCmd.MarkFlagRequired("message")

	rootCmd.AddCommand(annotateCmd)
	rootCmd.AddCommand(versionCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) {
	tags := strings.Split(annotationTags, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	var client grafana.GrafanaClient = grafana.NewClient(
		cfg.GrafanaURL,
		cfg.GrafanaServiceAccountToken,
		time.Duration(cfg.Timeout)*time.Second,
	)

	validator := validation.NewValidator()

	svc := service.NewAnnotatorService(client, log, validator)

	if err := svc.CreateAnnotations(ctx, tags, annotationText); err != nil {
		log.Fatalf("Failed to create annotations: %v", err)
	}

	log.Info("Successfully created annotations")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
