package model

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"html/template"
	"math"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"ezcoin.cc/ezcoin-go/server/pkg/email"
	devisecrypto "github.com/consyse/go-devise-encryptor"

	uuid "github.com/satori/go.uuid"
	"github.com/theckman/go-securerandom"
	"gorm.io/gorm"
)

type User struct {
	EZ_MODEL
	Username            string      `json:"username" gorm:"comment:用户名"`
	Email               string      `json:"email" gorm:"comment:電子郵件"`
	Password            string      `json:"-" gorm:"-"`
	EncryptedPassword   string      `json:"-" gorm:"comment:用戶密碼"`
	ResetPasswordToken  string      `json:"-" gorm:"<-:update,index,comment:重設密碼令牌"`
	ResetPasswordSentAt time.Time   `json:"-" gorm:"<-:update,comment:重設密碼寄送時間"`
	RememberCreatedAt   time.Time   `json:"-" gorm:"<-:update,comment:登入記住帳號建立時間"`
	SignInCount         uint        `json:"-" gorm:"comment:登入次數"`
	CurrentSignInIp     string      `json:"-" gorm:"comment:當前登入IP"`
	CurrentSignInAt     time.Time   `json:"-" gorm:"comment:當前登入時間"`
	LastSignInAt        time.Time   `json:"-" gorm:"comment:最後登入時間"`
	ConfirmationToken   string      `json:"-" gorm:"index,comment:確認帳號令牌"`
	ConfirmationSentAt  time.Time   `json:"-" gorm:"comment:確認帳號令牌寄送時間"`
	ConfirmedAt         time.Time   `json:"confirmed_at" gorm:"<-:update,comment:帳號確認時間"`
	FailedAttempts      uint        `json:"-" gorm:"<-:update,comment:登入失敗次數"`
	ApiKey              string      `json:"apiKey" gorm:"<-:update,comment:Bitfinex API Key"`
	ApiSecret           string      `json:"apiSecret" gorm:"<-:update,comment:Bitfinex API Secret"`
	ExpiredAt           time.Time   `json:"expiredAt" gorm:"<-:update,comment:過期時間"`
	PlanId              uint        `json:"planId" gorm:"comment:用戶註冊 Plan"`
	Plan                Plan        `json:"plan"` //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserRobotsCount     uint        `json:"userRobotsCount" gorm:"comment:用戶 Robot 記錄數"`
	UserRobots          []UserRobot `json:"-"`
	Jti                 string      `json:"-" gorm:"comment:JTI"`

	// Extra Attributes
	Machine                      string `json:"-" gorm:"-"`
	ConfirmationRawToken         string `json:"-" gorm:"-"`
	ResetPasswordRawToken        string `json:"-" gorm:"-"`
	SendConfirmationInstruction  bool   `json:"-" gorm:"-"`
	SendResetPasswordInstruction bool   `json:"-" gorm:"-"`
}

func (u User) Login(db *gorm.DB) (*User, error) {
	var user User
	pepper := ""
	err := db.Where("username = ?", u.Username).Or("email = ?", u.Email).Preload("Plan").Preload("UserRobots").First(&user).Error
	if err != nil {
		return nil, err
	}
	valid := devisecrypto.Compare(u.Password, pepper, user.EncryptedPassword)
	if !valid {
		return nil, errors.New("用戶不存在或密碼有誤")
	}

	if err := db.Model(&u).Where("id = ?", user.ID).Updates(User{
		CurrentSignInIp: u.CurrentSignInIp,
		CurrentSignInAt: time.Now(),
	}).Error; err != nil {
		return &user, err
	}

	user.Machine = u.Machine
	go user.sendSingedInEmail()

	return &user, nil
}

func (u User) List(db *gorm.DB, param request.TableQuery) (*response.PageResult, error) {
	var list []*User
	var result response.PageResult
	var count int64
	var query *gorm.DB

	query = db.Model(User{})

	// Filter
	for k, v := range param.Filter {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	query = query.Count(&count)
	global.GVA_LOG.Debug(fmt.Sprintf("err %v", query.Error))
	if query.Error != nil {
		return &result, query.Error
	}
	query = query.Limit(param.PageSize)
	if param.Sort != "" && param.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", param.Sort, param.Order))
	}
	query = query.Offset(app.GetPageOffset(param.Page, param.PageSize))
	err := query.Find(&list).Error
	if err != nil {
		return &result, err
	}
	result.List = list
	result.Meta = response.PageMeta{
		Total:     count,
		TotalPage: int(math.Ceil(float64(count) / float64(param.PageSize))),
		Page:      param.Page,
		PageSize:  param.PageSize,
	}

	return &result, err
}

func (u User) Create(db *gorm.DB) (*User, error) {
	var count int64
	db.Model(&User{}).Where("username = ?", u.Username).Or("email = ?", u.Email).Count(&count)
	if count != 0 {
		return nil, errors.New("用戶已存在")
	}

	if err := db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (u User) Update(db *gorm.DB) (*User, error) {
	if err := db.Where("id = ?", u.ID).Updates(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (u User) UpdateProfile(db *gorm.DB) (*User, error) {
	if err := db.Model(&u).Updates(User{ApiKey: u.ApiKey, ApiSecret: u.ApiSecret}).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (u User) Get(db *gorm.DB) (User, error) {
	var user User
	err := db.Where("id = ?", u.ID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, nil
}

func (u User) Delete(db *gorm.DB) error {
	if err := db.Where("id = ? AND deleted_at IS NULL", u.ID).Update("deleted_at = ?", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

func (u User) Robots(db *gorm.DB) ([]*UserRobot, error) {
	var robots []*UserRobot
	if err := db.Where("user_id = ?", u.ID).Preload("Currency").Preload("LendingStrategy").Find(&robots).Error; err != nil {
		return robots, err
	}
	return robots, nil
}

func (u User) ResendConfirmation(db *gorm.DB) (bool, error) {
	err := db.Where("email = ?", u.Email).First(&u).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}

	if u.Confirmed() {
		return false, errors.New("帳戶已開通")
	}

	token, err := securerandom.Base64InBytes(20)
	if err != nil {
		return false, err
	}

	u.ConfirmationRawToken = token
	u.SendConfirmationInstruction = true
	if err := db.Model(&u).Updates(
		User{
			ConfirmationToken:  fmt.Sprintf("%x", sha256.Sum256([]byte(token))),
			ConfirmationSentAt: time.Now(),
		}).Error; err != nil {
		return false, err
	}
	return true, nil
}

//@function: Confirm
//@description: Confirm account via token
//@param: db *gorm.DB
//@return: bool, error
func (u User) Confirm(db *gorm.DB) (bool, error) {
	token := fmt.Sprintf("%x", sha256.Sum256([]byte(u.ConfirmationRawToken)))
	global.GVA_LOG.Debug(fmt.Sprintf("token: %s %s", u.ConfirmationRawToken, token))
	err := db.Where("confirmation_token = ?", token).First(&u).Error
	if err != nil {
		return false, errors.New("token invalid")
	}

	if err := db.Model(&u).Updates(
		map[string]interface{}{"confirmed_at": time.Now(), "confirmation_token": gorm.Expr("NULL")},
	).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (u User) SendForgetPassword(db *gorm.DB) (bool, error) {
	err := db.Where("email = ?", u.Email).First(&u).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}
	token, err := securerandom.Base64InBytes(20)
	if err != nil {
		return false, err
	}

	u.ResetPasswordRawToken = token
	u.SendResetPasswordInstruction = true
	if err := db.Model(&u).Updates(
		User{
			ResetPasswordToken:  fmt.Sprintf("%x", sha256.Sum256([]byte(token))),
			ResetPasswordSentAt: time.Now(),
		}).Error; err != nil {
		return false, err
	}
	return true, nil
}

//@function: ResetPassword
//@description: Reset password via token
//@param: db *gorm.DB
//@return: bool, error
func (u User) ResetPassword(db *gorm.DB) (bool, error) {
	token := fmt.Sprintf("%x", sha256.Sum256([]byte(u.ResetPasswordRawToken)))
	err := db.Where("reset_password_token = ?", token).First(&u).Error
	if err != nil {
		return false, errors.New("token invalid")
	}
	pepper := ""
	encryptedPassword, err := devisecrypto.Digest(u.Password, 0, pepper)
	if err != nil {
		return false, err
	}

	if err := db.Model(&u).Updates(
		User{
			EncryptedPassword:  encryptedPassword,
			ResetPasswordToken: "",
		}).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (u User) Confirmed() bool {
	return !u.ConfirmedAt.Equal(time.Time{})
}

func (u *User) generateConfirmToken() error {
	if u.Confirmed() {
		return nil
	}

	token, err := securerandom.Base64InBytes(20)
	if err != nil {
		return err
	}
	u.SendConfirmationInstruction = true
	u.ConfirmationRawToken = token
	u.ConfirmationToken = fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	u.ConfirmationSentAt = time.Now()
	return nil
}

func (u *User) SendMail(token string) (err error) {
	// Send confirmation email
	global.GVA_LOG.Debug("SendMail")
	return
}

// Hook Events
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	global.GVA_LOG.Debug("BeforeCreate user")
	pepper := ""
	encryptedPassword, err := devisecrypto.Digest(u.Password, 0, pepper)
	if err != nil {
		return err
	}

	u.EncryptedPassword = encryptedPassword
	u.Jti = uuid.NewV4().String()
	u.PlanId = 1
	u.generateConfirmToken()
	return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// Send confirmation email
	global.GVA_LOG.Debug("AfterCreate user")
	// if err := u.sendConfirmationInstructions(); err != nil {
	// 	return err
	// }

	return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	// Send confirmation email
	global.GVA_LOG.Debug(fmt.Sprintf("AfterSave SendResetPassword : %v", u.SendResetPasswordInstruction))
	// if err := u.sendConfirmationInstructions(); err != nil {
	// 	return err
	// }

	if u.SendResetPasswordInstruction {
		go u.sendResetPasswordInstructions()
	}

	if u.SendConfirmationInstruction {
		go u.sendConfirmationInstructions()
	}
	return
}

func (u *User) sendConfirmationInstructions() error {
	if u.ConfirmationRawToken == "" {
		if err := u.generateConfirmToken(); err != nil {
			return err
		}
	}

	files := []string{
		"app/views/confirmation.html.tmpl",
		"app/views/layout/email.html.tmpl",
	}
	var tmpl = template.Must(template.ParseFiles(files...))
	var body bytes.Buffer
	confirmationUrl := fmt.Sprintf("%s/confirmation?confirmation_token=%s", global.GVA_CONFIG.System.FrontUrl, u.ConfirmationRawToken)
	tmpl.Execute(&body, struct {
		Username        string
		ConfirmationUrl string
	}{
		Username:        u.Username,
		ConfirmationUrl: confirmationUrl,
	})
	if err := email.Email(u.Email, "[EZCoin] 帳號驗證步驟", body.Bytes()); err != nil {
		return err
	}
	return nil
}

func (u *User) sendResetPasswordInstructions() error {
	if u.ResetPasswordToken == "" {
		return nil
	}
	files := []string{
		"app/views/reset_password.html.tmpl",
		"app/views/layout/email.html.tmpl",
	}
	var tmpl = template.Must(template.ParseFiles(files...))
	var body bytes.Buffer
	resetPasswordUrl := fmt.Sprintf("%s/password?reset_password_token=%s", global.GVA_CONFIG.System.FrontUrl, u.ResetPasswordRawToken)
	tmpl.Execute(&body, struct {
		Username         string
		ResetPasswordUrl string
	}{
		Username:         u.Username,
		ResetPasswordUrl: resetPasswordUrl,
	})
	if err := email.Email(u.Email, "[EZCoin] 重設密碼", body.Bytes()); err != nil {
		return err
	}
	return nil
}

func (u *User) sendSingedInEmail() error {
	files := []string{
		"app/views/signed_in.html.tmpl",
		"app/views/layout/email.html.tmpl",
	}
	var tmpl = template.Must(template.ParseFiles(files...))
	var body bytes.Buffer
	tmpl.Execute(&body, struct {
		Machine string
		Email   string
		LoginAt string
		IP      string
	}{
		Machine: u.Machine,
		Email:   u.Email,
		LoginAt: u.LastSignInAt.String(),
		IP:      u.CurrentSignInIp,
	})

	if err := email.Email(u.Email, "[EZCoin] 您已登入.", body.Bytes()); err != nil {
		return err
	}
	return nil
}
