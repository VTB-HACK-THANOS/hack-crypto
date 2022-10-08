package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/VTB-HACK-THANOS/hack-crypto/api"
	"github.com/VTB-HACK-THANOS/hack-crypto/services/auth"
	crytowallet "github.com/VTB-HACK-THANOS/hack-crypto/services/cryto-wallet"
	questions "github.com/VTB-HACK-THANOS/hack-crypto/services/questions"
	userManagement "github.com/VTB-HACK-THANOS/hack-crypto/services/user-management"
	"github.com/VTB-HACK-THANOS/hack-crypto/storage"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

var appName = "emlvid-back"

func main() {
	cfg := configuration{}
	yamlFeeder := feeder.Yaml{Path: "./config/config.yml"}
	if err := config.
		New().
		AddFeeder(yamlFeeder).
		AddStruct(&cfg).
		Feed(); err != nil {
		panic(err)
	}

	store, err := storage.New(
		cfg.Postgres.Addr,
		cfg.Postgres.Name,
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.Pool,
		cfg.App.Environment,
		appName)
	if err != nil {
		panic(err)
	}

	ver, err := store.Migrate()
	if err != nil {
		panic(err)
	}
	log.Printf("migrated to version %v\n", ver)

	//Services
	authSvc, err := auth.New(store)
	if err != nil {
		panic(err)
	}

	question, err := questions.New(store)
	if err != nil {
		panic(err)
	}

	cryptoWallet, err := crytowallet.New(cfg.CryptoWallet.Endpoint)
	if err != nil {
		panic(err)
	}

	userManagementSvc, err := userManagement.New(store, cryptoWallet)
	if err != nil {
		panic(err)
	}

	api := api.New(authSvc, userManagementSvc, question, cfg.RateLimiter.Rate, cfg.RateLimiter.ExpiresIn, cfg.RateLimiter.Burst)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	go func() {
		if err := api.Start(ctx, cfg.App.Addr); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

}

type configuration struct {
	App struct {
		Environment string
		Addr        string
	}
	Postgres struct {
		Addr string
		Name string
		User string
		Pass string
		Pool int
	}
	CryptoWallet struct {
		Endpoint string
	} `yaml:"cryptoWallet"`
	RateLimiter struct {
		Rate      time.Duration
		Burst     int
		ExpiresIn time.Duration
	}
}
