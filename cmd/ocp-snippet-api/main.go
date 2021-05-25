package main

import (
	"fmt"

	"errors"
	"os"
	"strconv"
	"time"
)

func task3() {
	// Task #3:
	// Открытие файла n раз.
	// В файле лежит цифра (с чсилом слишком сложно без ioutil :-), а суть задания не в этом ).
	// При каждом открытии программа считывает текущую цифру, и увеличивает её (происходит перезапись)
	getLastValue := func(file *os.File) (int, error) {
		if file == nil {
			return -1, errors.New("file is nil!")
		}

		b := make([]byte, 1)
		bSz, err := file.Read(b)

		if err != nil {
			return -1, err
		}

		return strconv.Atoi(string(b[:bSz]))
	}
	updateLastValue := func(file *os.File, newVal int) error {
		if file == nil {
			return errors.New("file is nil!")
		}

		_, err := file.WriteAt([]byte(strconv.Itoa(newVal)), 0)
		return err
	}
	readFile := func(path string) (int, error) {
		file, err := os.OpenFile(path, os.O_RDWR, 0666)
		defer file.Close()

		cur, err := getLastValue(file)

		if err != nil {
			return -1, nil
		}

		err = updateLastValue(file, cur+1)

		return cur, err
	}

	for range []int{1, 2, 3, 4, 5, 6, 7} {
		res, err := readFile("test.txt")
		if err == nil {
			fmt.Printf("file updated %d times\n", res)
		} else {
			fmt.Printf("error while updating file: %s", err.Error()) // TO BE FIXED: заменить на fmt.Errorf!
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("DONE!")
}

func main() {
	fmt.Println("ocp-snippet-api by Oleg Usov")

	task3()
}
