package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Alonza0314/it-system/controller/backend/config"
	"github.com/Alonza0314/it-system/controller/backend/internal"
	"github.com/Alonza0314/it-system/controller/backend/logger"

	loggergoUtil "github.com/Alonza0314/logger-go/v2/util"
	"github.com/free-ran-ue/util"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use: "controller",
	Run: controllerFunc,
}

func init() {
	controllerCmd.Flags().StringP("config", "c", "config.yaml", "Path to the configuration file")
	if err := controllerCmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}
}

func controllerFunc(cmd *cobra.Command, args []string) {
	controllerConfigFilePath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	controllerConfig := config.Config{}
	if err := util.LoadFromYaml(controllerConfigFilePath, &controllerConfig); err != nil {
		panic(err)
	}

	logger := logger.NewBackendLogger(loggergoUtil.LogLevelString(controllerConfig.Logger.Level), "", true)

	var discordWebhookURL string
	if controllerConfig.Backend.Discord.Enabled {
		if url, err := os.ReadFile(controllerConfig.Backend.Discord.WebhookUrlPath); err == nil {
			discordWebhookURL = strings.TrimSpace(string(url))
			if discordWebhookURL == "" {
				panic(fmt.Errorf("discord webhook URL is empty in %s", controllerConfig.Backend.Discord.WebhookUrlPath))
			}
		} else {
			panic(err)
		}
	}

	controller := internal.NewBackend(&controllerConfig, discordWebhookURL, logger)
	if controller == nil {
		panic("failed to initialize the backend")
	}

	controller.Start()
	defer controller.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}

func Execute() {
	if err := controllerCmd.Execute(); err != nil {
		panic(err)
	}
}
