package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ListenPort string
	MonitorInterval time.Duration
	WriterInterval time.Duration
	SerialDevice string
	DbPath string
}

var config Config

func Read() {
	config = Config{
		ListenPort: getEnvString("LISTEN_PORT", "3000"),
		MonitorInterval: time.Duration(getEnvInt("MONITOR_INTERVAL", 37000)),
		WriterInterval: time.Duration(getEnvInt("WRITER_INTERVAL", 97000)),
		SerialDevice: getEnvString("SERIAL_DEVICE", "/dev/ttyUSB0"),
		DbPath: getEnvString("DB_PATH", "/app/db/ups.sqlite"),
	}
}

func getEnvString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}		
	}
	return defaultVal
}

func GetListenPort() string {
	return config.ListenPort
}
func GetMonitorInterval() time.Duration {
	return config.MonitorInterval
}
func GetWriterInterval() time.Duration {
	return config.WriterInterval
}
func GetSerialDevice() string {
	return config.SerialDevice
}
func GetDbPath() string {
	return config.DbPath
}