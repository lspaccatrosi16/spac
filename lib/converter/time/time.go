package time

import (
	"fmt"
	"strconv"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type ConvTime struct {
	seconds int
}

type convType int

const numTypes = 2

type hmsType struct {
}

func (*hmsType) Input() (int, error) {
	hS := input.GetInput("Hours")
	mS := input.GetInput("Minutes")
	sS := input.GetInput("Seconds")

	h, err := strconv.Atoi(hS)
	if err != nil {
		return 0, err
	}
	m, err := strconv.Atoi(mS)
	if err != nil {
		return 0, err
	}
	s, err := strconv.Atoi(sS)
	if err != nil {
		return 0, err
	}

	return 3600*h + 60*m + s, err
}

func (*hmsType) Output(secs int) string {
	hInt := secs / 3600

	mInt := secs/60 - 60*hInt

	sInt := secs - 3600*hInt - 60*mInt

	hStr := fmt.Sprint(hInt)
	mStr := fmt.Sprint(mInt)
	sStr := fmt.Sprint(sInt)

	if hInt <= 9 {
		hStr = "0" + hStr
	}

	if mInt <= 9 {
		mStr = "0" + mStr
	}
	if sInt <= 9 {
		sStr = "0" + sStr
	}

	return fmt.Sprintf("%s:%s:%s", hStr, mStr, sStr)
}

type decimalType struct{}

func (d *decimalType) Input() (int, error) {
	fStr := input.GetInput("Decimal String")
	fInt, err := strconv.ParseFloat(fStr, 64)
	if err != nil {
		return 0, err
	}

	secs := fInt * 3600
	return int(secs), nil
}

func (d *decimalType) Output(secs int) string {
	return fmt.Sprintf("%.3f", float64(secs)/3600)
}

const (
	hms convType = iota
	decimal
)

func (c convType) String() string {
	switch c {
	case hms:
		return "HH:MM:SS"
	case decimal:
		return "Decimal"
	default:
		return "Invalid"
	}
}

func (c convType) FromString(s string) convType {
	switch s {
	case "HH:MM:SS":
		return hms
	case "Decimal":
		return decimal
	}
	return -1
}

func (c *ConvTime) Name() string {
	return "Time"
}

func (c *ConvTime) InputFormat() []string {
	arr := []string{}
	for i := 0; i < numTypes; i++ {
		arr = append(arr, convType(i).String())
	}
	return arr
}

func (c *ConvTime) OutputFormat() []string {
	return c.InputFormat()
}

func (c *ConvTime) Convert(inputFormat, outputFormat string) string {
	input := new(convType).FromString(inputFormat)
	output := new(convType).FromString(outputFormat)

	switch input {
	case hms:
		s, err := new(hmsType).Input()
		if err != nil {
			return err.Error()
		}
		c.seconds = s
	case decimal:
		s, err := new(decimalType).Input()
		if err != nil {
			return err.Error()
		}
		c.seconds = s
	}

	o := ""

	switch output {
	case hms:
		o = new(hmsType).Output(c.seconds)
	case decimal:
		o = new(decimalType).Output(c.seconds)
	}

	return o
}
