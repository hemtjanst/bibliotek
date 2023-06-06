package main

import (
	"context"
	"flag"
	"html/template"
	"lib.hemtjan.st/client"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"lib.hemtjan.st/server"
	"lib.hemtjan.st/transport/mqtt"
)

func main() {
	mCfg := mqtt.MustFlags(flag.String, flag.Bool)
	flgDelete := flag.String("delete", "", "Delete device (by topic)")
	flgHttpBind := flag.String("http.bind", ":8933", "Bind HTTP to host:port")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mq, err := mqtt.New(ctx, mCfg())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stop()
		for {
			ok, err := mq.Start()
			if err != nil {
				log.Printf("Error from MQTT: %s", err)
			}
			// If ok=true then the error is recoverable
			if !ok {
				break
			}
			log.Printf("Disconnected from MQTT, retrying in 3 seconds")
			time.Sleep(3 * time.Second)
		}
	}()

	srv := server.New(mq)

	l := sync.RWMutex{}
	devices := map[string]server.Device{}
	getDevices := func() (out []server.Device) {
		l.RLock()
		for _, v := range devices {
			out = append(out, v)
		}
		l.RUnlock()
		sort.Slice(out, func(i, j int) bool {
			return out[i].Id() < out[j].Id()
		})
		return
	}
	waitDelete := map[string][]chan struct{}{}

	ch := make(chan server.Update, 1)
	srv.SetUpdateChannel(ch)
	go func() {
		for {
			upd, open := <-ch
			if !open {
				return
			}
			dev := upd.Device
			switch upd.Type {
			case server.AddedDevice:
				if *flgDelete != "" && dev.Id() == *flgDelete {
					_ = client.DeleteDevice(dev.Info(), mq)
				}
				log.Printf("Device Created: %s (\"%s\" with serial %s model %s made by %s)",
					dev.Id(), dev.Name(), dev.SerialNumber(), dev.Model(), dev.Manufacturer())
				l.Lock()
				devices[dev.Id()] = dev
				l.Unlock()
				for _, ft := range dev.Features() {
					log.Printf("                * %s (min/max/step: %d/%d/%d", ft.Name(), ft.Min(), ft.Max(), ft.Step())
					go func(ft server.Feature) {
						ch, _ := ft.OnUpdate()
						for {
							d, open := <-ch
							if !open {
								log.Printf("Device: %s feature: %s disappeared", dev.Id(), ft.Name())
								return
							}
							log.Printf("Device: %s feature: %s is now %s", dev.Id(), ft.Name(), d)
						}
					}(ft)
				}
			case server.UpdatedDevice:
				if dev != nil {
					l.Lock()
					devices[dev.Id()] = dev
					l.Unlock()
				}
				for _, upd := range upd.Changes {
					log.Printf("[%s] %s changed \"%s\" -> \"%s\" (%+v)",
						dev.Id(),
						upd.Field,
						upd.Old,
						upd.New,
						upd.FeatureInfo,
					)
				}
			case server.RemovedDevice:
				if dev != nil {
					l.Lock()
					delete(devices, dev.Id())
					v, ok := waitDelete[dev.Id()]
					if ok {
						delete(waitDelete, dev.Id())
					}
					l.Unlock()
					if ok {
						for _, ch := range v {
							ch <- struct{}{}
						}
					}
				}
				log.Printf("Device Removed: %+v", dev)
			}

		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := srv.Start(ctx)
		if err != nil {
			stop()
			log.Print(err)
		}
	}()

	web := http.Server{Addr: *flgHttpBind}

	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		web.Handler = mux

		mux.HandleFunc("/", indexFunc(getDevices))
		mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
			qs := r.URL.Query()
			id := qs.Get("id")
			l.RLock()
			dev := devices[id]
			l.RUnlock()
			if dev != nil {
				ch := make(chan struct{}, 1)
				l.Lock()
				if _, ok := waitDelete[dev.Id()]; !ok {
					waitDelete[dev.Id()] = []chan struct{}{}
				}
				waitDelete[dev.Id()] = append(waitDelete[dev.Id()], ch)
				l.Unlock()
				_ = client.DeleteDevice(dev.Info(), mq)
				select {
				case <-ch:
				case <-time.After(2 * time.Second):
				}
			}
			w.Header().Add("Location", "/")
			w.WriteHeader(302)
		})
		err := web.ListenAndServe()
		if err != nil {
			stop()
			log.Print(err)
		}
	}()

	<-ctx.Done()
	shutctx, shutcancel := context.WithTimeout(context.Background(), 10*time.Second)
	_ = web.Shutdown(shutctx)
	wg.Wait()
	shutcancel()
}

func indexFunc(getDevices func() []server.Device) http.HandlerFunc {
	tplIndex, err := template.New("index").Parse(htmlIndex)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		devs := getDevices()
		err := tplIndex.Execute(w, devs)
		if err != nil {
			w.Write([]byte(err.Error()))
			log.Print(err)
		}
	}
}

const htmlIndex = `<!DOCTYPE html>
<html>
<head>
<body>

<table>
<tr>
<th>Action</th>
<th>ID</th>
<th>Name</th>
<th>Info</th>
</tr>
{{ range $d := . }}
<tr>
  <td><a href="/delete?id={{$d.Id}}">Delete</a></td>
  <td>{{$d.Id}}</td>
  <td>{{$d.Name}}</td>
  <td>{{$d.Manufacturer}} {{$d.Model}}</td>
</tr>

{{ end }}
</table>
`
