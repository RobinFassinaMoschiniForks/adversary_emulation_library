package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"attackevals.mitre.org/control_server/config"
	"attackevals.mitre.org/control_server/display"
	"attackevals.mitre.org/control_server/handlers/util"
	"attackevals.mitre.org/control_server/logger"
	"attackevals.mitre.org/control_server/sessions"
	"github.com/gorilla/mux"
)

var Server *http.Server

// Start enables the REST API server
func Start(listenAddress, payloadDir string) {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1.0/version", GetVersion).Methods("GET")
	r.HandleFunc("/api/v1.0/config", GetConfig).Methods("GET")
	r.HandleFunc("/api/v1.0/sessions", GetSessions).Methods("GET")
	r.HandleFunc("/api/v1.0/session", CreateSession).Methods("POST")
	r.HandleFunc("/api/v1.0/session/{guid}", GetSessionByName).Methods("GET")
	r.HandleFunc("/api/v1.0/session/delete/{guid}", RemoveSession).Methods("DELETE")
	r.HandleFunc("/api/v1.0/task/{guid}", GetTask).Methods("GET")
	r.HandleFunc("/api/v1.0/task/{guid}", SetTask).Methods("POST")
	r.HandleFunc("/api/v1.0/task/{guid}", RemoveTask).Methods("DELETE")
	r.HandleFunc("/api/v1.0/bootstraptask/{handler}", GetBootstrapTask).Methods("GET")
	r.HandleFunc("/api/v1.0/bootstraptask/{handler}", SetBootstrapTask).Methods("POST")
	r.HandleFunc("/api/v1.0/bootstraptask/{handler}", RemoveBootstrapTask).Methods("DELETE")
	r.HandleFunc("/api/v1.0/task/output/{guid}", GetTaskOutput).Methods("GET")
	r.HandleFunc("/api/v1.0/task/output/{guid}", SetTaskOutput).Methods("POST")
	r.HandleFunc("/api/v1.0/task/output/{guid}", RemoveTaskOutput).Methods("DELETE")
	r.HandleFunc("/api/v1.0/upload/{file}", UploadFile).Methods("POST")

	// serve files located in specified payload directory
	r.PathPrefix("/api/v1.0/files/").Handler(http.StripPrefix("/api/v1.0/files/", http.FileServer(http.Dir(payloadDir))))

	Server = &http.Server{
		Addr:         listenAddress,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	// start rest api in goroutine so it doesn't block
	go func() {
		err := Server.ListenAndServe()
		if err != nil && err.Error() != "http: Server closed" {

			if strings.Contains(err.Error(), fmt.Sprintf("listen tcp %s: bind: address already in use", listenAddress)) {
				logger.Warning(err)
				return
			} else {
				logger.Error(err)
			}
		}
	}()
}

// Stop disables the REST API server
func Stop() {
	emptyContext := context.Background()
	err := Server.Shutdown(emptyContext)
	if err != nil {
		logger.Error(fmt.Sprintf("REST server failed to shut down: %s", err.Error()))
	}
}

// GetVersion returns the current application version
func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ATT&CK Evaluations Control Server 1.0\n"))
}

// GetConfig returns the current server configuration
func GetConfig(w http.ResponseWriter, r *http.Request) {
	configuration, err := config.GetRestAPIConfig()
	if err != nil {
		e := err.Error()
		w.Write([]byte(e))
	}
	w.Write(configuration)
}

// GetStagers returns a list of files that are used to stage implants on target
// for example, droppers, download-exec 1-liners, etc.
func GetStagers() error {
	return errors.New("not implemented")
}

// GetStagerOptions returns detailed information and options for a given stager
func GetStagerOptions() error {
	return errors.New("not implemented")
}

// SetStager returns a stager configured with options specified by the user
func SetStager() error {
	return errors.New("not implemented")
}

// CreateSession registers a new C2 session to the control server
func CreateSession(w http.ResponseWriter, r *http.Request) {
	// read POST request body from client
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	// convert JSON data into a session struct
	var session sessions.Session
	err = json.Unmarshal(req, &session)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	// pass session struct to handler
	err = sessions.AddSession(session)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	// send success message to client
	_, err = fmt.Fprint(w, "successfully added session")
	if err != nil {
		logger.Error(err)
		return
	}
	// display new session on console
	display.PrintSession(session)
}

// GetSessions returns a list of all current C2 sessions
func GetSessions(w http.ResponseWriter, r *http.Request) {
	sessionList := sessions.GetSessionList()
	sessionJson, err := json.Marshal(sessionList)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	// send session list to client
	fmt.Fprintf(w, "%s", sessionJson)

	// display session list on console
	display.PrintSessionList(sessionList)
}

// GetSessionByName returns detailed info for the specified session
func GetSessionByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	theSession, err := sessions.GetSessionByName(guid)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	sessionJson, err := json.Marshal(theSession)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}
	// send session list to client
	fmt.Fprintf(w, "%s", sessionJson)
}

// RemoveSession terminates the specified C2 session
func RemoveSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	err := sessions.RemoveSession(guid)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v%v", "successfully removed session: ", guid)

}

// SetTask allows users to issue commands to the specified session
func SetTask(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received SetTask request")
	vars := mux.Vars(r)
	guid := vars["guid"]
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	task := string(req)
	err = sessions.SetTask(guid, task)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v", "successfully set task")
	logger.Success("Successfully set task")
}

// GetTask allows users to issue commands to the specified session
func GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	task, err := sessions.GetTask(guid)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v", task)

	// clear task so we don't keep sending it
	err = sessions.RemoveTask(guid)
	if err != nil {
		logger.Error(err)
	}
}

// RemoveTask deletes the currently queued tasking
func RemoveTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	err := sessions.RemoveTask(guid)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v%v", "successfully removed task for session: ", guid)
}

// SetBootstrapTask allows users to set a default bootstrap task for sessions that support this feature.
func SetBootstrapTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	handlerName := strings.ToLower(vars["handler"])
	logger.Info(fmt.Sprintf("Received SetBootstrapTask request for handler %s", handlerName))
	if _, ok := util.RunningHandlers[handlerName]; !ok {
		errMsg := fmt.Sprintf("Handler %s is not currently running. Failed to set bootstrap task.", handlerName)
		logger.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	task := string(req)
	sessions.SetBootstrapTask(handlerName, task)
	fmt.Fprintf(w, "successfully set bootstrap task for handler %s", handlerName)
	logger.Success(fmt.Sprintf("Successfully set bootstrap task for handler %s", handlerName))
}

// GetBootstrapTask returns the currently set bootstrap task.
func GetBootstrapTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	handlerName := strings.ToLower(vars["handler"])
	task := sessions.GetBootstrapTask(handlerName)
	fmt.Fprintf(w, "%v", task)
}

// RemoveBootstrapTask deletes the currently set bootstrap task
func RemoveBootstrapTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	handlerName := strings.ToLower(vars["handler"])
	if _, ok := util.RunningHandlers[handlerName]; !ok {
		errMsg := fmt.Sprintf("Handler %s is not currently running. Cannot manage bootstrap tasks.", handlerName)
		logger.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errMsg))
		return
	}
	sessions.RemoveBootstrapTask(handlerName)
	logger.Success(fmt.Sprintf("Successfully removed bootstrap task for handler %s", handlerName))
	fmt.Fprintf(w, "successfully removed bootstrap task for handler %s", handlerName)
}

// SetTaskOutput sets the output following task execution
func SetTaskOutput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	taskOutput := string(req)
	err = sessions.SetTaskOutput(guid, taskOutput)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v", "successfully set task output")
	logger.Info("Received task output for session: ", guid)
	logger.Task(taskOutput)
}

// GetTaskOutput returns the last task's output
func GetTaskOutput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	taskOutput, err := sessions.GetTaskOutput(guid)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v", taskOutput)
}

// RemoveTaskOutput clears the current task's output
func RemoveTaskOutput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid := vars["guid"]
	err := sessions.DeleteTaskOutput(guid)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%v%v", "successfully deleted task output for session: ", guid)

}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// get uploaded file name
	vars := mux.Vars(r)
	fileName := vars["file"]

	// read file data from HTTP POST request
	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	// write file data to disk
	filePath := "./files/" + fileName
	err = ioutil.WriteFile(filePath, req, 0444)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "Successfully uploaded file to control server at './files/%v'\n", fileName)
}
