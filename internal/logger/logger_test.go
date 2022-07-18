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
)

const testMessage = "Test message"

// Test that the default Streams are the stardard output and error of the os
func TestDefaultStreams(t *testing.T) {
	defaultStream := DefaultStreams()

	if defaultStream.OutStream != os.Stdout {
		t.Fatal("Unexpected default out stream")
	}

	if defaultStream.ErrStream != os.Stderr {
		t.Fatal("Unexpected default error stream")
	}
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

	logger.Warnf("%s", testMessage)
	logger.Warn(testMessage)
	if !bytes.Equal(errBuffer.Bytes(), []byte(expectedMessage)) {
		t.Fatalf("Unexpected message on buffer %s", stdBuffer.String())
	}

	errBuffer.Reset()
	// default logger has 0 level set as default so info messages are logged
	logger.V(0).Info(testMessage)
	logger.V(0).Infof("%s", testMessage)
	if !bytes.Equal(stdBuffer.Bytes(), []byte(expectedMessage)) {
		t.Fatalf("Unexpected message on buffer %s", stdBuffer.String())
	}

	stdBuffer.Reset()
	// 1 or more level are grater than the default so no log is written
	logger.V(1).Info(testMessage)
	logger.V(1).Infof("%s", testMessage)
	if !bytes.Equal(stdBuffer.Bytes(), []byte("")) {
		t.Fatalf("Unexpected message on buffer %s", stdBuffer.String())
	}
}

func TestChangingLogLevel(t *testing.T) {
	stdBuffer := new(bytes.Buffer)
	streams := LogStreams{
		ErrStream: stdBuffer,
		OutStream: stdBuffer,
	}
	logger := NewLogger(streams)

	// 1 or more level are grater than the default so no log is written
	logger.V(1).Info(testMessage)
	if !bytes.Equal(stdBuffer.Bytes(), []byte("")) {
		t.Fatalf("Unexpected message on buffer %s", stdBuffer.String())
	}

	logger.SetLogLevel(3)

	// now new info logger must log from level 3 and below
	logger.V(2).Info(testMessage)
	if !bytes.Equal(stdBuffer.Bytes(), []byte(fmt.Sprintf("%s\n", testMessage))) {
		t.Fatalf("Unexpected message on buffer %s", stdBuffer.String())
	}

	// Logger can be set as LogInterface variable without issues
	var logInterface LogInterface = NewLogger(streams)
	logInterface.SetLogLevel(3)
}

func TestDisabledLogger(t *testing.T) {
	logger := DisabledLogger{}
	// Disabled logger implement standard methods
	logger.Warn(testMessage)
	logger.Warnf("%s", testMessage)
	logger.V(1).Info(testMessage)
	logger.V(1).Infof("%s", testMessage)
	logger.SetLogLevel(3)

	// DisabledLogger can be set as a LogInterface variable without issues
	var logInterface LogInterface = DisabledLogger{}
	logInterface.SetLogLevel(3)
}
