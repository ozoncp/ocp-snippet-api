package main

import (
	"fmt"
	"io"

	"errors"
	"os"
	"strconv"
	"time"
)

// task3 открывает файл openCount раз.
// Функция принимает на вход путь к файлу (path) и количество требуемых открытий (openCount).
// Функция ожидает, что в файле лежит целое число-счётчик. Если в файле лежит не число (не только число) на экран будет выведена ошибка функции Atoi.
// По ходу выполнения показывает сколько раз был открыт файл.
// openCount раз: открывает файл, считывает и выводит на экран текущее значение счётчика, увеличивает счётчик, закрывает файл.
func task3(path string, openCount uint) {
	getLastValue := func(file *os.File) (int, error) {
		if file == nil {
			return -1, errors.New("file is nil!")
		}

		data := make([]byte, 0, 64)

		var err error
		for cur, n := make([]byte, 10), 0; err == nil; {
			n, err = file.Read(cur)

			data = append(data, cur[:n]...)
		}

		if err != io.EOF {
			return -1, err
		}

		return strconv.Atoi(string(data))
	}
	updateLastValue := func(file *os.File, newVal int) error {
		if file == nil {
			return errors.New("file is nil!")
		}

		_, err := file.WriteAt([]byte(strconv.Itoa(newVal)), 0)
		return err
	}
	readFile := func(path string) (int, error) {
		file, err := os.OpenFile(path, os.O_RDWR, 0600)
		defer file.Close()

		cur, err := getLastValue(file)

		if err != nil {
			return -1, err
		}

		err = updateLastValue(file, cur+1)

		return cur, err
	}

	for i := uint(0); i < openCount; i++ {
		res, err := readFile(path)
		if err == nil {
			fmt.Printf("file updated %d times\n", res)
		} else {
			fmt.Printf("error while updating file: %s", err.Error())
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("DONE!")
}

func main() {
	fmt.Println("ocp-snippet-api by Oleg Usov")

	task3("test.txt", 10)
}
