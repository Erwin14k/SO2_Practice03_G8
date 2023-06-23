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

// MemoryRow represents a row of memory information
type MemoryRow struct {
	InitialAddress string   `json:"initial_address"`
	FinalAddress   string   `json:"final_address"`
	Permissions    string   `json:"permissions"`
	Device         string   `json:"device"`
	File           string   `json:"file"`
	RSS            int      `json:"rss_mb"`
	Size           int      `json:"size_mb"`
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

// handleMemory handles the route "/memory" by reading the request body, extracting the PID (process ID), and retrieving the memory information of the process.
// It executes the "cat /proc/pid/maps" command, parses the output, generates JSON data, and sends it as the response.
// w - http.ResponseWriter: The response writer used to write the HTTP response.
// r - *http.Request: The HTTP request received.
func handleMemory(w http.ResponseWriter, r *http.Request) {
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

	cmd := exec.Command("sudo", "cat", fmt.Sprintf("/proc/%d/smaps", pid)) // Create a command to execute "cat /proc/pid/smaps"
	output, err := cmd.Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error executing the command
		fmt.Fprintln(w, "Error reading process memory")
		return
	}

	rows := parseMemoryRows(string(output)) // Parse the memory rows from the command output
	jsonData, err := json.Marshal(rows)     // Generate JSON data from the parsed rows
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error generating JSON
		fmt.Fprintln(w, "Error generating JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)              // Set HTTP 200 OK status code
	fmt.Println("Information: Process Memory PID", pid)
	fmt.Fprintln(w, string(jsonData))         // Write the process data to the response
}


// parseMemoryRows parses the memory rows from the output of cat /proc/pid/maps command and returns an array of MemoryRow objects.
// output - string: The output of cat /proc/pid/maps command.
// Returns []MemoryRow: An array of MemoryRow objects representing the parsed memory rows.
func parseMemoryRows(output string) []MemoryRow {
	rows := []MemoryRow{}

	lines := strings.Split(output, "\n")
	var row MemoryRow
	for _, line := range lines {
		if line == "" {
			if row.File != "" {
				rows = append(rows, row)
			}
			row = MemoryRow{}
			continue
		}

		if strings.Contains(line, "-") {
			fields := strings.Fields(line)
			if len(fields) < 6 {
				continue
			}

			row.InitialAddress = strings.Split(fields[0], "-")[0]
			row.FinalAddress = strings.Split(fields[0], "-")[1]
			row.Permissions = fields[1]
			row.Device = fields[3]
			row.File = fields[5]
		} else if strings.HasPrefix(line, "Rss:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			value, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			row.RSS = value / 1024 // Convert KB to MB
		} else if strings.HasPrefix(line, "Size:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			value, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			row.Size = value / 1024 // Convert KB to MB
		}
	}

	return rows
}

// calculateSizeKB calculates the size in kilobytes from the initial and final addresses.
// initialAddress - string: The initial address.
// finalAddress - string: The final address.
// Returns int: The size in kilobytes.
func calculateSizeKB(initialAddress, finalAddress string) int {
	initial, err := strconv.ParseUint(initialAddress, 16, 64)
	if err != nil {
		return 0
	}

	final, err := strconv.ParseUint(finalAddress, 16, 64)
	if err != nil {
		return 0
	}

	size := final - initial
	sizeKB := int(size) / 1024

	return sizeKB
}

// parsePermissions parses the permissions string and returns a list of permissions.
// permissions - string: The permissions string.
// Returns []string: A list of permissions.
func parsePermissions(permissions string) []string {
	perms := []string{}

	if strings.Contains(permissions, "r") {
		perms = append(perms, "Lectura")
	}
	if strings.Contains(permissions, "w") {
		perms = append(perms, "Escritura")
	}
	if strings.Contains(permissions, "x") {
		perms = append(perms, "Ejecucion")
	}

	return perms
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
	router.HandleFunc("/memory", handleMemory).Methods("POST")

	handler := cors.Default().Handler(router) // Create a new CORS handler with default settings
	log.Fatal(http.ListenAndServe(":8080", handler)) // Start the HTTP server and listen on port 8080

	fmt.Println("Server on port 8080")
}
