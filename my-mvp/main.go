package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"google.golang.org/api/gmail/v1"
)

var me = "me"
var service = newGmailService()

var (
	flagListMaxSenders = flag.Int("list", 10, "used to list the max common senders")
	flagSenderEmail    = flag.String("trash", "", "flagOfEmailToBeTrashed")
	flagMaxPage        = flag.Int("page", -1, "-1 denotes infinity ie no page limit")
)

func trashMessagesFromSendersMailAndDisplaySpinner(sender string) {
	// Collect all the message ids
	var ids []string
	s := spinner.New(spinner.CharSets[57], 100*time.Millisecond) // Build our new spinner
	s.Color("green")
	s.Suffix = " finding mails"
	s.Start()

	service.Users.Messages.
		List(me).
		Q(fmt.Sprintf("from:%v", sender)).
		Pages(
			context.TODO(),
			func(lmr *gmail.ListMessagesResponse) error {
				for _, message := range lmr.Messages {
					ids = append(ids, message.Id)
				}

				return nil
			},
		)

	s.Color("red")

	for idx, id := range ids {
		_, err := service.Users.Messages.Trash(me, id).Do()
		suffix := fmt.Sprintf(" Trashing %d of %d\n", idx+1, len(ids))

		if err != nil {
			suffix = fmt.Sprintf(" failed to trash %q\n", id)
		}

		s.Suffix = suffix
	}
	s.Stop()
	fmt.Println("All trashed!")
}

func listCommonSender(limit int) {
	fmt.Printf("Listing %d common senders\n", limit)

	limit = min(safeMap.Length(), limit)

	// Lipgloss table renderer
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render
	t := table.New()

	type EmailMessageidsPair struct {
		email  string
		length int
	}

	var container []EmailMessageidsPair

	for key, value := range safeMap.GetMapCopy() {
		container = append(
			container,
			EmailMessageidsPair{email: key, length: len(value)},
		)
	}

	// Sort the value based upon length
	sort.Slice(container, func(i, j int) bool {
		return container[i].length > container[j].length
	})

	for idx, pair := range container[:limit] {
		t.Row(
			s(fmt.Sprint(idx+1)),
			pair.email,
			fmt.Sprint(pair.length),
		)
	}

	fmt.Println(t.Render())

}

func main() {
	flag.Parse()

	if *flagSenderEmail != "" {
		fmt.Println("--trash", *flagSenderEmail)
		trashMessagesFromSendersMailAndDisplaySpinner(*flagSenderEmail)
		return
	}

	maxPage := *flagMaxPage

	listSendersConcurrentlyAndDisplayUi(maxPage)
	listCommonSender(*flagListMaxSenders)
}
