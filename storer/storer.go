package storer

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/volatiletech/authboss"
	aboauth "github.com/volatiletech/authboss/oauth2"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // The sqlite adapter for gorm.
)

// DBStorer stores users in a database.
type DBStorer struct {
	DB *gorm.DB
}

// NewDBStorer can create a new database storer.
func NewDBStorer() *DBStorer {
	return &DBStorer{}
}

// Connect connects to the database.
func (m *DBStorer) Connect() error {
	db, err := gorm.Open("sqlite3", "vspt.db")
	if err != nil {
		return err
	}

	db.LogMode(true)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})

	m.DB = db
	log.Println("Connected to database")
	return nil
}

// Save saves a user to the database.
func (m *DBStorer) Save(ctx context.Context, user authboss.User) error {
	u := user.(*User)

	if err := m.DB.Save(u).Error; err != nil {
		return err
	}

	log.Println("Saved user:", u.Name)
	return nil
}

// Load loads a user from the database.
func (m *DBStorer) Load(ctx context.Context, key string) (user authboss.User, err error) {
	var u User

	// Check to see if the key is an oauth2 pid.
	provider, uid, err := authboss.ParseOAuth2PID(key)
	if err == nil {
		if err = m.DB.Where(&User{OAuth2Provider: provider, OAuth2UID: uid}).First(&u).Error; err != nil {
			return nil, authboss.ErrUserNotFound
		}

		log.Println("Loaded oauth2 user:", u.Name)
		return &u, nil
	}

	if err = m.DB.Where(&User{Email: key}).First(&u).Error; err != nil {
		return nil, authboss.ErrUserNotFound
	}

	log.Println("Loaded user:", u.Name)
	return &u, nil
}

// New will just instantiate a new user.
func (m *DBStorer) New(ctx context.Context) authboss.User {
	return &User{}
}

// Create creates a user in the database.
func (m *DBStorer) Create(ctx context.Context, user authboss.User) error {
	u := user.(*User)

	if err := m.DB.Where(&User{Email: u.Email}).First(&u).Error; err == nil {
		return authboss.ErrUserFound
	}
	if err := m.DB.Create(u).Error; err != nil {
		return err
	}
	log.Println("Created new user:", u.Name)
	return nil
}

// LoadByConfirmSelector looks a user up by a confirmation token.
func (m *DBStorer) LoadByConfirmSelector(ctx context.Context, selector string) (user authboss.ConfirmableUser, err error) {
	var u User
	if err := m.DB.Where(&User{ConfirmSelector: selector}).First(&u).Error; err != nil {
		return nil, authboss.ErrUserNotFound
	}

	log.Println("Loaded user by confirm token:", u.Name)
	return &u, nil
}

// LoadByRecoverSelector looks a user up by a recover selector.
func (m *DBStorer) LoadByRecoverSelector(ctx context.Context, selector string) (user authboss.RecoverableUser, err error) {
	var u User
	if err := m.DB.Where(&User{RecoverSelector: selector}).First(&u).Error; err != nil {
		return nil, authboss.ErrUserNotFound
	}

	log.Println("Loaded user by recover selector:", u.Name)
	return &u, nil
}

// AddRememberToken adds a remember token to a user.
func (m *DBStorer) AddRememberToken(ctx context.Context, pid, token string) error {
	if err := m.DB.Create(&Token{PID: pid, Value: token}).Error; err != nil {
		return err
	}

	log.Printf("Added rm token to %s: %s\n", pid, token)
	return nil
}

// DelRememberTokens removes all tokens for the given pid.
func (m *DBStorer) DelRememberTokens(ctx context.Context, pid string) error {
	if err := m.DB.Delete(&Token{}, &Token{PID: pid}).Error; err != nil {
		return err
	}

	log.Printf("Deleted rm tokens from: %s\n", pid)
	return nil
}

// UseRememberToken finds the pid-token pair and deletes it.
// If the token could not be found it returns ErrTokenNotFound.
func (m *DBStorer) UseRememberToken(ctx context.Context, pid, token string) error {
	q := m.DB.Delete(&Token{}, &Token{PID: pid, Value: token})
	if q.Error != nil {
		return q.Error
	}
	if q.RowsAffected < 1 {
		return authboss.ErrTokenNotFound
	}

	log.Printf("Used remember for %s: %s\n", pid, token)
	return nil
}

// NewFromOAuth2 creates an oauth2 user (not in the database, just a blank one to be saved later).
func (m *DBStorer) NewFromOAuth2(ctx context.Context, provider string, details map[string]string) (authboss.OAuth2User, error) {
	switch provider {
	case "google":
		email := details[aboauth.OAuth2Email]

		var user *User
		if err := m.DB.Where(&User{Email: email}).First(&user).Error; err == nil {
			return user, nil
		}

		// Google OAuth2 doesn't allow us to fetch real name without more complicated API calls
		// in order to do this properly in your own app, look at replacing the authboss oauth2.GoogleUserDetails
		// method with something more thorough.
		user.Name = "Unknown"
		user.Email = details[aboauth.OAuth2Email]
		user.OAuth2UID = details[aboauth.OAuth2UID]
		user.Confirmed = true

		return user, nil
	}

	return nil, errors.Errorf("unknown provider %s", provider)
}

// SaveOAuth2 user
func (m *DBStorer) SaveOAuth2(ctx context.Context, user authboss.OAuth2User) error {
	u := user.(*User)
	if err := m.DB.Create(u).Error; err != nil {
		return err
	}

	log.Println("Saved oauth2 user:", u.Name)
	return nil
}
