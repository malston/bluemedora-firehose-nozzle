/**Copyright Blue Medora Inc. 2016**/

package bluemedorafirehosenozzle

import (
    "crypto/tls"
    "fmt"
    "time"
    //REST server
    
    "github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
    "github.com/cloudfoundry/noaa/consumer"
    "github.com/cloudfoundry/sonde-go/events"
    "github.com/cloudfoundry/gosteno"
    "github.com/cloudfoundry-incubator/uaago"
)

//BlueMedoraFirehoseNozzle consuems data from fire hose and exposes it via REST
type BlueMedoraFirehoseNozzle struct {
    config      *nozzleconfiguration.NozzleConfiguration
    errs        <-chan error
    messages    <-chan *events.Envelope
    logger      *gosteno.Logger
}

//New BlueMedoraFirhoseNozzle
func New(config *nozzleconfiguration.NozzleConfiguration, logger *gosteno.Logger) *BlueMedoraFirehoseNozzle {
    return &BlueMedoraFirehoseNozzle {
        config:     config,
        logger:     logger,
    }
}

//Start starts consuming events from firehose
func (nozzle *BlueMedoraFirehoseNozzle) Start() {
    nozzle.logger.Info("Starting Blue Medora Firehose Nozzle")
    
    var authToken string
    if !nozzle.config.DisableAccessControl {
        authToken = nozzle.fetchUAAAuthToken()
    }
    
    nozzle.collectFromFirehose(authToken)
    err := nozzle.logMessages()
    if err != nil {
        nozzle.logger.Errorf("Error while reading from firehose: %s", err)
    }
    
    nozzle.logger.Info("Closing Blue Medora Firehose Nozzle")
}

func (nozzle *BlueMedoraFirehoseNozzle) fetchUAAAuthToken() string {
    nozzle.logger.Debug("Fetching UAA authenticaiton token")
    
    UAAClient, err := uaago.NewClient(nozzle.config.UAAURL)
    if err != nil {
        nozzle.logger.Fatalf("Error creating UAA client: %s", err.Error())
    }   
    
    var token string
    token, err = UAAClient.GetAuthToken(nozzle.config.UAAUsername, nozzle.config.UAAPassword, nozzle.config.InsecureSSLSkipVerify)
    if err != nil {
        nozzle.logger.Fatalf("Failed to get oauth token: %s.", err.Error())
    }
    
    nozzle.logger.Debug(fmt.Sprintf("Successfully fetched UAA authentication token <%s>", token))
    return token
}

func (nozzle *BlueMedoraFirehoseNozzle) collectFromFirehose(authToken string) {
    consumer := consumer.New(nozzle.config.TrafficControllerURL, &tls.Config{InsecureSkipVerify: nozzle.config.InsecureSSLSkipVerify}, nil)
    
    debugPrinter := &BMDebugPrinter{nozzle.logger}
    consumer.SetDebugPrinter(debugPrinter)
    consumer.SetIdleTimeout(time.Duration(nozzle.config.IdleTimeoutSeconds) * time.Second)
    nozzle.messages, nozzle.errs = consumer.Firehose("bluemedora-nozzle", authToken)
}

func (nozzle *BlueMedoraFirehoseNozzle) logMessages() error {
    for {
        select {
            case envelope := <-nozzle.messages:
                nozzle.logEnvelope(envelope)
            case err := <-nozzle.errs:
                return err
        }
    }
}

func (nozzle *BlueMedoraFirehoseNozzle) logEnvelope(envelope *events.Envelope) {
    nozzle.logger.Debug(fmt.Sprintf("Received Envelope with Origin <%v>, EventType <%v>, Deployment <%v>, Job <%v>, Index <%v>, IP <%v>", envelope.Origin, envelope.EventType, envelope.Deployment, envelope.Job, envelope.Index, envelope.Ip))
}
