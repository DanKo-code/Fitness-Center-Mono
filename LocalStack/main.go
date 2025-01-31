package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	// Путь к файлу
	inputFile := "./aws/init_bucket.sh"

	// Чтение содержимого файла
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	// Удаление символов \r
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "\r", "")

	// Перезапись содержимого файла с удаленными символами \r
	err = ioutil.WriteFile(inputFile, []byte(contentStr), 0644)
	if err != nil {
		log.Fatalf("Ошибка при записи файла: %v", err)
	}

	fmt.Println("Файл успешно обработан и сохранен.")
}
