package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	listenPort      string
	monitorInterval time.Duration
	writerInterval  time.Duration
	serialDevice    string
	dbPath          string

	chargeManagementEnabled bool

	configUpdatedTrigger chan struct{}
}

func (c *Config)Read() {
	c.listenPort = getEnvString("LISTEN_PORT", "3000")
	c.monitorInterval = time.Duration(getEnvInt("MONITOR_INTERVAL", 37000))
	c.writerInterval = time.Duration(getEnvInt("WRITER_INTERVAL", 97000))
	c.serialDevice = getEnvString("SERIAL_DEVICE", "/dev/ttyUSB0")
	c.dbPath = getEnvString("DB_PATH", "/app/db/ups.sqlite")

	c.chargeManagementEnabled = true

	c.triggerUpdate()
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

func (c *Config)triggerUpdate() {
	prevTrigger := c.configUpdatedTrigger;
	c.configUpdatedTrigger = make(chan struct{})
	if (prevTrigger != nil) {
		close(prevTrigger)
	}
}

func (c *Config)ConfigUpdatedTriggerGet() chan struct{} {
	return c.configUpdatedTrigger
}

func (c *Config)ListenPortGet() string {
	return c.listenPort
}

func (c *Config)MonitorIntervalGet() time.Duration {
	return c.monitorInterval
}
func (c *Config)MonitorIntervalSet(interval time.Duration) {
	c.monitorInterval = interval
	c.triggerUpdate()
}

func (c *Config)WriterIntervalGet() time.Duration {
	return c.writerInterval
}
func (c *Config)SerialDeviceGet() string {
	return c.serialDevice
}
func (c *Config)DbPathGet() string {
	return c.dbPath
}

func (c *Config)ChargeManagementEnabledGet() bool {
	return c.chargeManagementEnabled
}
func (c *Config)ChargeManagementEnabledSet(chargeManagementEnabled bool) {
	c.chargeManagementEnabled = chargeManagementEnabled
	c.triggerUpdate()
}
