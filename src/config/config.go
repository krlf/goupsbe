package config

import (
	"encoding/json"
	"log"
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

	startChargingVoltage uint32
	stopChargingVoltage  uint32
	shutdownVoltage      uint32
}

func (c *Config) Read() {
	c.listenPort = getEnvString("LISTEN_PORT", "3000")
	c.monitorInterval = time.Duration(getEnvInt("MONITOR_INTERVAL", 37000))
	c.writerInterval = time.Duration(getEnvInt("WRITER_INTERVAL", 97000))
	c.serialDevice = getEnvString("SERIAL_DEVICE", "/dev/ttyUSB0")
	c.dbPath = getEnvString("DB_PATH", "/app/db/ups.sqlite")

	c.chargeManagementEnabled = true

	c.startChargingVoltage = uint32(getEnvInt("START_CHARGING_VOLTAGE", 7200))
	c.stopChargingVoltage = uint32(getEnvInt("STOP_CHARGING_VOLTAGE", 7600))
	c.shutdownVoltage = uint32(getEnvInt("SHUTDOWN_VOLTAGE", 6200))

	if cnf, err := json.Marshal(c); err == nil {
		log.Print("Configuration:")
		log.Print(cnf)
	}

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

func (c *Config) triggerUpdate() {
	log.Print("[triggerUpdate]")
	prevTrigger := c.configUpdatedTrigger
	c.configUpdatedTrigger = make(chan struct{})
	if prevTrigger != nil {
		close(prevTrigger)
	}
}

func (c *Config) ConfigUpdatedTriggerGet() chan struct{} {
	return c.configUpdatedTrigger
}

func (c *Config) ListenPortGet() string {
	return c.listenPort
}

func (c *Config) MonitorIntervalGet() time.Duration {
	return c.monitorInterval
}
func (c *Config) MonitorIntervalSet(interval time.Duration) {
	c.monitorInterval = interval
}

func (c *Config) WriterIntervalGet() time.Duration {
	return c.writerInterval
}
func (c *Config) SerialDeviceGet() string {
	return c.serialDevice
}
func (c *Config) DbPathGet() string {
	return c.dbPath
}

func (c *Config) ChargeManagementEnabledGet() bool {
	return c.chargeManagementEnabled
}
func (c *Config) ChargeManagementEnabledSet(chargeManagementEnabled bool) {
	c.chargeManagementEnabled = chargeManagementEnabled
}

func (c *Config) StartChargingVoltageGet() uint32 {
	return c.startChargingVoltage
}
func (c *Config) StartChargingVoltageSet(startChargingVoltage uint32) {
	c.startChargingVoltage = startChargingVoltage
}

func (c *Config) StopChargingVoltageGet() uint32 {
	return c.stopChargingVoltage
}
func (c *Config) StopChargingVoltageSet(stopChargingVoltage uint32) {
	c.stopChargingVoltage = stopChargingVoltage
}

func (c *Config) ShutdownVoltageGet() uint32 {
	return c.shutdownVoltage
}
func (c *Config) ShutdownVoltageSet(shutdownVoltage uint32) {
	c.shutdownVoltage = shutdownVoltage
}

func (c *Config) Validate() error {
	if c.chargeManagementEnabled {
		if c.startChargingVoltage >= c.stopChargingVoltage {
			return ErrStartChargeVoltageMoreStop
		}
		if c.startChargingVoltage >= 8400 {
			return ErrStartChargeVoltage
		}
		if c.stopChargingVoltage > 8400 {
			return ErrStopChargeVoltage
		}
		if c.startChargingVoltage <= c.shutdownVoltage {
			return ErrStartChargeVoltageLessShutdown
		}
	}
	if c.shutdownVoltage < 6100 {
		return ErrShutdownVoltage
	}
	return nil
}

func (c *Config) Apply(newConfig *Config) error {

	if err := newConfig.Validate(); err != nil {
		return err
	}

	changed := false

	if newConfig.chargeManagementEnabled != c.chargeManagementEnabled {
		c.chargeManagementEnabled = newConfig.chargeManagementEnabled
		changed = true
	}
	if c.chargeManagementEnabled {
		if newConfig.startChargingVoltage != c.startChargingVoltage {
			c.startChargingVoltage = newConfig.startChargingVoltage
			changed = true
		}
		if newConfig.stopChargingVoltage != c.stopChargingVoltage {
			c.stopChargingVoltage = newConfig.stopChargingVoltage
			changed = true
		}
	}
	if newConfig.shutdownVoltage != c.shutdownVoltage {
		c.shutdownVoltage = newConfig.shutdownVoltage
		changed = true
	}

	if changed {
		c.triggerUpdate()
	}

	return nil
}
