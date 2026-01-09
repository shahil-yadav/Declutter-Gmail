package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/progress"
	"google.golang.org/api/gmail/v1"
)

type SafeMap struct {
	mu sync.Mutex
	v  map[string][]string
}

var safeMap = SafeMap{v: make(map[string][]string)}

func (c *SafeMap) Length() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.v)
}

func (c *SafeMap) GetMapCopy() map[string][]string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.v
}

func (c *SafeMap) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] = append(c.v[key], value)
}

func (c *SafeMap) Get(key string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Lock so only one goroutine at a time can access the map c.v.
	return c.v[key]
}

// Utils
func parseSenderAddressByMsgId(id string) (string, error) {
	email, err := service.Users.Messages.Get(me, id).Format("metadata").Do()
	if err != nil {
		log.Fatalf("Failed to retrieve the message ~ %v", id)
	}

	var fromHeader string

	// get the `from` header in email
	for _, header := range email.Payload.Headers {
		if header.Name == "From" {
			fromHeader = header.Value
			break // no need to go further in search of headers, i got what i needed
		}
	}

	// parse the `from` header
	parsed, err := mail.ParseAddress(fromHeader)
	if err != nil {
		return "", fmt.Errorf("failed to parse email from\n> %v", email.Payload.Headers)
	}

	return parsed.Address, nil
}

// scanner tracker
func trackScanning(pw progress.Writer, emailIds []string) {
	// create a tracker
	tracker := progress.Tracker{
		Message: fmt.Sprintf("Fetching %d emails", len(emailIds)),
		Total:   int64(len(emailIds)),
		Units:   progress.UnitsDefault,
	}

	// append the tracker to the pw
	pw.AppendTracker(&tracker)

	// Do expensive API calls
	for _, id := range emailIds {
		address, err := parseSenderAddressByMsgId(id)
		if err != nil {
			pw.SetPinnedMessages(
				fmt.Sprintf("%v", err),
			)
		} else {
			// store the ids mapped to email address as key
			safeMap.Set(address, id)
		}

		// update tracker
		tracker.Increment(1)
	}

	// out of the loop means we are done
	tracker.MarkAsDone()
}

// progress writer to manage trackers
func NewProgressWriter() progress.Writer {
	pw := progress.NewWriter()
	pw.SetNumTrackersExpected(10)
	pw.SetAutoStop(true)
	pw.SetMessageLength(24)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerLength(25)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	pw.Style().Visibility.ETAOverall = true
	pw.Style().Visibility.Percentage = true
	pw.Style().Visibility.Pinned = true
	pw.Style().Visibility.Speed = true
	pw.Style().Visibility.TrackerOverall = true
	pw.Style().Visibility.Value = true

	// call Render() in async mode; yes we don't have any trackers at the moment
	go pw.Render()

	return pw
}

// this is strictly binded to progress.writer
func listSendersConcurrentlyAndDisplayUi(maxPage int) {
	pw := NewProgressWriter()

	// create an arrays of message ids
	var ids []string

	s := spinner.New(spinner.CharSets[57], 100*time.Millisecond) // Build our new spinner
	s.Prefix = ""
	s.Suffix = " Scanning the mailbox" // Append text after the spinner
	s.Start()

	page := 1
	s.Color("green")

	//[Expensive] collect all ids in container
	service.Users.Messages.List(me).Pages(
		context.TODO(),
		func(lmr *gmail.ListMessagesResponse) error {
			for _, message := range lmr.Messages {
				ids = append(ids, message.Id)
				s.Suffix = fmt.Sprintln(
					" ",
					"📬 Scanning your mailbox... This might take a while! ☕️",
					"\n   Hang tight and grab a coffee while we work our magic! 🚀",
					"\n   Fetching mails in pages of 100 — total count unknown.",
					fmt.Sprintf("\n   currently reading %d page", page),
				)
			}

			if maxPage != -1 && page == maxPage {
				// Returned error stops the fetching immediately
				return errors.New("reached max page limit")
			}

			page++

			// denotes successfull fetch and collection
			return nil
		},
	)

	s.Stop()

	// create a batches that contains 5 partitions running in independent thread
	var batches [][]string

	bLen := 10
	bRng := len(ids) / bLen

	fmt.Printf("Initiating parsing in batches of %d \n\n", bLen)

	start, end := 0, bRng
	for range bLen {
		batches = append(batches, ids[start:end])
		// fmt.Printf("batch %d [%d:%d] prepared\n", i+1, start, end)
		start = end
		end = start + bRng
	}

	if len(ids)%bLen != 0 {
		batches = append(batches, ids[start:])
	}

	// assert if batches total length is equal to ids arrays
	totalLenOfBatch := 0
	for _, batch := range batches {
		totalLenOfBatch += len(batch)
	}

	if totalLenOfBatch != len(ids) {
		log.Fatalln("creation of batch logic is incorrect")
	}

	//  make a concurrent call in the batches
	for _, batch := range batches {
		go trackScanning(pw, batch)
	}

	for pw.IsRenderInProgress() {
	}

	fmt.Println("All done!")
}
