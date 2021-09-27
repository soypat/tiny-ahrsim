package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/soypat/ahrs"
	"github.com/spf13/pflag"
)

func main() {
	var (
		tty  = "/dev/ttyUSB0"
		port = ":8080"
	)
	pflag.StringVarP(&tty, "ttl", "f", tty, "Device USB port/file. Has form COM1 on windows. Has form /dev/ttyACM on linux.")
	pflag.StringVarP(&port, "port", "p", port, "TCP port on which to serve HTTP.")

	pflag.Parse()
	fp, err := os.Open(tty)
	must(err)
	rd := bufio.NewReader(fp)
	imu := &fieldReaderIMU{r: rd}
	estimater := ahrs.NewXioARS(5, imu)
	tLast := time.Now()
	// IMU reader goroutine
	go func() {
		tick := time.NewTicker(5 * time.Millisecond)
		for tm := range tick.C {
			err = imu.Update(fp)
			if err != nil {
				log.Print(err)
			}
			estimater.Update(time.Since(tLast).Seconds())
			tLast = tm
		}
	}()

	// Server handler
	http.HandleFunc("/attitude", func(rw http.ResponseWriter, r *http.Request) {
		log.Print("serving attitude")
		if allowCORS(rw, r) {
			return
		}
		q := estimater.GetQuaternion()
		rot := ahrs.RotationMatrixFromQuat(q)
		angles := rot.TaitBryan(ahrs.OrderXYZ)
		tLast = time.Now()
		e := json.NewEncoder(rw)
		must(e.Encode(&angles))
	})
	srv := http.Server{
		Addr:    port,
		Handler: http.DefaultServeMux,
	}
	log.Print("serving on port ", port)
	srv.ListenAndServe()
}

func allowCORS(rw http.ResponseWriter, r *http.Request) (aborted bool) {
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Max-Age", "999999")             // Allow javascript requests from localhost
	rw.Header().Set("Access-Control-Allow-Methods", "GET,POST")     // We will use GET and POST methods
	rw.Header().Set("Access-Control-Allow-Headers", "content-type") // allow JSON interchange
	rw.Header().Set("Access-Control-Allow-Origin", "*")             // Allow javascript requests from localhost
	if r.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

// must asserts error is nil. If error is not nil program ends with error message.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

type fieldReaderIMU struct {
	sync.Mutex
	r *bufio.Reader
	a [3]int32
	g [3]int32
}

func (imu *fieldReaderIMU) Update(r io.Reader) error {
	rd := imu.r
	rd.Reset(r)        // discard old data
	rd.ReadBytes('\n') // Discard all data until next newline
	s, err := rd.ReadString('\n')
	if err != nil {
		return err
	}
	values := strings.Fields(s)
	if len(values) < 6 {
		return errors.New("not enough fields read")
	}
	log.Print("got values: ", values)
	ival := make([]int32, 6)
	for i := range ival {
		I, err := strconv.Atoi(values[i])
		if err != nil {
			return err
		}
		ival[i] = int32(I)
	}
	copy(imu.a[:], ival)
	copy(imu.g[:], ival[3:])
	return nil
}

func (imu *fieldReaderIMU) Acceleration() (ax, ay, az int32) {
	imu.Lock()
	defer imu.Unlock()
	return imu.a[0], imu.a[1], imu.a[2]
}
func (imu *fieldReaderIMU) AngularVelocity() (gx, gy, gz int32) {
	imu.Lock()
	defer imu.Unlock()
	return imu.g[0], imu.g[1], imu.g[2]
}
