// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package main

import (
    "context"
    "flag"
    "fmt"
    "net"
    "net/http"
    "time"
)

const (
    initialWaitSec = 2
    imageWaitSec = 2
    dnsResolverTimeoutMs = 5000 // Timeout (ms) for the DNS resolver
)

func main() {
    serverName, serverPort, dnsResolver := getArgs()

    ch := make(chan byte, 1)

    go resolveLoop(serverName, serverPort, dnsResolver, ch)
    <- ch

}

func resolveLoop(serverName, serverPort, dnsResolver string, ch chan byte) {
    fmt.Printf("Entered resolveLoop\n")
    counter := 0
    for {
        time.Sleep(1 * time.Second)
        counter++

        fmt.Printf("%d: ", counter)
        ip, err := lookupHost(serverName, dnsResolver)
        if err != nil {
            continue
        }

        serverURL := fmt.Sprintf("http://%s:%s?image=%d", ip, serverPort, counter)
        fmt.Printf("%s\n", serverURL)
        client := &http.Client{}
        resp, err1 := client.Get(serverURL)
        if err1 != nil {
           fmt.Printf("   Error fetchingurl: %s\n", err1.Error())
           continue
        }
        if resp != nil && resp.Body != nil {
            defer resp.Body.Close()
        }

    }
    ch <- 1
}

func lookupHost(serverName, dnsResolver string)  (string, error){
    r := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{
                Timeout: time.Millisecond * time.Duration(10000),
            }
            return d.DialContext(ctx, network, dnsResolver)
        },
    }

    ip, err := r.LookupHost(context.Background(), serverName)
    if err != nil {
       fmt.Printf("lookup failed. Error: %s\n", err.Error())
       return "", err
    }

    //fmt.Printf("First IP address: %s\n", ip[0])
    return ip[0], nil
}

func getArgs()  (string, string, string) {
    var serverName, serverPort, dnsServer string

    flag.StringVar(&serverName, "server", "webserver.demo.com", "HTTP server name")
    flag.StringVar(&serverPort, "port", "30000", "HTTP server port")
    flag.StringVar(&dnsServer, "dns", "0.0.0.0:53", "Custom DNS server with port")
    flag.Parse()

    fmt.Printf("Will connect to %s:%s using DNS %s\n", serverName, serverPort, dnsServer)

    return serverName, serverPort, dnsServer
}
