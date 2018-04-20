//
// Copyright (c) 2017
// Cavium
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package distro

import (
	"github.com/edgexfoundry/edgex-go/core/domain/models"
	export "github.com/edgexfoundry/edgex-go/export"
)

// Sender - Send interface
type Sender interface {
	Send(event *models.Event, data []byte)
}

// Formater - Format interface
type Formater interface {
	Format(event *models.Event) []byte
}

// Transformer - Transform interface
type Transformer interface {
	Transform(data []byte) []byte
}

// Filter - Filter interface
type Filterer interface {
	Filter(event *models.Event) (bool, *models.Event)
}

// RegistrationInfo - registration info
type registrationInfo struct {
	registration export.Registration
	format       Formater
	compression  Transformer
	encrypt      Transformer
	sender       Sender
	filter       []Filterer

	chRegistration chan *export.Registration
	chEvent        chan *models.Event

	deleteMe bool
}
