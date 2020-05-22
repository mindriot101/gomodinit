package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	username string
	hostname string
	verbose  bool
)

func calculateModStub() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}

func gomodinit(cmd *cobra.Command, args []string) {
	// Set up logging
	config := zap.NewDevelopmentConfig()
	if verbose {
		config.Level.SetLevel(zap.InfoLevel)
	} else {
		config.Level.SetLevel(zap.WarnLevel)
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	modStub, err := calculateModStub()
	if err != nil {
		sugar.Fatal(err)
	}

	modName := fmt.Sprintf("%s/%s/%s", hostname, username, modStub)
	sugar.Infof("Creating module under %s", modName)

	var out bytes.Buffer
	goCmd := exec.Command("go", "mod", "init", modName)
	goCmd.Stderr = &out

	if err = goCmd.Run(); err != nil {
		sugar.Fatalf("Error %v: %s\n", err, out.String())
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "gomodinit",
		Short: "",
		Long:  "",
		Run:   gomodinit,
	}

	rootCmd.PersistentFlags().StringVarP(&username, "user", "u", "mindriot101", "username to assign package to")
	rootCmd.PersistentFlags().StringVarP(&hostname, "host", "H", "github.com", "host to assign package to")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
