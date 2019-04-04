package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/video-sharing-platform-thingy/backend/storer"
	"github.com/video-sharing-platform-thingy/backend/util"
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss/confirm"
	"github.com/volatiletech/authboss/defaults"
	"github.com/volatiletech/authboss/otp/twofactor"
	"github.com/volatiletech/authboss/otp/twofactor/totp2fa"
	"github.com/volatiletech/authboss/remember"

	_ "github.com/volatiletech/authboss/auth"
	_ "github.com/volatiletech/authboss/logout"
	_ "github.com/volatiletech/authboss/recover"
	_ "github.com/volatiletech/authboss/register"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	abclientstate "github.com/volatiletech/authboss-clientstate"
	abrenderer "github.com/volatiletech/authboss-renderer"
	aboauth "github.com/volatiletech/authboss/oauth2"
)

var (
	ab        = authboss.New()
	schemaDec = schema.NewDecoder()
	database  = storer.NewDBStorer()

	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer

	loadedConfig config
)

// init just sets up the config variable and
// connects to the database when the program
// starts.
func init() {
	database.Connect()
	loadedConfig = loadConfig()
}

// setupAuth sets up authentication with authboss.
func setupAuth() {
	// Set the base url for authboss.
	ab.Config.Paths.RootURL = loadedConfig.BaseURL

	// Set various request methods so you can actually
	// specify JSON bodies, because authboss is annoying
	// and expects JSON bodies even for GET requests.
	ab.Config.Modules.LogoutMethod = http.MethodPost
	ab.Config.Modules.ConfirmMethod = http.MethodPost
	ab.Config.Modules.MailRouteMethod = http.MethodPost
	ab.Config.Modules.TwoFactorEmailAuthRequired = false

	// Set up session and cookie storage mechanisms.
	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	// Set the renderer to a JSON renderer.
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

	// We render mail with the authboss-renderer but we use a LogMailer
	// which simply sends the e-mail to stdout.
	// TODO: Change this to something better.
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")

	// TOTP2FAIssuer is the name of the issuer we use for totp 2fa.
	ab.Config.Modules.TOTP2FAIssuer = loadedConfig.Name

	// This automatically sets up a bunch of useful
	// default configuration so it doesn't all have
	// to be written out.
	defaults.SetCore(&ab.Config, true, false)

	// Set up some validation options and such.
	// TODO: Figure out what fields we need.
	// FIXME: Add a better regex and more secure options.
	emailRule := defaults.Rules{
		FieldName:  "email",
		Required:   true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]{1,}`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password",
		Required:  true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName:       "name",
		Required:        true,
		AllowWhitespace: true,
	}

	// Set up the authboss request body reader.
	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: true,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": []string{"email", "name", "password"},
		},
	}

	// Set up totp 2fa.
	// TODO: Do we want other types of 2fa?
	twofaRecovery := &twofactor.Recovery{Authboss: ab}
	err := twofaRecovery.Setup()
	util.CheckError(err)
	totp := &totp2fa.TOTP{Authboss: ab}
	err = totp.Setup()
	util.CheckError(err)

	// Set up various oauth2 providers.
	ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{
		"google": authboss.OAuth2Provider{
			OAuth2Config: &oauth2.Config{
				ClientID:     loadedConfig.Oauth.Google.ClientID,
				ClientSecret: loadedConfig.Oauth.Google.ClientSecret,
				Scopes:       []string{`profile`, `email`},
				Endpoint:     google.Endpoint,
			},
			FindUserDetails: aboauth.GoogleUserDetails,
		},
	}

	// Actually initialize authboss.
	// I actually forgot to do this for a while,
	// and was confused why nothing was working.
	err = ab.Init()
	util.CheckError(err)
}

func main() {
	// Initialize session and cookie store keys.
	// FIXME: THESE ARE INSECURE, CHANGE THEM!
	cookieStoreKey, err := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	util.CheckError(err)
	sessionStoreKey, err := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	util.CheckError(err)

	// Initialize the actual session and cookie stores.
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(loadedConfig.SessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false

	// Set up authboss.
	setupAuth()

	// Create out router.
	schemaDec.IgnoreUnknownKeys(true)
	mux := chi.NewRouter()

	// Set up middlewares.
	// - logger will log requests and some debug info
	// - LoadClientStateMiddleware makes session/cookie stuff work
	// - remember middleware logs users in if they have a remember token
	// FIXME: Add CSRF protection.
	mux.Use(logger, ab.LoadClientStateMiddleware, remember.Middleware(ab))

	// Set up authed routes, currently just a test route.
	// FIXME: Actually add needed routes.
	mux.Group(func(mux chi.Router) {
		// Set up group-specific middlewares.
		// - authboss middleware2 prevents unauthed users from accessing the routes.
		// - confirm middleware makes sure user accounts are confirmed.
		// TODO: What is full auth in this context?
		mw2 := authboss.Middleware2(ab, authboss.RequireFullAuth, authboss.RespondUnauthorized)
		mux.Use(mw2, confirm.Middleware(ab))

		// Set up the actual routes.
		mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ping":"pong"}`)
		})
	})

	// Set up auth routes (not authed routes).
	mux.Group(func(mux chi.Router) {
		mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})

	// Set up unauthed routes, currently just a test route.
	// FIXME: Actually add needed routes.
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"home":true}`)
	})

	// Listen on a port.
	log.Printf("Listening on localhost:%s\n", loadedConfig.Port)
	err = http.ListenAndServe("localhost:"+loadedConfig.Port, mux)
	util.CheckError(err)
}
