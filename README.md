netrasp
=======

Netrasp is a package that communicates to network devices over SSH. It takes
care of handling the pty terminal of network devices giving you an API with
common actions such as executing commands and configuring devices.

Warning
-------
Netrasp is in pre release mode so some parts of the API might change before
the initial version is released.

Example
-------

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/networklore/netrasp/pkg/netrasp"
)

func main() {
	device, err := netrasp.New("switch1",
		netrasp.WithUsernamePassword("my_user", "my_password123"),
		netrasp.WithDriver("ios"),
	)
	if err != nil {
		log.Fatalf("unable to create client: %v", err)
	}

	ctx, cancelOpen := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancelOpen()
	err = device.Dial(ctx)
	if err != nil {
		fmt.Printf("unable to connect: %v\n", err)

		return
	}
	defer device.Close(context.Background())

	ctx, cancelRun := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancelRun()
	output, err := device.Run(ctx, "show running")
	if err != nil {
		fmt.Printf("unable to run command: %v\n", err)

		return
	}
	fmt.Println(output)
}
```

Network Device Support
----------------------
The initial release of Netrasp comes with support for the following platforms:

* Cisco IOS: netrasp.WithDriver("ios")
* Cisco NXOS: netrasp.WithDriver("nxos")
* Cisco ASA: netrasp.WithDriver("asa")
* Nokia SR OS: netrasp.WithDriver("sros")

Use cases
---------

You can use Netrasp as a package as in the example above of combine it with
something like [Gornir](https://github.com/nornir-automation/gornir) to get
the same type of experience you'd have from using [Netmiko](https://github.com/ktbyers/netmiko)
and [Nornir](https://github.com/nornir-automation/nornir) in the Python world.

Blog Posts
----------

* [Hello Netrasp](https://networklore.com/hello-netrasp/)

Credits
-------

Netrasp was created by [Patrick Ogenstad](https://github.com/ogenstad). Special
thanks to [David Barroso](https://github.com/dbarrosop) for providing feedback
and recommendations of the structure and code.
