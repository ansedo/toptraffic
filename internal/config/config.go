package config

import (
	"context"
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/ansedo/toptraffic/internal/logger"
)

type Config struct {
	ServerPort      string `env:"SERVER_PORT,notEmpty"`
	AdvDomainString string `env:"ADV_DOMAINS,notEmpty"`
	AdvDomains      []string
}

func New(ctx context.Context) *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.FromCtx(ctx).Fatal("construct Config: parse flags", zap.Error(err))
	}
	flag.StringVar(
		&cfg.ServerPort,
		"p",
		cfg.ServerPort,
		`server port to listen on`,
	)
	flag.StringVar(
		&cfg.AdvDomainString,
		"d",
		cfg.AdvDomainString,
		`string of adv domains should contains from 1 to 10 domains comma separated`,
	)
	flag.Parse()

	if err := cfg.Validate(); err != nil {
		logger.FromCtx(ctx).Fatal("construct Config: parse flags", zap.Error(err))
	}
	return &cfg
}

func (c *Config) Validate() error {
	// Add colon to server port if it does not exist.
	if c.ServerPort[0] != ':' {
		c.ServerPort = ":" + c.ServerPort
	}

	// Get adv domains from `ADV_DOMAINS` string and check the count.
	c.AdvDomains = strings.Split(c.AdvDomainString, ",")
	if len(c.AdvDomains) < 1 || len(c.AdvDomains) > 10 {
		return errors.New(
			`env var "ADV_DOMAINS" should contains 1 to 10 domains (now: ` + strconv.Itoa(len(c.AdvDomains)) + `)`,
		)
	}
	return nil
}
