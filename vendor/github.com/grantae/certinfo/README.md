# Certificate Information for Go

A golang tool for printing x509 TLS certificates in a format similar to OpenSSL.

## Installation

``` bash
go get github.com/grantae/certinfo
```

## Usage

### Print a certificate from a website

``` go
package main

import (
  "crypto/tls"
  "fmt"
  "github.com/grantae/certinfo"
  "log"
)

func main() {
  // Connect to google.com
  cfg := tls.Config{}
  conn, err := tls.Dial("tcp", "google.com:443", &cfg)
  if err != nil {
    log.Fatalln("TLS connection failed: " + err.Error())
  }
  // Grab the last certificate in the chain
  certChain := conn.ConnectionState().PeerCertificates
  cert := certChain[len(certChain)-1]

  // Print the certificate
  result, err := certinfo.CertificateText(cert)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Print(result)
}
```

### Print a PEM-encoded certificate from a file

``` go
package main

import (
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "github.com/grantae/certinfo"
  "io/ioutil"
  "log"
)

func main() {
  // Read and parse the PEM certificate file
  pemData, err := ioutil.ReadFile("cert.pem")
  if err != nil {
    log.Fatal(err)
  }
  block, rest := pem.Decode([]byte(pemData))
  if block == nil || len(rest) > 0 {
    log.Fatal("Certificate decoding error")
  }
  cert, err := x509.ParseCertificate(block.Bytes)
  if err != nil {
    log.Fatal(err)
  }

  // Print the certificate
  result, err := certinfo.CertificateText(cert)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Print(result)
}
```

## Testing

``` bash
go test github.com/grantae/certinfo
```

This compares several PEM-encoded certificates with their expected outputs.

## License

MIT -- see `LICENSE` for more information.

