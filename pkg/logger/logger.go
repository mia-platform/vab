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
	"fmt"
	"io"
	"os"
)

// A logger that writes messages on different writers
type Logger struct {
	writer         io.Writer
	errorWriter    io.Writer
	verbosityLevel LogLevel
}

// The streams used by the logger for writing messages
type LogStreams struct {
	OutStream io.Writer
	ErrStream io.Writer
}

const defaultLogLevel = 0

// DefaultStreams for writing log messages
func DefaultStreams() LogStreams {
	return LogStreams{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}
}

// Return a new logger with the verbosity LogLevel that will write on the writer and that will use
// the stdErr writer for the warn messages if useStdErr is set to true
func NewLogger(streams LogStreams) *Logger {
	return &Logger{
		writer:         streams.OutStream,
		errorWriter:    streams.ErrStream,
		verbosityLevel: defaultLogLevel,
	}
}

func (l *Logger) Warn() LogWriterInterface {
	return LogWriter{
		enabled: true,
		writer:  l.errorWriter,
	}
}

func (l *Logger) V(logLevel LogLevel) LogWriterInterface {
	return LogWriter{
		enabled: l.verbosityLevel >= logLevel,
		writer:  l.writer,
	}
}

func (l *Logger) SetLogLevel(logLevel LogLevel) {
	l.verbosityLevel = logLevel
}

type LogWriter struct {
	enabled bool
	writer  io.Writer
}

func (l LogWriter) Write(message string) {
	if !l.enabled {
		return
	}

	fmt.Fprintln(l.writer, message)
}

func (l LogWriter) Writef(format string, args ...interface{}) {
	if !l.enabled {
		return
	}

	message := fmt.Sprintf(format, args...)
	fmt.Fprintln(l.writer, message)
}
