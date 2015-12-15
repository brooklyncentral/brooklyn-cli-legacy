package command
import (
    "fmt"
    "github.com/brooklyncentral/brooklyn-cli/net"
    "net/url"
    "errors"
)

type Verifier struct {
    network *net.Network
}

func (v Verifier) SetNetwork(nw *net.Network) {
    v.network = nw
}

func (v Verifier) Verify() error {
    fmt.Println("Verifying....")
    return VerifyLoginURL(v.network)
}

func VerifyLoginURL(network *net.Network) error {
    url, err := url.Parse(network.BrooklynUrl)
    if err != nil {
        return err
    }
    if url.Scheme != "http" && url.Scheme != "https" {
        return errors.New("Brooklyn URL must have a scheme of \"http\" or \"https\"")
    }
    if url.Host == "" {
        return errors.New("Brooklyn URL must have a valid host")
    }
    return nil
}