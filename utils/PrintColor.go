package utils

import (
	"os"
	"fmt"
	"log"
	logmki "github.com/mkideal/log"
)

const (
	PFCLR_NONE      = "\033[0m"
	PFCLR_BLACK     = "\033[0;30m"
	PFCLR_L_BLACK   = "\033[1;30m"
	PFCLR_RED       = "\033[0;31m"
	PFCLR_L_RED     = "\033[1;31m"
	PFCLR_GREEN     = "\033[0;32m"
	PFCLR_L_GREEN   = "\033[1;32m"
	PFCLR_BROWN     = "\033[0;33m"
	PFCLR_YELLOW    = "\033[1;33m"
	PFCLR_BLUE      = "\033[0;34m"
	PFCLR_L_BLUE    = "\033[1;34m"
	PFCLR_PURPLE    = "\033[0;35m"
	PFCLR_L_PURPLE  = "\033[1;35m"
	PFCLR_CYAN      = "\033[0;36m"
	PFCLR_L_CYAN    = "\033[1;36m"
	PFCLR_GRAY      = "\033[0;37m"
	PFCLR_WHITE     = "\033[1;37m"
	PFCLR_BOLD      = "\033[1m"
	PFCLR_UNDERLINE = "\033[4m"
	PFCLR_BLINK     = "\033[5m"
	PFCLR_REVERSE   = "\033[7m"
	PFCLR_HIDE      = "\033[8m"
	PFCLR_CLEAR     = "\033[2J"
	PFCLR_CLRLINE   = "\r\033[K" //or "\033[1K\r"
)

func SprintRed(v ...interface{}) (out string) {
	out += PFCLR_L_RED
	out += fmt.Sprint(v...)
	out += PFCLR_NONE
	return
}


func LogFatalf(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf(PFCLR_L_RED+"--ERROR-- "+format+PFCLR_NONE+"\n", v...))
	logmki.Printf(2, logmki.LvFATAL, format,v...)
	os.Exit(0)
}
func LogErrorf(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf(PFCLR_L_RED+"--ERROR-- "+format+PFCLR_NONE+"\n", v...))
	logmki.Printf(2, logmki.LvERROR, format,v...)
}
func LogWarnf(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf(PFCLR_RED+"--WARN-- "+format+PFCLR_NONE+"\n", v...))
	logmki.Printf(2, logmki.LvWARN, format,v...)
}
func LogInfof(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf(PFCLR_L_GREEN+"--INFO-- "+format+PFCLR_NONE+"\n", v...))
	logmki.Printf(2, logmki.LvINFO, format,v...)
}
func LogDepthInfof(depth int,format string, v ...interface{}) {
	log.Output(2+depth, fmt.Sprintf(PFCLR_L_GREEN+"--INFO-- "+format+PFCLR_NONE+"\n", v...))
}
func LogDebugf(format string, v ...interface{}) {
	log.Output(2, fmt.Sprintf("--DEBUG-- "+format+"\n", v...))
	logmki.Printf(2, logmki.LvDEBUG, format,v...)
}
func PrintErrorf(format string,v ...interface{}) {
	fmt.Fprintf(os.Stderr,PFCLR_L_RED+"--ERROR-- "+format+PFCLR_NONE+"\n", v...)
}
func PrintErrorFatalf(format string,v ...interface{}) {
	fmt.Fprintf(os.Stderr,PFCLR_L_RED+"--ERROR-- "+format+PFCLR_NONE+"\n", v...)
	os.Exit(0)
}

func PrintGreenf(format string,v ...interface{}) {
	fmt.Fprintf(os.Stdout,PFCLR_L_GREEN+format+PFCLR_NONE+"\n", v...)
}
func PrintRedf(format string,v ...interface{}) {
	fmt.Print(PFCLR_L_RED)
	fmt.Print(v...)
	fmt.Println(PFCLR_NONE)
}
func SprintGreen(v ...interface{}) (out string) {
	out += PFCLR_L_GREEN
	out += fmt.Sprint(v...)
	out += PFCLR_NONE
	return
}
func PrintGreen(v ...interface{}) {
	fmt.Print(PFCLR_L_GREEN)
	fmt.Print(v...)
	fmt.Println(PFCLR_NONE)
}
func SprintBlue(v ...interface{}) (out string) {
	out += PFCLR_L_BLUE
	out += fmt.Sprint(v...)
	out += PFCLR_NONE
	return
}
func PrintBlue(v ...interface{}) {
	fmt.Print(PFCLR_L_BLUE)
	fmt.Print(v...)
	fmt.Println(PFCLR_NONE)
}
