// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

// The noisiness level of the logger, we recommend to not use a level greater than 10
// because already 10 level of diverse logs are silly, using all 255 available can be too much
type LogLevel uint8

// LogInterface is a simple interface that must be adopted by a logger
type LogInterface interface {
	// Warn should be used to return a LogWriterInterface that can always write to the user,
	// it may use the stdErr stream
	Warn() LogWriterInterface

	// V() returns an LogWriterInterface for a given verbosity LogLevel
	//
	// Normal verbosity levels:
	// V(0): normal user facing messages go to V(0)
	// V(1): debug messages start when V(N > 0), these should be high level
	// V(2+): trace level logging, in increasing "noisiness"
	//
	// It is expected that the returned LogWriterInterface will implement a near to noop
	// function for LogLevel major than the enabled one
	V(LogLevel) LogWriterInterface

	// Change the LogLevel set for the logger to the new logLevel
	SetLogLevel(logLevel LogLevel)
}

// LogWriterInterface defines a logger interface for wrinting messages
type LogWriterInterface interface {
	// Write is used to write a user facing status message
	Write(message string)
	// Writef is used to write a Printf style user facing status message
	Writef(format string, args ...interface{})
}
