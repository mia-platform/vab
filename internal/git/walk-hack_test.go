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

// WARNING This file is a copy of the Walk implementation in billy package because is not publicly available. Remove
// the file and use the official one once is made available with a new version https://github.com/go-git/go-billy/releases

package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"

	. "gopkg.in/check.v1" //revive:disable-line:dot-imports
)

type WalkSuite struct{}

func TestWalk(t *testing.T) { TestingT(t) }

var _ = Suite(&WalkSuite{})

func (s *WalkSuite) TestWalkCanSkipTopDirectory(c *C) {
	filesystem := memfs.New()
	c.Assert(Walk(filesystem, "/root/that/does/not/exist", func(_ string, _ os.FileInfo, _ error) error { return filepath.SkipDir }), IsNil)
}

func (s *WalkSuite) TestWalkReturnsAnErrorWhenRootDoesNotExist(c *C) {
	filesystem := memfs.New()
	c.Assert(Walk(filesystem, "/root/that/does/not/exist", func(_ string, _ os.FileInfo, err error) error { return err }), NotNil)
}

func (s *WalkSuite) TestWalkOnPlainFile(c *C) {
	filesystem := memfs.New()
	createFile(c, filesystem, "./README.md")
	discoveredPaths := []string{}
	c.Assert(Walk(filesystem, "./README.md", func(path string, _ os.FileInfo, _ error) error {
		discoveredPaths = append(discoveredPaths, path)
		return nil
	}), IsNil)
	c.Assert(discoveredPaths, DeepEquals, []string{"./README.md"})
}

func (s *WalkSuite) TestWalkOnExistingFolder(c *C) {
	filesystem := memfs.New()
	createFile(c, filesystem, "path/to/some/subfolder/that/contain/file")
	createFile(c, filesystem, "path/to/some/file")
	discoveredPaths := []string{}
	c.Assert(Walk(filesystem, "path", func(path string, _ os.FileInfo, _ error) error {
		discoveredPaths = append(discoveredPaths, path)
		return nil
	}), IsNil)
	c.Assert(discoveredPaths, Contains, "path")
	c.Assert(discoveredPaths, Contains, "path/to")
	c.Assert(discoveredPaths, Contains, "path/to/some")
	c.Assert(discoveredPaths, Contains, "path/to/some/file")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder/that")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder/that/contain")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder/that/contain/file")
}

func (s *WalkSuite) TestWalkCanSkipFolder(c *C) {
	filesystem := memfs.New()
	createFile(c, filesystem, "path/to/some/subfolder/that/contain/file")
	createFile(c, filesystem, "path/to/some/file")
	discoveredPaths := []string{}
	c.Assert(Walk(filesystem, "path", func(path string, _ os.FileInfo, _ error) error {
		discoveredPaths = append(discoveredPaths, path)
		if path == "path/to/some/subfolder" {
			return filepath.SkipDir
		}
		return nil
	}), IsNil)
	c.Assert(discoveredPaths, Contains, "path")
	c.Assert(discoveredPaths, Contains, "path/to")
	c.Assert(discoveredPaths, Contains, "path/to/some")
	c.Assert(discoveredPaths, Contains, "path/to/some/file")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that/contain")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that/contain/file")
}

func (s *WalkSuite) TestWalkStopsOnError(c *C) {
	filesystem := memfs.New()
	createFile(c, filesystem, "path/to/some/subfolder/that/contain/file")
	createFile(c, filesystem, "path/to/some/file")
	discoveredPaths := []string{}
	c.Assert(Walk(filesystem, "path", func(path string, _ os.FileInfo, _ error) error {
		discoveredPaths = append(discoveredPaths, path)
		if path == "path/to/some/subfolder" {
			return errors.New("uncaught error")
		}
		return nil
	}), NotNil)
	c.Assert(discoveredPaths, Contains, "path")
	c.Assert(discoveredPaths, Contains, "path/to")
	c.Assert(discoveredPaths, Contains, "path/to/some")
	c.Assert(discoveredPaths, Contains, "path/to/some/file")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that/contain")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that/contain/file")
}

func (s *WalkSuite) TestWalkForwardsStatErrors(c *C) {
	memFilesystem := memfs.New()
	filesystem := &fnFs{
		Filesystem: memFilesystem,
		lstat: func(path string) (os.FileInfo, error) {
			if path == "path/to/some/subfolder" {
				return nil, errors.New("uncaught error")
			}
			return memFilesystem.Lstat(path)
		},
	}

	createFile(c, filesystem, "path/to/some/subfolder/that/contain/file")
	createFile(c, filesystem, "path/to/some/file")
	discoveredPaths := []string{}
	c.Assert(Walk(filesystem, "path", func(path string, _ os.FileInfo, err error) error {
		discoveredPaths = append(discoveredPaths, path)
		if path == "path/to/some/subfolder" {
			c.Assert(err, NotNil)
		}
		return err
	}), NotNil)
	c.Assert(discoveredPaths, Contains, "path")
	c.Assert(discoveredPaths, Contains, "path/to")
	c.Assert(discoveredPaths, Contains, "path/to/some")
	c.Assert(discoveredPaths, Contains, "path/to/some/file")
	c.Assert(discoveredPaths, Contains, "path/to/some/subfolder")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that/contain")
	c.Assert(discoveredPaths, NotContain, "path/to/some/subfolder/that/contain/file")
}

func createFile(c *C, filesystem billy.Filesystem, path string) {
	fd, err := filesystem.Create(path)
	c.Assert(err, IsNil)
	if err != nil {
		fd.Close()
	}
}

type fnFs struct {
	billy.Filesystem
	lstat func(path string) (os.FileInfo, error)
}

func (f *fnFs) Lstat(path string) (os.FileInfo, error) {
	if f.lstat != nil {
		return f.lstat(path)
	}
	return nil, errors.New("not implemented")
}

type containsChecker struct {
	*CheckerInfo
}

func (checker *containsChecker) Check(params []interface{}, _ []string) (result bool, err string) {
	defer func() {
		if v := recover(); v != nil {
			result = false
			err = fmt.Sprint(v)
		}
	}()

	value := reflect.ValueOf(params[0])
	result = false
	err = fmt.Sprintf("%v does not contain %v", params[0], params[1])
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			r := reflect.DeepEqual(value.Index(i).Interface(), params[1])
			if r {
				result = true
				err = ""
			}
		}
	default:
		return false, "obtained value type is not iterable"
	}
	return
}

var Contains Checker = &containsChecker{
	&CheckerInfo{Name: "Contains", Params: []string{"obtained", "expected"}},
}

var NotContain Checker = Not(Contains)
