# Keypresslog

Sample usage:
```go
package main
 
import (
    "fmt"
    "strconv"
    
    "github.com/marekfilip/keypresslog"
)
 
func main() {
    var inputSign string
    devs := keypresslog.Find()
    
    for _, val := range devs {
        fmt.Println("Id->", val.GetId(), "Device->", val.GetName())
    }
    
    fmt.Println("Choose device")
    fmt.Scanln(&inputSign)
    intId, _ := strconv.Atoi(inputSign)
    
    in, err := devs[intId].Read()
    
    if err != nil {
        fmt.Println(err)
        return
    }
    
    for i := range in {
        //listen only key stroke event
        if i.Type == keypresslog.EV_KEY {
            if i.Value == 1 {
                fmt.Println("Click", i.ToString())
            }
        }
    }
}
 ```
