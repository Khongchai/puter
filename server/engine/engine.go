package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"puter/interpreter"
	"puter/logging"
	lsproto "puter/lsp"
	"puter/utils"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Reader interface {
	Read() ([]byte, error)
}

type Writer interface {
	Write(msg []byte) error
}

type pendingClientRequest struct {
	req    *lsproto.RequestMessage
	cancel context.CancelFunc
}

type Engine struct {
	ctx                     context.Context
	reader                  Reader
	writer                  Writer
	requestQueue            chan *lsproto.RequestMessage
	outgoingQueue           chan *lsproto.Message
	pendingServerRequests   map[lsproto.ID]chan *lsproto.ResponseMessage
	pendingServerRequestsMu sync.Mutex
	pendingClientRequests   map[lsproto.ID]pendingClientRequest
	pendingClientRequestsMu sync.Mutex
	logger                  logging.Logger
	initComplete            bool
	interpreter             *interpreter.Interpreter
}

func NewEngine(
	ctx context.Context,
	reader Reader,
	writer Writer,
	logger logging.Logger,
	interpreter *interpreter.Interpreter,
) *Engine {
	return &Engine{
		ctx:           ctx,
		reader:        reader,
		writer:        writer,
		requestQueue:  make(chan *lsproto.RequestMessage, 100),
		outgoingQueue: make(chan *lsproto.Message, 100),
		logger:        logger,
		initComplete:  false,
		interpreter:   interpreter,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// request response queues
	// starts a worker that writes to outgoing queue.
	g.Go(func() error { return e.writeLoop(ctx) })
	// Don't run readLoop in the group, as it blocks on stdin read and cannot be cancelled.
	readLoopErr := make(chan error, 1)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-readLoopErr:
			return err
		}
	})
	// listens to incoming io events and handle in dispatchLoop
	go func() { readLoopErr <- e.readLoop(ctx) }()

	// handler queue
	// handles incoming request
	g.Go(func() error { return e.dispatchLoop(ctx) })

	if err := g.Wait(); err != nil && !errors.Is(err, io.EOF) && ctx.Err() != nil {
		return err
	}

	return nil
}

func (e *Engine) readLoop(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		data, err := e.reader.Read()
		msg := &lsproto.Message{}
		if err := json.Unmarshal(data, msg); err != nil {
			return fmt.Errorf("%w: %w", lsproto.ErrorCodeInvalidRequest, err)
		}

		e.logger.Info("read %+v", msg)

		if err != nil {
			if errors.Is(err, lsproto.ErrorCodeInvalidRequest) {
				e.sendError(nil, err)
				continue
			}
			return err
		}

		if !e.initComplete && msg.Kind == lsproto.MessageKindRequest {
			req := msg.AsRequest()
			if req.Method == lsproto.MethodInitialize {
				resp, err := e.handleInitialize(ctx, req.Params.(*lsproto.InitializeParams), req)
				if err != nil {
					return err
				}
				e.sendResult(req.ID, resp)
			} else {
				e.sendError(req.ID, lsproto.ErrorCodeServerNotInitialized)
			}
			continue
		}

		if msg.Kind == lsproto.MessageKindResponse {
			resp := msg.AsResponse()
			e.pendingServerRequestsMu.Lock()
			if respChan, ok := e.pendingServerRequests[*resp.ID]; ok {
				respChan <- resp
				close(respChan)
				delete(e.pendingServerRequests, *resp.ID)
			}
			e.pendingServerRequestsMu.Unlock()
		} else {
			req := msg.AsRequest()
			if req.Method == lsproto.MethodCancelRequest {
				e.cancelRequest(req.Params.(*lsproto.CancelParams).Id)
			} else {
				e.requestQueue <- req
			}
		}
	}
}

func (e *Engine) cancelRequest(rawID lsproto.IntegerOrString) {
	id := lsproto.NewID(rawID)
	e.pendingClientRequestsMu.Lock()
	defer e.pendingClientRequestsMu.Unlock()
	if pendingReq, ok := e.pendingClientRequests[*id]; ok {
		pendingReq.cancel()
		delete(e.pendingClientRequests, *id)
	}
}

func (e *Engine) writeLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data := <-e.outgoingQueue:
			bytes, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("%w: %w", lsproto.ErrorCodeInvalidRequest, err)
			}
			if err := e.writer.Write(bytes); err != nil {
				return fmt.Errorf("failed to write message: %w", err)
			}
		}
	}
}

func (e *Engine) dispatchLoop(ctx context.Context) error {
	ctx, lspExit := context.WithCancel(ctx)
	defer lspExit()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-e.requestQueue:
			if req.ID != nil {
				var cancel context.CancelFunc
				e.pendingClientRequestsMu.Lock()
				e.pendingClientRequests[*req.ID] = pendingClientRequest{
					req:    req,
					cancel: cancel,
				}
				e.pendingClientRequestsMu.Unlock()
			}

			handle := func() {
				if err := e.handleRequestOrNotification(ctx, req); err != nil {
					if errors.Is(err, context.Canceled) {
						e.sendError(req.ID, lsproto.ErrorCodeRequestCancelled)
					} else if errors.Is(err, io.EOF) {
						lspExit()
					} else {
						e.sendError(req.ID, err)
					}
				}

				if req.ID != nil {
					e.pendingClientRequestsMu.Lock()
					delete(e.pendingClientRequests, *req.ID)
					e.pendingClientRequestsMu.Unlock()
				}
			}

			if isBlockingMethod(req.Method) {
				handle()
			} else {
				go handle()
			}
		}
	}
}

func (e *Engine) handleRequestOrNotification(ctx context.Context, req *lsproto.RequestMessage) error {
	if handler := handlers()[req.Method]; handler != nil {
		start := time.Now()
		err := handler(e, ctx, req)
		e.logger.Info("handled method '", req.Method, "' in ", time.Since(start))
		return err
	}
	e.logger.Warn("unknown method '", req.Method, "'")
	if req.ID != nil {
		e.sendError(req.ID, lsproto.ErrorCodeInvalidRequest)
	}
	return nil
}

type handlerMap map[lsproto.Method]func(*Engine, context.Context, *lsproto.RequestMessage) error

var handlers = sync.OnceValue(func() handlerMap {
	handlers := make(handlerMap)

	registerRequestHandler(handlers, lsproto.InitializeInfo, (*Engine).handleInitialize)
	registerNotificationHandler(handlers, lsproto.InitializedInfo, (*Engine).handleInitialized)

	registerNotificationHandler(handlers, lsproto.TextDocumentDidChangeInfo, (*Engine).handleTextDocumentDidChange)

	return handlers
})

func registerNotificationHandler[Req any](handlers handlerMap, info lsproto.NotificationInfo[Req], fn func(*Engine, context.Context, Req) error) {
	handlers[info.Method] = func(e *Engine, ctx context.Context, req *lsproto.RequestMessage) error {
		var params Req
		// Ignore empty params; all generated params are either pointers or any.
		if req.Params != nil {
			params = req.Params.(Req)
		}
		if err := fn(e, ctx, params); err != nil {
			return err
		}
		return ctx.Err()
	}
}

func registerRequestHandler[Req, Resp any](
	handlers handlerMap,
	info lsproto.RequestInfo[Req, Resp],
	fn func(*Engine, context.Context, Req, *lsproto.RequestMessage) (Resp, error),
) {
	handlers[info.Method] = func(e *Engine, ctx context.Context, req *lsproto.RequestMessage) error {
		var params Req
		// Ignore empty params.
		if req.Params != nil {
			params = req.Params.(Req)
		}
		resp, err := fn(e, ctx, params, req)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return e.sendResult(req.ID, resp)
	}
}

func isBlockingMethod(method lsproto.Method) bool {
	switch method {
	case lsproto.MethodInitialize,
		lsproto.MethodInitialized,
		lsproto.MethodTextDocumentDidOpen,
		lsproto.MethodTextDocumentDidChange,
		lsproto.MethodTextDocumentDidSave,
		lsproto.MethodTextDocumentDidClose,
		lsproto.MethodWorkspaceDidChangeWatchedFiles,
		lsproto.MethodWorkspaceDidChangeConfiguration,
		lsproto.MethodWorkspaceConfiguration:
		return true
	}
	return false
}

func (e *Engine) sendError(id *lsproto.ID, err error) {
	code := lsproto.ErrorCodeInternalError
	if errCode := lsproto.ErrorCode(0); errors.As(err, &errCode) {
		code = errCode
	}
	e.send((&lsproto.ResponseMessage{
		ID: id,
		Error: &lsproto.ResponseError{
			Code:    int32(code),
			Message: err.Error(),
		},
	}).Message())
}

func (e *Engine) sendResult(id *lsproto.ID, result any) error {
	return e.send((&lsproto.ResponseMessage{
		ID:     id,
		Result: result,
	}).Message())
}

func (e *Engine) send(resp *lsproto.Message) error {
	select {
	case e.outgoingQueue <- resp:
		return nil
	case <-e.ctx.Done():
		return e.ctx.Err()
	}
}

func (e *Engine) handleInitialize(ctx context.Context, params *lsproto.InitializeParams, _ *lsproto.RequestMessage) (lsproto.InitializeResponse, error) {
	response := &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name:    "puter",
			Version: utils.PointerTo("0.0.1"),
		},
		Capabilities: &lsproto.ServerCapabilities{
			TextDocumentSync: &lsproto.TextDocumentSyncOptionsOrKind{
				Options: &lsproto.TextDocumentSyncOptions{
					OpenClose: utils.PointerTo(true),
					Change:    utils.PointerTo(lsproto.TextDocumentSyncKindFull),
					Save: &lsproto.BooleanOrSaveOptions{
						Boolean: utils.PointerTo(true),
					},
				},
			},
			// HoverProvider: &lsproto.BooleanOrHoverOptions{
			// 	Boolean: utils.PointerTo(true),
			// },
			// DefinitionProvider: &lsproto.BooleanOrDefinitionOptions{
			// 	Boolean: utils.PointerTo(true),
			// },
			// TypeDefinitionProvider: &lsproto.BooleanOrTypeDefinitionOptionsOrTypeDefinitionRegistrationOptions{
			// 	Boolean: utils.PointerTo(true),
			// },
			// ReferencesProvider: &lsproto.BooleanOrReferenceOptions{
			// 	Boolean: utils.PointerTo(true),
			// },
			// DiagnosticProvider: &lsproto.DiagnosticOptionsOrRegistrationOptions{
			// 	Options: &lsproto.DiagnosticOptions{
			// 		InterFileDependencies: true,
			// 	},
			// },
			// RenameProvider: &lsproto.BooleanOrRenameOptions{
			// 	Boolean: utils.PointerTo(true),
			// },
		},
	}

	return response, nil
}

func (e *Engine) handleInitialized(ctx context.Context, params *lsproto.InitializedParams) error {
	e.initComplete = true
	return nil
}

func (e *Engine) handleTextDocumentDidChange(ctx context.Context, params *lsproto.DidChangeTextDocumentParams) error {
	// uri := params.TextDocument.Uri
	for _, change := range params.ContentChanges {
		interpretations := e.interpreter.Interpret(
			change.WholeDocument.Text,
		)
		response := &lsproto.RequestMessage{
			Method: "custom/evaluationResult",
			Params: interpretations,
		}
		e.send(response.Message())
	}
	return nil
}
