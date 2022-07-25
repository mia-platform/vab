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

// The DisabledLogger implement the LogInterface but all the functions are noop
type DisabledLogger struct{}

// implement the LogInterface

// Warn meets the LogInterface interface but does nothing
func (logger DisabledLogger) Warn() LogWriterInterface {
	return DisabledLogWriter{}
}

// V meets the Logger interface but does nothing
func (logger DisabledLogger) V(level LogLevel) LogWriterInterface {
	return DisabledLogWriter{}
}

// SetLogLevel meets the Logger interface but does nothing
func (logger DisabledLogger) SetLogLevel(logLevel LogLevel) {}

// The DisabledLogWriter implement the LogWriterInterface but all the functions are noop
type DisabledLogWriter struct{}

// Info meets the LogWriterInterface interface but does nothing
func (logger DisabledLogWriter) Write(message string) {}

// Writef meets the LogWriterInterface interface but does nothing
func (logger DisabledLogWriter) Writef(message string, args ...interface{}) {}
