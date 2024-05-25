package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createInt() int {
	return 42
}

type Logger interface {
	Error(message string)
}

type LoggerType interface {
	LoggerTypeId() int
}

type MyPersonalLoggerType struct {
}

func (m MyPersonalLoggerType) LoggerTypeId() int {
	return 1
}

func createLoggerType() LoggerType {
	return MyPersonalLoggerType{}
}

type MyPersonalLogger1 struct {
}

func (m MyPersonalLogger1) Error(message string) {
	fmt.Printf("[MyLogger1] %s\n", message)
}

type MyPersonalLogger2 struct {
}

func (m MyPersonalLogger2) Error(message string) {
	fmt.Printf("[MyLogger2] %s\n", message)
}

func createLogger(t LoggerType) Logger {
	if t.LoggerTypeId() == 1 {
		return MyPersonalLogger1{}
	} else {
		return MyPersonalLogger2{}
	}
}

type ManyLoggers struct {
	L1 Logger `autowired:""`
	L2 Logger `autowired:""`
	L3 Logger
}

func setUp() {
	CleanRegistry()
}

func TestWireIntWithFactory(t *testing.T) {
	setUp()
	err := WireFactory(createInt)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	number := *Construct[int]()
	assert.Equal(t, number, 42)
}

func TestWireLoggerWithFactory(t *testing.T) {
	setUp()
	err := WireFactory(createLogger)
	assert.Nil(t, err)
	err = WireFactory(createLoggerType)
	assert.Nil(t, err)

	l := *Construct[Logger]()
	l.Error("My error!")
}

func TestAutoWireLoggersWithFactory(t *testing.T) {
	setUp()
	err := WireFactory(createLogger)
	assert.Nil(t, err)
	err = WireFactory(createLoggerType)
	assert.Nil(t, err)
	m := ManyLoggers{}
	AutoWire(&m)
	m.L1.Error("My error L1!")
	m.L2.Error("My error L2!")
	assert.Nil(t, m.L3)
}

func TestWireWithType(t *testing.T) {
	setUp()
	err := Wire(new(Logger), new(MyPersonalLogger2))
	assert.Nil(t, err)
	l := *Construct[Logger]()
	l.Error("Error type wired")
}
