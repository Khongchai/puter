package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	lsproto "puter/lsp"
	"puter/utils"

	"golang.org/x/sync/errgroup"
)

type Reader interface {
	Read() ([]byte, error)
}

type Writer interface {
	Write(msg []byte) error
}

type Engine struct {
	ctx           context.Context
	reader        Reader
	writer        Writer
	requestQueue  chan *lsproto.RequestMessage
	outgoingQueue chan *lsproto.Message
}

func NewEngine(
	ctx context.Context,
	reader Reader,
	writer Writer,
) *Engine {
	return &Engine{
		ctx,
		reader,
		writer,
		make(chan *lsproto.RequestMessage, 100),
		make(chan *lsproto.Message, 100),
	}
}

func (e *Engine) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return e.dispatchLoop(ctx) })
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
	go func() { readLoopErr <- e.readLoop(ctx) }()

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

		if err != nil {
			if errors.Is(err, lsproto.ErrorCodeInvalidRequest) {
				e.sendError(nil, err)
				continue
			}
			return err
		}

		if msg.Kind == lsproto.MessageKindRequest {
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

func (e *Engine) sendError(id *lsproto.ID, err error) {
	code := lsproto.ErrorCodeInternalError
	if errCode := lsproto.ErrorCode(0); errors.As(err, &errCode) {
		code = errCode
	}
	e.outgoingQueue <- &lsproto.ResponseMessage{
		ID: id,
		Error: &lsproto.ResponseError{
			Code:    int32(code),
			Message: err.Error(),
		},
	}
}

func (e *Engine) handleInitialize(ctx context.Context, params *lsproto.InitializeParams, _ *lsproto.RequestMessage) (lsproto.InitializeResponse, error) {
	response := &lsproto.InitializeResult{
		ServerInfo: &lsproto.ServerInfo{
			Name:    "puter",
			Version: utils.PointerTo("0.0.1"),
		},
		Capabilities: &lsproto.ServerCapabilities{
			PositionEncoding: ptrTo(e.positionEncoding),
			TextDocumentSync: &lsproto.TextDocumentSyncOptionsOrKind{
				Options: &lsproto.TextDocumentSyncOptions{
					OpenClose: ptrTo(true),
					Change:    ptrTo(lsproto.TextDocumentSyncKindIncremental),
					Save: &lsproto.BooleanOrSaveOptions{
						Boolean: ptrTo(true),
					},
				},
			},
			HoverProvider: &lsproto.BooleanOrHoverOptions{
				Boolean: ptrTo(true),
			},
			DefinitionProvider: &lsproto.BooleanOrDefinitionOptions{
				Boolean: ptrTo(true),
			},
			TypeDefinitionProvider: &lsproto.BooleanOrTypeDefinitionOptionsOrTypeDefinitionRegistrationOptions{
				Boolean: ptrTo(true),
			},
			ReferencesProvider: &lsproto.BooleanOrReferenceOptions{
				Boolean: ptrTo(true),
			},
			ImplementationProvider: &lsproto.BooleanOrImplementationOptionsOrImplementationRegistrationOptions{
				Boolean: ptrTo(true),
			},
			DiagnosticProvider: &lsproto.DiagnosticOptionsOrRegistrationOptions{
				Options: &lsproto.DiagnosticOptions{
					InterFileDependencies: true,
				},
			},
			CompletionProvider: &lsproto.CompletionOptions{
				TriggerCharacters: &ls.TriggerCharacters,
				ResolveProvider:   ptrTo(true),
				// !!! other options
			},
			SignatureHelpProvider: &lsproto.SignatureHelpOptions{
				TriggerCharacters: &[]string{"(", ","},
			},
			DocumentFormattingProvider: &lsproto.BooleanOrDocumentFormattingOptions{
				Boolean: ptrTo(true),
			},
			DocumentRangeFormattingProvider: &lsproto.BooleanOrDocumentRangeFormattingOptions{
				Boolean: ptrTo(true),
			},
			DocumentOnTypeFormattingProvider: &lsproto.DocumentOnTypeFormattingOptions{
				FirstTriggerCharacter: "{",
				MoreTriggerCharacter:  &[]string{"}", ";", "\n"},
			},
			WorkspaceSymbolProvider: &lsproto.BooleanOrWorkspaceSymbolOptions{
				Boolean: ptrTo(true),
			},
			DocumentSymbolProvider: &lsproto.BooleanOrDocumentSymbolOptions{
				Boolean: ptrTo(true),
			},
			FoldingRangeProvider: &lsproto.BooleanOrFoldingRangeOptionsOrFoldingRangeRegistrationOptions{
				Boolean: ptrTo(true),
			},
			RenameProvider: &lsproto.BooleanOrRenameOptions{
				Boolean: ptrTo(true),
			},
			DocumentHighlightProvider: &lsproto.BooleanOrDocumentHighlightOptions{
				Boolean: ptrTo(true),
			},
			SelectionRangeProvider: &lsproto.BooleanOrSelectionRangeOptionsOrSelectionRangeRegistrationOptions{
				Boolean: ptrTo(true),
			},
			InlayHintProvider: &lsproto.BooleanOrInlayHintOptionsOrInlayHintRegistrationOptions{
				Boolean: ptrTo(true),
			},
			CodeLensProvider: &lsproto.CodeLensOptions{
				ResolveProvider: ptrTo(true),
			},
			CodeActionProvider: &lsproto.BooleanOrCodeActionOptions{
				CodeActionOptions: &lsproto.CodeActionOptions{
					CodeActionKinds: &[]lsproto.CodeActionKind{
						lsproto.CodeActionKindQuickFix,
					},
				},
			},
			CallHierarchyProvider: &lsproto.BooleanOrCallHierarchyOptionsOrCallHierarchyRegistrationOptions{
				Boolean: ptrTo(true),
			},
		},
	}

	return response, nil
}
