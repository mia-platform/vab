//go:build !e2e
// +build !e2e

// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMessage = "Test message"

// Test that the default Streams are the stardard output and error of the os
func TestDefaultStreams(t *testing.T) {
	defaultStream := DefaultStreams()

	assert.Equal(t, defaultStream.OutStream, os.Stdout, "Unexpected default out stream")
	assert.Equal(t, defaultStream.ErrStream, os.Stderr, "Unexpected default error stream")
}

// Test that a newly created logger will log warning messages to buffer if is set to not using stderr
func TestNewLogger(t *testing.T) {
	expectedMessage := fmt.Sprintf("%s\n%s\n", testMessage, testMessage)
	stdBuffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	streams := LogStreams{
		ErrStream: errBuffer,
		OutStream: stdBuffer,
	}
	logger := NewLogger(streams)

	logger.Warn().Writef("%s", testMessage)
	logger.Warn().Write(testMessage)
	assert.Equal(t, errBuffer.String(), expectedMessage)

	errBuffer.Reset()
	// default logger has 0 level set as default so info messages are logged
	logger.V(0).Write(testMessage)
	logger.V(0).Writef("%s", testMessage)
	assert.Equal(t, stdBuffer.String(), expectedMessage)

	stdBuffer.Reset()
	// 1 or more level are grater than the default so no log is written
	logger.V(1).Write(testMessage)
	logger.V(1).Writef("%s", testMessage)
	assert.Equal(t, stdBuffer.String(), "")
}

func TestChangingLogLevel(t *testing.T) {
	stdBuffer := new(bytes.Buffer)
	streams := LogStreams{
		ErrStream: stdBuffer,
		OutStream: stdBuffer,
	}
	logger := NewLogger(streams)

	// 1 or more level are grater than the default so no log is written
	logger.V(1).Write(testMessage)
	assert.Equal(t, stdBuffer.String(), "")

	logger.SetLogLevel(3)

	// now new info logger must log from level 3 and below
	logger.V(2).Write(testMessage)
	assert.Equal(t, stdBuffer.String(), fmt.Sprintf("%s\n", testMessage))

	// Logger can be set as LogInterface variable without issues
	var logInterface LogInterface = NewLogger(streams)
	logInterface.SetLogLevel(3)
}

func TestDisabledLogger(t *testing.T) {
	logger := DisabledLogger{}
	// Disabled logger implement standard methods
	logger.Warn().Write(testMessage)
	logger.Warn().Writef("%s", testMessage)
	logger.V(1).Write(testMessage)
	logger.V(1).Writef("%s", testMessage)
	logger.SetLogLevel(3)

	// DisabledLogger can be set as a LogInterface variable without issues
	var logInterface LogInterface = DisabledLogger{}
	logInterface.SetLogLevel(3)
}
