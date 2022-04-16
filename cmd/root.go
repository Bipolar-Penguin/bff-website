package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
	"github.com/Bipolar-Penguin/bff-website/pkg/repository"
	"github.com/Bipolar-Penguin/bff-website/pkg/service"
	"github.com/Bipolar-Penguin/bff-website/pkg/transport/amqp"
	httptransport "github.com/Bipolar-Penguin/bff-website/pkg/transport/http"
)

const (
	deafaultHTTPPort int = 8000
)

var (
	cfgHTTPPort   int
	cfgMongoURL   string
	cfgRabbimqURL string
)

var rootCmd = &cobra.Command{
	Use:   "bff-website",
	Short: "bff-website",
	Long:  "Backend for frontend",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVar(&cfgHTTPPort, "http-port", deafaultHTTPPort, "http port to connect to")
	rootCmd.PersistentFlags().StringVar(&cfgMongoURL, "mongo-url", "", "mongo URL")
	rootCmd.PersistentFlags().StringVar(&cfgMongoURL, "rabbitmq-url", "", "rabbitmq URL")

	viper.BindPFlag("http-port", rootCmd.PersistentFlags().Lookup("http-port"))
	viper.BindPFlag("mongo-url", rootCmd.PersistentFlags().Lookup("mongo-url"))
	viper.BindPFlag("rabbitmq-url", rootCmd.PersistentFlags().Lookup("rabbitmq-url"))
}

func initConfig() {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
}

func run() { // nolint:funlen
	var err error

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger.Log("app", os.Args[0], "event", "starting")
	}
	// Parameters declaration & validation
	httpPort := viper.GetInt("http-port")
	if httpPort == 0 {
		logger.Log("error", "http-port argument was not provided")
		os.Exit(1)
	}
	mongoURL := viper.GetString("mongo-url")
	if mongoURL == "" {
		logger.Log("error", "mongo-url argument was not provided")
		os.Exit(1)
	}
	rabbitmqURL := viper.GetString("rabbitmq-url")
	if rabbitmqURL == "" {
		logger.Log("error", "rabbitmq-url argument was not provided")
		os.Exit(1)
	}

	// Broker declaration
	var amqpBroker *amqp.RabbitBroker
	{
		logger := log.With(logger, "module", "transport.amqp")

		amqpBroker = amqp.NewRabbitBroker(rabbitmqURL, logger)
		amqpBroker.PublishEvent(domain.Event{GUID: "foobar"})
	}

	// Repositories declaration
	var rep *repository.Repositories
	{
		logger := log.With(logger, "module", "repository")

		rep, err = repository.MakeRepositories(mongoURL, logger)
		if err != nil {
			os.Exit(1)
		}
	}

	// Services declaration
	var svc *service.Service
	{
		svc = service.NewService(rep, amqpBroker)
	}

	// HTTP server declaration
	var srv *httptransport.HTTPServer
	{
		logger := log.With(logger, "module", "http.transport")

		srv = httptransport.NewHttpServer(httpPort, logger, svc)
	}
	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Log("event", "got os shutdown signal")

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log("error", err)
			os.Exit(1)
		}

		close(idleConnsClosed)
		logger.Log("event", "server stopped")
	}()

	logger.Log("event", "server started")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Log("error", err)
		os.Exit(1)
	}

	<-idleConnsClosed

	logger.Log("event", "server exited normally")
}
