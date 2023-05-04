package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func createEC2Instance(w http.ResponseWriter, r *http.Request) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("AWS_KEY", "AWS_SECRET_KEY", ""),
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}
	svc := ec2.New(sess)

	params := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-007855ac798b5175e"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	}
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	resp, err := svc.RunInstances(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	instanceID := *resp.Instances[0].InstanceId
	fmt.Fprintf(w, "Created instance %s", instanceID)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func terminateEC2Instance(w http.ResponseWriter, r *http.Request) {

	instanceID := r.FormValue("instance_id")

	if instanceID == "" {
		http.Error(w, "instance_id is empty", http.StatusBadRequest)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("AWS_KEY", "AWS_SECRET_KEY", ""),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	svc := ec2.New(sess)

	params := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	_, err = svc.TerminateInstances(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Terminated instance %s", instanceID)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	name := r.URL.Query().Get("name")

	message := fmt.Sprintf("Hello, %s!", name)
	fmt.Fprintln(w, message)
}

func main() {

	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/ec2/create", createEC2Instance)
	http.HandleFunc("/api/ec2/terminate", terminateEC2Instance)

	fmt.Println("Starting server on :80")
	http.ListenAndServe(":80", nil)
}
