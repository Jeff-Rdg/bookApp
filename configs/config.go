package configs

import (
	"BookApp/author"
	"BookApp/book"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type conf struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}

func LoadDatabase(conf *conf) (*gorm.DB, error) {
	var conn gorm.Dialector

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)

	if conf != nil {
		conn = mysql.Open(dsn)
	} else {
		conn = sqlite.Open("test.db")
	}

	db, err := gorm.Open(conn, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(author.Author{}, book.Book{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
