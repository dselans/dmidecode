dmidecode
=========

`dmidecode` is a Go library that parses the output of the `dmidecode` command
and makes it accessible via a simple map data structure.

In addition, it exposes a helper method for quickly looking up specific
records.

## Usage

```go
import (
    dmidecode "github.com/dselans/dmidecode"
)

dmi := dmidecode.NewDMI()

if err := dmi.Run(); err != nil {
    fmt.Printf("Unable to get dmidecode information. Error: %v\n", err)
}

sysInfo, err := dmi.Search("System Information")
if err != nil {
    fmt.Println("Unable to find System Information")
}

for k, v := range sysInfo {
    fmt.Printf("Key: %v Value: %v\n", key, value)
}
```
