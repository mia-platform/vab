// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

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
	// Warn should be used to write user facing warnings, it can be outputted to stderr
	// and it cannot be turn off via LogLevel
	Warn(message string)
	// Warnf should be used to write Printf style user facing warnings, it can be outputted to stderr
	// and it cannot be turn off via LogLevel
	Warnf(format string, args ...interface{})

	// V() returns an InfoLogWriterInterface for a given verbosity LogLevel
	//
	// Normal verbosity levels:
	// V(0): normal user facing messages go to V(0)
	// V(1): debug messages start when V(N > 0), these should be high level
	// V(2+): trace level logging, in increasing "noisiness"
	//
	// It is expected that the returned InfoLogWriterInterface will implement a near to noop
	// function for LogLevel major than the enabled one
	V(LogLevel) InfoLogWriterInterface

	// Change the LogLevel set for the logger to the new logLevel
	SetLogLevel(logLevel LogLevel)
}

// InfoLogWriterInterface defines a logger interface for wrinting messages if the set LogLevel
// is minor than the enabled one
type InfoLogWriterInterface interface {
	// Info is used to write a user facing status message
	Info(message string)
	// Infof is used to write a Printf style user facing status message
	Infof(format string, args ...interface{})
}
