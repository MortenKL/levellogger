package levellogger_test

import (
	//"log"

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	log "github.com/MortenKL/levellogger"
)

//Some times a logfile is deleted, and the logger should then be able to recreate the logfile
//We have a flag for that
func Test_LogFileDisappears(t *testing.T) {

	log.FatalCausesExit = false
	log.SetLogFile(log.LINFO|log.LDEBUG, "log.2.txt")
	log.SetLogFile(log.LWARN|log.LERROR|log.LFATAL, "log.3.txt")
	log.Print("Test logging 01")
	log.Debug("Test logging 02")
	os.Remove("log.2.txt")
	log.Print("Test logging 03")
	log.Debug("Test logging 04")
	log.Error("Test logging 05")
	log.Warn("Test logging 06")
	os.Remove("log.3.txt")
	log.Error("Test logging 07")
	log.Warn("Test logging 08")
	log.Fatal("Test logging 09")

}

//Test of combined flags when setting log files
func Test_BitMasks(t *testing.T) {

	log.FatalCausesExit = false
	log.SetLogFile(log.LWARN|log.LERROR|log.LFATAL, "WEF.txt")
	log.SetLogFile(log.LDEBUG|log.LINFO, "DI.txt")

	log.Debug("Log test 01")
	log.Info("Log test 02")
	log.Warn("log test 03")
	log.Error("Log test 04")
	log.Fatal("Log test 05")

	log.SetLogFile(log.LALL, "ALL.txt")

	log.Debug("Log test 01")
	log.Info("Log test 02")
	log.Warn("log test 03")
	log.Error("Log test 04")
	log.Fatal("Log test 05")

}

//Testing the functionality for multiple log outputs
//The test should output line 01 to stdout line 02 and 07 log.txt and lines 03,04 and 05 to errors.txt
//OBS! errors.txt should be deleted after test
func Test_LogOutput(t *testing.T) {

	defer func() {
		//This defer is used to recover from testing the panic button below
		if err := recover(); err != nil {
			t.Log("Panic recovered:", err)
			log.Info("And We're done 07")

			if f2, err3 := ioutil.ReadFile("errors.txt"); err3 == nil {
				lst := strings.Split(string(f2), "\n")
				if len(lst) != 5 {
					t.Log("Too many or few log lines in errors.txt")
					t.FailNow()
				}
			}

			if f1, err2 := ioutil.ReadFile("log.txt"); err2 == nil {
				lst := strings.Split(string(f1), "\n")
				if len(lst) != 3 {
					t.Log("Too many or few log lines in log.txt")
					t.Fail()
					t.FailNow()
				}

			}

		}
	}()

	log.FatalCausesExit = false //Don't want to exit on Fataæ logs
	//Starting with output to stderr using Print, Print is INFO Level
	log.Print("Test logging 01")

	//Setting general output to the file log.txt
	f, _ := os.Create("log.txt")
	log.SetOutput(f)
	log.Print("Test logging 02")

	//Setting all logs from ERROR level to a specific file
	c := log.SetLogFile(log.LERROR, "errors.txt")
	//I want WARN and FATAL in same output as Error, use SetLogOutput to same filedescriptor
	log.SetLogOutput(log.LWARN, c)
	log.SetLogOutput(log.LFATAL, c)

	log.Warn("Test logging 03, this is a warning")
	log.Error("Test logging 04, this is an error")

	//This will make the test fail..
	log.Fatal("Test logging 05, this is fatal")
	log.Panic("Total panic 06") //Panic will cause a fatal log line and panic away

}

//Testing the exclusive level setting,
//The test should give an output of log lines 02 and 07
func Test_LogLevel(t *testing.T) {

	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetLevel("Info")

	log.Debug("Test logging 01")
	log.Print("Test logging 02")

	log.SetLoglevel(log.LERROR)

	log.Debug("Test logging 03")
	log.Println("Test logging 04")
	log.Info("Test logging 05")
	log.Warn("Test logging 06")
	log.Error("Test logging 07")

	lst := strings.Split(buf.String(), "\n")
	if len(lst) > 3 {
		t.Log("Too many lines in output")
		t.FailNow()
	}
	if strings.HasSuffix(lst[0], "02") && strings.HasSuffix(lst[1], "07") {
		t.Log("All fine")
		fmt.Print(buf.String()) //Printing the log bugger for visual confirmation
	} else {
		t.Log("Unexpected lines in output")
		t.FailNow()
	}

}
