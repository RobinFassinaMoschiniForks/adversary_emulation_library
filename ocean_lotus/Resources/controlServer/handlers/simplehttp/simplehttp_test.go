package simplehttp_test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"attackevals.mitre.org/control_server/config"
	"attackevals.mitre.org/control_server/handlers/simplehttp"
	"attackevals.mitre.org/control_server/restapi"
	"attackevals.mitre.org/control_server/sessions"
	"attackevals.mitre.org/control_server/tasks"
)

var restAPIlistenHost = "127.0.0.1:9998"
var registerImplantURL = "http://127.0.0.1:8080/register"
var getTaskURL = "http://127.0.0.1:8080/task"
var postOutputURL = "http://127.0.0.1:8080/output"
var getFileURL = "http://127.0.0.1:8080/getFile"
var putFileURL = "http://127.0.0.1:8080/putFile"

// exampleBeacon shows how to construct a well formed beacon
var exampleBeacon = sessions.Session{
	GUID:     "simplehttp-implant",
	IPAddr:   "127.0.0.1",
	HostName: "myHostName",
	User:     "myUserName",
	Dir:      "C:\\MyDir\\",
	PID:      1234,
	PPID:     4,
	Task:     &exampleTask,
}

var exampleTask = tasks.Task{
	GUID     : "example-task",
	Command  : "whoami",
	Output   : "",
	Status   : tasks.TASK_STATUS_NEW,
	ExitCode : -1,
}

var handler = &simplehttp.HttpHandler{}
var configEntry = config.HandlerConfigEntry{
	"host": "127.0.0.1",
	"port": "8080",
}

func createTestTaskForSession(t *testing.T) {
	// Create Test Task
	createTaskUrl := fmt.Sprintf("http://%s/api/v1.0/session/%s/task", restAPIlistenHost, exampleBeacon.GUID)
	// setup HTTP POST request
	req, err := http.NewRequest("POST", createTaskUrl, bytes.NewBufferString(exampleTask.Command))

	req.Header.Set("X-Task-Guid", exampleTask.GUID)
	req.Close = true
	if err != nil {
		t.Error(err)
	}

	// execute HTTP POST and read response
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		t.Errorf("Expected error code 200, got %v", response.StatusCode)
	}
}

func startSimpleHTTPHandler(t *testing.T) {
	if err := handler.StartHandler(restAPIlistenHost, configEntry); err != nil {
		t.Errorf("Error when starting HTTP handler: %s", err.Error())
	}
	time.Sleep(50 * time.Millisecond)
}

func stopSimpleHTTPHandler(t *testing.T) {
	if err := handler.StopHandler(); err != nil {
		t.Errorf("Error when stopping HTTP handler: %s", err.Error())
	}
	time.Sleep(50 * time.Millisecond)
}

func startRESTAPI() {
	restapi.Start(restAPIlistenHost, "./test_payloads")
	time.Sleep(50 * time.Millisecond)
}

func stopRESTAPI() {
	restapi.Stop()
	time.Sleep(50 * time.Millisecond)
}

func registerImplant(registerURL string) (string, error) {

	// convert exampleBeacon data into JSON
	exampleBeaconJSON, err := json.Marshal(exampleBeacon)
	if err != nil {
		return "", err
	}

	// setup HTTP POST request
	req, err := http.NewRequest("POST", registerURL, bytes.NewBuffer(exampleBeaconJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// execute HTTP POST request and read response
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}
func TestStartStopTrickBotHandler(t *testing.T) {
	startSimpleHTTPHandler(t)
	stopSimpleHTTPHandler(t)
}

func TestResgisterImplant(t *testing.T) {

	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)

	// test implant registration against handler
	response, err := registerImplant(registerImplantURL)
	if err != nil {
		t.Error(err)
	}

	expectedOutput := "Successfully added session."
	if string(response) != expectedOutput {
		t.Errorf("Got '%v' expected '%v'", response, expectedOutput)
	}

	// delete implant?
}

func TestGetTask(t *testing.T) {
	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)

	url := getTaskURL + "/" + exampleBeacon.GUID
	response, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}

	expectedOutput := "whoami"
	if string(body) != expectedOutput {
		t.Errorf("Got '%v' expected '%v'", response, expectedOutput)
	}
}

func TestPostTaskOutput(t *testing.T) {

	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)

	createTestTaskForSession(t)

	url := postOutputURL + "/" + exampleBeacon.GUID
	output := "this is example task output"

	// setup HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(output))
	if err != nil {
		t.Error(err)
	}

	// execute HTTP POST and read response
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Errorf("Expected error code 200, got %v", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	returnedOutput := string(body)

	expectedOutput := "Successfully set task output."
	if returnedOutput != expectedOutput {
		t.Errorf("Expected \"%v\", got \"%v\"", expectedOutput, returnedOutput)
	}
}

func TestForwardGetFileFromServer(t *testing.T) {

	// set current working directory to wizard_spider/Resources/control_server
	// this is needed so that the unit tests can find ./config/config.yml
	cwd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(cwd) // restore cwd at end of test

	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)

	resp, err := simplehttp.ForwardGetFileFromServer("hello_world.elf")
	if err != nil {
		t.Error(err)
	}
	fileData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	goodHash := "fe7c47d38224529c7d8f9a11a62cdd7a"
	h := md5.Sum([]byte(fileData))
	actualHash := hex.EncodeToString(h[:])
	if goodHash != actualHash {
		t.Errorf("Expected %v, got %v", goodHash, actualHash)
	}

	// test non-existent file
	resp2, err := simplehttp.ForwardGetFileFromServer("does_not_exist.txt")
	if err != nil {
		t.Error(err)
	}

	expectedErr := "404 Not Found"
	actualErr := resp2.Status
	if expectedErr != actualErr {
		t.Errorf("expected '%v' got '%v'", expectedErr, actualErr)
	}
}

func TestGetFileFromServer(t *testing.T) {
	// set current working directory to wizard_spider/Resources/control_server
	// this is needed so that the unit tests can find ./config/config.yml
	cwd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(cwd) // restore cwd at end of test

	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)

	// download the test file
	url := getFileURL + "/hello_world.elf"
	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	// Create the file
	filepath := "/tmp/hello_world.elf"
	out, err := os.Create(filepath)
	if err != nil {
		t.Error(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		t.Error(err)
	}

	// compare file hashes
	goodHash := "fe7c47d38224529c7d8f9a11a62cdd7a"
	fileData, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Error(err)
	}
	h := md5.Sum(fileData)
	actualHash := hex.EncodeToString(h[:])
	if goodHash != actualHash {
		t.Errorf("Expected %v, got %v", goodHash, actualHash)
	}

	// test non-existent file
	url = getFileURL + "/does_not_exist.txt"
	resp2, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}

	expectedErr := "404 Not Found"
	actualErr := resp2.Status
	if expectedErr != actualErr {
		t.Errorf("expected '%v' got '%v'", expectedErr, actualErr)
	}
}

func TestForwardPostFileToServer(t *testing.T) {
	// set current working directory to wizard_spider/Resources/control_server
	// this is needed so that the unit tests can find ./config/config.yml
	cwd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(cwd) // restore cwd at end of test

	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)

	testFile := "./test_payloads/hello_world.elf"
	fileNameOnUpload := "test_binary.elf"
	fileData, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Error(err)
	}

	// invoke file upload
	got, err := simplehttp.ForwardPostFileToServer(fileNameOnUpload, fileData)
	if err != nil {
		t.Error(err)
	}
	// validate response
	want := "Successfully uploaded file to control server at './files/test_binary.elf'"
	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	// confirm file made it to disk properly
	// read test file bytes
	uploadedFile := "./files/" + fileNameOnUpload
	uploadedFileData, err := ioutil.ReadFile(uploadedFile)
	if err != nil {
		t.Error(err)
	}
	// compare file hashes
	goodHash := "fe7c47d38224529c7d8f9a11a62cdd7a"
	h := md5.Sum(uploadedFileData)
	actualHash := hex.EncodeToString(h[:])
	if goodHash != actualHash {
		t.Errorf("Expected %v, got %v", goodHash, actualHash)
	}

	// clean up test file
	err = os.Remove(uploadedFile)
	if err != nil {
		t.Error(err)
	}
}

func TestPostFileToServer(t *testing.T) {
	// set current working directory to wizard_spider/Resources/control_server
	// this is needed so that the unit tests can find ./config/config.yml
	cwd, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(cwd) // restore cwd at end of test

	// start REST API
	startRESTAPI()
	defer stopRESTAPI()

	// start handler
	startSimpleHTTPHandler(t)
	defer stopSimpleHTTPHandler(t)
	// read file to upload
	// pass file data to function
	// validate response

	testFile := "./test_payloads/hello_world.elf"
	testFileName := "test_file.elf"
	fileData, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Error(err)
	}

	url := putFileURL + "/" + testFileName
	response, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(fileData))
	if err != nil {
		t.Error(err)
	}
	s, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	temp := string(s)
	got := strings.TrimSuffix(temp, "\n")
	want := "Successfully uploaded file to control server at './files/test_file.elf'"
	if got != want {
		t.Errorf("Expected: \n%v \ngot: \n%v", want, got)
	}
	// confirm file made it to disk properly
	// read test file bytes
	uploadedFile := "./files/" + testFileName
	uploadedFileData, err := ioutil.ReadFile(uploadedFile)
	if err != nil {
		t.Error(err)
	}
	// compare file hashes
	goodHash := "fe7c47d38224529c7d8f9a11a62cdd7a"
	h := md5.Sum(uploadedFileData)
	actualHash := hex.EncodeToString(h[:])
	if goodHash != actualHash {
		t.Errorf("Expected %v, got %v", goodHash, actualHash)
	}

	// clean up test file

	err = os.Remove(uploadedFile)
	if err != nil {
		t.Error(err)
	}
}
