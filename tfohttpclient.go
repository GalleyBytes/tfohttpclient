package tfohttpclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func envOrDefault(s, defaultval string) string {
	env := os.Getenv(s)
	if env != "" {
		return env
	}
	return defaultval
}

func Resource() ([]byte, error) {
	var (
		host      = envOrDefault("KUBERNETES_SERVICE_HOST", "kubernetes.default.svc")
		certFile  = envOrDefault("CERTFILE", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
		tokenFile = envOrDefault("TOKENFILE", "/var/run/secrets/kubernetes.io/serviceaccount/token")
		group     = envOrDefault("TFO_GROUP", "tf.isaaguilar.com/v1alpha2")
		namespace = envOrDefault("TFO_NAMESPACE", "")
		resource  = envOrDefault("TFO_RESOURCE", "")
	)
	url := fmt.Sprintf("https://%s/apis/%s/namespaces/%s/terraforms/%s", host, group, namespace, resource)

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(certFile)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to append %q to RootCAs: %v", certFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return []byte{}, fmt.Errorf("no certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return []byte{}, err
	}

	authToken, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to load token from tokenFile '%s': %s", tokenFile, err.Error())
	}
	req.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", string(authToken))},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("errored when sending request to the server: %s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
