package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/StudioSol/set"
)

const baseURL = "https://registro.br/dominio"

func refreshDomainsList() (err error) {
	now := time.Now()
	outputFile := fmt.Sprintf("history/release-%d-%d.txt", now.Month(), now.Year())

	err = downloadFile(fmt.Sprintf("%s/lista-processo-liberacao.txt", baseURL), outputFile)
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(source, output string) (err error) {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error on download file. Status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func readDomainFile(fileName string, skipComments bool) (lines []string, err error) {
	sourceFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()

	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if skipComments && strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, scanner.Text())
	}

	return
}

func diffFiles(source, comparision string) (newDomains []string, removedDomains []string, err error) {
	domains, err := readDomainFile(comparision, true)
	if err != nil {
		return nil, nil, err
	}
	oldList := set.NewLinkedHashSetString(domains...)

	domains, err = readDomainFile(source, true)
	if err != nil {
		return nil, nil, err
	}
	currentList := set.NewLinkedHashSetString(domains...)

	removedList := set.NewLinkedHashSetString()
	for domain := range oldList.Iter() {
		if !currentList.InArray(domain) {
			removedList.Add(domain)
		}
	}

	newList := set.NewLinkedHashSetString(domains...)
	for domain := range oldList.Iter() {
		newList.Remove(domain)
	}

	return newList.AsSlice(), removedList.AsSlice(), nil
}
