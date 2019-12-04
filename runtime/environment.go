package runtime

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-pg/pg/v9"

	"github.com/jacob-ebey/golang-ecomm/apis"
)

type BraintreeConfig struct {
	MerchantID string
	PublicKey  string
	PrivateKey string
}

func BaseUrl() string { return os.Getenv("BASE_URL") }

func IsDevelopment() bool { return os.Getenv("ENVIRONMENT") == "development" }

func ShippoPrivateToken() string { return os.Getenv("SHIPPO_PRIVATE_TOKEN") }

func ShouldServeStaticFiles() bool { return os.Getenv("GO_SERVES_STATIC") == "true" }

func Braintree() BraintreeConfig {
	return BraintreeConfig{
		MerchantID: os.Getenv("BRAINTREE_MERCHANT_ID"),
		PublicKey:  os.Getenv("BRAINTREE_PUBLIC_KEY"),
		PrivateKey: os.Getenv("BRAINTREE_PRIVATE_KEY"),
	}
}

func ZeitToken() string { return os.Getenv("ZEIT_TOKEN") }

type SmtpConfig struct {
	From     string
	Username string
	Password string
	Host     string
	Port     int
}

func Smtp() SmtpConfig {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return SmtpConfig{
		From:     os.Getenv("SMTP_FROM"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
	}
}

func GetAddress() string {
	addr := ":" + os.Getenv("PORT")
	// Prevents windows dialog in dev
	if IsDevelopment() && !IsDocker() {
		addr = "localhost" + addr
	}

	return addr
}

func JwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		panic(fmt.Errorf("No JWT_SECRET environment variable provided."))
	}

	return []byte(secret)
}

func IsDocker() bool { return os.Getenv("IS_DOCKER") == "true" }

func IsHeroku() bool { return os.Getenv("IS_HEROKU") == "true" }

func GetPgOptions() *pg.Options {
	databaseURL := os.Getenv("DATABASE_URL")

	var options *pg.Options
	if strings.Contains(databaseURL, "@") {
		parsed, _ := pg.ParseURL(databaseURL)
		options = parsed
	} else {
		options = &pg.Options{
			Addr:     os.Getenv("POSTGRESS_ADDRESS"),
			Database: os.Getenv("POSTGRESS_DATABASE"),
			User:     os.Getenv("POSTGRESS_USER"),
			Password: os.Getenv("POSTGRESS_PASSWORD"),
		}
	}

	if IsHeroku() {
		if options.TLSConfig == nil {
			options.TLSConfig = &tls.Config{}
		}
		options.TLSConfig.InsecureSkipVerify = true
	}
	return options
}

func GetAvatax() *apis.Avatax {
	return &apis.Avatax{
		BearerToken: os.Getenv("AVATAX_BEARER_TOKEN"),
		Username:    os.Getenv("AVATAX_USERNAME"),
		Password:    os.Getenv("AVATAX_PASSWORD"),
		AccountID:   os.Getenv("AVATAX_ACCOUNTID"),
		LicenseKey:  os.Getenv("AVATAX_LICENSEKEY"),
	}
}
