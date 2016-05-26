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
    "github.com/BlueMedora/bluemedora-firehose-nozzle/webserver"
)

//BlueMedoraFirehoseNozzle consuems data from fire hose and exposes it via REST
type BlueMedoraFirehoseNozzle struct {
    config      *nozzleconfiguration.NozzleConfiguration
    errs        <-chan error
    messages    <-chan *events.Envelope
    serverErrs  <-chan error
    logger      *gosteno.Logger
    server      *webserver.WebServer
}

//New BlueMedoraFirhoseNozzle
func New(config *nozzleconfiguration.NozzleConfiguration, server *webserver.WebServer, logger *gosteno.Logger) *BlueMedoraFirehoseNozzle {
    return &BlueMedoraFirehoseNozzle {
        config:     config,
        logger:     logger,
        server:     server,
    }
}

//Start starts consuming events from firehose
func (nozzle *BlueMedoraFirehoseNozzle) Start() error {
    nozzle.logger.Info("Starting Blue Medora Firehose Nozzle")
    
    var authToken string
    if !nozzle.config.DisableAccessControl {
        authToken = nozzle.fetchUAAAuthToken()
    }
    
    nozzle.serverErrs = nozzle.server.Start(webserver.DefaultKeyLocation, webserver.DefaultCertLocation)
    
    nozzle.collectFromFirehose(authToken)
    err := nozzle.processMessages()
    
    nozzle.logger.Info("Closing Blue Medora Firehose Nozzle")
    return err
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

//Method blocks until error occurs
func (nozzle *BlueMedoraFirehoseNozzle) processMessages() error {
    flushTicker := time.NewTicker(time.Duration(nozzle.config.MetricCacheDurationSeconds) * time.Second)
    for {
        select {
            case <-flushTicker.C:
                nozzle.flushMetricCaches()
            case envelope := <-nozzle.messages:
                nozzle.cacheEnvelope(envelope)
            case err := <-nozzle.serverErrs:
                if err != nil {
                    nozzle.logger.Errorf("Error while running webserver: %s", err)
                    return err
                }
            case err := <-nozzle.errs:
                if err != nil {
                    nozzle.logger.Errorf("Error while reading from firehose: %s", err)
                    return err
                }
        }
    }
}

func (nozzle *BlueMedoraFirehoseNozzle) cacheEnvelope(envelope *events.Envelope) {
    nozzle.logger.Debug("Cache envelope in server")
}

func (nozzle *BlueMedoraFirehoseNozzle) flushMetricCaches() {
    nozzle.logger.Debug("Flushing metric caches")
}
