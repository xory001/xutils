package xlogs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xory001/xutils/xsys"
	"github.com/xory001/xutils/xzip"
)

type CXLogFile struct {
	logFileName string
	logFile     *os.File
	logFileList []string

	maxFileSize     int64 //file max size, bytes
	currentFileSize int64 //

	logFileIndex int //from 0 - 9999
	logFileDir   string
	//newFile      bool

	//for new log
	day    int
	hour   int
	min    int
	second int
}

func NewXLogFile(logDir string, maxFileSize int64, zipLog bool) *CXLogFile {
	if 0 == len(logDir) {
		strExeFullName, err := os.Executable()
		if nil != err {
			fmt.Println("NewXLogFile, Executable err =", err)
			return nil
		}

		logDir = filepath.Join(filepath.Dir(strExeFullName), "logs")
		_, err = os.Stat(logDir)
		if os.IsNotExist(err) {
			os.MkdirAll(logDir, 0644)
		}
	}
	if 0 == maxFileSize {
		maxFileSize = 100 //50 * 1024 * 1024
	}
	t := time.Now().UTC()
	xlog := &CXLogFile{logFileDir: logDir, maxFileSize: maxFileSize, day: t.Day(), hour: t.Hour(), min: t.Minute(), second: t.Second()}
	if zipLog {
		go xlog.processHistoryLogFile()
	}

	return xlog
}

func (x *CXLogFile) Write(p []byte) (n int, err error) {
	x.createNewLogFileIfNeed()
	if nil == x.logFile {
		return 0, errors.New("can't create log file")
	}
	bytesWrite, err := x.logFile.Write(p)
	if nil != err {
		fmt.Println("write log err =", err)
	}
	x.currentFileSize += int64(bytesWrite)
	return bytesWrite, err
}

func (x *CXLogFile) createNewLogFileIfNeed() error {
	var err error

	//check file size
	if x.currentFileSize > x.maxFileSize {
		x.logFile.Close()
		x.logFile = nil
	}

	//check time
	t := time.Now().UTC()
	if x.day != t.Day() { //another day, need create new
		x.day = t.Day()
		x.hour, x.min, x.second = t.Clock()
		if nil != x.logFile {
			x.logFile.Close()
			x.logFile = nil
		}
		//err = x.createNewLogFileIfNeed()
		//if x.newFile { //need create new log file
		//	x.newFile = false
		//	err = x.createNewLogFile()
		//}
	}

	if nil == x.logFile {
		err = x.createNewLogFile()
	}
	if nil != x.logFile {
		fmt.Println("log file name =", x.logFile.Name())
	}

	return err
}

func (x *CXLogFile) createNewLogFile() error {
	var err error
	strFileName, err := x.generateFileName()
	if nil != err {
		return err
	}
	fmt.Println("new log file name =", strFileName)
	x.logFile, err = os.OpenFile(strFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if nil != err {
		x.logFile = nil
		fmt.Println("open file err, err =", err, ", file =", strFileName)
	} else {
		x.logFileList = append(x.logFileList, strFileName)
	}
	x.currentFileSize = 0
	return err
}

func (x *CXLogFile) generateFileName() (string, error) {
	strExeFullName, err := os.Executable()
	if nil != err {
		log.Println("Executable err =", err)
		return "", err
	}
	strExeName := filepath.Base(strExeFullName)
	strNamePrefix := filepath.Join(x.logFileDir, strExeName)
	strLogFileName := fmt.Sprintf("%s_info_%s_%04d_%d.log", strNamePrefix, time.Now().UTC().Format("20060102"), x.logFileIndex, os.Getpid())
	x.logFileIndex++
	for {
		_, err := os.Stat(strLogFileName)
		if nil != err {
			if os.IsNotExist(err) {
				break
			}
		}
		strLogFileName = fmt.Sprintf("%s_info_%s_%04d_%d.log", strNamePrefix, time.Now().UTC().Format("20060102"), x.logFileIndex, os.Getpid())
		x.logFileIndex++
	}

	return strLogFileName, nil
}

func (x *CXLogFile) processHistoryLogFile() {
	moduleName, _ := os.Executable()
	execName := filepath.Base(moduleName)
	tickerCheck := time.NewTicker(time.Second)
	for timeNow := range tickerCheck.C {
		sliceFiles, err := xsys.ReadDirAscByTime(x.logFileDir)
		if nil == err {
			var sliceSrcFile []string
			var fileDay = 0
			var fileHour = 0
			var lastTime time.Time
			for _, file := range sliceFiles {
				fmt.Println("file =", file.Name())
				if file.IsDir() {
					//delete expire data file
				} else {
					if strings.Contains(file.Name(), execName) && (".log" == filepath.Ext(file.Name())) {
						fileInfo, _ := file.Info()
						fileTime := fileInfo.ModTime().UTC()
						if 0 == fileDay {
							fileDay = fileTime.Day()
						}
						if fileTime.Day() != timeNow.UTC().Day() { //need zip
							if fileDay == fileTime.Day() {
								if 0 == fileHour {
									fileHour = fileTime.Hour()
								}

								if fileHour == fileTime.Hour() {
									sliceSrcFile = append(sliceSrcFile, filepath.Join(x.logFileDir, file.Name()))
									lastTime = fileTime
									continue
								} else {
									x.zipFiles(sliceSrcFile, execName, lastTime)
									sliceSrcFile = nil

									fileHour = fileTime.Hour()
									sliceSrcFile = append(sliceSrcFile, filepath.Join(x.logFileDir, file.Name()))
									lastTime = fileTime
									continue
								}
							} else {
								x.zipFiles(sliceSrcFile, execName, lastTime)
								sliceSrcFile = nil

								fileDay = fileTime.Day()
								fileHour = fileTime.Hour()
								sliceSrcFile = append(sliceSrcFile, filepath.Join(x.logFileDir, file.Name()))
								lastTime = fileTime
								continue
							}
						} else { //if fileTime.Day() != timeNow.UTC().Day() { //need zip
							if len(sliceSrcFile) > 0 {
								x.zipFiles(sliceSrcFile, execName, lastTime)
								sliceSrcFile = nil
							}
							break
						}
					}

				}
			}
		}
	}
}

func (x *CXLogFile) zipFiles(sliceSrcFile []string, execName string, lastTime time.Time) {
	destDir := filepath.Join(x.logFileDir, lastTime.Format("20060102"))
	os.MkdirAll(destDir, 0777)
	destFile := filepath.Join(destDir, fmt.Sprintf("%s_info_%02d.tar.gz", execName, lastTime.Hour()))
	xzip.ZipFilesToTarGz(sliceSrcFile, destFile)
	sliceSrcFile = nil
}
