package managers

import (
	"database/sql/driver"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getAuthManager(ctrl *gomock.Controller) (*AuthManagerImpl, sqlmock.Sqlmock, error) {
    db, mock, err := sqlmock.New()
    if err != nil {
        ctrl.T.Errorf("failed to init sqlmock %s", err)
        return nil, nil, err
    }
    dialector := postgres.New(postgres.Config{
        DSN:                  "sqlmock",
        DriverName:           "postgres",
        Conn:                 db,
        PreferSimpleProtocol: true,
    })
    orm, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        ctrl.T.Errorf("failed to init gorm %s", err)
        return nil, nil, err
    }
    return &AuthManagerImpl {
        &BaseServices {
            Logger: log.Default(),
            Database: orm,
        },
        NewMockEmailService(ctrl),
        NewMockTokenService(ctrl),
    }, mock, nil
}

type Anything struct {}
func (Anything) Match(v driver.Value) bool {
	return true
}

func getUserRows(t *testing.T) (*sqlmock.Rows) {
    rows := sqlmock.NewRows(GetColumns(reflect.TypeOf(User{})))
    pass, err := bcrypt.GenerateFromPassword([]byte("zaq1@WSX"), 10)
    if err != nil {
        t.Errorf("bcrypt returned an error %s", err)
    }

    rows.AddRow(1, 0, 0, 0, "test", "test", "test", "test@example.com", pass, true)
    return rows
}

type UserMatcher struct {
    ID uint
}

func (matcher UserMatcher) Matches(x interface{}) bool {
    return x.(User).ID == matcher.ID
}

func (UserMatcher) String() string {
    return "UserMatcher"
}

type TimeGreaterThanNow struct {
    Time time.Time
}

func (matcher TimeGreaterThanNow) Matches(x interface{}) bool {
    return matcher.Time.Before(x.(time.Time))
}

func (TimeGreaterThanNow) String() string {
    return "TimeGreaterThanNow"
}

func TestAuthManagerCreate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, dbmock, _ := getAuthManager(ctrl)

    dbmock.ExpectExec("INSERT INTO \"users\" .*").
		WithArgs(getUserRows(t))

    //TODO subject and body probably should be tested too
    manager.emailService.(*MockEmailService).
        EXPECT().
        Send(
            gomock.Eq("test@example.com"), 
            gomock.Any(), 
            gomock.Any())

    manager.tokenService.(*MockTokenService).
        EXPECT().
        Create(
            &UserMatcher{1},
            gomock.Eq(EmailTokenPurpose),
            &TimeGreaterThanNow{time.Now().Add(24*time.Hour)},
        )

    user, err := manager.Create(
        UserDetails{
            Email: "test@example.com",
            Password: "zaq1@WSX",
        },
    )
    if err != nil {
        t.Errorf("Expected no errors but received: %s", err)
    }

    if user != nil {
        assert.Equal(t, user.Email, "test@test.com", "User email is expected")
        err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX"))
        if err != nil {
            t.Errorf("CompareHashAndPassword retruned an error %s", err)
        }
    } else {
        t.Errorf("User is nil")
    }

	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthManagerCreateWithInvalidEmail(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, _, _ := getAuthManager(ctrl)

    user, err := manager.Create(
        UserDetails{
            Email: "test",
            Password: "zaq1@WSX",
        },
    )
    if err != InvalidEmail {
        t.Errorf("Expected an InvalidEmail error")
    }
    if user != nil {
        t.Errorf("User is not nil")
    } 
}

func GetColumns(t reflect.Type) []string {
    fields := t.NumField()
    var result []string
    for i := 0; i < fields; i++ {
        field := t.Field(i)
        if reflect.ValueOf(field).Kind() == reflect.Struct && field.Anonymous {
            tmp := GetColumns(field.Type)
            result = append(result, tmp...)
        } else {
            gormField := field.Tag.Get("gorm")
            if gormField == "" || !strings.Contains(gormField, "foreignkey:") {
                result = append(result, strings.ToLower(field.Name))
            }
        }
    }
    return result
}

func getTokenRows(t *testing.T, hash string) (*sqlmock.Rows) {
    rows := sqlmock.NewRows(GetColumns(reflect.TypeOf(Token{})))
    pass, err := bcrypt.GenerateFromPassword([]byte(hash), 10)
    if err != nil {
        t.Errorf("bcrypt returned an error %s", err)
    }

    rows.AddRow(1, 0, 0, 0, 1, "test", "test", pass, time.Now().Add(time.Hour), SessionTokenPurpose, false, false)
    return rows
}

func TestAuthManagerLogin(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, dbmock, _ := getAuthManager(ctrl)

    dbmock.ExpectQuery("SELECT \\* FROM \"users\" .*").
		WithArgs("test@example.com").
		WillReturnRows(getUserRows(t))

    user, token, err := manager.Login("test@example.com", "zaq1@WSX")
    if err != nil {
        t.Errorf("Error is not nil %s", err)
    } 
    if user == nil {
        t.Errorf("User is nil")
    } else {
        assert.Equal(t, user.Email, "test@example.com", "Invalid user email")
        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX")); err != nil {
            t.Errorf("bcrypt did not return nil %s", err)
        }
    }
    if token == nil {
        t.Errorf("User is nil")
    }  else {
        assert.Equal(t, token.OwnerId, user.ID, "Invalid token owner id")
        assert.Equal(t, token.TokenPurpose, SessionTokenPurpose, "Invalid token purpose")
    }
}

func assertInvalidLogin(t *testing.T, dbmock sqlmock.Sqlmock, user *User, token *Token, err error) {
    if err != InvalidLogin {
        t.Errorf("Error is not InvalidLogin %s", err)
    } 
    if user != nil {
        t.Errorf("User is not nil")
    } 
    if token != nil {
        t.Errorf("Token is not nil")
    }  
	if err := dbmock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAuthManagerInvalidPassword(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, dbmock, _ := getAuthManager(ctrl)

    dbmock.ExpectQuery("SELECT \\* FROM \"users\" .*").
		WithArgs("test@example.com").
		WillReturnRows(getUserRows(t))

    user, token, err := manager.Login("test@example.com", "invalid_password")
    assertInvalidLogin(t, dbmock, user, token, err)
}

func TestAuthManagerInvalidEmail(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, dbmock, _ := getAuthManager(ctrl)

    dbmock.ExpectQuery("SELECT \\* FROM \"users\" .*").
		WithArgs("unknown@example.com").
		WillReturnRows(getUserRows(t))

    user, token, err := manager.Login("unknown@example.com", "invalid_password")
    assertInvalidLogin(t, dbmock, user, token, err)
}

//TODO should i actually mock the database?
func TestAuthManagerLogout(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, dbmock, _ := getAuthManager(ctrl)

    //user, err := manager.Create(
    //    UserDetails{
    //        Email: "test@example.com",
    //        Password: "zaq1@WSX",
    //    },
    //)
    //if err != nil {
    //    t.Errorf("Failed to create user %s", err)
    //}

    //dbmock.ExpectQuery("SELECT \\* FROM \"users\" .*").
	//	WithArgs("test@example.com").
	//	WillReturnRows(getUserRows(t))

    //user, token, err := manager.Login("test@example.com", "zaq1@WSX")
    //if err != nil {
    //    t.Errorf("Failed to login %s", err)
    //}

    dbmock.ExpectQuery("SELECT \\* FROM \"users\" .*").
		WithArgs("test").
		WillReturnRows(getTokenRows(t, "test"))

    dbmock.ExpectQuery("SELECT \\* FROM \"token\" .*").
		WithArgs("test@example.com").
		WillReturnRows(getUserRows(t))

    logoutUser, logoutToken, err := manager.Logout("test", "test")
    if err != nil {
        t.Errorf("Logout returned an error %s", err)
    }
    if logoutUser == nil {
        t.Errorf("logoutUser is nil")
    } else {
        assert.Equal(t, logoutUser.ID, 1, "Logout user does not match")
    }
    if logoutToken == nil {
        t.Errorf("logoutToken is nil")
    } else {
        assert.Equal(t, logoutToken.ID, 1, "Logout token does not match")
        assert.Equal(t, logoutToken.Recalled, true, "Logout token does not match")
    }
}

func TestAuthManagerConfirmEmail(t *testing.T) {
    t.Fail()
}

func TestAuthManagerChangePassword(t *testing.T) {
    t.Fail()
}

func TestAuthManagerChangeEmail(t *testing.T) {
    t.Fail()
}
