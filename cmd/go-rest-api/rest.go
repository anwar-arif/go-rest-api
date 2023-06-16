package main

import (
	"context"
	"go-rest-api/api/controller"
	"go-rest-api/infra/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"go-rest-api/api"
	"go-rest-api/config"
	infraMongo "go-rest-api/infra/mongo"
	infraRedis "go-rest-api/infra/redis"
	infraSentry "go-rest-api/infra/sentry"
	"go-rest-api/logger"
)

// srvCmd is the serve sub command to start the api server
var srvCmd = &cobra.Command{
	Use:   "serve-rest",
	Short: "serve serves the api server",
	RunE:  serve,
}

func setEnvPath() {
	if env == DevEnv {
		envPath = "dev.env"
	} else if env == TestEnv {
		envPath = "test.env"
	} else if env == LocalEnv {
		envPath = "local.env"
	} else if env == ProdEnv {
		envPath = "prod.env"
	} else {
		envPath = ""
	}
}

func init() {
	srvCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "", "config file path")
	srvCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "environment type: dev, prod")
}

func serve(cmd *cobra.Command, args []string) error {
	setEnvPath()
	cfgApp := config.GetApp(cfgPath)
	cfgMongo := config.GetMongo(cfgPath, env, envPath)
	cfgSentry := config.GetSentry(cfgPath)
	cfgDBTable := config.GetTable(cfgPath)
	cfgRedis := config.GetRedis(cfgPath, envPath)

	ctx := context.Background()

	lgr := logger.DefaultOutStructLogger

	mgo, err := infraMongo.New(ctx, cfgMongo)
	if err != nil {
		return err
	}
	defer mgo.Close(ctx)

	rds, err := infraRedis.New(ctx, cfgRedis, lgr)
	if err != nil {
		return err
	}
	defer rds.Close()

	db := infra.NewDB(mgo, rds)

	err = infraSentry.NewInit(cfgSentry.URL)
	if err != nil {
		return err
	}

	api.SetLogger(logger.DefaultOutLogger)

	errChan := make(chan error)
	go func() {
		if err := startHealthServer(cfgApp, db); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := startApiServer(cfgApp, cfgDBTable, db, lgr); err != nil {
			errChan <- err
		}
	}()
	return <-errChan

}

func startHealthServer(cfg *config.Application, db *infra.DB) error {
	log.Println("startHealthServer")
	sc := controller.NewSystemController(db)
	r := chi.NewMux()
	r.Mount("/system/v1", api.NewSystemRouter(sc))

	srvr := http.Server{
		Addr:    getAddressFromHostAndPort(cfg.Host, cfg.SystemServerPort),
		Handler: r,
		//ErrorLog: logger.DefaultErrLogger,
		//WriteTimeout: cfg.WriteTimeout,
		//ReadTimeout:  cfg.ReadTimeout,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	graceful := func() error {
		log.Println("To shutdown immediately press again")

		return nil
	}

	errCh := make(chan error)
	forced := func() error {
		log.Println("Shutting down server forcefully")
		return nil
	}
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}

	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()

	return <-errCh
}

func startApiServer(cfgApp *config.Application, cfgDBTable *config.Table, db *infra.DB, lgr logger.StructLogger) error {

	r := chi.NewMux()
	r.Mount("/api/v1", api.NewApiRouter(cfgDBTable, db, lgr))
	r.Mount("/", api.NewPingRouter(lgr))

	srvr := http.Server{
		Addr:    getAddressFromHostAndPort(cfgApp.Host, cfgApp.Port),
		Handler: r,
		//ErrorLog: logger.DefaultErrLogger,
		//WriteTimeout: cfg.WriteTimeout,
		//ReadTimeout:  cfg.ReadTimeout,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return ManageServer(&srvr, 30*time.Second)
}

func ManageServer(srvr *http.Server, gracePeriod time.Duration) error {
	errCh := make(chan error)

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, os.Interrupt}

	graceful := func() error {
		log.Println("Shutting down server gracefully with in", gracePeriod)
		log.Println("To shutdown immediately press again")

		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		return srvr.Shutdown(ctx)
	}

	forced := func() error {
		log.Println("Shutting down server forcefully")
		return srvr.Close()
	}

	go func() {
		log.Println("Starting server on", srvr.Addr)
		if err := srvr.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()

	return <-errCh
}

// HandleSignals listen on the registered signals and fires the gracefulHandler for the
// first signal and the forceHandler (if any) for the next this function blocks and
// return any error that returned by any of the handlers first
func HandleSignals(sigs []os.Signal, gracefulHandler, forceHandler func() error) error {
	sigCh := make(chan os.Signal)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, sigs...)
	defer signal.Stop(sigCh)

	grace := true
	for {
		select {
		case err := <-errCh:
			return err
		case <-sigCh:
			if grace {
				grace = false
				go func() {
					errCh <- gracefulHandler()
				}()
			} else if forceHandler != nil {
				err := forceHandler()
				errCh <- err
			}
		}
	}
}

func getAddressFromHostAndPort(host string, port int) string {
	addr := host
	if port != 0 {
		addr = addr + ":" + strconv.Itoa(port)
	}
	return addr
}
