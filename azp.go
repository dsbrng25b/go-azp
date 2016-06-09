package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"errors"
)

type WorkUnit struct {
	start_time time.Time
	end_time   time.Time
	break_time time.Duration
	work_time  time.Duration
	comment    string
}

func ParseLine(line string) (WorkUnit, error) {
	var w WorkUnit
	var year int = time.Now().Year()
	var field []string = strings.Fields(line)

	if len(field) < 4 {
		return w, errors.New("wrong line format")
	}
	//work start time
	start_time, err := time.Parse("2.1.2006 15:04", fmt.Sprintf("%s.%.4d %s", field[0], year, field[1]))
	if err != nil {
		return w, err
	}

	//work finish time
	end_time, err := time.Parse("2.1.2006 15:04", fmt.Sprintf("%s.%.4d %s", field[0], year, field[2]))
	if err != nil {
		return w, err
	}

	//move end_time to next day if before start time
	if end_time.Before(start_time) {
		end_time = end_time.Add( time.Duration( time.Hour ) * 24 )
	}

	//break time
	var break_time time.Duration

	//if last byte is numeric treat as minutes
	var last_byte = rune(field[3][len(field[3])-1])
	if '0' <= last_byte && last_byte <= '9' {
		break_time_int, err := strconv.Atoi(field[3])
		if err != nil {
			return w, err
		}
		break_time = time.Duration(break_time_int) * time.Minute
	} else {
		break_time, err = time.ParseDuration(field[3])
	}

	//work time
	work_time := end_time.Sub(start_time) - break_time

	//comment
	var comment string = ""
	if len(field) > 4 {
		comment = strings.Join(field[4:], " ")
	}
	w = WorkUnit{start_time, end_time, break_time, work_time, comment}
	return w, nil
}

func GetWorkUnits(r io.Reader) ([]WorkUnit, error) {
	var WorkUnits []WorkUnit
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		//skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		WorkUnit, err := ParseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "faild to parse '%s': %s\n", line, err)
			continue
		}
		WorkUnits = append(WorkUnits, WorkUnit)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return WorkUnits, nil
}

func PrintSummary(WorkUnits *[]WorkUnit) {
	for _, w := range *WorkUnits {
		fmt.Println(w.start_time)
		fmt.Println(w.end_time)
		fmt.Println(w.break_time)
		fmt.Println(w.work_time)
		fmt.Println(w.comment)
		fmt.Println("------")
	}
}

func main() {
	current_user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	file_path := filepath.Join(current_user.HomeDir, "azp", "worktime.txt")
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	WorkUnits, err := GetWorkUnits(file)
	if err != nil {
		log.Fatal(err)
	}
	PrintSummary(&WorkUnits)
}
