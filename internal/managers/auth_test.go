package managers

import (
	"database/sql/driver"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	. "github.com/StampWallet/backend/internal/database"
	. "github.com/StampWallet/backend/internal/database/mocks"
	. "github.com/StampWallet/backend/internal/services"
	. "github.com/StampWallet/backend/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func getAuthManager(ctrl *gomock.Controller) (*AuthManagerImpl, error) {
    //db, mock, err := sqlmock.New()
    //if err != nil {
    //    ctrl.T.Errorf("failed to init sqlmock %s", err)
    //    return nil, nil, err
    //}
    return &AuthManagerImpl {
        &BaseServices {
            Logger: log.Default(),
            Database: NewMockGormDB(ctrl),
        },
        NewMockEmailService(ctrl),
        NewMockTokenService(ctrl),
    }, nil
}

type Anything struct {}
func (Anything) Match(v driver.Value) bool {
	return true
}

type UserMatcher struct {
    ID *uint
    Email *string
    PasswordHash *string
    FirstName *string
    LastName *string
    EmailVerified *bool
}

func matchEntities(matcher interface{}, obj interface{}) bool {
    o := reflect.ValueOf(obj)
    if o.Kind() == reflect.Pointer {
        return matchEntities(matcher, o.Elem().Interface())
    } else {
        m := reflect.ValueOf(matcher)
        mt := reflect.TypeOf(matcher)
        for i := 0; i < mt.NumField(); i++ {
            mtf := mt.Field(i)
            of := o.FieldByName(mtf.Name)
            mf := m.FieldByName(mtf.Name)
            if !mf.IsNil() && !of.Equal(mf.Elem()) {
                return false
            }
        }
        return true
    }
}

type StructMatcher struct {
    obj interface{} 
}

func (matcher StructMatcher) Matches(x interface{}) bool {
    return matchEntities(matcher.obj, x)
}

func (StructMatcher) String() string {
    return "StructMatcher"
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

func String(s string) *string {
    return &s
}

func Bool(b bool) *bool {
    return &b
}

func TestAuthManagerCreate(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, _ := getAuthManager(ctrl)

    userMatcher := &StructMatcher{UserMatcher{
        Email: String("test@example.com"),
        FirstName: String("first"),
        LastName: String("last"),
        EmailVerified: Bool(false),
    }} 

    manager.baseServices.Database.(*MockGormDB).
        EXPECT().
        Find(gomock.Any(), &StructMatcher{UserMatcher{ 
            Email: String("test@example.com"),
        }})

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
            gomock.Eq(EmailTokenPurpose),
            &TimeGreaterThanNow{time.Now().Add(24*time.Hour)},
        ).
        Return(Token{
            TokenPurpose: EmailTokenPurpose,
            Used: false,
            Recalled: false,
        }, nil)

    manager.tokenService.(*MockTokenService).
        EXPECT().
        Create(
            userMatcher,
            gomock.Eq(SessionTokenPurpose),
            &TimeGreaterThanNow{time.Now().Add(time.Hour)},
        ).
        Return(Token{
            TokenPurpose: SessionTokenPurpose,
            Used: true,
            Recalled: false,
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
            Email: "test@example.com",
            Password: "zaq1@WSX",
            FirstName: "first",
            LastName: "last",
        },
    )
    if err != nil {
        t.Errorf("Expected no errors but received: %s", err)
    }

    if user != nil {
        assert.Equal(t, user.Email, "test@example.com", "User email is expected")
        err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("zaq1@WSX"))
        if err != nil {
            t.Errorf("CompareHashAndPassword retruned an error %s", err)
        }
    } else {
        t.Errorf("User is nil")
    }

    if token != nil {
        assert.Equal(t, token.TokenPurpose, SessionTokenPurpose)
        assert.Equal(t, token.OwnerId, user.ID)
        assert.Equal(t, token.Used, true)
        assert.Equal(t, token.Recalled, false)
    } else {
        t.Errorf("Token is nil")
    }
}

func TestAuthManagerCreateWithInvalidEmail(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, _ := getAuthManager(ctrl)

    user, _, err := manager.Create(
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

func TestAuthManagerLogin(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, _ := getAuthManager(ctrl)

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
    manager, _ := getAuthManager(ctrl)

    user, token, err := manager.Login("test@example.com", "invalid_password")
    assertInvalidLogin(t, user, token, err)
}

func TestAuthManagerInvalidEmail(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, _ := getAuthManager(ctrl)

    user, token, err := manager.Login("unknown@example.com", "invalid_password")
    assertInvalidLogin(t, user, token, err)
}

//TODO should i actually mock the database?
func TestAuthManagerLogout(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    manager, _ := getAuthManager(ctrl)

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
