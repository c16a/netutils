package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var HttpUrl string
var HttpMethod string
var HttpTimeout int
var HttpVersion int
var HttpProxy string
var HttpHeaders map[string]string

func init() {
	httpCommand.Flags().StringVarP(&HttpUrl, "url", "u", "", "set the URL")
	httpCommand.Flags().StringVarP(&HttpMethod, "method", "m", "", "set the method")
	httpCommand.Flags().IntVarP(&HttpTimeout, "timeout", "t", 5000, "set the timeout in milliseconds")
	httpCommand.Flags().StringToStringVarP(&HttpHeaders, "headers", "", nil, "set the headers")
	httpCommand.Flags().IntVarP(&HttpVersion, "version", "v", 1, "set the HTTP version")
	httpCommand.Flags().StringVarP(&HttpProxy, "proxy", "p", "", "set the HTTP proxy")
	rootCmd.AddCommand(httpCommand)
}

var httpCommand = &cobra.Command{
	Use:   "http",
	Short: "Perform HTTP operations",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{
			Timeout: time.Duration(HttpTimeout) * time.Millisecond,
		}

		caCertPool, err := x509.SystemCertPool()
		if err != nil {
			return err
		}

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}

		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy: func(r *http.Request) (*url.URL, error) {
				u, err := url.Parse(HttpProxy)
				if err != nil {
					return nil, err
				}
				return u, nil
			},
		}

		if HttpVersion == 2 {
			t, err := http2.ConfigureTransports(transport)
			if err != nil {
				return err
			}
			client.Transport = t
		} else {
			client.Transport = transport
		}

		switch HttpVersion {
		case 1:
			client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
		case 2:
			client.Transport = &http2.Transport{TLSClientConfig: tlsConfig}
		}

		request, err := http.NewRequest(HttpMethod, HttpUrl, nil)
		if err != nil {
			return err
		}

		response, err := client.Do(request)
		if err != nil {
			return err
		}

		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(bodyBytes))
		return nil
	},
}
