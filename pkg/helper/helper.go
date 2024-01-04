package helper

import (
	cfg "WatchHive/pkg/config"
	interfaces "WatchHive/pkg/helper/interface"
	"WatchHive/pkg/utils/errmsg"
	"WatchHive/pkg/utils/models"
	"reflect"
	"strconv"
	"unicode"

	"errors"
	"fmt"
	"mime/multipart"

	"regexp"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
	"golang.org/x/crypto/bcrypt"
)

type helper struct {
	cfg cfg.Config
}

func NewHelper(config cfg.Config) interfaces.Helper {
	return &helper{
		cfg: config,
	}
}

type AuthCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func (helper *helper) GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error) {

	accesTokenClaims := &AuthCustomClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCustomClaims{
		Id:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	cfg, _ := cfg.LoadConfig()
	accesToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accesTokenClaims)
	accessTokenString, err := accesToken.SignedString([]byte(cfg.AdminAccessKey))
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.AdminRefreshKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil

}
func (h *helper) GenerateTokenClients(user models.UserDetailsResponse) (string, error) {
	claims := &AuthCustomClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  "client",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	cfg, _ := cfg.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.UserAccessKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (h *helper) PasswordHashing(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New(errmsg.ErrServer)
	}
	hash := string(hashedPassword)
	return hash, nil
}

func (h *helper) CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}
	return nil
}

func (h *helper) Copy(udr *models.UserDetailsResponse, usr *models.UserSignInResponse) (models.UserDetailsResponse, error) {
	err := copier.Copy(udr, usr)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}
	return *udr, nil
}

var client *twilio.RestClient

func (h *helper) TwilioSetup(username string, password string) {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})

}

func (h *helper) TwilioSendOTP(phone string, serviceID string) (string, error) {
	fmt.Println("ghhkkk", phone)
	to := "+91" + phone
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {
		return " ", err
	}
	return *resp.Sid, nil

}

func (h *helper) TwilioVerifyOTP(serviceID string, code string, phone string) error {
	to := "+91" + phone
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(serviceID, params)
	if err != nil {
		return err
	}

	if *resp.Status == errmsg.StatusApprove {
		return nil
	}

	return errors.New(errmsg.ErrOtpValidate)

}

func (h *helper) ValidatePhoneNumber(phone string) bool {
	phoneNumber := phone
	pattern := `^\d{10}$`
	regex := regexp.MustCompile(pattern)
	value := regex.MatchString(phoneNumber)
	return value
}

func (h *helper) ValidatePin(pin string) bool {

	match, _ := regexp.MatchString(`^\d{4}(\d{2})?$`, pin)
	return match

}
func (h *helper) ValidateDatatype(data, intOrString string) (bool, error) {

	switch intOrString {
	case "int":
		if _, err := strconv.Atoi(data); err != nil {
			return false, errors.New(errmsg.ErrDataIsNot + "integer")
		}
		return true, nil
	case "string":
		kind := reflect.TypeOf(data).Kind()
		return kind == reflect.String, nil
	default:
		return false, errors.New(errmsg.ErrDataIsNot + intOrString)
	}

}

func (h *helper) AddImageToAwsS3(file *multipart.FileHeader) (string, error) {

	f, openErr := file.Open()
	if openErr != nil {
		return "", openErr
	}
	defer f.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(h.cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			h.cfg.AWSAccesskeyID,
			h.cfg.AWSSecretaccesskey,
			"",
		),
	})

	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(sess)
	bucketName := "watch-hive"

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Filename),
		Body:   f,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, file.Filename)
	return url, nil
}

func (h *helper) GetTimeFromPeriod(timePeriod string) (time.Time, time.Time) {

	endDate := time.Now()

	if timePeriod == "week" {
		startDate := endDate.AddDate(0, 0, -6)
		return startDate, endDate
	}

	if timePeriod == "month" {
		startDate := endDate.AddDate(0, -1, 0)
		return startDate, endDate
	}

	if timePeriod == "year" {
		startDate := endDate.AddDate(-1, 0, 0)
		return startDate, endDate
	}

	return endDate.AddDate(0, 0, -6), endDate

}

func (h *helper) ValidateDate(dateString string) bool {

	// dateLayout := "2006-01-02"
	dateLayout := "02-01-2006"

	_, err := time.Parse(dateLayout, dateString)

	// if err != nil {
	// 	return false
	// }

	return err == nil
}

func (h *helper) ValidateAlphabets(data string) (bool, error) {
	for _, char := range data {
		if !unicode.IsLetter(char) {
			return false, errors.New(errmsg.ErrAlphabet)
		}
	}
	return true, nil
}

func (h *helper) ConvertToExel(sales []models.OrderDetailsAdmin) (*excelize.File, error) {

	filename := "salesReport/sales_report.xlsx"
	file := excelize.NewFile()

	file.SetCellValue("Sheet1", "A1", "Product")
	file.SetCellValue("Sheet1", "B1", "Amount Sold")

	// Bold style for headings
	boldStyle, err := file.NewStyle(`{"font":{"bold":true}}`)
	if err != nil {
		return nil, err
	}

	file.SetCellStyle("Sheet1", "A1", "B1", boldStyle)

	var Total float64
	var Limit int
	for i, sale := range sales {
		col1 := fmt.Sprintf("A%d", i+2)
		col2 := fmt.Sprintf("B%d", i+2)

		file.SetCellValue("Sheet1", col1, sale.ProductName)
		file.SetCellValue("Sheet1", col2, sale.TotalAmount)
		Limit = i + 3
		Total += sale.TotalAmount

	}
	col1 := fmt.Sprintf("A%d", Limit)
	file.SetCellValue("Sheet1", col1, "Final Total")
	col2 := fmt.Sprintf("B%d", Limit)
	file.SetCellValue("Sheet1", col2, Total)

	// Larger font size for 'Final Total'
	largerFontStyle, err := file.NewStyle(`{"font":{"size":10}}`)
	if err != nil {
		return nil, err
	}
	file.SetCellStyle("Sheet1", col1, col2, largerFontStyle)

	if err := file.SaveAs(filename); err != nil {
		return nil, err
	}

	return file, nil
}
