package main

import (
	"fmt"
	"net"
	"strconv"
	"testing"
)

func Test_isTcpPortOpen(t *testing.T) {

	listenPort := 34582
	ln1, err := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort))
	if err != nil {
		fmt.Println("ERROR: failed to bind to tcp port ", listenPort, "and encountered error")
		t.Errorf("unable to tcp listen on port %v", strconv.Itoa(listenPort))
	}
	defer func(ln1 net.Listener) {
		ln1.Close()
	}(ln1)

	listenPort2 := 34583
	ln2, err2 := net.Listen("tcp", "localhost:"+strconv.Itoa(listenPort2))
	if err2 != nil {
		fmt.Println("ERROR: failed to bind to tcp port ", listenPort2, "and encountered error")
		t.Errorf("unable to tcp listen on port %v", strconv.Itoa(listenPort2))
		return
	}
	err = ln2.Close()
	if err != nil {
		t.Errorf("unable to close tcp port port %v", strconv.Itoa(listenPort2))
		return
	}

	type args struct {
		targetHostWithPort string
		portTimeout        int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should be able to connect to tcp port " + strconv.Itoa(listenPort),
			args: args{
				targetHostWithPort: "localhost:" + strconv.Itoa(listenPort),
				portTimeout:        0,
			},
			want: true,
		},
		{
			name: "should not be able to connect to tcp port " + strconv.Itoa(listenPort2),
			args: args{
				targetHostWithPort: "localhost:" + strconv.Itoa(listenPort2),
				portTimeout:        0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTcpPortOpen(tt.args.targetHostWithPort, tt.args.portTimeout); got != tt.want {
				t.Errorf("isTcpPortOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkForInvalidHostname(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test localhost",
			args: args{
				hostname: "localhost",
			},
			wantErr: false,
		},
		{
			name: "test bad hostname",
			args: args{
				hostname: "badHostname",
			},
			wantErr: true,
		},
		{
			name: "test ip 127.0.0.1",
			args: args{
				hostname: "127.0.0.1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkForInvalidHostname(tt.args.hostname); (err != nil) != tt.wantErr {
				t.Errorf("checkForInvalidHostname() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
