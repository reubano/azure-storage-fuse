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

type blockTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *blockTestSuite) SetupTest() {
}

func (suite *blockTestSuite) cleanupTest() {
}

func (suite *blockTestSuite) TestAllocate() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(0)
	suite.assert.Nil(b)
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid size")

	b, err = AllocateBlock(10)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	suite.assert.NotNil(b.data)
	suite.assert.Equal(b.closed, true)
	suite.assert.Equal(b.id, uint64(0))

	_ = b.Delete()
}

func (suite *blockTestSuite) TestAllocateBig() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(100 * 1024 * 1024)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	suite.assert.NotNil(b.data)
	suite.assert.Equal(b.closed, true)
	suite.assert.Equal(cap(b.data), 100*1024*1024)

	b.Delete()
}

func (suite *blockTestSuite) TestAllocateHuge() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(50 * 1024 * 1024 * 1024)
	suite.assert.Nil(b)
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "mmap error")
}

func (suite *blockTestSuite) TestFreeInvalid() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(1)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	b.data = nil

	err = b.Delete()
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid buffer")
}

func (suite *blockTestSuite) TestFreeInvalid2() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(1)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	b.data = make([]byte, 1)

	err = b.Delete()
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid argument")
}

func (suite *blockTestSuite) TestResuse() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(1)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	suite.assert.Nil(b.state)

	b.ReUse()
	suite.assert.NotNil(b.state)
	suite.assert.Equal(b.closed, false)

	_ = b.Delete()
}

func (suite *blockTestSuite) TestReadyForReading() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(1)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	suite.assert.Nil(b.state)

	err = b.ReadyForReading()
	suite.assert.NotNil(err)
	suite.assert.Nil(b.state)
	suite.assert.Equal(b.closed, true)
	suite.assert.Contains(err.Error(), "block was never used")

	b.ReUse()
	suite.assert.NotNil(b.state)
	suite.assert.Equal(b.closed, false)

	err = b.ReadyForReading()
	suite.assert.Nil(err)
	suite.assert.NotNil(b.state)
	suite.assert.Equal(b.closed, false)
	suite.assert.Equal(len(b.state), 2)

	<-b.state
	suite.assert.Equal(len(b.state), 1)

	err = b.ReadyForReading()
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid state to mark it ready")

	<-b.state
	suite.assert.Equal(len(b.state), 0)

	_ = b.Delete()
}

func (suite *blockTestSuite) TestUnBlock() {
	suite.assert = assert.New(suite.T())

	b, err := AllocateBlock(1)
	suite.assert.NotNil(b)
	suite.assert.Nil(err)
	suite.assert.Nil(b.state)

	b.ReUse()
	suite.assert.NotNil(b.state)
	suite.assert.Equal(b.closed, false)

	err = b.ReadyForReading()
	suite.assert.Nil(err)
	suite.assert.NotNil(b.state)
	suite.assert.Equal(b.closed, false)
	suite.assert.Equal(len(b.state), 2)

	err = b.Unblock()
	suite.assert.Nil(err)
	suite.assert.Equal(b.closed, true)

	<-b.state
	suite.assert.Equal(len(b.state), 1)
	<-b.state
	suite.assert.Equal(len(b.state), 0)

	<-b.state
	suite.assert.Equal(len(b.state), 0)

	err = b.Unblock()
	suite.assert.NotNil(err)
	suite.assert.Equal(err.Error(), "invalid state of block to unblock")

	_ = b.Delete()
}

func TestBlockSuite(t *testing.T) {
	suite.Run(t, new(blockTestSuite))
}