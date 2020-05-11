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
	/*
	configUpdatedTriggerStream chan types.TriggerStream
	updateConfigTrigger types.TriggerStream
	*/
}

func (c *Config)Read() {
	c.listenPort = getEnvString("LISTEN_PORT", "3000")
	c.monitorInterval = time.Duration(getEnvInt("MONITOR_INTERVAL", 37000))
	c.writerInterval = time.Duration(getEnvInt("WRITER_INTERVAL", 97000))
	c.serialDevice = getEnvString("SERIAL_DEVICE", "/dev/ttyUSB0")
	c.dbPath = getEnvString("DB_PATH", "/app/db/ups.sqlite")
	/*
	if (c.configUpdatedTriggerStream == nil) {
		c.configUpdatedTriggerStream = make(chan types.TriggerStream)
		c.updateConfigTrigger = types.TriggerStreamCreate()
		c.configUpdatedTriggerStream <- c.updateConfigTrigger
	}
	*/
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

func (c *Config)ListenPortGet() string {
	return c.listenPort
}
func (c *Config)MonitorIntervalGet() time.Duration {
	return c.monitorInterval
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

/*
func (c *Config)MonitorIntervalSet(interval time.Duration) {
	c.monitorInterval = interval
	c.triggerUpdate()
}

func (c *Config)triggerUpdate() {
	close(c.updateConfigTrigger.Flag)
	c.updateConfigTrigger = types.TriggerStreamCreate()
	c.configUpdatedTriggerStream <- c.updateConfigTrigger
}

func (c *Config)ConfigUpdatedTriggerGet() chan types.TriggerStream {
	return c.configUpdatedTriggerStream
}
*/