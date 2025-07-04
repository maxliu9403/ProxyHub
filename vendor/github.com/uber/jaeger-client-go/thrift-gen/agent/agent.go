// Code generated by Thrift Compiler (0.14.1). DO NOT EDIT.

package agent

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/uber/jaeger-client-go/thrift"
	"github.com/uber/jaeger-client-go/thrift-gen/jaeger"
	"github.com/uber/jaeger-client-go/thrift-gen/zipkincore"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal

var _ = jaeger.GoUnusedProtection__
var _ = zipkincore.GoUnusedProtection__

type Agent interface {
	// Parameters:
	//  - Spans
	EmitZipkinBatch(ctx context.Context, spans []*zipkincore.Span) (_err error)
	// Parameters:
	//  - Batch
	EmitBatch(ctx context.Context, batch *jaeger.Batch) (_err error)
}

type AgentClient struct {
	c    thrift.TClient
	meta thrift.ResponseMeta
}

func NewAgentClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *AgentClient {
	return &AgentClient{
		c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
	}
}

func NewAgentClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *AgentClient {
	return &AgentClient{
		c: thrift.NewTStandardClient(iprot, oprot),
	}
}

func NewAgentClient(c thrift.TClient) *AgentClient {
	return &AgentClient{
		c: c,
	}
}

func (p *AgentClient) Client_() thrift.TClient {
	return p.c
}

func (p *AgentClient) LastResponseMeta_() thrift.ResponseMeta {
	return p.meta
}

func (p *AgentClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
	p.meta = meta
}

// Parameters:
//   - Spans
func (p *AgentClient) EmitZipkinBatch(ctx context.Context, spans []*zipkincore.Span) (_err error) {
	var _args0 AgentEmitZipkinBatchArgs
	_args0.Spans = spans
	p.SetLastResponseMeta_(thrift.ResponseMeta{})
	if _, err := p.Client_().Call(ctx, "emitZipkinBatch", &_args0, nil); err != nil {
		return err
	}
	return nil
}

// Parameters:
//   - Batch
func (p *AgentClient) EmitBatch(ctx context.Context, batch *jaeger.Batch) (_err error) {
	var _args1 AgentEmitBatchArgs
	_args1.Batch = batch
	p.SetLastResponseMeta_(thrift.ResponseMeta{})
	if _, err := p.Client_().Call(ctx, "emitBatch", &_args1, nil); err != nil {
		return err
	}
	return nil
}

type AgentProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      Agent
}

func (p *AgentProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *AgentProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *AgentProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewAgentProcessor(handler Agent) *AgentProcessor {

	self2 := &AgentProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self2.processorMap["emitZipkinBatch"] = &agentProcessorEmitZipkinBatch{handler: handler}
	self2.processorMap["emitBatch"] = &agentProcessorEmitBatch{handler: handler}
	return self2
}

func (p *AgentProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
	if err2 != nil {
		return false, thrift.WrapTException(err2)
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(ctx, thrift.STRUCT)
	iprot.ReadMessageEnd(ctx)
	x3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
	x3.Write(ctx, oprot)
	oprot.WriteMessageEnd(ctx)
	oprot.Flush(ctx)
	return false, x3

}

type agentProcessorEmitZipkinBatch struct {
	handler Agent
}

func (p *agentProcessorEmitZipkinBatch) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := AgentEmitZipkinBatchArgs{}
	var err2 error
	if err2 = args.Read(ctx, iprot); err2 != nil {
		iprot.ReadMessageEnd(ctx)
		return false, thrift.WrapTException(err2)
	}
	iprot.ReadMessageEnd(ctx)

	tickerCancel := func() {}
	_ = tickerCancel

	if err2 = p.handler.EmitZipkinBatch(ctx, args.Spans); err2 != nil {
		tickerCancel()
		return true, thrift.WrapTException(err2)
	}
	tickerCancel()
	return true, nil
}

type agentProcessorEmitBatch struct {
	handler Agent
}

func (p *agentProcessorEmitBatch) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := AgentEmitBatchArgs{}
	var err2 error
	if err2 = args.Read(ctx, iprot); err2 != nil {
		iprot.ReadMessageEnd(ctx)
		return false, thrift.WrapTException(err2)
	}
	iprot.ReadMessageEnd(ctx)

	tickerCancel := func() {}
	_ = tickerCancel

	if err2 = p.handler.EmitBatch(ctx, args.Batch); err2 != nil {
		tickerCancel()
		return true, thrift.WrapTException(err2)
	}
	tickerCancel()
	return true, nil
}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//   - Spans
type AgentEmitZipkinBatchArgs struct {
	Spans []*zipkincore.Span `thrift:"spans,1" db:"spans" json:"spans"`
}

func NewAgentEmitZipkinBatchArgs() *AgentEmitZipkinBatchArgs {
	return &AgentEmitZipkinBatchArgs{}
}

func (p *AgentEmitZipkinBatchArgs) GetSpans() []*zipkincore.Span {
	return p.Spans
}
func (p *AgentEmitZipkinBatchArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *AgentEmitZipkinBatchArgs) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*zipkincore.Span, 0, size)
	p.Spans = tSlice
	for i := 0; i < size; i++ {
		_elem4 := &zipkincore.Span{}
		if err := _elem4.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem4), err)
		}
		p.Spans = append(p.Spans, _elem4)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *AgentEmitZipkinBatchArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "emitZipkinBatch_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *AgentEmitZipkinBatchArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "spans", thrift.LIST, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:spans: ", p), err)
	}
	if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Spans)); err != nil {
		return thrift.PrependError("error writing list begin: ", err)
	}
	for _, v := range p.Spans {
		if err := v.Write(ctx, oprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
		}
	}
	if err := oprot.WriteListEnd(ctx); err != nil {
		return thrift.PrependError("error writing list end: ", err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:spans: ", p), err)
	}
	return err
}

func (p *AgentEmitZipkinBatchArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AgentEmitZipkinBatchArgs(%+v)", *p)
}

// Attributes:
//   - Batch
type AgentEmitBatchArgs struct {
	Batch *jaeger.Batch `thrift:"batch,1" db:"batch" json:"batch"`
}

func NewAgentEmitBatchArgs() *AgentEmitBatchArgs {
	return &AgentEmitBatchArgs{}
}

var AgentEmitBatchArgs_Batch_DEFAULT *jaeger.Batch

func (p *AgentEmitBatchArgs) GetBatch() *jaeger.Batch {
	if !p.IsSetBatch() {
		return AgentEmitBatchArgs_Batch_DEFAULT
	}
	return p.Batch
}
func (p *AgentEmitBatchArgs) IsSetBatch() bool {
	return p.Batch != nil
}

func (p *AgentEmitBatchArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRUCT {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *AgentEmitBatchArgs) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	p.Batch = &jaeger.Batch{}
	if err := p.Batch.Read(ctx, iprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", p.Batch), err)
	}
	return nil
}

func (p *AgentEmitBatchArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "emitBatch_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *AgentEmitBatchArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "batch", thrift.STRUCT, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:batch: ", p), err)
	}
	if err := p.Batch.Write(ctx, oprot); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", p.Batch), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:batch: ", p), err)
	}
	return err
}

func (p *AgentEmitBatchArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AgentEmitBatchArgs(%+v)", *p)
}
