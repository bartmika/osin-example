package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openshift/osin"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"github.com/bartmika/osin-example/internal/controllers"
	repo "github.com/bartmika/osin-example/internal/repositories"
	"github.com/bartmika/osin-example/internal/session"
	"github.com/bartmika/osin-example/internal/utils"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the JSON API over HTTP",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runServeCmd()
	},
}

func doRunServe() {
	fmt.Println("Server started")
}

func runServeCmd() {
	// // Load up the S3.
	// key := os.Getenv("OSIN_AWS_S3_ACCESS_KEY")
	// secret := os.Getenv("OSIN_AWS_S3_SECRET_KEY")
	// endpoint := os.Getenv("OSIN_AWS_S3_ENDPOINT")
	// region := os.Getenv("OSIN_AWS_S3_REGION")
	// bucketName := os.Getenv("OSIN_AWS_S3_BUCKET_NAME")

	// Load up our database.
	db, err := utils.ConnectDB(
		databaseHost,
		databasePort,
		databaseUser,
		databasePassword,
		databaseName,
		"public",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load up our repositories.
	tr := repo.NewTenantRepo(db)
	ur := repo.NewUserRepo(db)

	// Open up our session handler, powered by redis and let's save the user
	// account with our ID
	sm := session.New()

	//
	// oAuth 2.0 Storage
	//

	oastore := controllers.NewOSINRedisStorage()
	oastore.CreateClient(&osin.DefaultClient{
		Id:          "1234",
		Secret:      "aabbccdd",
		RedirectUri: "http://localhost:8000/appauth/code",
	})

	//
	// oAuth 2.0
	//

	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
		osin.PASSWORD,
		osin.CLIENT_CREDENTIALS,
		osin.ASSERTION,
	}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true
	oas := osin.NewServer(sconfig, oastore)

	//
	// Background worker.
	//

	//TODO: Impl.

	//
	// Controller server.
	//

	c := &controllers.Controller{
		SecretSigningKeyBin: []byte(applicationSigningKey),
		OAuthServer:         oas,
		OAuthStorage:        oastore,
		TenantRepo:          tr,
		UserRepo:            ur,
		SessionManager:      sm,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", c.AttachMiddleware(c.HandleRequests))

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation via `https://github.com/rs/cors` for more options.
	handler := cors.AllowAll().Handler(mux)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "localhost", "8000"),
		Handler: handler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go runMainRuntimeLoop(srv)

	log.Print("Server Started")

	// Run the main loop blocking code.
	<-done

	stopMainRuntimeLoop(srv)
}

func runMainRuntimeLoop(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func stopMainRuntimeLoop(srv *http.Server) {
	log.Printf("Starting graceful shutdown now...")

	// Execute the graceful shutdown sub-routine which will terminate any
	// active connections and reject any new connections.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Printf("Graceful shutdown finished.")
	log.Print("Server Exited")
}
