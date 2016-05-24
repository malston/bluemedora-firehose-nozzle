/**Copyright Blue Medora Inc. 2016**/

package bluemedorafirehosenozzle

import (
    "fmt"
    
    "github.com/cloudfoundry/gosteno"
)

type BMDebugPrinter struct {
	logger *gosteno.Logger
}

type bmDebugPrinterMessage struct {
	Title, Body string
}

func (p *BMDebugPrinter) Print(title, body string) {
	p.logger.Debug(fmt.Sprintf("BMPrinter message %s: <%s>", title, body))
}