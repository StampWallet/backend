package managers

import (
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/mocks"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	. "github.com/StampWallet/backend/internal/testutils"
)

func GetAuthManager(ctrl *gomock.Controller) (*AuthManagerImpl, error) {
	//db, mock, err := sqlmock.New()
	//if err != nil {
	//    ctrl.T.Errorf("failed to init sqlmock %s", err)
	//    return nil, nil, err
	//}
	return &AuthManagerImpl{
		&BaseServices{
			Logger:   log.Default(),
			Database: NewMockGormDB(ctrl),
		},
		NewMockEmailService(ctrl),
		NewMockTokenService(ctrl),
	}, nil
}

type UserMatcher struct {
	ID            *uint
	Email         *string
	PasswordHash  *string
	FirstName     *string
	LastName      *string
	EmailVerified *bool
}

type TokenMatcher struct {
	OwnerId      *uint
	TokenId      *string
	Expires      *time.Time
	TokenPurpose *TokenPurposeEnum
	Used         *bool
	Recalled     *bool
}

func TestAuthManagerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	userMatcher := &StructMatcher{UserMatcher{
		Email:         Ptr("test@example.com"),
		FirstName:     Ptr("first"),
		LastName:      Ptr("last"),
		EmailVerified: Ptr(false),
	}}

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		First(gomock.Any(), &StructMatcher{UserMatcher{
			Email: Ptr("test@example.com"),
		}}).
		DoAndReturn(func(user *User, cond interface{}) GormDB {
			manager.baseServices.Database.(*MockGormDB).
				EXPECT().
				GetError().
				Return(gormlogger.ErrRecordNotFound)
			return manager.baseServices.Database
		})

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Create(userMatcher)

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Commit()

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(
			userMatcher,
			gomock.Eq(TokenPurposeEmail),
			&TimeGreaterThanNow{time.Now().Add(24 * time.Hour)},
		).
		Return(&Token{
			TokenPurpose: TokenPurposeEmail,
			Used:         false,
			Recalled:     false,
		}, nil)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(
			userMatcher,
			gomock.Eq(TokenPurposeSession),
			&TimeGreaterThanNow{time.Now().Add(time.Hour)},
		).
		Return(&Token{
			TokenPurpose: TokenPurposeSession,
			Used:         true,
			Recalled:     false,
		}, nil)

	//TODO subject and body probably should be tested too
	manager.emailService.(*MockEmailService).
		EXPECT().
		Send(
			gomock.Eq("test@example.com"),
			gomock.Any(),
			gomock.Any())

	user, token, err := manager.Create(
		UserDetails{
			Email:     "test@example.com",
			Password:  "zaq1@WSX",
			FirstName: "first",
			LastName:  "last",
		},
	)
	if err != nil {
		t.Errorf("Expected no errors but received: %s", err)
	}

	if user != nil {
		assert.Equal(t, "test@example.com", user.Email, "User email is expected")
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX"))
		if err != nil {
			t.Errorf("CompareHashAndPassword retruned an error %s", err)
		}
	} else {
		t.Errorf("User is nil")
	}

	if token != nil {
		assert.Equal(t, TokenPurposeSession, token.TokenPurpose)
		assert.Equal(t, user.ID, token.OwnerId)
		assert.Equal(t, true, token.Used)
		assert.Equal(t, false, token.Recalled)
	} else {
		t.Errorf("Token is nil")
	}
}

func TestAuthManagerCreateWithInvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	user, _, err := manager.Create(
		UserDetails{
			Email:    "test",
			Password: "zaq1@WSX",
		},
	)

	if err != InvalidEmail {
		t.Errorf("Expected an InvalidEmail error %s", err)
	}
	if user != nil {
		t.Errorf("User is not nil")
	}
}

func TestAuthManagerCreateWithExistingEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		First(gomock.Any(), &StructMatcher{UserMatcher{
			Email: Ptr("test@example.com"),
		}}).
		DoAndReturn(func(arg *User) *GormDB {
			arg.ID = 1
			arg.Email = "test@example.com"
			return &manager.baseServices.Database
		})

	user, _, err := manager.Create(
		UserDetails{
			Email:    "test@example.com",
			Password: "zaq1@WSX",
		},
	)

	if err != EmailExists {
		t.Errorf("Expected an ExistingEmail error %s", err)
	}
	if user != nil {
		t.Errorf("User is not nil")
	}
}

func GetExampleUser() User {
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

func mockExampleUser(database *MockGormDB) User {
	user := GetExampleUser()
	database.
		EXPECT().
		First(gomock.Any(), &StructMatcher{UserMatcher{
			Email: Ptr("test@example.com"),
		}}).
		Do(func(arg *User, conds ...interface{}) GormDB {
			*arg = user
			return database
		})
	return user
}

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

func mockExampleUserEmailVerificationToken(database *MockGormDB) Token {
	token := createExampleToken("test_email", TokenPurposeEmail)
	database.
		EXPECT().
		Find(gomock.Any(), &StructMatcher{TokenMatcher{
			TokenId: Ptr("test_email"),
		}}).
		Do(func(arg *Token) GormDB {
			*arg = token
			return database
		})
	return token
}

func mockExampleUserLogin(tokenService *MockTokenService) Token {
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
	user := GetExampleUser()
	tokenService.
		EXPECT().
		Check("test_login", "test_hash").
		Return(&user, &token, nil)
	return token
}

func TestAuthManagerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	mockUser := mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(UserMatcher{
			ID: &mockUser.ID,
		}, TokenPurposeSession, TimeGreaterThanNow{time.Now().Add(time.Hour)}).
		DoAndReturn(ReturnArg)

	user, token, err := manager.Login("test@example.com", "zaq1@WSX")
	if err != nil {
		t.Errorf("Error is not nil %s", err)
	}
	if user == nil {
		t.Errorf("User is nil")
	} else {
		assert.Equal(t, "test@example.com", user.Email, "Invalid user email")
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX")); err != nil {
			t.Errorf("bcrypt did not return nil %s", err)
		}
	}
	if token == nil {
		t.Errorf("User is nil")
	} else {
		assert.Equal(t, user.ID, token.OwnerId, "Invalid token owner id")
		assert.Equal(t, TokenPurposeSession, token.TokenPurpose, "Invalid token purpose")
	}
}

func assertInvalidLogin(t *testing.T, user *User, token *Token, err error) {
	if err != InvalidLogin {
		t.Errorf("Error is not InvalidLogin %s", err)
	}
	if user != nil {
		t.Errorf("User is not nil")
	}
	if token != nil {
		t.Errorf("Token is not nil")
	}
}

func TestAuthManagerInvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)
	mockExampleUser(manager.baseServices.Database.(*MockGormDB))

	user, token, err := manager.Login("test@example.com", "invalid_password")
	assertInvalidLogin(t, user, token, err)
}

func TestAuthManagerInvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)
	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		First(gomock.Any(), StructMatcher{&UserMatcher{
			Email: Ptr("unknown@example.com"),
		}})

	user, token, err := manager.Login("unknown@example.com", "invalid_password")
	assertInvalidLogin(t, user, token, err)
}

func TestAuthManagerLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	token := mockExampleUserLogin(manager.tokenService.(*MockTokenService))

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Invalidate(token.TokenId)

	logoutUser, logoutToken, err := manager.Logout("test", "test_hash")
	if err != nil {
		t.Errorf("Logout returned an error %s", err)
	}
	if logoutUser == nil {
		t.Errorf("logoutUser is nil")
	} else {
		assert.Equal(t, 1, logoutUser.ID, "Logout user does not match")
	}
	if logoutToken == nil {
		t.Errorf("logoutToken is nil")
	} else {
		assert.Equal(t, 1, logoutToken.ID, "Logout token id does not match")
		assert.Equal(t, true, logoutToken.Recalled, "Logout token recalled does not match")
	}
}

func TestAuthManagerInvalidLogoutHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("test", "test_hah").
		Return(nil, ErrUnknownToken)

	logoutUser, logoutToken, err := manager.Logout("test", "test_hah")
	if err != InvalidToken {
		t.Errorf("Logout did not return InvalidToken %s", err)
	}
	if logoutUser != nil {
		t.Errorf("logoutUser is not nil")
	}
	if logoutToken != nil {
		t.Errorf("logoutUser is not nil")
	}
}

func TestAuthManagerInvalidLogoutValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("test_invalid", "test_hash").
		Return(nil, ErrUnknownToken)

	logoutUser, logoutToken, err := manager.Logout("test_invalid", "test_hash")
	if err != InvalidToken {
		t.Errorf("Logout did not return InvalidToken %s", err)
	}
	if logoutUser != nil {
		t.Errorf("logoutUser is not nil")
	}
	if logoutToken != nil {
		t.Errorf("logoutUser is not nil")
	}
}

func TestAuthManagerConfirmEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	user := mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	token := mockExampleUserEmailVerificationToken(manager.baseServices.Database.(*MockGormDB))

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("test_email", "test_hash").
		Return(&token, nil)

	manager.tokenService.(*MockTokenService).
		EXPECT().
		Invalidate(StructMatcher{&TokenMatcher{TokenId: Ptr("test_email")}}).
		Return(&token, nil)

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Save(StructMatcher{&UserMatcher{
			ID:            Ptr(user.ID),
			EmailVerified: Ptr(true),
		}})

	changedUser, err := manager.ConfirmEmail("test_email", "test_hash")
	if err != nil {
		t.Errorf("ConfirmEmail did not return nil %s", err)
	}
	if changedUser != nil {
		if changedUser.EmailVerified {
			t.Errorf("User email is not verified")
		}
	} else {
		t.Errorf("ChangedUser is nil")
	}
}

func TestAuthManagerConfirmEmailInvalidId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("invalid_id", "test_hash").
		Return(nil, ErrUnknownToken)

	_, err := manager.ConfirmEmail("invalid_id", "test_hash")
	if err != InvalidToken {
		t.Errorf("ConfirmEmail did not return InvalidToken %s", err)
	}
}

func TestAuthManagerConfirmEmailInvalidHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	mockExampleUserLogin(manager.tokenService.(*MockTokenService))
	mockExampleUserEmailVerificationToken(manager.baseServices.Database.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Check("test_email", "invalid_hash").
		Return(nil, ErrUnknownToken)

	_, err := manager.ConfirmEmail("test_email", "invalid_hash")
	if err != InvalidToken {
		t.Errorf("ConfirmEmail did not return InvalidToken %s", err)
	}
}

func TestAuthManagerChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	var hash string
	user := mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Save(StructMatcher{&UserMatcher{
			ID: Ptr(user.ID),
		}}).
		DoAndReturn(func(value *User) *GormDB {
			hash = value.PasswordHash
			return &manager.baseServices.Database
		})

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Commit()

	_, err := manager.ChangePassword(user, "zaq1@WSX", "nu9AhYoo")
	if err != nil {
		t.Errorf("ChangePassword did not return nil %s", err)
	}
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(hash), []byte("nu9AhYoo"))
	if bcryptErr != nil {
		t.Errorf("bcrypt did not return nil %s", bcryptErr)
	}
}

func TestAuthManagerChangePasswordInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	user := mockExampleUser(manager.baseServices.Database.(*MockGormDB))

	_, err := manager.ChangePassword(user, "test", "nu9AhYoo")
	if err == nil {
		t.Errorf("ChangePassword returned nil %s", err)
	}
}

func TestAuthManagerChangeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	manager, _ := GetAuthManager(ctrl)

	user := mockExampleUser(manager.baseServices.Database.(*MockGormDB))
	manager.tokenService.(*MockTokenService).
		EXPECT().
		Create(
			&StructMatcher{UserMatcher{
				Email:         Ptr("test@example.com"),
				EmailVerified: Ptr(true),
			}},
			gomock.Eq(TokenPurposeEmail),
			&TimeGreaterThanNow{time.Now().Add(24 * time.Hour)},
		).
		Return(&Token{
			TokenPurpose: TokenPurposeEmail,
			Used:         false,
			Recalled:     false,
		}, nil)

	manager.emailService.(*MockEmailService).
		EXPECT().
		Send(
			gomock.Eq("test@example.com"),
			gomock.Any(),
			gomock.Any())

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Save(StructMatcher{&UserMatcher{
			ID:            Ptr(user.ID),
			EmailVerified: Ptr(false),
		}})

	manager.baseServices.Database.(*MockGormDB).
		EXPECT().
		Commit()

	changedUser, err := manager.ChangeEmail(user, "test2@example.com")

	if err != nil {
		t.Errorf("ChangeEmail returned an error %s", err)
	}
	if changedUser != nil {
		if changedUser.EmailVerified {
			t.Errorf("changedUser still has verified email")
		}
	} else {
		t.Errorf("changedUser is nil")
	}
}
