# tfohttpclient

This is a package for tfo-plugins intended to read the terraform resource spec of the terraform resource that it belongs to. All envs required are provided to the pod by the terraform-operator that provision the job/pod.

**Example usage**

```bash
export TFO_NAMESPACE=default
export TFO_RESOURCE=hello-tfo
```

```go
package main

import (
	"log"
	"github.com/galleybytes/tfohttpclient"
)

func main() {
	b, err := tfohttpclient.Resource()
	if err != nil {
		panic(err)
	}
	log.Print(string(b))
}
```

Or unmarshal into the api

```go
import "github.com/isaaguilar/terraform-operator/pkg/apis/tf/v1alpha2"
```

```go
	b, _ := tfohttpclient.Resource()
	//Extract the spec
	var tf v1alpha2.Terraform
	err = json.Unmarshal(b, &tf)
	// handle err
	log.Printf("%+v", tf.Spec)
```
