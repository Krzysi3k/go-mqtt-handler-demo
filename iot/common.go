package iot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"mqtt-handler-go/addr"
)

type MacAddress [6]byte

type MagicPacket struct {
	header  [6]byte
	payload [16]MacAddress
}

type GrafanaAnnotation struct {
	dashboardId int64
	panelId     int64
	msg         string
}

func logEvent(topic, msg string) {
	log.Printf("topic: %v, payload: %v\n", topic, msg)
}

func logError(err error, msg string) {
	if err != nil {
		log.Printf("ERROR: %v, ", msg)
		log.Println(err)
	}
}

func httpPost(url string, contentType string, body []byte) (string, bool) {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	logError(err, "http request faled")
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	logError(err, "cannot read body response")
	if err != nil {
		return "", false
	}
	return string(content), true
}

func httpGet(url string) (string, bool) {
	resp, err := http.Get(url)
	logError(err, "http request faled")
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	logError(err, "cannot read body response")
	if err != nil {
		return "", false
	}
	return string(content), true
}

func turnAllDevices() {
	devices := []string{"dekoder", "tv", "creative"}
	for _, i := range devices {
		httpGet(fmt.Sprintf("http://%v/%v", addr.NodeMcuIP, i))
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(10 * time.Second)
	httpGet(fmt.Sprintf("http://%v/dekoder-back", addr.NodeMcuIP))
}

func powerOffLamps() {
	lamps := []string{addr.LampaSalonIP, addr.LampaSwieczkiIP}
	for _, lamp := range lamps {
		go httpGet(fmt.Sprintf("http://%v/cm?cmnd=Power%%20Off", lamp))
	}
}

func postTelegramMsg(msg string, includeSticker bool) {
	API_KEY := os.Getenv("TELEGRAM_KEY")
	CHAT_ID := os.Getenv("TELEGRAM_CHAT_ID")
	STICKER_ID := os.Getenv("SPIDER_STICKER")
	if includeSticker {
		data := url.Values{
			"chat_id": {CHAT_ID},
			"sticker": {STICKER_ID},
		}
		_, err := http.PostForm("https://api.telegram.org/bot"+API_KEY+"/sendSticker", data)
		logError(err, "cannot POST request")
		time.Sleep(time.Second * 1)
	}
	currTime := time.Now().Format("2006-01-02 15:04:05")
	msgPayload := currTime + "\n" + msg
	data := `{"text":"` + msgPayload + `","chat_id":"` + CHAT_ID + `"}`
	httpPost(addr.ApiUrl+API_KEY+"/sendMessage", "application/json", []byte(data))
}

func postGrafanaReleaseAnnotation(ga GrafanaAnnotation) {
	client := &http.Client{}
	data := fmt.Sprintf("{\"dashboardId\":%v, \"panelId\":%v, \"text\": \"%v\"}", ga.dashboardId, ga.panelId, ga.msg)
	url := fmt.Sprintf("http://%v:3000/api/annotations", addr.WyseIP)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth("admin", os.Getenv("GRAFANA_PASS"))
	resp, _ := client.Do(req)
	resp.Body.Close()
}

func dockerInfo(event Event) {
	if resp, ok := httpGet(fmt.Sprintf("http://%v:5001/docker-info?items=%v", addr.WyseIP, event.msg)); ok {
		event.mqttClient.Publish("tele/rpi/dockerinfo", 0, false, resp)
	}
}

func redisInfo(event Event) {
	if resp, ok := httpGet(fmt.Sprintf("http://%v:5001/redis-info", addr.WyseIP)); ok {
		event.mqttClient.Publish("tele/rpi/redis-info", 0, false, resp)
	}
}

func sendPing(event Event, device, publishTopic string) {
	cnt := "-c"
	substring := " 0% packet loss"
	if runtime.GOOS == "windows" {
		cnt = "-n"
		substring = "Lost = 0"
	}
	out, _ := exec.Command("ping", device, cnt, "1").Output()
	outstr := string(out)
	if strings.Contains(outstr, "host unreachable") {
		event.mqttClient.Publish(publishTopic, 0, true, "OFF")
	} else if strings.Contains(outstr, substring) {
		event.mqttClient.Publish(publishTopic, 0, true, "ON")
	} else {
		event.mqttClient.Publish(publishTopic, 0, true, "OFF")
	}
}

func pingStatus(event Event) {
	switch event.topic {
	case "cmnd/pingstatus/Ryzen":
		sendPing(event, addr.RyzenIP, "tele/pingstatus/Ryzen")
	case "cmnd/pingstatus/Tv":
		sendPing(event, addr.TvIP, "tele/pingstatus/Tv")
	case "cmnd/pingstatus/Hypervisor":
		sendPing(event, addr.HypervisorIP, "tele/pingstatus/Hypervisor")
	}
}

func wakeOnLan(event Event) {
	if event.msg == "ON" {
		if ok, _ := regexp.Match("^([0-9a-fA-F]{2}[:-]){5}([0-9a-fA-F]{2})$", []byte(addr.HypervisorMAC)); ok {
			var packet MagicPacket
			var macAddr MacAddress
			hwAddr, err := net.ParseMAC(addr.HypervisorMAC)
			if err != nil {
				logError(err, "cannot parse mac addr")
				return
			}
			for idx := range macAddr {
				macAddr[idx] = hwAddr[idx]
			}
			for idx := range packet.header {
				packet.header[idx] = 0xFF
			}
			for idx := range packet.payload {
				packet.payload[idx] = macAddr
			}
			var buf bytes.Buffer
			binary.Write(&buf, binary.BigEndian, packet)
			udpAddr, err := net.ResolveUDPAddr("udp", "ipaddr:9")
			if err != nil {
				logError(err, "unable to get UDP address")
				return
			}
			connection, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				logError(err, "Unable to dial UDP address")
				return
			}
			defer connection.Close()
			bytesWritten, err := connection.Write(buf.Bytes())
			if err != nil {
				logError(err, "Unable to write packet to connection")
				return
			} else if bytesWritten != 102 {
				logEvent(event.topic, fmt.Sprintf("Warning: %v bytes written 102 expected!", bytesWritten))
			}
		}
	}
}
