package main

import (
	"flag"
	"net/url"
	"os"
	"strings"

	"github.com/j4ng5y/argo-inventory/pkg/argocd"
	"github.com/rs/zerolog"
)

var (
	outfile, logLevel, argoUrl, argoUsername, argoPassword, kubeTargetVersion string
)

func init() {
	flag.StringVar(&argoUsername, "argo-username", "", "The ArgoCD username to use.")
	flag.StringVar(&argoPassword, "argo-password", "", "The ArgoCD password to use.")
	flag.StringVar(&argoUrl, "argo-url", "", "The ArgoCD URL to use.")
	flag.StringVar(&kubeTargetVersion, "kube-target-version", "1.30.1", "The Version of Kubernetes to ")
	flag.StringVar(&logLevel, "log-level", "info", "The log level to use.")
	flag.StringVar(&outfile, "out", "report.csv", "The output file to write the report to.")
	flag.Parse()
}

func main() {
	var (
		argourl *url.URL
		err     error
	)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	loglvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		logger.Warn().Err(err).Msg("log level failed to parse, defaulting to info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(loglvl)
	}

	if argoUrl == "" {
		logger.Fatal().Msg("argo-url flag was not set")
	} else {
		argourl, err = url.Parse(argoUrl)
		if err != nil {
			logger.Fatal().Err(err).Send()
		} else {
			logger.Debug().Msgf("argo-url set to '%s'", argourl.String())
		}
	}
	if argoUsername == "" {
		logger.Fatal().Msg("argo-username flag was not set")
	} else {
		logger.Debug().Msgf("argo-username set to '%s'", argoUsername)
	}
	if argoPassword == "" {
		logger.Fatal().Msg("argo-password flag was not set")
	} else {
		logger.Debug().Msgf("argo-password is set: %s", strings.Repeat("*", len(argoPassword)))
	}

	if argourl == nil {
		logger.Fatal().Msg("argourl parsed to nil")
	}

	client, err := argocd.NewArgoClient(argourl, argoUsername, argoPassword, outfile, &logger)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	if err := client.FetchApplications(); err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msg("done")
}
