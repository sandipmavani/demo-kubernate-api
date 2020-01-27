package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Payload struct {
	FilterBy string `json:"filter_by"`
}

// Ping is only for testing purposes to make sure the service is up and responding
func Ping(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "OK")
}

func EnumerateRBAC(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read request body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "body is empty", http.StatusBadRequest)
		return
	}

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "can't parse body to JSON", http.StatusBadRequest)
		return
	}

	// TODO: sanitize payload.FilterBy
	roles := getListOfRoles(payload.FilterBy)

	jsonBytes, err := json.Marshal(roles)
	if err != nil {
		http.Error(w, "can't make response in JSON format", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(jsonBytes))
}

// this code should work only inside a Kubernetes cluster
func getListOfRoles(filterBy string) []v1.Subject {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	roleBindings, err := clientset.RbacV1().RoleBindings("").List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error loading cluster role bindings, reason is " + err.Error())
	}

	var prefixRegexp = regexp.MustCompile(`^` + filterBy)
	arrRolesFiltered := make([]v1.Subject, 0)

	for _, roleBinding := range roleBindings.Items {
		for _, subject := range roleBinding.Subjects {
			if prefixRegexp.MatchString(subject.Name) {
				arrRolesFiltered = append(arrRolesFiltered, subject)
			}
		}
	}

	return arrRolesFiltered
}
