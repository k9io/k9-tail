/*
 * Copyright (C) 2024 Key9 Identity, Inc <k9.io>
 * Copyright (C) 2024 Champ Clark III <cclark@k9.io>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License Version 2 as
 * published by the Free Software Foundation.  You may not use, modify or
 * distribute this program under any other version of the GNU General
 * Public License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA 02111-1307, USA.
 */

package main

/*****************************************************************************/
/* k9-tail								     */
/*									     */
/* This routine "follows" a file and send data to the Key9 "client logging"  */
/* API.  In the configuration file, there is a "tail_file" option. This is   */
/* the file to "follow".  This is typically the "auth.log" file.  We only    */
/* send authentication related logs (i.e - "sshd" logs).  The position of    */
/* file is tracked by the "waldo_file", which stores the last seek position. */
/*									     */
/* Champ Clark III (cclark@k9.io)					     */
/* Version 1.0 - 20240305						     */
/*****************************************************************************/

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/clarketm/json"
	"github.com/nxadm/tail"
)

func main() {

	var waldo_position int64 /* Storage for seek position */

	JSON := JSON_F{} /* See struct.go */

	debug := flag.Bool("debug", false, "Debug option")
	flag.Parse()

	/* Signal processing. */

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {

		sig := <-signalChannel

		switch sig {

		case syscall.SIGINT:

			if *debug == true {
				log.Printf("Caught SIGINT\n")
			}

			Catch_Signal()

		case syscall.SIGTERM:

			if *debug == true {
				log.Printf("Caught SIGTERM\n")
			}
			Catch_Signal()

		case syscall.SIGABRT:

			if *debug == true {
				log.Printf("Caught SIGABRT\n")
			}
			Catch_Signal()
		}

	}()

	/* Load in configuration */

	Config := LoadConfig("/opt/k9/etc/k9.yaml")

	API_KEY := fmt.Sprintf("%s:%s", Config.Authentication.Company_UUID, Config.Authentication.Api_Key)

	/* Grab local hostname */

	hostname, err := os.Hostname()

	if err != nil {

		log.Fatalf("Cannot determine local hostname (%s).\n", err)

	}

	/* No waldo file,  create one */

	waldo_data, err := os.ReadFile(Config.Tail.Waldo_File)

	if err != nil {

		err := os.WriteFile(Config.Tail.Waldo_File, []byte("{0 0}"), 0644)

		if err != nil {

			log.Fatalf("Can't write to waldo file %s. (%s)", Config.Tail.Waldo_File, err)

		}

		waldo_position = 0

	} else {

		if string(waldo_data) != "" {

			/* Waldo is stored like this: "{2227597 0}".  We only wan the
			"2227597" seek position.  We carve this out of the string for
			later use */

			splitme := strings.Split(string(waldo_data)[1:], " ")
			waldo_position, _ = strconv.ParseInt(splitme[0], 10, 64)

			if *debug {
				log.Printf("| Tail File: %s\n", Config.Tail.Tail_File)
				log.Printf("| Data: %s\n", string(waldo_data))
				log.Printf("| Waldo Split: %s\n", splitme[0])
				log.Printf("| Pre-Position: %v\n", waldo_position)
			}

		} else {

			log.Fatalf("Waldo is corrupt/incomplete!\n")

		}
	}

	if *debug {
		log.Printf("| Position: %v\n", waldo_position)
	}

	/* Open tail file and seek to position */

	t, err := tail.TailFile(Config.Tail.Tail_File, tail.Config{
		Follow:   true,
		ReOpen:   true, /* Handles log rotate, etc */
		Location: &tail.SeekInfo{Offset: waldo_position, Whence: io.SeekStart},
	})

	if err != nil {
		log.Fatalf("Can't tall file %s. (%s)\n", Config.Tail.Tail_File, err)
	}

	/* Loop for reading data as it comes in */

	for line := range t.Lines {

		/* Only process data related to sshd, please. Try and exclude audit logs. */

		if strings.Contains(line.Text, "sshd") == true && strings.Contains(line.Text, " audit[") == false && strings.Contains(line.Text, " audit:") == false {

			if *debug {
				log.Printf("| %v|%d| %s\n", line.SeekInfo, line.Num, line.Text)
			}

			JSON.Log = line.Text
			JSON.Hostname = hostname

			JSON_OUT, err := json.Marshal(JSON)

			if err != nil {
				log.Fatalf("Can't decode JSON - '%s'\n", JSON.Log)
			}

			waldo_data := fmt.Sprintf("%v", line.SeekInfo) /* Write position */

			/* POST data and look for errors */

			Status := Post_Log(API_KEY, JSON_OUT)

			for Status != "200 OK" {

				log.Printf("Got '%s' instead of '200 OK'. Retrying.....", Status)
				Status = Post_Log(API_KEY, JSON_OUT)

				time.Sleep(2 * time.Second) /* So we don't overwhelm network/CPU */

			}

			err = os.WriteFile(Config.Tail.Waldo_File, []byte(waldo_data), 0644) /* Store new position */

			if err != nil {

				log.Fatalf("Can't update waldo %s (%s).\n", Config.Tail.Waldo_File, err)

			}
		}
	}
}

/*****************************************************************/
/* Post_Log - Send data via HTTP POST request (over TLS) to Key9 */
/*****************************************************************/

func Post_Log(API_KEY string, JSON_OUT []byte) string {

	req, err := http.NewRequest("POST", Config.Tail.Client_Logging_URL, bytes.NewBuffer(JSON_OUT))

	if err != nil {
		log.Fatalf("Can't build http.NewRequest (%s)\n", err)
	}

	req.Header.Set("API_KEY", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("| Failed to make HTTP request (%s). Waiting to retry....\n", err)
		time.Sleep(2 * time.Second)
		return "Connection Error"
	}

	defer resp.Body.Close()

	return resp.Status

}

/****************************************************************/
/* Catch_Signal - Simple "signal" catcher.  May want this to do */
/* more in the future						*/
/****************************************************************/

func Catch_Signal() {

	os.Exit(0)

}
