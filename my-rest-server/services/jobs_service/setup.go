package jobs_service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"my-gmail-server/settings"
	"net/http"

	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikestefanello/backlite"
	"github.com/mikestefanello/backlite/ui"
)

var Client *backlite.Client
var sqliteDb *sql.DB

func RunWebUi(port int) {
	mux := http.DefaultServeMux
	h, err := ui.NewHandler(ui.Config{
		DB:       sqliteDb,
		BasePath: "/ui",
	})
	if err != nil {
		log.Fatal(err)
	}

	h.Register(mux)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(err)
	}
}

// initialise my jobs queues to do some fancy sse stuff
func Setup() {
	var err error

	// todo: add cancel support
	ctx := context.Background()

	sqliteDb, err = sql.Open("sqlite3", "data.db?_journal=WAL&_timeout=5000")
	if err != nil {
		log.Fatal(err)
	}

	Client, err = backlite.NewClient(backlite.ClientConfig{
		DB:              sqliteDb,
		Logger:          slog.Default(),
		ReleaseAfter:    10 * time.Minute,
		NumWorkers:      10,
		CleanupInterval: time.Hour,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err = Client.Install(); err != nil {
		log.Fatal(err)
	}

	for _, q := range []backlite.Queue{backlite.NewQueue(ScanJobWork), backlite.NewQueue(TrashJobWork)} {
		Client.Register(q)
	}

	Client.Start(ctx)

	// Start the web server on port 9000 to monitor
	go RunWebUi(settings.ServerSetting.JobsPort)
}
