# tfohttpclient

This is a package for tfo-plugins intended to read the terraform resource spec of the terraform resource that it belongs to. All envs required are provided to the pod by the terraform-operator that provision the job/pod.

**Example usage**

```go
package main

import (
	"log"
	"tfohttpclient"
)

func main() {
	b, err := tfohttpclient.ResourceSpec()
	if err != nil {
		panic(err)
	}
	log.Print(string(b))
}
```
