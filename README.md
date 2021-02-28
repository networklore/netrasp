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

	"github.com/networklore/netrasp/pkg/netrasp"
)


func main() {
	// Request a new client using the ios driver.
	client, err := netrasp.New("switch1", "ios", "my_user", "my_password123")
	if err != nil {
		log.Fatalf("unable to create client: %w", err)
	}

	// Connect to the device
	device, err := client.Connect(context.Background())
	defer device.Disconnect()
	if err != nil {
		log.Fatalf("unable to connect: %w", err)
	}

	// Run command against the device.
	output, err := device.Run(context.Background(), "show inventory")
	if err != nil {
		log.Fatalf("unable to run command: %w", err)
	}

	// Display the results.
	fmt.Println(output)
}
```

Use cases
---------

You can use Netrasp as a package as in the example above of combine it with
something like [Gornir](https://github.com/nornir-automation/gornir) to get
the same type of experience you'd have from using [Netmiko](https://github.com/ktbyers/netmiko)
and [Nornir](https://github.com/nornir-automation/nornir) in the Python world.

