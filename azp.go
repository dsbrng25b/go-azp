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

func (w *WorkUnit) String() string {
	return fmt.Sprintf("from %s to %s, break %d, work_time %s, comment: %s", w.from, w.to, w.brk_min, w.work_time, w.comment)
}

func ParseTime (timestring string) (hour, minute int) {
	fmt.Println("timestring:", timestring)
	str_hour, str_minute := strings.Split(timestring, ":")
	hour, err := strconv.Atoi(str_hour)
	minute, err := strconv.Atoi(str_minute)
	return hour, minute
}
	

func ParseLine (line string) (WorkUnit, error) {
	year := time.Now().Year()
	field := strings.Fields(line)
	//date
	day, month = strings.Split(field[0],".")
	day, err := strconv.Atoi(day)
	month, err := strconv.Atoi(month)
	//from time
	f_hour, f_min := ParseTime(field[1])
	from_time, err = time.Date(year, time.Month(month), day, f_hour, f_min, 0, 0, time.UTC)
	//to time
	t_hour, t_min := ParseTime(field[2])
	to_time, err = time.Date(year, time.Month(month), day, t_hour, t_min, 0, 0, time.UTC)

	brk_time_int, err := strconv.Atoi(field[3])
	brk_time := time.Duration( brk_time_int * time.Minute )

	work_time = to_time.Sub(from_time) - brk_time

	w := WorkUnit{ from_time, to_time, brk_time, work_time, "foobar" }
	if 1 == 0 {
		return w, errors.New("error parsing line")
	}
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
		fmt.Printf("'%v' %T\n%v %T\n", line, line, fields, fields)
		//fmt.Printf(" -->> %v\n", fields)
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

	fmt.Println("start get workunits")
	WorkUnits, err := GetWorkUnits(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(WorkUnits)

}
