package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type DateResponse struct {
	Date string `json:"date"`
}

type ShellResponse struct {
	KernelVersion string `json:"kernel_version"`
}

type PodInfoResponse struct {
	PodName string `json:"pod_name"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func getDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DateResponse{Date: time.Now().Format(time.RFC3339)})
}

func printRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var data interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getShell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cmd := exec.Command("uname", "-r")
	out, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShellResponse{KernelVersion: string(out)})
}

func getPodName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	podName := os.Getenv("HOSTNAME") // In Kubernetes, HOSTNAME env var is usually the pod name
	json.NewEncoder(w).Encode(PodInfoResponse{PodName: podName})
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Simple Web App</title>
		<style>
			body { font-family: sans-serif; margin: 20px; }
			button { padding: 10px 20px; margin: 5px; cursor: pointer; }
			.response { margin-top: 20px; padding: 10px; border: 1px solid #ccc; background-color: #f9f9f9; white-space: pre-wrap; }
			textarea { width: 80%; height: 100px; margin-top: 10px; padding: 5px; }
		</style>
	</head>
	<body>
		<h1>Welcome to Simple Web App</h1>
		<p>Click the buttons below to interact with the services:</p>

		<button onclick="getHealth()">Health Check (GET)</button>
		<button onclick="getDate()">Get Date (GET)</button>

		<p>Enter text for /print (POST):</p>
		<textarea id="printRequestBody">Hello World!</textarea>
		<button onclick="postPrint()">Print Request Body (POST)</button>
		<button onclick="postShell()">Get Kernel Version (POST)</button>

		<button onclick="getPodName()">Get Pod Name (GET)</button>

		<div class="response" id="response"></div>

		<script>
			async function getHealth() {
				const responseDiv = document.getElementById('response');
				try {
					const response = await fetch('/health');
					const result = await response.json();
					responseDiv.textContent = 'Response from /health:\n' + JSON.stringify(result, null, 2);
				} catch (error) {
					responseDiv.textContent = 'Error calling /health:\n' + error;
				}
			}

			async function getDate() {
				const responseDiv = document.getElementById('response');
				try {
					const response = await fetch('/date');
					const result = await response.json();
					responseDiv.textContent = 'Response from /date:\n' + JSON.stringify(result, null, 2);
				} catch (error) {
					responseDiv.textContent = 'Error calling /date:\n' + error;
				}
			}

			async function postPrint() {
                const responseDiv = document.getElementById('response');
                try {
                    const textarea = document.getElementById('printRequestBody');
                    const rawData = textarea.value;
                    const dataToSend = { "input": rawData }; // Wrap the string in a JSON object

                    const response = await fetch('/print', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(dataToSend)
                    });
                    const result = await response.json();
                    responseDiv.textContent = 'Response from /print:\n' + JSON.stringify(result, null, 2);
                } catch (error) {
                    responseDiv.textContent = 'Error calling /print:\n' + error;
                }
            }

			async function postShell() {
				const responseDiv = document.getElementById('response');
				try {
					const response = await fetch('/shell', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ action: "get_kernel" })
					});
					const result = await response.json();
					responseDiv.textContent = 'Response from /shell:\n' + JSON.stringify(result, null, 2);
				} catch (error) {
					responseDiv.textContent = 'Error calling /shell:\n' + error;
				}
			}

			async function getPodName() {
				const responseDiv = document.getElementById('response');
				try {
					const response = await fetch('/pod-name');
					const result = await response.json();
					responseDiv.textContent = 'Response from /pod-name:\n' + JSON.stringify(result, null, 2);
				} catch (error) {
					responseDiv.textContent = 'Error calling /pod-name:\n' + error;
				}
			}
		</script>
	</body>
	</html>
	`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/date", getDate)
	http.HandleFunc("/print", printRequest)
	http.HandleFunc("/shell", getShell)
	http.HandleFunc("/pod-name", getPodName)

	fmt.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}