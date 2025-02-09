package proxyfilter

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/exp/rand"
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

type ProxyStatus struct {
	URL     string `json:"url"`
	Working bool   `json:"working"`
}

func checkProxy(proxyURL string) bool {
	uri, err := url.Parse(proxyURL)
	if err != nil {
		return false
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(uri),
		// Opcional: Configura un tiempo de espera para evitar bloqueos
		Dial: (&net.Dialer{Timeout: 5 * time.Second}).Dial,
	}

	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	// Intenta realizar una solicitud a un sitio web de prueba
	resp, err := client.Get("http://www.google.com") // Puedes usar cualquier sitio confiable
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func GetRandomWorkingProxy(jsonFile string) (string, error) {
	// 1. Read JSON file
	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		return "", fmt.Errorf("error reading JSON file: %v", err)
	}

	// 2. Unmarshal JSON data
	var proxyStatuses []ProxyStatus
	err = json.Unmarshal(jsonData, &proxyStatuses)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// 3. Filter working proxies
	workingProxies := []string{}
	for _, status := range proxyStatuses {
		if status.Working {
			workingProxies = append(workingProxies, status.URL)
		}
	}

	// 4. Check if any working proxies exist
	if len(workingProxies) == 0 {
		return "", fmt.Errorf("no working proxies found in JSON file")
	}

	// 5. Seed the random number generator
	rand.Seed(uint64(time.Now().Unix())) // Important for actual randomness
	// 6. Generate a random index
	//randomIndex := rand.Intn(len(workingProxies))

	// 7. Return the randomly selected proxy
	//return workingProxies[randomIndex], nil
	return "https://13.56.192.187:80", nil
}

func main() {
	proxyStatuses := []ProxyStatus{}
	workingProxies := []string{}

	for _, proxy := range proxyList {
		working := checkProxy(proxy)
		status := ProxyStatus{URL: proxy, Working: working}
		proxyStatuses = append(proxyStatuses, status)

		if working {
			workingProxies = append(workingProxies, proxy)
			fmt.Println(proxy, " - Funciona")
		} else {
			fmt.Println(proxy, " - No funciona")
		}
	}

	fmt.Println("\nProxies funcionales:")
	for _, proxy := range workingProxies {
		fmt.Println(proxy)
	}

	// Convertir a JSON
	jsonData, err := json.MarshalIndent(proxyStatuses, "", "  ") // Indentación para mejor legibilidad
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	// Imprimir JSON a la consola
	fmt.Println("\nJSON:")
	fmt.Println(string(jsonData))

	// Guardar JSON a un archivo (opcional)
	err = os.WriteFile("proxy_status.json", jsonData, 0644) // Permisos de lectura/escritura para el usuario
	if err != nil {
		fmt.Println("Error al guardar JSON en archivo:", err)
		return
	}
	fmt.Println("Información de proxies guardada en proxy_status.json")

}
