package main

// Config represents configuration for the main process.
type Config struct {
	Host         string `yaml:"host"`
	MotionPinNo  int    `yaml:"motion_pin_no"`
	LEDPinNo     int    `yaml:"led_pin_no"`
	httpHostPort string `yaml:"http_host_port"`
}
