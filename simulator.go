package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/soypat/ahrs"
	"github.com/spf13/pflag"
	"github.com/tarm/serial"
)

var dbgPrintDiv = 10

func main() {
	var (
		tty          = "/dev/ttyUSB0"
		port         = ":8080"
		baud    uint = 115200
		monitor      = false
	)
	pflag.StringVarP(&tty, "ttl", "F", tty, "Device USB port/file. Has form COM1 on windows. Has form /dev/ttyACM on linux.")
	pflag.StringVarP(&port, "port", "p", port, "TCP port on which to serve HTTP.")
	pflag.UintVarP(&baud, "baud", "", baud, "Serial (USB) baudrate.")
	pflag.IntVarP(&dbgPrintDiv, "debugDiv", "", dbgPrintDiv, "Debug print divider. Higher this number, the less prints.")
	pflag.BoolVarP(&monitor, "monitor", "m", monitor, "Just print port output.")
	pflag.Parse()

	dbgPrintDiv++ // can not be zero
	fp, err := serial.OpenPort(&serial.Config{
		Name:        tty,
		Baud:        int(baud),
		ReadTimeout: time.Second,
	})
	must(err)
	if monitor {
		_, err := io.Copy(os.Stdout, fp)
		must(err)
		log.Fatal("monitor ended")
	}

	rd := bufio.NewReader(fp)
	imu := &fieldReaderIMU{r: rd}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go imu.Run(ctx)
	estimater := ahrs.NewXioARS(0.3, imu)

	// Attitude update goroutine
	go func() {
		i := 0
		tick := time.NewTicker(10 * time.Millisecond)
		tlast := time.Now()
		for tickTime := range tick.C {
			select {
			case <-ctx.Done():
				log.Println("attitude update done")
				return
			default:
			}
			i++
			estimater.Update(time.Since(tlast).Seconds())
			tlast = tickTime
			if i%int(dbgPrintDiv) == 1 {
				log.Println("updated estimate")
			}
		}
	}()
	// Server handler
	attCounter := 0
	http.HandleFunc("/attitude", func(rw http.ResponseWriter, r *http.Request) {
		attCounter++
		if allowCORS(rw, r) {
			return
		}
		q := estimater.GetQuaternion()
		rot := ahrs.RotationMatrixFromQuat(q)
		angles := rot.TaitBryan(ahrs.OrderXYZ)
		e := json.NewEncoder(rw)
		must(e.Encode(&angles))
		if attCounter%int(dbgPrintDiv) == 1 {
			log.Print("attitude served")
		}
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
	r   *bufio.Reader
	val [6]int32
}

func (imu *fieldReaderIMU) Run(ctx context.Context) {
	defer log.Println("imu.Run finished")
	rd := imu.r
	j := 0
	for {
		j++
		select {
		case <-ctx.Done():
			return
		default:
		}
		line, err := rd.ReadBytes('\n') // Discard all data until next newline
		if err != nil {
			log.Println(err)
			continue
		}
		values := bytes.Fields(line)
		if len(values) != len(imu.val) {
			// expected exactly 6 accelerometer+gyro values
			log.Println("expected 6 values, found ", len(values))
			continue
		}
		imu.Lock()
		for i := range values {
			v, err := strconv.ParseInt(string(values[i]), 10, 32)
			if err != nil {
				log.Println("error reading value")
				break // if one value is corrupted/bad it is unlikely others will be fine.
			}
			imu.val[i] = int32(v)
		}
		imu.Unlock()
		if j%dbgPrintDiv == 1 {
			log.Println("processed values ", imu.val)
		}
	}
}

var sensCount = 0

func (imu *fieldReaderIMU) Acceleration() (ax, ay, az int32) {
	imu.Lock()
	sensCount++
	defer imu.Unlock()
	if sensCount%dbgPrintDiv == 1 {
		log.Println("imu accel read: ", imu.val)
	}
	return imu.val[0], imu.val[1], imu.val[2]
}
func (imu *fieldReaderIMU) AngularVelocity() (gx, gy, gz int32) {
	imu.Lock()
	defer imu.Unlock()
	return imu.val[3], imu.val[4], imu.val[5]
}
