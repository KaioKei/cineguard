package shell

import (
	"cineguard/internal/api/rest"
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var db *sql.DB

var router *gin.Engine

var connexion *pgx.Conn

// DatabaseType helps to handle an enumeration of database types supported by cineguard
type DatabaseType string

const (
	Postgres DatabaseType = "postgres"
)

const DatabaseName = "cineguard"

// ViperFlagsServe defines a struct to hold the values of cobra CLI flags and use viper to populate them
type ViperFlagsServe struct {
	// cineguard server parameters
	Address  string             `mapstructure:"address"`
	Port     string             `mapstructure:"port"`
	Database ViperFlagsDatabase `mapstructure:"database"`
}

type ViperFlagsDatabase struct {
	Type     string `mapstructure:"type"`
	Address  string `mapstructure:"address"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Declare the viper CLI flag values buffer
var vprFlgsServe ViperFlagsServe

var serveCmd = &cobra.Command{
	Use:              "serve",
	Short:            "Start the cineguard server",
	Long:             "Start the cineguard server with a REST API with the given address and port",
	TraverseChildren: true, // ensure local flags do not spread to sub commands

	// Initialize and populate cobra CLI flags values with viper during the Persistent pre-run
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitViperSubCmdE(viper.GetViper(), cmd, &vprFlgsServe); err != nil {
			logrus.WithField("cobra-cmd", cmd.Use).WithError(err).Error("Error initializing Viper")
			return err
		}
		return nil
	},

	// Run the command
	Run: func(cmd *cobra.Command, args []string) {
		initDatabaseConnection()
		// defer close connexion to database at the end of the program
		defer func(connexion *pgx.Conn, ctx context.Context) {
			logrus.Info("Closing database connection")
			err := connexion.Close(ctx)
			if err != nil {
				logrus.Error("Error closing connection to database")
			}
		}(connexion, context.Background())

		initServer()

		runServer()

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// server
	serveCmd.Flags().StringP("address", "a", "127.0.0.1", "Address to bind the server")
	serveCmd.Flags().StringP("port", "p", "8080", "Port to bind the server")

	// database
	// the following flags are given by database.<key> because it concerns the nested part of the
	// yaml dedicated to the database configuration.
	// The dynamic reference to the nested database configuration in the yaml file is given by the
	// type structure 'ViperFlagsDatabase', nested inside the ViperFlagsServe type structure.
	serveCmd.Flags().String("database.type", "postgres", "Type of the cineguard database. Possible values: [postrgres]")
	serveCmd.Flags().String("database.address", "0.0.0.0", "IP of the cineguard database")
	serveCmd.Flags().String("database.port", "5432", "Port number of the cineguard database")
	serveCmd.Flags().String("database.username", "", "Username of the cineguard database")
	serveCmd.Flags().String("database.password", "", "Password of the cineguard database")
}

func initDatabaseConnection() {

	switch vprFlgsServe.Database.Type {
	case string(Postgres):
		// todo connexion to the postgressql database
		url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			vprFlgsServe.Database.Username,
			vprFlgsServe.Database.Password,
			vprFlgsServe.Database.Address,
			vprFlgsServe.Database.Port,
			DatabaseName,
		)
		var err error
		connexion, err = pgx.Connect(context.Background(), url)
		if err != nil {
			logrus.Fatalf("Error connecting to Postgres database '%s' at '%s:%s': %v", DatabaseName, vprFlgsServe.Database.Address, vprFlgsServe.Database.Port, err.Error())
		}

	default:
		logrus.Fatalf("Unsupported database type: %s", vprFlgsServe.Database.Type)
	}
}

func initServer() {

	// Initialize a Gin router using Default.
	router = gin.Default()

	// CORS config
	// CONFIGURE IT BEFORE ROUTES !
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// paths declarations
	v1 := router.Group("/api/v1")
	v1.GET("/health", rest.Health)
	//v1.POST("/ssla", restapi.PostSSLA)

}

func runServer() {
	logrus.Info("Starting Cineguard server")

	err := router.Run(fmt.Sprintf("%s:%s", vprFlgsServe.Address, vprFlgsServe.Port))
	if err != nil {
		logrus.Fatalf("Error starting server: %s", err.Error())
	}
}
