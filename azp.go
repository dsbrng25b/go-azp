package main

import (
	"fmt"
	"time"
	"bufio"
        "os"
	"os/user"
	"log"
	"path/filepath"
	"strings"
	"errors"
	"io"
	"strconv"
)


type WorkUnit struct {
	from      time.Time	 //work start time
	to        time.Time	 //work end time
        brk_min   time.Duration  //break
	work_time time.Duration  //to - from - brk_min
	comment   string
}

func ParseTime (t string) (int, int, error){
	var err error
	var hour int
	var minute int
	parts := strings.Split(t, ":")	
	if len(parts) != 2 {
		return hour, minute, errors.New("error parsing time")
	}
	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return hour, minute, err
	}
	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return hour, minute, err
	}
	return hour, minute, nil
}

func ParseDate (d string) (int, int, error){
	var err error
	var day int
	var month int
	parts := strings.Split(d, ".")	
	if len(parts) != 2 {
		return day, month, errors.New("error parsing date")
	}
	day, err = strconv.Atoi(parts[0])
	if err != nil {
		return day, month, err
	}
	month, err = strconv.Atoi(parts[1])
	if err != nil {
		return day, month, err
	}
	return day, month, nil
}
	

func ParseLine (line string) (WorkUnit, error) {
	var w WorkUnit
	year := time.Now().Year()
	field := strings.Fields(line)
	day, month, err := ParseDate(field[0])
	if err != nil {
		return w, err
	}
	f_hour, f_min, _ := ParseTime(field[1])
	//from_time := time.Date(year, time.Month(month), day, f_hour, f_min, 0, 0, time.UTC)
	from_time, err_fd := time.Parse("2006-1-2 15:04", fmt.Sprintf("%.4d-%.2d-%.2d %.2d:%.2d", year, month, day, f_hour, f_min))
	if err_fd != nil {
		log.Fatal(err_fd)
	}
	t_hour, t_min, _ := ParseTime(field[2])
	to_time, err_td := time.Parse("2006-01-02 15:04", fmt.Sprintf("%.4d-%.2d-%.2d %.2d:%.2d", year, month, day, t_hour, t_min))
	if err_td != nil {
		log.Fatal(err_td)
	}
	brk_time_int, _ := strconv.Atoi(field[3])
	brk_time := time.Duration( brk_time_int ) * time.Minute 

	work_time := to_time.Sub(from_time) - brk_time

	w = WorkUnit{ from_time, to_time, brk_time, work_time, "foobar" }
	return w, nil
}

func GetWorkUnits(r io.Reader) ([]WorkUnit, error){
	var WorkUnits []WorkUnit
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		//skip empty lines and comments
		if line == "" || strings.HasPrefix(line,"#") {
			continue
		}
		WorkUnit, err := ParseLine(line)
		if err != nil {
			log.Fatal(err)
		}
		WorkUnits = append(WorkUnits, WorkUnit)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return WorkUnits, nil
}

func main() {

	current_user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	file_path := filepath.Join( current_user.HomeDir, "azp", "worktime.txt")
	fmt.Printf("%v %T\n", file_path, file_path)
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	WorkUnits, err := GetWorkUnits(file)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range WorkUnits {
		fmt.Println(w.from)
		fmt.Println(w.to)
		fmt.Println(w.brk_min)
		fmt.Println(w.work_time)
		fmt.Println(w.comment)
		fmt.Println("------")
	}
}
