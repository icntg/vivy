package initialize

import (
	"app/core/global"
	"app/core/utility/common"
	"app/core/utility/crypto"
	"app/core/utility/errno"
	"app/core/utility/terminal/qrcode"
	"app/core/web/model/system"
	"bufio"
	"database/sql"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
)

func InitDatabase() {
	var (
		gCfg = global.GetConfig()
		//gDb = global.GetGORM()
	)
	conn, err := sql.Open("mysql", gCfg.DataSource.MySQL.GetDSNWithOutDatabase())
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	if nil != err {
		common.ErrPrintf("sql cannot open database with [%s]: %v\n", gCfg.DataSource.MySQL.GetMaskedDSNWithOutDatabase(), err)
		os.Exit(errno.ErrorConnectDatabase)
	}
	common.OutPrintf("conn = %v\n", conn)

	// create database
	{
		_sql := fmt.Sprintf("CREATE DATABASE `%s` /*!40100 COLLATE 'utf8mb4_general_ci' */", gCfg.DataSource.MySQL.Database)
		r, err := conn.Exec(_sql)
		if nil != err {
			common.ErrPrintf("sql cannot create database [%s] with [%s]: %v\n", gCfg.DataSource.MySQL.Database, _sql, err)
		} else {
			common.OutPrintf("sql created database [%s] successfully: %v\n", gCfg.DataSource.MySQL.Database, r)
		}
	}

	// select to database
	{
		_sql := fmt.Sprintf("USE `%s`", gCfg.DataSource.MySQL.Database)
		r, err := conn.Exec(_sql)
		if nil != err {
			common.ErrPrintf("sql cannot use database [%s] with [%s]: %v\n", gCfg.DataSource.MySQL.Database, _sql, err)
			os.Exit(errno.ErrorConnectDatabase)
		} else {
			common.OutPrintf("sql use database [%s] successfully: %v\n", gCfg.DataSource.MySQL.Database, r)
		}
	}

	// wrap connection with GORM
	var gormDB *gorm.DB
	{
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn: conn,
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if nil != err {
			common.ErrPrintf("gorm cannot open database with [sql connection]: %v\n", err)
			os.Exit(errno.ErrorConnectDatabase)
		}
	}

	// create tables
	{
		err = gormDB.AutoMigrate(system.Department{})
		err = gormDB.AutoMigrate(system.Resource{})
		err = gormDB.AutoMigrate(system.Role{})
		err = gormDB.AutoMigrate(system.RoleResource{})
		err = gormDB.AutoMigrate(system.User{})
		err = gormDB.AutoMigrate(system.UserRole{})
	}

	admin := initAdmin()
	r := gormDB.Create(&admin)
	err = r.Error
	return
}

func initAdmin() system.User {
	admin := system.User{}
	admin.Service.Id = common.ObjectIdB32x()
	admin.Code = "P00000000"
	admin.Name = "系统管理员"
	admin.EMail = fmt.Sprintf("admin@%s.com", strings.ToLower(global.ProductName))

	input := bufio.NewScanner(os.Stdin)

	common.OutPrintf("===== CREATE ADMIN Begin =====\n")

	common.EndlessFunc(func() bool {
		common.OutPrintf("Please input the username of admin(default: 'admin', char must in `[A-Za-z0-9_]{4,20}`): ")
		input.Scan()
		adminName := strings.TrimSpace(input.Text())
		if len(adminName) == 0 {
			admin.LoginName = "admin"
			return true
		} else if common.IsAlphaDigitsBaseline(adminName) {
			admin.LoginName = adminName
			return true
		}
		common.ErrPrintf("ERROR: username is too short or contains illegal characters.")
		return false
	})
	common.OutPrintf("To use username: [%s]\n", admin.LoginName)

	common.EndlessFunc(func() bool {
		common.OutPrintf("Do you want to use Google Token? (Y/n): ")
		input.Scan()
		answer := strings.ToLower(strings.TrimSpace(input.Text()))
		if answer == "n" {
			return true
		}

		token := crypto.Rand(15, true)
		b32loToken := strings.ToLower(base32.StdEncoding.EncodeToString(token))
		otp := common.Format(`otpauth://totp/{{.issuer}}:{{.account}}?algorithm=SHA1&digits=6&period=30&issuer={{.issuer}}&secret={{.secret}}`).Exec(map[string]interface{}{
			"issuer":  global.ProductName,
			"account": admin.LoginName,
			"secret":  b32loToken,
		})
		qrCodes := qrcode.ConsoleQRCode(otp)
		common.OutPrintf("Now several QRCode(s) will be draw on terminal. You can choose one which can be recognized by Google Authenticator.\n")
		mode := ""
		for _, qr := range qrCodes {
			common.EndlessFunc(func() bool {
				common.OutPrintf(qr)
				common.OutPrintf("\nIs this one ok? Y(es) / N(o) / T(ext)): ")
				input.Scan()
				mode = strings.ToLower(strings.TrimSpace(input.Text()))
				if mode == "y" || mode == "n" || mode == "t" {
					return true
				}
				return false
			})
			if mode == "y" || mode == "t" {
				break
			}
		}
		if mode != "y" {
			buffer := strings.Builder{}
			for i := 0; i < len(b32loToken); i++ {
				buffer.WriteByte(b32loToken[i])
				if (i+1)%4 == 0 {
					buffer.WriteString(" ")
				}
			}
			outB32loToken := strings.TrimSpace(buffer.String())
			common.OutPrintf("Google Token in TextMode: [%s]\n", outB32loToken)
		}
		admin.Token = b32loToken
		return true
	})

	{
		clearPassword := hex.EncodeToString(crypto.Rand(10, true))
		admin.Salt = strings.ToLower(base32.StdEncoding.EncodeToString(crypto.Rand(10, true)))
		admin.Password, _ = crypto.EncPassword(clearPassword, admin.Salt)
		common.OutPrintf("Initial password of [%s] is: [%s]\n", admin.LoginName, clearPassword)
	}
	common.OutPrintf("===== CREATE ADMIN End =====\n")
	return admin
}
