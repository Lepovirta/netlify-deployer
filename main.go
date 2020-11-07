package main

import (
	"context"
	"fmt"

	oapiclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/kelseyhightower/envconfig"
	netlify "github.com/netlify/open-api/go/porcelain"
	ooapicontext "github.com/netlify/open-api/go/porcelain/context"
	"github.com/sirupsen/logrus"
)

type config struct {
	AuthToken     string `envconfig:"auth_token" required:"true"`
	SiteID        string `envconfig:"site_id" required:"true"`
	Directory     string `required:"true"`
	Draft         bool   `default:"true"`
	DeployMessage string `default:""`
	LogLevel      string `default:"warn"`
	LogFormat     string `default:"text"`
}

func (c *config) readFromEnv() error {
	return envconfig.Process("netlify", c)
}

func main() {
	// Base logger
	logger := logrus.New()

	// Read config
	var c config
	if err := c.readFromEnv(); err != nil {
		logger.Fatalf("failed to read config: %s", err)
	}

	// Setup logging
	if err := setupLogging(&c, logger); err != nil {
		logger.Fatal(err)
	}

	// Netlify setup
	client := setupNetlifyClient(&c)
	ctx := setupContext(&c, logger)

	// Deploy site
	resp, err := client.DoDeploy(ctx, &netlify.DeployOptions{
		SiteID:  c.SiteID,
		Dir:     c.Directory,
		IsDraft: c.Draft,
		Title:   c.DeployMessage,
	}, nil)
	if err != nil {
		logger.Fatalf("failed to deploy site: %s", err)
	}

	// Print the site URL
	if resp.DeploySslURL != "" {
		fmt.Println(resp.DeploySslURL)
	} else if resp.DeployURL != "" {
		fmt.Println(resp.DeployURL)
	}
}

func setupLogging(c *config, logger *logrus.Logger) error {
	logLevel, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %s", err)
	}
	logger.SetLevel(logLevel)

	switch c.LogFormat {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.Warnf("invalid log format: %s", c.LogFormat)
	}
	return nil
}

func setupNetlifyClient(c *config) *netlify.Netlify {
	formats := strfmt.NewFormats()
	return netlify.NewHTTPClient(formats)
}

func setupContext(c *config, logger *logrus.Logger) ooapicontext.Context {
	ctx := ooapicontext.WithLogger(context.Background(), logger.WithFields(logrus.Fields{
		"source": "netlify",
	}))
	return ooapicontext.WithAuthInfo(ctx, oapiclient.BearerToken(c.AuthToken))
}
