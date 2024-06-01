package configs

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func getEnv(name string, defaultValue string) string {
	res, ok := os.LookupEnv(name)

	if !ok || res == "" {
		logrus.Info(fmt.Sprintf("ENV Variable doesn't have '%v' key, using '%v' instead", name, defaultValue))
		res = defaultValue
	}

	return res
}

func GetPort() string {
	return getEnv("PORT", "5000")
}

func GetDBHOST() string {
	return getEnv("DB_HOST", "localhost")
}

func GetDBPORT() string {
	return getEnv("DB_PORT", "3306")
}

func GetDBUSER() string {
	return getEnv("DB_USER", "root")
}

func GetDBPASS() string {
	return getEnv("DB_PASS", "")
}

func GetDBNAME() string {
	return getEnv("DB_NAME", "simple")
}
