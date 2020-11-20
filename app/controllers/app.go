package controllers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/revel/revel"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var pix = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{1, 1}})
var img image.Image = pix
var buffer = new(bytes.Buffer)
var foo = jpeg.Encode(buffer, img, nil)
var content_length = strconv.Itoa(len(buffer.Bytes()))

type Pixel string

func (p Pixel) Apply(req *revel.Request, resp *revel.Response) {
	resp.Out.Header().Add("Content-Type", "image/jpeg")
	resp.Out.Header().Add("Content-Length", content_length)
	resp.GetWriter().Write(buffer.Bytes())
}

type Hermes struct {
	*revel.Controller
}

func (c Hermes) Index() revel.Result {
	return c.Render()
}

func send_notification(title, recipient, token string) {
	http.PostForm("https://gotify.yigitcolakoglu.com/message?token="+token,
		url.Values{"title": {"E-mail Seen"}, "message": {fmt.Sprintf("To: %s \n Title: %s", recipient, title)}, "priority": {"5"}})
}

func (h Hermes) Read(title, recipient, user, provided_hash string) revel.Result {
	ip := strings.Split(h.Request.RemoteAddr, ":")[0]
	raw_userdata, _ := ioutil.ReadFile("storage/userdata")
	userdata := strings.Split(string(raw_userdata), "\n")
	key := ""
	gotify_token := ""
	for i := 0; i < len(userdata); i++ {
		if strings.Split(userdata[i], " ")[0] == user {
			key = strings.Split(userdata[i], " ")[1]
			gotify_token = strings.Split(userdata[i], " ")[2]
		}
	}
	if key == "" {
		return Pixel("No dice")
	}
	identity := title + recipient + user + key
	fmt.Printf(identity)
	hash_string := fmt.Sprintf("%x", sha256.Sum256([]byte(identity)))
	if hash_string != provided_hash {
		fmt.Printf("Hashes don't match!\n")
		return Pixel("You think you can fool me?")
	}
	if _, err := os.Stat("storage/ipdata/" + hash_string); os.IsNotExist(err) {
		os.Create("storage/ipdata/" + hash_string)
	}
	raw_ipdata, _ := ioutil.ReadFile("storage/ipdata/" + hash_string)
	ipdata := strings.Split(string(raw_ipdata), "\n")
	for i := 0; i < len(ipdata)-1; i++ {
		if ipdata[i] == ip {
			fmt.Printf("IP collision detected!\n")
			return Pixel("I see what you are trying to do!")
		}
	}

	f, _ := os.OpenFile("storage/ipdata/"+hash_string, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(fmt.Sprintf("%s\n", ip))
	f.Close()
	send_notification(title, recipient, gotify_token)
	return Pixel("Success")
}
