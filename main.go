package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	username string
	hostname string
)

func calculateModStub() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}

func gomodinit(cmd *cobra.Command, args []string) {
	modStub, err := calculateModStub()
	if err != nil {
		log.Fatal(err)
	}

	modName := fmt.Sprintf("%s/%s/%s", hostname, username, modStub)
	log.Printf("Creating module under %s", modName)

	var out bytes.Buffer
	goCmd := exec.Command("go", "mod", "init", modName)
	goCmd.Stderr = &out

	if err = goCmd.Run(); err != nil {
		log.Fatalf("Error %v: %s\n", err, out.String())
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

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
