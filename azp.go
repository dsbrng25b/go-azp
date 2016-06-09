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
	"unicode"
)


type WorkUnit struct {
	start_time   time.Time
	end_time     time.Time
        break_time   time.Duration
	work_time    time.Duration
	comment      string
}

func ParseLine (line string) (WorkUnit, error) {
	year := time.Now().Year()
	field := strings.Fields(line)

	//work start time
	from_time, err := time.Parse("2006-1-2 15:04",fmt.Sprintf("%s.%.4d %s", field[0], year, field[1]))
	if err != nil {
		log.Fatal(err_fd)
	}

	//work finish time
	from_time, err := time.Parse("2006-1-2 15:04",fmt.Sprintf("%s.%.4d %s", field[0], year, field[2]))
	if err != nil {
		log.Fatal(err_td)
	}

	//break time
	var break_time time.Duration

	//if last byte is numeric treat as minutes
	var last_byte = rune(field[3][len(field[3])-1])
	if '0' <= last_byte && last_byte <= '9' {
		break_time_int, err := strconv.Atoi(field[3])
		if err != nil {
			return nil, err
		}
		break_time = time.Duration( break_time_int) * time.Minute
	} else {
		break_time, err = time.ParseDuration(field[3])
	}

	//work time
	work_time := to_time.Sub(from_time) - break_time

	//comment
	var comment string = ""
	if len(field) > 4 {
		comment = field[4]
	}	
	w := WorkUnit{ from_time, to_time, break_time, work_time, comment}
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
