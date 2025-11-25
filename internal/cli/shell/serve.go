package shell

import (
	"cineguard/internal/api/rest"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var router *gin.Engine

// ViperFlagsServe defines a struct to hold the values of cobra CLI flags and use viper to populate them
type ViperFlagsServe struct {
	// cineguard server parameters
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
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
		initServer()

		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// server
	serveCmd.Flags().StringP("address", "a", "127.0.0.1", "Address to bind the server")
	serveCmd.Flags().StringP("port", "p", "8080", "Port to bind the server")
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
