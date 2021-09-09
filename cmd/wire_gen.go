// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/navidrome/navidrome/core"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/core/agents/lastfm"
	"github.com/navidrome/navidrome/core/scrobbler"
	"github.com/navidrome/navidrome/core/transcoder"
	"github.com/navidrome/navidrome/db"
	"github.com/navidrome/navidrome/persistence"
	"github.com/navidrome/navidrome/scanner"
	"github.com/navidrome/navidrome/server"
	"github.com/navidrome/navidrome/server/events"
	"github.com/navidrome/navidrome/server/nativeapi"
	"github.com/navidrome/navidrome/server/subsonic"
	"sync"
)

// Injectors from wire_injectors.go:

func CreateServer(musicFolder string) *server.Server {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	serverServer := server.New(dataStore)
	return serverServer
}

func CreateNativeAPIRouter() *nativeapi.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	broker := events.GetBroker()
	share := core.NewShare(dataStore)
	router := nativeapi.New(dataStore, broker, share)
	return router
}

func CreateSubsonicAPIRouter() *subsonic.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	artworkCache := core.GetImageCache()
	artwork := core.NewArtwork(dataStore, artworkCache)
	transcoderTranscoder := transcoder.New()
	transcodingCache := core.GetTranscodingCache()
	mediaStreamer := core.NewMediaStreamer(dataStore, transcoderTranscoder, transcodingCache)
	archiver := core.NewArchiver(dataStore)
	players := core.NewPlayers(dataStore)
	agentsAgents := agents.New(dataStore)
	externalMetadata := core.NewExternalMetadata(dataStore, agentsAgents)
	scanner := GetScanner()
	broker := events.GetBroker()
	playTracker := scrobbler.GetPlayTracker(dataStore, broker)
	router := subsonic.New(dataStore, artwork, mediaStreamer, archiver, players, externalMetadata, scanner, broker, playTracker)
	return router
}

func CreateLastFMRouter() *lastfm.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	router := lastfm.NewRouter(dataStore)
	return router
}

func createScanner() scanner.Scanner {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	artworkCache := core.GetImageCache()
	artwork := core.NewArtwork(dataStore, artworkCache)
	cacheWarmer := core.NewCacheWarmer(artwork, artworkCache)
	broker := events.GetBroker()
	scannerScanner := scanner.New(dataStore, cacheWarmer, broker)
	return scannerScanner
}

// wire_injectors.go:

var allProviders = wire.NewSet(core.Set, subsonic.New, nativeapi.New, persistence.New, lastfm.NewRouter, events.GetBroker, db.Db)

// Scanner must be a Singleton
var (
	onceScanner     sync.Once
	scannerInstance scanner.Scanner
)

func GetScanner() scanner.Scanner {
	onceScanner.Do(func() {
		scannerInstance = createScanner()
	})
	return scannerInstance
}
