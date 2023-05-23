package managers

import (
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/mocks"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
)

// Subset of database.User that allows to check if some keys match using StructMatcher
// nil == ignore key
type userMatcher struct {
	ID            *uint
	Email         *string
	PasswordHash  *string
	FirstName     *string
	LastName      *string
	EmailVerified *bool
}

// Subset of database.Token that allows to check if some keys match using StructMatcher
// nil == ignore key
type tokenMatcher struct {
	OwnerId      *uint
	TokenId      *string
	Expires      *time.Time
	TokenPurpose *TokenPurposeEnum
	Used         *bool
	Recalled     *bool
}

func getAuthManager(ctrl *gomock.Controller) (*AuthManagerImpl, error) {
	return &AuthManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: NewMockGormDB(ctrl),
		},
		NewMockEmailService(ctrl),
		NewMockTokenService(ctrl),
	}, nil
}

// Returns example user model
func getExampleUser() User {
	hash, err := bcrypt.GenerateFromPassword([]byte("zaq1@WSX"), 10)
	if err != nil {
		panic(err)
	}
	return User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: false,
			},
		},
		PublicId:      "Es3Aepo7",
		FirstName:     "test_first_name",
		LastName:      "test_last_name",
		Email:         "test@example.com",
		PasswordHash:  string(hash),
		EmailVerified: false,
	}
}

// Mocks user in the database
func mockExampleUser(database *MockGormDB) User {
	user := getExampleUser()
	database.
		EXPECT().
		First(gomock.Any(), &StructMatcher{userMatcher{
			Email: Ptr("test@example.com"),
		}}).
		DoAndReturn(func(arg *User, conds ...interface{}) GormDB {
			database.
				EXPECT().
				GetError().
				AnyTimes().
				Return(nil)
			*arg = user
			return database
		})
	return user
}

// Retuns example Token model
func createExampleToken(tokenId string, tokenPurpose TokenPurposeEnum) Token {
	hash, err := bcrypt.GenerateFromPassword([]byte("test_hash"), 10)
	if err != nil {
		panic(err)
	}
	token := Token{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: false,
			},
		},
		TokenId:      tokenId,
		TokenHash:    string(hash),
		Expires:      time.Now().Add(24 * time.Hour),
		TokenPurpose: tokenPurpose,
		Used:         false,
		Recalled:     false,
	}
	return token
}

// Mocks example Token model in the database
func mockExampleUserEmailVerificationToken(database *MockGormDB) Token {
	token := createExampleToken("test_email", TokenPurposeEmail)
	database.
		EXPECT().
		Find(gomock.Any(), &StructMatcher{tokenMatcher{
			TokenId: Ptr("test_email"),
		}}).
		Do(func(arg *Token) GormDB {
			*arg = token
			return database
		})
	return token
}

// Returns example session token for example user
func getExampleUserLogin() Token {
	hash, err := bcrypt.GenerateFromPassword([]byte("test_hash"), 10)
	if err != nil {
		panic(err)
	}
	token := Token{
		Model: gorm.Model{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: false,
			},
		},
		TokenId:      "test_login",
		TokenHash:    string(hash),
		Expires:      time.Now().Add(time.Hour),
		TokenPurpose: TokenPurposeSession,
		Used:         true,
		Recalled:     false,
	}
	return token
}

// Mocks example session token in the database
func mockExampleUserLogin(tokenService *MockTokenService) Token {
	user := getExampleUser()
	token := getExampleUserLogin()
	tokenService.
		EXPECT().
		Check("test_login", "test_hash").
		Return(&user, &token, nil)
	return token
}

// Mocks transaction begin
func mockBegin(db GormDB) {
	db.(*MockGormDB).
		EXPECT().
		Begin().
		Return(db)
}

// Mocks transaction rollback
func mockRollback(db GormDB) {
	db.(*MockGormDB).
		EXPECT().
		Rollback().
		Return(db)
}

// Mocks transaction commit
func mockCommit(db GormDB) {
	db.(*MockGormDB).
		EXPECT().
		Commit().
		DoAndReturn(returnError0(db, nil))
}

// Utility functions used to return an error in DoAndReturn, each for different amount of method arguments
func returnError0(db GormDB, err error) func() GormDB {
	return (func() GormDB {
		db.(*MockGormDB).
			EXPECT().
			GetError().
			Return(err)
		return db
	})
}

func returnError1(db GormDB, err error) func(arg any) GormDB {
	return (func(arg any) GormDB {
		db.(*MockGormDB).
			EXPECT().
			GetError().
			Return(err)
		return db
	})
}

func returnError2(db GormDB, err error) func(arg any, arg2 any) GormDB {
	return (func(arg any, arg2 any) GormDB {
		db.(*MockGormDB).
			EXPECT().
			GetError().
			Return(err)
		return db
	})
}

// Tests

// Tests if AuthManagerImpl.Create works correctly on the happy path
func TestAuthManagerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	mainUserMatcher := &StructMatcher{userMatcher{
		Email:         Ptr("test@example.com"),
		FirstName:     Ptr("first"),
		LastName:      Ptr("last"),
		EmailVerified: Ptr(false),
	}}

	mockBegin(db)

	db.(*MockGormDB).
		EXPECT().
		First(gomock.Any(), &StructMatcher{userMatcher{
			Email: Ptr("test@example.com"),
		}}).
		DoAndReturn(returnError2(db, gorm.ErrRecordNotFound))

	db.(*MockGormDB).
		EXPECT().
		Create(mainUserMatcher).
		DoAndReturn(returnError1(db, nil))

	mockCommit(db)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(
			mainUserMatcher,
			gomock.Eq(TokenPurposeEmail),
			&TimeGreaterThanNow{time.Now().Add(24 * time.Hour)},
		).
		Return(&Token{
			TokenPurpose: TokenPurposeEmail,
			Used:         false,
			Recalled:     false,
		}, "emailSecret", nil)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(
			mainUserMatcher,
			gomock.Eq(TokenPurposeSession),
			&TimeGreaterThanNow{time.Now().Add(time.Hour)},
		).
		Return(&Token{
			TokenPurpose: TokenPurposeSession,
			Used:         true,
			Recalled:     false,
		}, "sessionSecret", nil)

	//TODO subject and body probably should be tested too
	manager.emailService.(*MockEmailService).
		EXPECT().
		Send(
			gomock.Eq("test@example.com"),
			gomock.Any(),
			gomock.Any())

	user, token, secret, err := manager.Create(
		UserDetails{
			Email:     "test@example.com",
			Password:  "zaq1@WSX",
			FirstName: "first",
			LastName:  "last",
		},
	)

	require.Nilf(t, err, "manager.Create should return a nil error")

	require.NotNilf(t, user, "user returned by login should not be nil")
	assert.Equal(t, "test@example.com", user.Email, "User email is expected")
	assert.Equal(t, "sessionSecret", secret, "Invalid session secret")
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX"))
	require.Nilf(t, bcryptErr, "bcrypt.CompareHashAndPassword should return a nil error")

	require.NotNilf(t, token, "token returned by login should not be nil")
	assert.Equal(t, TokenPurposeSession, token.TokenPurpose)
	assert.Equal(t, user.ID, token.OwnerId)
	assert.Equal(t, true, token.Used)
	assert.Equal(t, false, token.Recalled)
}

// Tests if AuthManagerImpl.Create works correctly if provided email is invalid
func TestAuthManagerCreateWithInvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)

	user, _, _, err := manager.Create(
		UserDetails{
			Email:    "test",
			Password: "zaq1@WSX",
		},
	)

	require.ErrorIsf(t, ErrInvalidEmail, err, "manager.Create should return InvalidEmail error")
	require.Nilf(t, user, "manager.Create should return nil user")
}

// Tests if AuthManagerImpl.Create works correctly if user with the same email exists
func TestAuthManagerCreateWithExistingEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	mockBegin(db)

	db.(*MockGormDB).
		EXPECT().
		First(gomock.Any(), &StructMatcher{userMatcher{
			Email: Ptr("test@example.com"),
		}}).
		DoAndReturn(func(arg *User, conds interface{}) GormDB {
			db.(*MockGormDB).
				EXPECT().
				GetError().
				Return(nil)
			arg.ID = 1
			arg.Email = "test@example.com"
			return db
		})

	mockRollback(db)

	user, _, _, err := manager.Create(
		UserDetails{
			Email:    "test@example.com",
			Password: "zaq1@WSX",
		},
	)

	require.ErrorIsf(t, ErrEmailExists, err, "manager.Create should return EmailExists error")
	require.Nilf(t, user, "manager.Create should return nil user")
}

// Tests if AuthManagerImpl.Login works correctly on the happy path
func TestAuthManagerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	mockUser := mockExampleUser(db.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(&StructMatcher{userMatcher{
			ID: &mockUser.ID,
		}}, TokenPurposeSession, TimeGreaterThanNow{time.Now().Add(time.Hour)}).
		DoAndReturn(func(user *User, arg1 interface{}, arg2 interface{}) (*Token, string, error) {
			return &Token{OwnerId: user.ID, TokenId: "test", TokenHash: "test", TokenPurpose: TokenPurposeSession}, "sessionSecret", nil
		})

	user, token, sessionSecret, err := manager.Login("test@example.com", "zaq1@WSX")
	require.Nilf(t, err, "manager.Login should return a nil error")
	require.NotNilf(t, user, "manager.Login should not return a nil user")
	require.NotNilf(t, token, "manager.Login should not return a nil token")
	assert.Equal(t, "test@example.com", user.Email, "Invalid user email")
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX")); err != nil {
		t.Errorf("bcrypt did not return nil %s", err)
	}
	assert.Equal(t, user.ID, token.OwnerId, "Invalid token owner id")
	assert.Equal(t, TokenPurposeSession, token.TokenPurpose, "Invalid token purpose")
	assert.Equal(t, "sessionSecret", sessionSecret, "Invalid session secret")
}

// Asserts that user, token and sesessionSecret are nil, error is InvalidLogin - user failed to login
func assertInvalidLogin(t *testing.T, user *User, token *Token, sessionSecret string, err error) {
	if err != ErrInvalidLogin {
		t.Errorf("Error is not InvalidLogin %s", err)
	}
	if user != nil {
		t.Errorf("User is not nil")
	}
	if token != nil {
		t.Errorf("Token is not nil")
	}
	if sessionSecret != "" {
		t.Errorf("Session secret is not empty")
	}
}

// Tests if AuthManagerImpl.Login works correctly if password is invalid
func TestAuthManagerLoginInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	mockExampleUser(db.(*MockGormDB))

	user, token, sessionSecret, err := manager.Login("test@example.com", "invalid_password")
	assertInvalidLogin(t, user, token, sessionSecret, err)
}

// Tests if AuthManagerImpl.Login works correctly if email is invalid
func TestAuthManagerLoginInvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	db.(*MockGormDB).
		EXPECT().
		First(gomock.Any(), &StructMatcher{userMatcher{
			Email: Ptr("unknown@example.com"),
		}}).
		DoAndReturn(returnError2(db, gorm.ErrRecordNotFound))

	user, token, sessionSecret, err := manager.Login("unknown@example.com", "invalid_password")
	assertInvalidLogin(t, user, token, sessionSecret, err)
}

// Tests if AuthManagerImpl.Logiut works correctly on the happy path
func TestAuthManagerLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	manager, _ := getAuthManager(ctrl)

	user := getExampleUser()
	token := getExampleUserLogin()
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check(token.TokenId, "test_hash").
		Return(&user, &token, nil)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Invalidate(&StructMatcher{tokenMatcher{
			TokenId: Ptr(token.TokenId),
		}}).
		DoAndReturn(func(token *Token) (*User, *Token, error) {
			token.Recalled = true
			return &user, token, nil
		})

	logoutUser, logoutToken, err := manager.Logout(token.TokenId, "test_hash")
	require.Nil(t, err)
	require.NotNil(t, logoutUser)
	require.NotNil(t, logoutToken)
	require.True(t, logoutToken.Recalled)
	require.Equal(t, token.TokenId, logoutToken.TokenId)
}

// Tests if AuthManagerImpl.Logout works correctly when provided session token is invalid
func TestAuthManagerLogoutInvalidPurpose(t *testing.T) {
	ctrl := gomock.NewController(t)
	manager, _ := getAuthManager(ctrl)

	user := getExampleUser()
	token := getExampleUserLogin()
	token.TokenPurpose = TokenPurposeEmail
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check(token.TokenId, "test_hash").
		Return(&user, &token, nil)

	logoutUser, logoutToken, err := manager.Logout(token.TokenId, "test_hash")
	require.ErrorIs(t, err, ErrInvalidTokenPurpose)
	require.Nil(t, logoutUser)
	require.Nil(t, logoutToken)
}

// Tests if AuthManagerImpl.ConfirmEmail works correctly on the happy path
func TestAuthManagerConfirmEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	user := getExampleUser()
	token := createExampleToken("test_email", TokenPurposeEmail)

	mockBegin(db)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("test_email", "test_hash").
		Return(&user, &token, nil)

	db.(*MockGormDB).
		EXPECT().
		Save(&StructMatcher{userMatcher{
			ID:            Ptr(user.ID),
			EmailVerified: Ptr(true),
		}}).
		DoAndReturn(returnError1(db, nil))

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Invalidate(&StructMatcher{tokenMatcher{TokenId: Ptr("test_email")}}).
		Return(&user, &token, nil)

	mockCommit(db)

	changedUser, err := manager.ConfirmEmail("test_email", "test_hash")
	require.Nilf(t, err, "confirmEmail returned not nil error")
	require.NotNilf(t, changedUser, "changedUser should not be nil")
	require.Truef(t, changedUser.EmailVerified, "user email should be verified")
}

// Tests if AuthManagerImpl.ConfirmEmail works correctly when token id is invalid
func TestAuthManagerConfirmEmailInvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	mockBegin(db)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("invalid_id", "test_hash").
		Return(nil, nil, ErrInvalidToken)

	mockRollback(db)

	_, err := manager.ConfirmEmail("invalid_id", "test_hash")
	require.ErrorIsf(t, ErrInvalidToken, err, "ConfirmEmail should return UnknownToken")
}

// Tests if AuthManagerImpl.ConfirmEmail works correctly when token secret is invalid
func TestAuthManagerConfirmEmailInvalidHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	mockBegin(db)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("test_email", "invalid_hash").
		Return(nil, nil, UnknownToken)

	mockRollback(db)

	_, err := manager.ConfirmEmail("test_email", "invalid_hash")
	if err != ErrInvalidToken {
		t.Errorf("ConfirmEmail did not return InvalidToken %s", err)
	}
}

// Tests if AuthManagerImpl.ChangePassword works correctly on the happy path
func TestAuthManagerChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	var hash string
	user := getExampleUser()

	db.(*MockGormDB).
		EXPECT().
		Save(&StructMatcher{userMatcher{
			ID: Ptr(user.ID),
		}}).
		DoAndReturn(func(value *User) GormDB {
			hash = value.PasswordHash
			db.(*MockGormDB).
				EXPECT().
				GetError().
				Return(nil)
			return db
		})

	_, err := manager.ChangePassword(user, "zaq1@WSX", "nu9AhYoo")
	require.Nilf(t, err, "manager.ChangePassword should return a nil error")
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(hash), []byte("nu9AhYoo"))
	require.Nilf(t, bcryptErr, "bcrypt.CompareHashAndPassword should return a nil error")
}

// Tests if AuthManagerImpl.ChangePassword works correctly when oldPassword is invalid
func TestAuthManagerChangePasswordInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)

	user := getExampleUser()

	_, err := manager.ChangePassword(user, "test", "nu9AhYoo")
	require.ErrorIsf(t, ErrInvalidOldPassword, err, "manager.ChangePassword should return a nil error")
}

// Tests if AuthManagerImpl.ChangeEmail works correctly on the happy path
func TestAuthManagerChangeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)
	db := manager.baseServices.Database

	user := getExampleUser()

	mockBegin(db)

	db.(*MockGormDB).
		EXPECT().
		Save(&StructMatcher{userMatcher{
			ID:            Ptr(user.ID),
			EmailVerified: Ptr(false),
		}}).
		DoAndReturn(returnError1(db, nil))

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(
			&StructMatcher{userMatcher{
				Email:         Ptr("test2@example.com"),
				EmailVerified: Ptr(false),
			}},
			gomock.Eq(TokenPurposeEmail),
			&TimeGreaterThanNow{time.Now().Add(24 * time.Hour)},
		).
		Return(&Token{
			TokenPurpose: TokenPurposeEmail,
			Used:         false,
			Recalled:     false,
		}, "test_secret", nil)

	manager.emailService.(*MockEmailService).
		EXPECT().
		Send(
			gomock.Eq("test2@example.com"),
			gomock.Any(),
			gomock.Any())

	mockCommit(db)

	changedUser, err := manager.ChangeEmail(user, "test2@example.com")

	require.Nilf(t, err, "ChangeEmail should return a nil error")
	require.NotNilf(t, changedUser, "ChangeEmail should not return a nil user")
	require.Falsef(t, changedUser.EmailVerified, "ChangeEmail should return a user with EmailVerified set to false")
}

// Tests if AuthManagerImpl.ChangeEmail works correctly when email is invalid
func TestAuthManagerChangeEmailInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := getAuthManager(ctrl)

	user := getExampleUser()

	changedUser, err := manager.ChangeEmail(user, "asd")

	require.ErrorIsf(t, ErrInvalidEmail, err, "error should be InvalidEmail")
	require.Nilf(t, changedUser, "changedUser should be nil")
}
