//go:build !authtest
// +build !authtest

/*
    _____           _____   _____   ____          ______  _____  ------
   |     |  |      |     | |     | |     |     | |       |            |
   |     |  |      |     | |     | |     |     | |       |            |
   | --- |  |      |     | |-----| |---- |     | |-----| |-----  ------
   |     |  |      |     | |     | |     |     |       | |       |
   | ____|  |_____ | ____| | ____| |     |_____|  _____| |_____  |_____


   Licensed under the MIT License <http://opensource.org/licenses/MIT>.

   Copyright © 2020-2023 Microsoft Corporation. All rights reserved.
   Author : <blobfusedev@microsoft.com>

   Permission is hereby granted, free of charge, to any person obtaining a copy
   of this software and associated documentation files (the "Software"), to deal
   in the Software without restriction, including without limitation the rights
   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
   copies of the Software, and to permit persons to whom the Software is
   furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all
   copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
   SOFTWARE
*/

package block_cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type threadPoolTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *threadPoolTestSuite) SetupTest() {
}

func (suite *threadPoolTestSuite) cleanupTest() {
}

func (suite *threadPoolTestSuite) TestCreate() {
	suite.assert = assert.New(suite.T())

	tp := newThreadPool(0, nil)
	suite.assert.Nil(tp)

	tp = newThreadPool(1, nil)
	suite.assert.Nil(tp)

	tp = newThreadPool(1, func(interface{}) {})
	suite.assert.NotNil(tp)
	suite.assert.Equal(tp.worker, uint32(1))
}

func (suite *threadPoolTestSuite) TestStartStop() {
	suite.assert = assert.New(suite.T())

	r := func(i interface{}) {
		suite.assert.Equal(i.(int), 1)
	}
	tp := newThreadPool(2, r)
	suite.assert.NotNil(tp)
	suite.assert.Equal(tp.worker, uint32(2))

	tp.Start()
	suite.assert.NotNil(tp.priorityCh)
	suite.assert.NotNil(tp.normalCh)

	tp.Stop()
}

func (suite *threadPoolTestSuite) TestSchedule() {
	suite.assert = assert.New(suite.T())

	r := func(i interface{}) {
		suite.assert.Equal(i.(int), 1)
	}

	tp := newThreadPool(2, r)
	suite.assert.NotNil(tp)
	suite.assert.Equal(tp.worker, uint32(2))

	tp.Start()
	suite.assert.NotNil(tp.priorityCh)
	suite.assert.NotNil(tp.normalCh)

	tp.Schedule(false, 1)
	tp.Schedule(true, 1)

	tp.Stop()
}

func TestThreadPoolSuite(t *testing.T) {
	suite.Run(t, new(threadPoolTestSuite))
}