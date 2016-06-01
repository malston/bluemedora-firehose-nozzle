// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package webserver

import (
    "fmt"
    
    "github.com/cloudfoundry/sonde-go/events"
    "github.com/cloudfoundry/gosteno"
)

//Resource represents cloud controller data
type Resource struct {
    Deployment      string
    Job             string
    Index           string
    IP              string
    ValueMetrics    map[string]float64
    CounterMetrics  map[string]float64
}

func createEnvelopeKey(envelope *events.Envelope) string {
	return fmt.Sprintf("%s | %s | %s | %s", envelope.GetDeployment(), envelope.GetJob(), envelope.GetIndex(), envelope.GetIp())
}

func addMetric(envelope *events.Envelope, valueMetricMap map[string]float64, counterMetricMap map[string]float64, logger *gosteno.Logger) {
    if envelope.GetEventType() == events.Envelope_ValueMetric {
        valueMetric := envelope.GetValueMetric()
		
		valueMetricMap[valueMetric.GetName()] = valueMetric.GetValue()
		logger.Debugf("Addeding Counter Event Name %s, Value %d", valueMetric.GetName(), valueMetric.GetValue())
    } else if envelope.GetEventType() == events.Envelope_CounterEvent {
        counterEvent := envelope.GetCounterEvent()
		
		counterMetricMap[counterEvent.GetName()] = float64(counterEvent.GetTotal())
		logger.Debugf("Addeding Counter Event Name %s, Value %d", counterEvent.GetName(), counterEvent.GetTotal())
    } else {
		logger.Errorf("Unkown event type %s", envelope.GetEventType())
	}
}