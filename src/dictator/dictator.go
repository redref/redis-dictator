package main

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func Run(conf DictatorConfiguration, stop <-chan bool, finished chan<- bool) {
	var master Redis // Dummy master node
	err := master.Initialize("", conf.MasterService, conf.Node.Port, conf.Node.LoadingTimeout)
	if err != nil {
		log.WithError(err).Warn("Fail to initialize Redis node")
		finished <- true
	}

	var re Redis // Create a Redis Node
	err = re.Initialize(conf.Node.Name, conf.Node.Host, conf.Node.Port, conf.Node.LoadingTimeout)
	if err != nil {
		log.WithError(err).Warn("Fail to initialize Redis node")
		finished <- true
	}

	// Set default to slave
	re.SetRole("SLAVE", &master)
	log.Info("Node started, new status: ", re.Role, " of ", master.Host)

	// http signals management
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, re.Role)
		log.Info("Call to node status, status was: ", re.Role)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := re.Connect()
		if err != nil {
			fmt.Fprintf(w, "OK")
			log.Info("Ping failed")
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Info("Ping success")
		}
	})
	http.HandleFunc("/promote", func(w http.ResponseWriter, r *http.Request) {
		re.SetRole("MASTER", nil)
		log.Info("Node promoted, new status: ", re.Role)
	})
	http.HandleFunc("/demote", func(w http.ResponseWriter, r *http.Request) {
		re.SetRole("SLAVE", &master)
		log.Info("Node demoted, new status: ", re.Role, " of ", master.Host)
	})
	go http.ListenAndServe(":"+strconv.Itoa(conf.HttpPort), nil)

	// Wait for the stop signal
Loop:
	for {
		select {
		case hasToStop := <-stop:
			if hasToStop {
				log.Debug("Close Signal Received!")
			} else {
				log.Debug("Close Signal Received (but a strange false one)")
			}
			break Loop
		}
	}

	finished <- true
}
