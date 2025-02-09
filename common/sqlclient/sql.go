package sqlclient

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var secretKey []byte = []byte("anhle_golang@!&*")

type ISqlClientConn interface {
	GetDB() *bun.DB
	Connect(secretKey string) (err error)
}

type SqlConfig struct {
	SecretKey string
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
}

type SqlClientConn struct {
	SqlConfig
	DB *bun.DB
}

func NewSqlClient(config SqlConfig) ISqlClientConn {
	client := &SqlClientConn{}
	client.SqlConfig = config
	if err := client.Connect(config.SecretKey); err != nil {
		log.Fatal(err)
		return nil
	}
	if err := client.DB.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}
	return client
}

func (c *SqlClientConn) Connect(secretKey string) (err error) {
	if ok, err := decrypt(secretKey); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("invalid secret key")
	}
	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", c.Host, c.Port)),
		pgdriver.WithUser(c.Username),
		pgdriver.WithPassword(c.Password),
		pgdriver.WithDatabase(c.Database),
	)
	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	c.DB = db
	return nil
}

func (c *SqlClientConn) GetDB() *bun.DB {
	return c.DB
}

func decrypt(encryptedText string) (bool, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return false, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return false, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return false, fmt.Errorf("invalid ciphertext")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, err
	}

	return string(plaintext) == "anhle_golang@!&*", nil
}
