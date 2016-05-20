/**Copyright Blue Medora Inc. 2016**/

package bluemedorafirehosenozzle

import (
    "crypto/tls"
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
func (nozzle *BlueMedoraFirehoseNozzle) Start() error {
    nozzle.logger.Info("Starting Blue Medora Firehose Nozzle")
    
    var authToken string
    if !nozzle.config.disableAccessControl {
        authToken = nozzle.FetchUAAAuthToken()
    }
    
    nozzle.logger.Info("Closing Blue Medora Firehose Nozzle")
}

func (nozzle *BlueMedoraFirehoseNozzle) FetchUAAAuthToken() string {
    nozzle.logger.Debug("Fetching UAA authenticaiton token")
    
    UAAClient, err := uaago.NewClient(nozzle.config.UAAURL)
    if err != nil {
        nozzle.logger.Fatalf("Error creating UAA client: %s", err.Error())
    }   
    
    var token string
    token, err = UAAClient.GetAuthToken(nozzle.config.UAAUsername, nozzle.config.UAAPassword, nozzle.config.useSSL)
    if err != nil{
        nozzle.logger.Fatalf("Failed to get oauth toke: %s. Verify username and password.", err.Error())
    }
    
    nozzle.logger.Debug("Successfully fetched UAA authentication token")
    return token
}

func (nozzle *BlueMedoraFirehoseNozzle) CollectFromFirehose(authToken string) {
    consumer := consumer.New(nozzle.config.trafficControllerURL, &tls.Config{InsecureSkipVerify: nozzle.config.useSSL}, nil)
    consumer.SetIdleTimeout(time.Duration(nozzle.config.IdleTimeout) * time.Second)
    nozzle.messages, nozzle.errs = consumer.Firehose("bluemedora-nozzle", authToken)
}
