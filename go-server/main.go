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

// Process represents a process with its attributes.
type Process struct {
	Pid     int    `json:"pid"`     // Pid represents the process ID.
	Nombre  string `json:"nombre"`  // Nombre represents the name of the process.
	Usuario string `json:"usuario"` // Usuario represents the user associated with the process.
	Estado  string `json:"estado"`  // Estado represents the state of the process.
	Ram     int    `json:"ram"`     // Ram represents the amount of RAM (in bytes) used by the process.
	Padre   int    `json:"padre"`   // Padre represents the parent process ID.
}

// CPUInfo represents CPU information and process tasks.
type CPUInfo struct {
	TotalCPU int       `json:"totalcpu"` // TotalCPU represents the total number of CPUs.
	Running  int       `json:"running"`  // Running represents the number of running processes.
	Sleeping int       `json:"sleeping"` // Sleeping represents the number of sleeping processes.
	Stopped  int       `json:"stopped"`  // Stopped represents the number of stopped processes.
	Zombie   int       `json:"zombie"`   // Zombie represents the number of zombie processes.
	Total    int       `json:"total"`    // Total represents the total number of processes.
	Tasks    []Process `json:"tasks"`    // Tasks represents a list of process tasks.
}

// RAMInfo represents RAM information.
type RAMInfo struct {
	TotalRAM    int `json:"totalram"`    // TotalRAM represents the total amount of RAM.
	RAMLibre    int `json:"ramlibre"`    // RAMLibre represents the amount of free RAM.
	RAMOcupada  int `json:"ramocupada"`  // RAMOcupada represents the amount of occupied RAM.
}

// General represents general system information.
type General struct {
	TotalRAM    int `json:"totalram"`    // TotalRAM represents the total amount of RAM.
	RAMLibre    int `json:"ramlibre"`    // RAMLibre represents the amount of free RAM.
	RAMOcupada  int `json:"ramocupada"`  // RAMOcupada represents the amount of occupied RAM.
	TotalCPU    int `json:"totalcpu"`    // TotalCPU represents the total number of CPUs.
}

// Counters represents process counters.
type Counters struct {
	Running  int `json:"running"`  // Running represents the number of running processes.
	Sleeping int `json:"sleeping"` // Sleeping represents the number of sleeping processes.
	Stopped  int `json:"stopped"`  // Stopped represents the number of stopped processes.
	Zombie   int `json:"zombie"`   // Zombie represents the number of zombie processes.
	Total    int `json:"total"`    // Total represents the total number of processes.
}

// AllData represents all system data.
type AllData struct {
	AllGenerales    []General       `json:"AllGenerales"`    // AllGenerales represents a list of general system information.
	AllTipoProcesos []Process       `json:"AllTipoProcesos"` // AllTipoProcesos represents a list of process information.
	AllProcesos     []Counters      `json:"AllProcesos"`     // AllProcesos represents a list of process counters.
}

// MemoryBlock represents a memory block.
type MemoryBlock struct {
	InitialAddress string   `json:"initial_address"` // InitialAddress represents the initial address of the memory block.
	FinalAddress   string   `json:"final_address"`   // FinalAddress represents the final address of the memory block.
	Permissions    []string `json:"permissions"`     // Permissions represents the permissions associated with the memory block.
	Device         string   `json:"device"`          // Device represents the device associated with the memory block.
	File           string   `json:"file"`            // File represents the file associated with the memory block.
	Size           float64  `json:"size"`            // Size represents the size of the memory block.
	Rss            float64  `json:"rss"`             // Rss represents the resident set size (RSS) of the memory block.
}

// MemoryResult represents the result of memory information.
type MemoryResult struct {
	TotalSize   float64        `json:"total_size"`   // TotalSize represents the total size of memory.
	TotalRss    float64        `json:"total_rss"`    // TotalRss represents the total resident set size (RSS) of memory.
	Blocks      []MemoryBlock  `json:"blocks"`       // Blocks represents a list of memory blocks.
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

// smapsHandler handles the HTTP request for retrieving process memory information.
// It expects a POST request with the process ID (PID) in the request body.
// It returns memory information for the specified process as a JSON response.
// w - http.ResponseWriter: The response writer used to write the HTTP response.
// r - *http.Request: The HTTP request received.
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

	blocks, totalSize, totalRss := parseSmapsOutput(string(output)) // Parse the smaps output and get blocks, total size, and total RSS
	result := MemoryResult{
		Blocks:    blocks,
		TotalSize: totalSize,
		TotalRss:  totalRss,
	}

	response, err := json.Marshal(result) // Convert the result to JSON format
	if err != nil {
		http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the response header as JSON
	w.Write(response) // Write the response to the HTTP response body
}

// parseSmapsOutput parses the output of the smaps command and extracts memory block information, total size, and total RSS.
// It takes a string parameter 'output' representing the output of the smaps command.
// It returns a slice of MemoryBlock representing individual memory blocks, and two float64 values for total size and total RSS.
func parseSmapsOutput(output string) ([]MemoryBlock, float64, float64) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanLines)
	// Data to send to frontend enviroment
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
			// parsing data fields
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
	// Returning data
	return blocks, totalSize, totalRss
}

// extractValue extracts the numerical value from a line of text.
// It takes a string parameter 'line' representing the line of text.
// It returns an integer value extracted from the line, or 0 if no value is found.
func extractValue(line string) int {
	re := regexp.MustCompile(`\d+`) // Regular expression to match numerical values
	match := re.FindString(line)   // Find the first numerical value in the line
	if match != "" {
		value := match
		size, _ := strconv.Atoi(value) // Convert the matched value to an integer
		return size
	}
	return 0
}

// mapPermissions maps the permissions string to a list of human-readable permission names.
// It takes a string parameter 'permissions' representing the permissions string.
// It returns a slice of strings representing the mapped permission names.
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
