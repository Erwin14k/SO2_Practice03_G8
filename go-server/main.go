package main

import (
	"github.com/gorilla/mux" // Package for creating HTTP routers
	"github.com/rs/cors" // Package for Cross-Origin Resource Sharing (CORS) support
	"encoding/json" // Package for JSON encoding and decoding
	"io/ioutil" // Package for reading and writing files and data streams
	"net/http" // Package for HTTP client and server implementations
	"os/exec" // Package for executing external commands
	"strconv" // Package for string conversions
	"strings" // Package for manipulating strings
	"regexp"
	"bufio"
	"log" // Package for logging
	"fmt" // Package for formatted I/O
)

// Process represents a process with its properties
type Process struct {
	Pid     int    `json:"pid"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
	Estado  string `json:"estado"`
	Ram     int    `json:"ram"`
	Padre   int    `json:"padre"`
}

// CPUInfo represents CPU information and process tasks
type CPUInfo struct {
	TotalCPU int       `json:"totalcpu"`
	Running  int       `json:"running"`
	Sleeping int       `json:"sleeping"`
	Stopped  int       `json:"stopped"`
	Zombie   int       `json:"zombie"`
	Total    int       `json:"total"`
	Tasks    []Process `json:"tasks"`
}

// RAMInfo represents RAM information
type RAMInfo struct {
	TotalRAM    int `json:"totalram"`
	RAMLibre    int `json:"ramlibre"`
	RAMOcupada  int `json:"ramocupada"`
}

// General represents general system information
type general struct {
	TotalRAM    int `json:"totalram"`
	RAMLibre    int `json:"ramlibre"`
	RAMOcupada  int `json:"ramocupada"`
	TotalCPU    int `json:"totalcpu"`
}

// Counters represents process counters
type counters struct {
	Running  int       `json:"running"`
	Sleeping int       `json:"sleeping"`
	Stopped  int       `json:"stopped"`
	Zombie   int       `json:"zombie"`
	Total    int       `json:"total"`
}

// AllData represents all system data
type AllData struct {
	AllGenerales    []general    `json:"AllGenerales"`
	AllTipoProcesos []Process  `json:"AllTipoProcesos"`
	AllProcesos     []counters   `json:"AllProcesos"`
}

type MemoryBlock struct {
	InitialAddress string   `json:"initial_address"`
	FinalAddress   string   `json:"final_address"`
	Permissions    []string `json:"permissions"`
	Device         string   `json:"device"`
	File           string   `json:"file"`
	Size           float64  `json:"size"`
	Rss            float64  `json:"rss"`
}

type MemoryResult struct {
	TotalSize   float64       `json:"total_size"`
	TotalRss    float64       `json:"total_rss"`
	Blocks      []MemoryBlock `json:"blocks"`
}

/* createData creates the system data by reading the contents of specific files, 
 processing the data, and returning a JSON representation of the system information.*/
func createData() (string, error) {
	// read /proc/mem_grupo8 file
	outRAM, err := ioutil.ReadFile("/proc/mem_grupo8")
	if err != nil {
		fmt.Println(err)
	}
	// read /proc/cpu_grupo8 file
	outCPU, err := ioutil.ReadFile("/proc/cpu_grupo8")
	if err != nil {
		fmt.Println(err)
	}
	// --------- PROCESS ---------
	var cpuInfo CPUInfo
	err = json.Unmarshal(outCPU, &cpuInfo)
	if err != nil {
		fmt.Println("Error: Cpu json unmarshal failed", err)
		return "", err
	}

	for i, task := range cpuInfo.Tasks {
		uid, err := strconv.Atoi(task.Usuario)
		if err != nil {
			fmt.Println("Error: Failed to convert UID to int", err)
			return "", err
		}

		// Sh -> Interpreter
		// -c -> Read the command from the argument string
		// grep -m 1 -> Show the first match
		// cut -d: -f1 -> Cuts the string at the first delimiter and displays the first field
		cmdUsr := exec.Command("sh", "-c", "grep -m 1 '"+strconv.Itoa(uid)+":' /etc/passwd | cut -d: -f1")
		

		outUsr, err := cmdUsr.Output()
		if err != nil {
			fmt.Println("Error: Failed to get username for UID ", task.Usuario, err)
			return "", err
		}
		username := strings.TrimSpace(string(outUsr))
		cpuInfo.Tasks[i].Usuario = username
	}

	// --------- RAM ---------
	var ramInfo RAMInfo
	err = json.Unmarshal(outRAM, &ramInfo)
	if err != nil {
		fmt.Println("Error: Ram json unmarshal failed", err)
		return "", err
	}
	
	// allData struct variable contains the data to send to frontend
	allData := AllData{
		AllGenerales: []general{
			{
				TotalRAM:     ramInfo.TotalRAM,
				RAMLibre:     ramInfo.RAMLibre,
				RAMOcupada:   ramInfo.RAMOcupada,
				TotalCPU:     cpuInfo.TotalCPU,
			},
		},
		AllTipoProcesos: cpuInfo.Tasks,
		AllProcesos: []counters{
			{
				Running: cpuInfo.Running,
				Sleeping: cpuInfo.Sleeping,
				Stopped: cpuInfo.Stopped,
				Zombie: cpuInfo.Zombie,
				Total: cpuInfo.Total,
			},
		},
	}

	// Marshal allData variable to json
	allDataJSON, err := json.Marshal(allData)
	if err != nil {
		fmt.Println("Error: AllData json marshal failed", err)
		return "", err
	}

	return string(allDataJSON), nil
}

// handleGet handles the GET request, receives a response writer and a request object as parameters, and returns no value. 
// It generates system data using the createData function and sends it as a JSON response.
// w - http.ResponseWriter: The response writer used to write the HTTP response.
// r - *http.Request: The HTTP request received.
func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET")

	allData, err := createData() // Get all process data in JSON format
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if data is empty
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, allData)
}

// handlePost handles the POST request by reading the request body, extracting the PID (process ID), and killing the corresponding process.
// If successful, it returns an HTTP 200 OK status code and a response message indicating that the process has been deleted.
// w - http.ResponseWriter: The response writer used to write the HTTP response.
// r - *http.Request: The HTTP request received.
func handlePost(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error reading the body
		return
	}

	pid, err := strconv.Atoi(string(body)) // Convert the body to an integer (PID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Return HTTP 400 Bad Request if the body is not a valid PID
		fmt.Fprintln(w, "Information: Invalid PID")
		return
	}

	cmd := exec.Command("sudo", "kill", strconv.Itoa(pid)) // Create a command to kill the process
	err = cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error killing the process
		fmt.Fprintln(w, "Error killing process")
		return
	}

	fmt.Println("Information: Process with PID", pid, "has been deleted") // Print the information about the deleted process
	w.WriteHeader(http.StatusOK)              // Set HTTP 200 OK status code
	fmt.Fprintln(w, "Process deleted")        // Write the response message to the response writer
}

// handleRoute handles the route "/", prints single a message only.
// w - http.ResponseWriter: The response writer used to write the HTTP response.
// r - *http.Request: The HTTP request received.
func handleRoute(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to my API :D")
}

func smapsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error reading the body
		return
	}

	pid, err := strconv.Atoi(string(body)) // Convert the body to an integer (PID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Return HTTP 400 Bad Request if the body is not a valid PID
		fmt.Fprintln(w, "Invalid PID")
		return
	}

	cmd := exec.Command("sudo", "cat", fmt.Sprintf("/proc/%d/smaps", pid)) // Create a command to execute "cat /proc/pid/maps"
	output, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error executing the command
		fmt.Fprintln(w, "Error reading process memory")
		return
	}

	blocks, totalSize, totalRss := parseSmapsOutput(string(output))
	result := MemoryResult{
		Blocks:    blocks,
		TotalSize: totalSize,
		TotalRss:  totalRss,
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Error al convertir los datos a JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func parseSmapsOutput(output string) ([]MemoryBlock, float64, float64) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanLines)

	var blocks []MemoryBlock
	var currentBlock MemoryBlock
	var totalSize, totalRss float64

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-") {
			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}

			address := fields[0]
			permissions := fields[1]
			device := fields[3]
			file := fields[len(fields)-1]

			addresses := strings.Split(address, "-")
			initialAddress := addresses[0]
			finalAddress := addresses[1]

			currentBlock = MemoryBlock{
				InitialAddress: initialAddress,
				FinalAddress:   finalAddress,
				Permissions:    mapPermissions(permissions),
				Device:         device,
				File:           file,
			}
		} else if strings.HasPrefix(line, "Size:") {
			size := extractValue(line)
			currentBlock.Size = float64(size) / 1024
			totalSize += currentBlock.Size
		} else if strings.HasPrefix(line, "Rss:") {
			rss := extractValue(line)
			currentBlock.Rss = float64(rss) / 1024
			totalRss += currentBlock.Rss
			blocks = append(blocks, currentBlock)
		}
	}

	return blocks, totalSize, totalRss
}

func extractValue(line string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(line)
	if match != "" {
		value := match
		size, _ := strconv.Atoi(value)
		return size
	}
	return 0
}

func mapPermissions(permissions string) []string {
	mappedPermissions := make([]string, 0)

	if strings.Contains(permissions, "r") {
		mappedPermissions = append(mappedPermissions, "Lectura")
	}
	if strings.Contains(permissions, "w") {
		mappedPermissions = append(mappedPermissions, "Escritura")
	}
	if strings.Contains(permissions, "x") {
		mappedPermissions = append(mappedPermissions, "Ejecucion")
	}
	return mappedPermissions
}

// main is the entry point of the application.
// It sets up the router, defines the route handlers, and starts the HTTP server.
// The server listens on port 8080 for incoming requests.
func main() {
	fmt.Println("************************************************************")
	fmt.Println("*                 SO2 Practica 3 - Grupo 8                 *")
	fmt.Println("************************************************************")

	router := mux.NewRouter().StrictSlash(true) // Create a new router instance
	router.HandleFunc("/", handleRoute) // Set the handler function for the root route ("/")
	router.HandleFunc("/tasks", handlePost).Methods("POST") // Set the handler function for the "/tasks" route with POST method
	router.HandleFunc("/tasks", handleGet).Methods("GET") // Set the handler function for the "/tasks" route with GET method
	router.HandleFunc("/memory", smapsHandler).Methods("POST")

	handler := cors.Default().Handler(router) // Create a new CORS handler with default settings
	log.Fatal(http.ListenAndServe(":8080", handler)) // Start the HTTP server and listen on port 8080

	fmt.Println("Server on port 8080")
}
