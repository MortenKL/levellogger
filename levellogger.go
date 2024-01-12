//Package levellogger a leveled logger, a log wrapper / drop in, using multiple std loggers
package levellogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var output io.Writer
var ll LLogger
var logFlags int

//Loglevel determines the loglevel
type Loglevel uint8

const (
	//LDEBUG the debug level
	LDEBUG Loglevel = 1 << iota
	//LINFO the info level
	LINFO
	//LWARN the warning level
	LWARN
	//LERROR the error level
	LERROR
	//LFATAL yeah well ..
	LFATAL
)

//LALL all log levels
const LALL Loglevel = LDEBUG | LINFO | LWARN | LERROR | LFATAL

//SplitLevels is used to determine if we need separate level files
//var SplitLevels bool

//RecreateLogfiles recreates the log file if it has dissapeared
var RecreateLogfiles bool

//FatalCausesExit determines if we need a log.Fatal to cause os.Exit(), folowing the dtandard is true
var FatalCausesExit bool

//Level is obviosly the LogLevel
var Level Loglevel

//Ldate the date in the local time zone: 2009/01/23
var Ldate int

//Ltime the time in the local time zone: 01:23:23
var Ltime int

//Lmicroseconds microsecond resolution: 01:23:23.123123.  assumes Ltime.
var Lmicroseconds int

//Llongfile full file name and line number: /a/b/c/d.go:23
var Llongfile int

//Lshortfile final file name element and line number: d.go:23. overrides Llongfile
var Lshortfile int

//LUTC if Ldate or Ltime is set, use UTC rather than the local time zone
var LUTC int

//LstdFlags initial values for the standard logger
var LstdFlags int

//var Lmsgprefix int    // move the "prefix" from the beginning of the line to before the message

func init() {
	FatalCausesExit = true
	output = os.Stderr
	Ldate = log.Ldate
	Ltime = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile = log.Llongfile
	Lshortfile = log.Lshortfile
	LUTC = log.LUTC
	LstdFlags = log.LstdFlags
	//Lmsgprefix =log.L
	ll = LLogger{}
	ll.filemap = make(map[string]*os.File)
	logFlags = Ldate | Ltime | Lshortfile
	log.SetPrefix("STD ")
}

//LLogger is short for LevelLogger.. figures
type LLogger struct {
	//std           *log.Logger
	info          *log.Logger
	infofilename  string
	debug         *log.Logger
	debugfilename string
	warn          *log.Logger
	warnfilename  string
	err           *log.Logger
	errfilename   string
	fatal         *log.Logger
	fatalfilename string
	//logfileName   string
	logoutput   io.Writer
	debugoutput io.Writer
	infooutput  io.Writer
	warnoutput  io.Writer
	erroutput   io.Writer
	fataloutput io.Writer
	filemap     map[string]*os.File
}

//Std logger interface

//Fatal logs fatal lines
func Fatal(v ...interface{}) {
	ll.GetLogger(LFATAL).Output(2, fmt.Sprint(v...))
	if FatalCausesExit {
		os.Exit(1)
	}

}

//Fatalf logs fatal lines
func Fatalf(format string, v ...interface{}) {
	ll.GetLogger(LFATAL).Output(2, fmt.Sprintf(format, v...))
	if FatalCausesExit {
		os.Exit(1)
	}

}

//Fatalln logs fatal lines
func Fatalln(v ...interface{}) {
	ll.GetLogger(LFATAL).Output(2, fmt.Sprintln(v...))
	if FatalCausesExit {
		os.Exit(1)
	}

}

//Flags gets the flags, doh!
func Flags() int {
	return log.Flags()
}

//Output hmm tricky bastard..
func Output(calldepth int, s string) error {
	return log.Output(calldepth, s)
}

//Panic but dont panic..
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	ll.GetLogger(LFATAL).Output(2, s)
	panic(s)
}

//Panicf but dont ..
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	ll.GetLogger(LFATAL).Output(2, s)
	panic(s)
}

//Panicln but dont ..
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	ll.GetLogger(LFATAL).Output(2, s)
	panic(s)

}

//Prefix returns, well the prefix
func Prefix() string {
	return log.Prefix()
}

//Print a log line
func Print(v ...interface{}) {
	if !(Level > LINFO) {
		ll.GetLogger(LINFO).Output(2, fmt.Sprint(v...))
	}

}

//Printf a log line
func Printf(format string, v ...interface{}) {
	if !(Level > LINFO) {
		ll.GetLogger(LINFO).Output(2, fmt.Sprintf(format, v...))
	}
}

//Println a log line
func Println(v ...interface{}) {
	if !(Level > LINFO) {
		ll.GetLogger(LINFO).Output(2, fmt.Sprintln(v...))
	}
}

//SetFlags sets some..flags i suppose ?!
func SetFlags(flag int) {
	ll.GetLogger(Level).SetFlags(flag)
}

//SetOutput sets the Output of the log
func SetOutput(w io.Writer) {
	Close()
	output = w
	ll.logoutput = output
	//log.SetOutput(w)

}

/*
func SetPrefix(prefix string) {
	log.SetPrefix(prefix)
}
*/

//Writer hmm what is it
func Writer() io.Writer {
	return output
}

/*
func NewLLogger(stdfile string) LLogger {
	var ll lLogger
	if len(stdfile) > 0 {
		var err error
		ll.logfile, err = os.OpenFile(stdfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panic("Unable to open output file")
		}
		ll.logfileName = stdfile
	} else {
		ll.logfile = os.Stderr
	}
	ll.std = log.New(ll.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return &ll
}
*/

//GetLogger gets as speperate logger instance for the specified level
func (ll *LLogger) GetLogger(level Loglevel) *log.Logger {

	switch level {
	case LINFO:
		if ll.info == nil {
			ll.info = log.New(output, "INFO ", logFlags)
		}
		checkLogLevelFilename(LINFO)
		return ll.info
	case LDEBUG:
		if ll.debug == nil {
			ll.debug = log.New(output, "DEBUG ", logFlags)
		}
		checkLogLevelFilename(LDEBUG)
		return ll.debug
	case LWARN:
		if ll.warn == nil {
			ll.warn = log.New(output, "WARN ", logFlags)
		}
		checkLogLevelFilename(LWARN)
		return ll.warn
	case LERROR:
		if ll.err == nil {
			ll.err = log.New(output, "ERROR ", logFlags)
		}
		checkLogLevelFilename(LERROR)
		return ll.err
	case LFATAL:
		if ll.fatal == nil {
			ll.fatal = log.New(output, "FATAL ", logFlags)
		}
		checkLogLevelFilename(LFATAL)
		return ll.fatal
	default:
		return nil
	}
}

//Debug writes Debug logs
func Debug(pattern string, args ...interface{}) {
	if !(Level > LDEBUG) {
		ll.GetLogger(LDEBUG).Output(2, fmt.Sprintf(pattern, args...))
	}
}

//Info writes Info logs
func Info(pattern string, args ...interface{}) {
	if !(Level > LINFO) {
		ll.GetLogger(LINFO).Output(2, fmt.Sprintf(pattern, args...))
	}
}

//Warn writes log lines
func Warn(pattern string, args ...interface{}) {
	if !(Level > LWARN) {
		ll.GetLogger(LWARN).Output(2, fmt.Sprintf(pattern, args...))
	}
}

//Error writes error lines
func Error(pattern string, args ...interface{}) {
	if !(Level > LERROR) {
		ll.GetLogger(LERROR).Output(2, fmt.Sprintf(pattern, args...))
	}
}
func checkLogLevelFilename(level Loglevel) {
	switch level {
	case LDEBUG:
		checkLogFilename(ll.debugfilename)
		break
	case LINFO:
		checkLogFilename(ll.infofilename)
		break
	case LWARN:
		checkLogFilename(ll.warnfilename)
		break
	case LERROR:
		checkLogFilename(ll.errfilename)
		break
	case LFATAL:
		checkLogFilename(ll.fatalfilename)
		break
	}
}

func checkLogFilename(fileName string) {
	if fileName == "" {
		return
	}
	if _, err := os.Stat(fileName); err != nil {
		log.Printf("ERROR Logfile disappeared, %v re-creating file %v", err.Error(), fileName)
		f, err2 := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err2 == nil {
			ll.filemap[fileName] = f
			if ll.debugfilename == fileName {
				ll.debug.SetOutput(f)
			}
			if ll.infofilename == fileName {
				ll.info.SetOutput(f)
			}
			if ll.warnfilename == fileName {
				ll.warn.SetOutput(f)
			}
			if ll.errfilename == fileName {
				ll.err.SetOutput(f)
			}
			if ll.fatalfilename == fileName {
				ll.fatal.SetOutput(f)
			}
		} else {
			log.Fatalf("FATAL Could not aquire logfile, %v", err2.Error())
		}
	}
}

//SetLogFile sets a logfile on a speific loglevel
func SetLogFile(level Loglevel, filename string) (closer io.WriteCloser) {
	/*if ll.errfilename == filename || ll.fatalfilename == filename || ll.infofilename == filename || ll.debugfilename == filename || ll.warnfilename == filename {
		//log.Fatalf("filename %v is already in use", filename)
	}*/

	if v, ok := ll.filemap[filename]; ok {
		SetLogOutput(level, v)
		closer = v
	} else {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicf("Could not open file %v,%v", filename, err.Error())
		}
		closer = f
		SetLogOutput(level, f)
	}

	if level&LDEBUG != 0 {
		ll.debugfilename = filename
	}
	if level&LINFO != 0 {
		ll.infofilename = filename
	}
	if level&LWARN != 0 {
		ll.warnfilename = filename
	}
	if level&LERROR != 0 {
		ll.errfilename = filename
	}
	if level&LFATAL != 0 {
		ll.fatalfilename = filename
	}

	return
}

//SetLogOutput for a specific level
func SetLogOutput(level Loglevel, w io.Writer) {

	if level&LINFO != 0 {
		if ll.info == nil {
			ll.info = log.New(w, "INFO ", logFlags)
		} else {
			ll.info.SetOutput(w)
		}
		ll.infooutput = w
	}
	if level&LDEBUG != 0 {
		if ll.debug == nil {
			ll.debug = log.New(w, "DEBUG ", logFlags)
		} else {
			ll.debug.SetOutput(w)
		}
		ll.debugoutput = w
	}
	if level&LWARN != 0 {
		if ll.warn == nil {
			ll.warn = log.New(w, "WARN ", logFlags)
		} else {
			ll.warn.SetOutput(w)
		}
		ll.warnoutput = w
	}
	if level&LERROR != 0 {
		if ll.err == nil {
			ll.err = log.New(w, "ERROR ", logFlags)
		} else {
			ll.err.SetOutput(w)
		}
		ll.erroutput = w
	}
	if level&LFATAL != 0 {
		if ll.fatal == nil {
			ll.fatal = log.New(w, "FATAL ", logFlags)
		} else {
			ll.fatal.SetOutput(w)
		}
		ll.fataloutput = w
	}

}

//SetLevel Sets the level of the logger, no logs below the Level which is set, will be printet
func SetLevel(level string) {
	l := strings.ToLower(level)
	switch l {
	case "debug":
		Level = LDEBUG
	case "info":
		Level = LINFO
	case "warn", "warning":
		Level = LWARN
	case "error":
		Level = LERROR
	case "fatal":
		Level = LFATAL
	default:
		Level = LINFO
	}
}

//SetLoglevel sets the level to the enum chosen, more failsafe than SetLevel
func SetLoglevel(level Loglevel) {
	Level = level
}

func closeWriter(w io.Writer) {
	if w != nil {
		if _, ok := w.(io.WriteCloser); ok {
			(w.(io.WriteCloser)).Close()
		}
	}
}

//Close all open files with Close if need be, and kill logger instances
func Close() {
	ll.debug = nil
	ll.info = nil
	ll.warn = nil
	ll.err = nil
	ll.fatal = nil
	closeWriter(ll.logoutput)
	closeWriter(ll.debugoutput)
	closeWriter(ll.infooutput)
	closeWriter(ll.warnoutput)
	closeWriter(ll.erroutput)
	closeWriter(ll.warnoutput)

}
