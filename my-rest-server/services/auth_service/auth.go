package auth_service

import (
	"context"
	"fmt"
	"my-gmail-server/models"
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"my-gmail-server/settings"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func Session(name string) gin.HandlerFunc {
	return sessions.Sessions(name, store)
}

func CreateGmailServiceWithTokens(tokens oauth2.Token) (*gmail.Service, error) {
	ctx := context.Background()

	client := conf.Client(ctx, &tokens)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	return service, err
}

func RegisterUser(c *gin.Context) {
	appG := app.Gin{C: c}

	// r.Use(auth_service.Session(settings.AuthSettings.SessionName))
	// retrieve the session used here
	mysqlAuthSession := sessions.Default(c)

	// Check the state from which the req was generated matches the redirect uri's state from google
	retrievedState := mysqlAuthSession.Get(SessionState)
	if retrievedState != c.Query(SessionState) {
		appG.Response(http.StatusUnauthorized, e.STATUS_UNAUTHORIZED, "क्या आप वास्तव में Google से रीडायरेक्ट किए गए थे, मुझे तुम पर संदेह है ।")
		return
	}

	user := models.User{}
	ctx := context.Background()

	// get my juicy oauth tokens
	tok, err := conf.Exchange(ctx, c.Query("code"))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, "Oauth टोकन के लिए कोड का आदान-प्रदान करने में विफल")
		return
	}

	client := conf.Client(ctx, tok)
	peopleService, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, "Google People API की सेवा बनाने में विफलता")
		return
	}

	person, err := peopleService.People.Get("people/me").PersonFields("names,emailAddresses,photos").Do()

	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, "People API अनुरोध विफल रहा: उपयोगकर्ता की पहचान नहीं हो सकी ।")
		return
	}

	user.Token = *tok
	var myName, myEmail, myCoverPhoto string
	if len(person.Names) > 0 {
		myName = person.Names[0].DisplayName
	}

	if len(person.EmailAddresses) > 0 {
		// todo: i can also get name from person.EmailAddresses[0].DisplayName
		// can i fix the scopes
		myEmail = person.EmailAddresses[0].Value
	}

	if len(person.Photos) > 0 {
		myCoverPhoto = person.Photos[0].Url
	}

	user.FullName = myName
	user.Email = myEmail
	user.CoverPhoto = myCoverPhoto

	err = models.AddUser(user)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_USER_FAIL, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, fmt.Sprintf("Welcome %v", myName))
}

// setup oauth2 google config
func Setup() {
	SetupGoogleConfig(
		settings.AuthSettings.RedirectUrl,
		"./conf/credentials.json",
		[]string{
			// gmail api scopes
			gmail.GmailModifyScope,

			// people api scopes
			people.UserinfoProfileScope,
			people.UserinfoEmailScope,
		},
		[]byte(settings.AuthSettings.Secret), // load from app.ini config
	)
}
