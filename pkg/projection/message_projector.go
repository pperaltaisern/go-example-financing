package projection

import (
	"errors"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pperaltaisern/financing/pkg/financing"
)

// MessageProjector is a layer on top of an EventProjector responsible of deserializing and routing events
type MessageProjector struct {
	eventProjector EventProjector
	eventMarshaler cqrs.CommandEventMarshaler
	logErr         func(*message.Message, error)
}

func NewMessageProjector(ep EventProjector, em cqrs.CommandEventMarshaler, logErr func(*message.Message, error)) *MessageProjector {
	if logErr == nil {
		logErr = func(m *message.Message, e error) {}
	}
	return &MessageProjector{
		eventProjector: ep,
		eventMarshaler: em,
		logErr:         logErr,
	}
}

func (p *MessageProjector) ProjectMessage(m *message.Message) {
	err := p.projectMessage(m)
	if err != nil {
		p.logErr(m, errUnknownMessage)
	}
	m.Ack()
}

func (p *MessageProjector) projectMessage(m *message.Message) error {
	switch p.eventMarshaler.NameFromMessage(m) {
	case "InvestorCreatedEvent":
		e := &financing.InvestorCreatedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvestorCreatedEvent(e); err != nil {
			return err
		}

	case "InvestorFundsAddedEvent":
		e := &financing.InvestorFundsAddedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvestorFundsAddedEvent(e); err != nil {
			return err
		}

	case "BidOnInvoicePlacedEvent":
		e := &financing.BidOnInvoicePlacedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectBidOnInvoicePlacedEvent(e); err != nil {
			return err
		}

	case "InvestorFundsReleasedEvent":
		e := &financing.InvestorFundsReleasedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvestorFundsReleasedEvent(e); err != nil {
			return err
		}

	case "InvestorFundsCommittedEvent":
		e := &financing.InvestorFundsCommittedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvestorFundsCommittedEvent(e); err != nil {
			return err
		}

	case "InvoiceCreatedEvent":
		e := &financing.InvoiceCreatedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvoiceCreatedEvent(e); err != nil {
			return err
		}

	case "InvoiceFinancedEvent":
		e := &financing.InvoiceFinancedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvoiceFinancedEvent(e); err != nil {
			return err
		}

	case "InvoiceReversedEvent":
		e := &financing.InvoiceReversedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvoiceReversedEvent(e); err != nil {
			return err
		}

	case "InvoiceApprovedEvent":
		e := &financing.InvoiceApprovedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectInvoiceApprovedEvent(e); err != nil {
			return err
		}

	case "IssuerCreatedEvent":
		e := &financing.IssuerCreatedEvent{}
		if err := p.eventMarshaler.Unmarshal(m, e); err != nil {
			return err
		}
		if err := p.eventProjector.ProjectIssuerCreatedEvent(e); err != nil {
			return err
		}

	default:
		return errUnknownMessage
	}
	return nil
}

var errUnknownMessage = errors.New("unknown message")
