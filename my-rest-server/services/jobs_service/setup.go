package jobs_service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikestefanello/backlite"
	"github.com/mikestefanello/backlite/ui"
)

var Client *backlite.Client
var sqliteDb *sql.DB

func RunWebUi(port int) error {
	mux := http.DefaultServeMux
	h, err := ui.NewHandler(ui.Config{
		DB:       sqliteDb,
		BasePath: "/ui",
	})
	if err != nil {
		return err
	}

	h.Register(mux)
	err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mux)
	if err != nil {
		return err
	}

	return nil
}

// initialise my jobs queues to do some fancy sse stuff
func Setup() {
	// todo: add cancel support
	ctx := context.Background()

	var err error
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
}
