package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ErwinSalas/pharma-scrapper/sitemap"
	"github.com/ErwinSalas/pharma-scrapper/types"
	"github.com/go-rod/rod"
)

var Products = []string{
	"https://sucreenlinea.com/0010100-amoebriz-x-2-tabs.html",
	"https://sucreenlinea.com/0012400-gyno-daktarin-dual-1-dia-otc.html",
}

var pharmacyIndex = map[string]int{
	"labomba":    0,
	"fishel":     1,
	"sucre":      2,
	"farmavalue": 3,
}

var linkIndex = map[string]int{
	"labomba":    0,
	"fishel":     0,
	"sucre":      0,
	"farmavalue": 0,
}

var pharmacySources = []types.Pharmacy{
	types.Pharmacy{
		ID:     "labomba",
		Domain: "farmacialabomba.com",
		ProductSelectors: types.Product{
			Name:   "h1",
			Price:  "h2.mb-0.fw-500",
			Active: "#flush-collapseThree > div",
			Lab:    "",
			Img:    "img.img-product-detail",
		},
	},
	types.Pharmacy{
		ID:     "fishel",
		Domain: "fischelenlinea.com",
		ProductSelectors: types.Product{
			Name:   "h1",
			Price:  "h2.mb-0.fw-500",
			Active: "#flush-collapseThree > div",
			Lab:    "",
			Img:    "img.img-product-detail",
		},
	},
	types.Pharmacy{
		ID:     "sucre",
		Domain: "sucreenlinea.com",
		ProductSelectors: types.Product{
			Name:   "span.base",
			Price:  "span.price",
			Active: "div.product-details-content.d-none.d-md-block > div:nth-child(1) > div:nth-child(3) > p",
			Lab:    "div.product-details-content.d-none.d-md-block > div:nth-child(2) > div:nth-child(2) > p",
			Img:    "img.main-product-photo",
		},
	},
	// types.Pharmacy{
	// 	ID:     "farmavalue",
	// 	Domain: "farmavalue.com",
	// 	ProductSelectors: types.Product{
	// 		Name:   "p.titulo-producto",
	// 		Price:  "p.cantidad-2",
	// 		Active: "p.principio-activo",
	// 		Lab:    "p.nombre-lab",
	// 		Img:    "img.product-image",
	// 	},
	// },
}

func getTextIfExists(page *rod.Page, selector string) string {
	if selector == "" || page.MustHas(selector) == false {
		return ""
	}
	return strings.TrimSpace(page.MustElement(selector).MustText())
}

func getAttributeIfExists(page *rod.Page, selector, attr string) string {
	if selector == "" || page.MustHas(selector) == false {
		return ""
	}
	attrPtr, _ := page.MustElement(selector).Attribute(attr)
	if attrPtr != nil {
		return *attrPtr
	}
	return ""
}

func getImageSrcIfExist(page *rod.Page, selector, attr string) string {

	timeout := time.Now().Add(10 * time.Second)
	placeholderImages := []string{
		"default-placeholder.jpg",
		"no-image.png",
	}

	var imgSrc string
	for time.Now().Before(timeout) {
		imgSrc = getAttributeIfExists(page, selector, "attr")

		// Si la imagen no es un placeholder, sal del bucle
		isPlaceholder := false
		for _, placeholder := range placeholderImages {
			if strings.Contains(imgSrc, placeholder) {
				isPlaceholder = true
				break
			}
		}
		if !isPlaceholder {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	return imgSrc
}

func extractProductInfo(page *rod.Page, config types.Pharmacy) (types.Product, error) {
	product := types.Product{}
	// Solo extraer datos si el selector existe
	product.Name = getTextIfExists(page, config.ProductSelectors.Name)
	product.Price = getTextIfExists(page, config.ProductSelectors.Price)
	product.Active = getTextIfExists(page, config.ProductSelectors.Active)
	product.Lab = getTextIfExists(page, config.ProductSelectors.Lab)
	product.Img = getAttributeIfExists(page, config.ProductSelectors.Img, "src")

	return product, nil
}

func validateURLHost(link string, host string) bool {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	return strings.HasSuffix(parsedURL.Host, host)
}
func findPharmacyConfig(url string) *types.Pharmacy {
	for _, pharmacy := range pharmacySources {
		if validateURLHost(url, pharmacy.Domain) {
			return &pharmacy
		}
	}

	return nil
}

func Crawl() {
	browser := rod.New().Trace(true).MustConnect()

	defer browser.MustClose()

	outputFile := "main.csv"

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}

	if fileInfo.Size() == 0 {
		writer.Write([]string{"Name", "Price", "Image", "Lab", "Active", "URL", "Pharmacy"})
	}

	productLinks, err := sitemap.LoadURLsFromFile(pharmacySources)
	if err != nil {
		fmt.Println("Error getting products links:", err)
		return
	}

	if len(productLinks) > 855 {
		productLinks = productLinks[855:]
	}

	for index, url := range productLinks {
		page := browser.MustPage(url)
		page.MustWaitDOMStable()
		config := findPharmacyConfig(url)
		if config == nil {
			continue
		}
		product, err := extractProductInfo(page, *config)
		if err != nil {
			fmt.Println("index", index)

			fmt.Println("Error extracting product info:", err)
			continue // Skip to the next product if extraction fails
		}

		fmt.Println("Name:", product.Name)
		fmt.Println("Price:", product.Price)
		fmt.Println("Image:", product.Img)
		fmt.Println("Lab:", product.Lab)
		fmt.Println("Active:", product.Active)
		fmt.Println("URL:", url)

		writer.Write([]string{product.Name, product.Price, product.Img, product.Lab, product.Active, url, config.ID})
	}

	fmt.Println("Scraping completed, data saved to labomba.csv")
}

func main() {
	command := os.Args[1]

	if command == "crawl" {
		Crawl()
	} else if command == "sitemaps" {
		config := os.Args[2]
		if config == "" {
			fmt.Println("Error: missing sitemap param")
		}

		pharmacyConfig := pharmacySources[pharmacyIndex[config]]
		err := sitemap.ExtractURLsFromSitemap(pharmacyConfig)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Error: not command:")

	}

}
