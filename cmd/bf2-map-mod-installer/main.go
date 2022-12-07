//go:generate goversioninfo

package main

import (
	"os"

	"github.com/cetteup/bf2-map-mod-installer/cmd/bf2-map-mod-installer/internal/gui"
	filerepo "github.com/cetteup/filerepo/pkg"
	"github.com/cetteup/joinme.click-launcher/pkg/registry_repository"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	registryRepository := registry_repository.New()
	fileRepository := filerepo.New()

	finder = software_finder.New(registryRepository, fileRepository)
}

var (
	finder *software_finder.SoftwareFinder
)

func main() {
	mw, err := gui.CreateMainWindow(finder)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create main window")
		os.Exit(1)
	}

	mw.Run()
}
