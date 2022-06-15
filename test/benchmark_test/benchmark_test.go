// +build !unittest

package benchmark_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/suite"
)

var mntPath = "benchmark"
var n = 5
var sizes = []float32{0.5, 1, 2, 3, 4}

type benchmarkSuite struct {
	suite.Suite
}

func createSingleFile(size float32, path string) (float64, error) {
	sizeInBytes := int(size * 1024 * 1024 * 1024)
	buffer := make([]byte, sizeInBytes)

	start := time.Now()

	err := ioutil.WriteFile(path, buffer, os.FileMode(0755))
	if err != nil {
		return 0, err
	}
	return float64(time.Now().Sub(start)), nil
}

func (suite *benchmarkSuite) TestCreateSingleFiles() {
	for _, size := range sizes {
		suite.T().Logf("\nCreating %.2fGB files: \n", size)
		durs := make([]float64, 0)
		for i := 0; i < n; i++ {
			suite.T().Logf("Round %d\n", i)
			fileName := filepath.Join(mntPath, "testfile"+strconv.Itoa(i))
			completionTime, err := createSingleFile(size, fileName)
			if err != nil {
				suite.T().Logf("error creating file of size %.2fGB for %d run [%s]\n", size, i, err)
			}
			durs = append(durs, completionTime)
			os.Remove(fileName)
		}
		mean, err := stats.Mean(durs)
		if err != nil {
			suite.T().Logf("error computing mean for %.2fGB file [%s]\n", size, err)
		}
		suite.T().Logf("\nMean create time for %.2fGB files=%s\n", size, time.Duration(mean))
		std, err := stats.StandardDeviation(durs)
		if err != nil {
			suite.T().Logf("error computing std for %.2fGB file [%s]\n", size, err)
		}
		suite.T().Logf("Standard Deviation of create files for %.2fGB files=%s\n", size, time.Duration(std))
	}
}

func TestBenchmarkSuite(t *testing.T) {
	suite.Run(t, new(benchmarkSuite))
}

func TestMain(m *testing.M) {
	pathFlag := flag.String("mnt-path", ".", "Mount Path of container")
	nFlag := flag.Int("n", 5, "Number of times to run a test.")
	fileSizesFlag := flag.String("sizes", "0.5,1,2,3,4", "List different sizes of uploads to run. All values are specified in GBs. Default sizes=0.5,1,2,3,4")

	flag.Parse()

	mntPath = filepath.Join(*pathFlag, mntPath)
	err := os.RemoveAll(mntPath)
	if err != nil {
		fmt.Printf("error cleaning up base dir %s [%s]", mntPath, err)
	}
	err = os.Mkdir(mntPath, os.FileMode(0755))
	if err != nil {
		fmt.Printf("error mkdir for base dir %s [%s]", mntPath, err)
	}

	n = *nFlag

	//Parse size flags
	sizeStrs := strings.Split(*fileSizesFlag, ",")
	tempSizes := make([]float32, 0)
	success := true
	for _, i := range sizeStrs {
		sizeStr := strings.TrimSpace(i)
		size, err := strconv.ParseFloat(sizeStr, 32)
		if err != nil {
			success = false
			break
		}
		tempSizes = append(tempSizes, float32(size))
	}
	if success {
		sizes = tempSizes
	}

	exit := m.Run()
	err = os.RemoveAll(mntPath)
	if err != nil {
		fmt.Printf("error post test clean up %s [%s]", mntPath, err)
	}
	os.Exit(exit)
}