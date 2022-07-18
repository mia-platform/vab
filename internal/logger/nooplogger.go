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

// The DisabledLogger implement the LogInterface but all the functions are noop
type DisabledLogger struct{}

// implement the LogInterface

// Warn meets the LogInterface interface but does nothing
func (logger DisabledLogger) Warn(message string) {}

// Warnf meets the LogInterface interface but does nothing
func (logger DisabledLogger) Warnf(format string, args ...interface{}) {}

// V meets the Logger interface but does nothing
func (logger DisabledLogger) V(level LogLevel) InfoLogWriterInterface {
	return DisabledInfoLogWriter{}
}

// SetLogLevel meets the Logger interface but does nothing
func (logger DisabledLogger) SetLogLevel(logLevel LogLevel) {}

// The DisabledInfoLogWriter implement the InfoLogWriterInterface but all the functions are noop
type DisabledInfoLogWriter struct{}

// Info meets the InfoLogWriterInterface interface but does nothing
func (logger DisabledInfoLogWriter) Info(message string) {}

// Infof meets the InfoLogWriterInterface interface but does nothing
func (logger DisabledInfoLogWriter) Infof(message string, args ...interface{}) {}
