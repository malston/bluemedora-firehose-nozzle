/**Copyright Blue Medora Inc. 2016**/

package logger

import (
    "testing"
    "os"
    "fmt"
    "io/ioutil"
    "strings"
)

var (
    logDirectory = "../logs"
    testLog = "test log"
)

func TestLogDirectoryCreation(t *testing.T) {
    //Setup Envrionment
    setupEnvrionment(t)
    
    New(logDirectory)
    
    //See if log file was created
    checkLogFileExists(t)
}

func TestLogFileContents(t *testing.T) {
    //Setup Enivronment
    setupEnvrionment(t)
    
    logger := New(logDirectory)
    logger.Info(testLog)
    
    //Test if log contents contains test string
    checkLogContents(t)
}

func setupEnvrionment(t *testing.T) {
    t.Log("Removing logs directory...")
    if _, err := os.Stat(logDirectory); err != nil {
        os.RemoveAll(logDirectory)
    }
    t.Log("Removed logs directory")
}

func checkLogFileExists(t *testing.T) {
    t.Log("Check if log file bm_nozzle.log exists...")
    if _,err := os.Stat(fmt.Sprintf("%s/bm_nozzle.log", logDirectory)); os.IsNotExist(err) {
        t.Fatalf("Log file bm_nozzle.log not created")
    }
}

func checkLogContents(t *testing.T) {
    fileBuffer, err := ioutil.ReadFile(fmt.Sprintf("%s/bm_nozzle.log", logDirectory))
    if err != nil {
        t.Fatalf("Failed to load log file: %s", err)
    }
    
    fileString := string(fileBuffer)
    
    t.Logf("Checking log contents... (expecting log contains: %s)", testLog)
    if !strings.Contains(fileString, testLog) {
        t.Errorf("Expected log contents to contain %s, string was not in log", testLog)
    }
}