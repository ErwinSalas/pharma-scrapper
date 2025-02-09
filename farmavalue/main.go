package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

var proxyList = []string{
	"http://103.152.112.120:80",
	"http://23.247.136.245:80",
	"http://23.247.136.248:80",
	"http://23.247.136.254:80",
	"http://103.152.112.157:80",
	"http://184.73.68.87:11",
	"http://47.88.59.79:82",
	"http://143.198.226.25:80",
	"http://98.80.66.1:10018",
	"http://23.94.136.205:80",
	"http://23.247.137.142:80",
	"http://143.42.191.48:80",
	"http://138.91.159.185:80",
	"http://3.136.29.104:80",
	"http://3.212.148.199:3128",
	"http://23.94.137.130:80",
	"http://204.236.137.68:80",
	"http://44.219.175.186:80",
	"http://47.251.43.115:33333",
	"http://50.174.7.156:80",
	"http://147.124.222.230:3128",
	"http://47.251.122.81:8888",
	"http://216.229.112.25:8080",
	"http://66.29.154.105:3128",
	"http://162.223.90.130:80",
	"http://50.223.246.237:80",
	"http://44.218.183.55:80",
	"http://50.207.199.80:80",
	"http://50.207.199.83:80",
}

func getRandomProxy() string {
	rand.Seed(time.Now().UnixNano())
	return proxyList[rand.Intn(len(proxyList))]
}

func scrapeData() {
	// Create a CSV file to store the scraped data
	file, err := os.Create("farmavalue.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV headers
	writer.Write([]string{"Name", "Price", "Image", "Lab", "URL"})

	// Initialize the collector with concurrency limits
	c := colly.NewCollector(
		colly.Async(true),
	)

	// Set a limit on the number of concurrent requests
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2, // Adjust this number based on system capacity
		Delay:       2 * time.Second,
	})

	// Set up proxy rotation
	c.OnRequest(func(r *colly.Request) {
		proxy := getRandomProxy()
		r.Ctx.Put("proxy", proxy)
		r.Headers.Set("Proxy", proxy)
		log.Printf("Using proxy: %s", proxy)
		log.Printf("site: %s", r.URL)

	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error occurred for %s: %v", r.Body, err)
		// Retry with new proxy
		proxy := getRandomProxy()
		log.Printf("Retrying with new proxy: %s", proxy)
		r.Ctx.Put("proxy", proxy)
		r.Headers.Set("Proxy", proxy)
	})

	c.OnHTML("#app > div.mainContainer > div > section > div.product-detail.grid-2.box-responsive > div.grid-1.center > div", func(e *colly.HTMLElement) {
		name := e.ChildText("p.titulo-producto")
		price := e.ChildText("p.cantidad-2")
		active := e.ChildText("p.principio-activo")
		image := e.ChildAttr("img.product-image", "src")
		lab := e.ChildText("p.nombre-lab")
		url := e.Request.URL.String()

		// Print the values to the console
		fmt.Println("Name:", name)
		fmt.Println("Price:", price)
		fmt.Println("Image:", image)
		fmt.Println("Lab:", lab)
		fmt.Println("URL:", url)
		fmt.Println("------")
		// Write the scraped data to CSV
		writer.Write([]string{name, price, image, lab, active, url})
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Loop through pages (adjust the range based on your scraping needs)
	for i := 1; i <= 5000; i++ {
		url := fmt.Sprintf("https://www.farmavalue.com/#/cr/products/%d", i)
		err := c.Visit(url)
		if err != nil {
			log.Printf("Failed to visit URL %s: %v", url, err)
		}
	}

	// Wait for all requests to finish
	c.Wait()
}

func main() {
	scrapeData()
	fmt.Println("Scraping completed, data saved to farmavalue.csv")
}
