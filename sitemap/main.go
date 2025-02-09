package sitemap

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/ErwinSalas/pharma-scrapper/types"
)

func ExtractURLsFromSitemap(pharmacyConfig types.Pharmacy) error {
	xmlFile, err := os.Open(fmt.Sprintf("./%s/sitemap.xml", pharmacyConfig.ID))
	if err != nil {
		return fmt.Errorf("error leyendo el archivo sitemap: %v", err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false

	var urlset types.URLSet
	err = decoder.Decode(&urlset)
	if err != nil {
		return fmt.Errorf("error decodificando el XML: %v", err)
	}

	outputFile := fmt.Sprintf("./%s/index.txt")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creando el archivo de salida: %v", err)
	}
	defer file.Close()

	for _, url := range urlset.URLs {
		_, err := file.WriteString(url.Loc + "\n")
		if err != nil {
			return fmt.Errorf("error escribiendo en el archivo: %v", err)
		}
	}

	fmt.Printf("URLs extraídas y guardadas en %s\n", outputFile)
	return nil
}

// Leer archivo y devolver array de strings
func ReadURLsFromFile(filename string) ([]string, error) {
	var urls []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error abriendo el archivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // Ignorar líneas vacías
			urls = append(urls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error leyendo el archivo: %v", err)
	}

	return urls, nil
}

func LoadURLsFromFile(pharmacys []types.Pharmacy) ([]string, error) {
	var urlsMatrix [][]string

	// Initialize the matrix with empty slices
	urlsMatrix = make([][]string, len(pharmacys))

	// Read URLs from each pharmacy file
	for index, pharmacy := range pharmacys {
		if pharmacy.ID == "farmavalue" {
			farmavalueURLS := []string{}
			for i := 1; i <= 5000; i++ {
				url := fmt.Sprintf("https://www.farmavalue.com/#/cr/products/%d", i)
				farmavalueURLS = append(farmavalueURLS, url)
			}
			urlsMatrix[index] = farmavalueURLS
			continue
		}
		filename := fmt.Sprintf("./%s/%s.txt", pharmacy.ID, pharmacy.ID)
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %v", filename, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" { // Ignore empty lines
				urlsMatrix[index] = append(urlsMatrix[index], line)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file %s: %v", filename, err)
		}
	}

	// Merge the matrix into a single array with interleaving
	return interleaveMatrix(urlsMatrix), nil
}

// interleaveMatrix takes a matrix of strings and interleaves elements into a single slice.
func interleaveMatrix(matrix [][]string) []string {
	var result []string
	indices := make([]int, len(matrix))
	activeLists := len(matrix)

	for activeLists > 0 {
		activeLists = 0 // Reset count for active lists

		for i := 0; i < len(matrix); i++ {
			if indices[i] < len(matrix[i]) {
				result = append(result, matrix[i][indices[i]])
				indices[i]++
				activeLists++
			}
		}
	}

	return result
}
