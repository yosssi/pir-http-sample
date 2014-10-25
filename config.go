package main

// Config represents configuration for the main process.
type Config struct {
	MotionPinNo int `yaml:"motion_pin_no"`
	LEDPinNo    int `yaml:"led_pin_no"`
}
