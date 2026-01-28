import * as path from "path";
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

const evalDecorationType = vscode.window.createTextEditorDecorationType({
  after: {
    color: "#637777",
    fontStyle: "italic",
    margin: "0 0 0 3em",
  },
});

export async function activate(context: vscode.ExtensionContext) {
  const serverOptions = (() => {
    const lsPath = path.resolve(context.extensionPath, "binaries");
    const exePath = path.join(
      lsPath,
      `puter${process.platform === "win32" ? ".exe" : ""}`,
    );
    const option: ServerOptions = {
      // Note: if we can't find the package during build, take a look at this
      // https://github.com/microsoft/typescript-go/blob/main/_packages/native-preview/lib/getExePath.js
      command: exePath,
      transport: TransportKind.stdio,
    };
    return option;
  })();

  const clientOptions: LanguageClientOptions = (() => {
    const outputChannel = vscode.window.createOutputChannel("puter");
    const traceOutputChannel =
      vscode.window.createOutputChannel("puter (trace)");
    return {
      outputChannel,
      traceOutputChannel,
      documentSelector: [
        {
          scheme: "file",
          language: "*",
        },
        {
          scheme: "untitled",
          language: "*",
        },
      ],
    };
  })();

  client = new LanguageClient(
    "puter",
    "puter language server",
    serverOptions,
    clientOptions,
  );

  await client.start();
  await client.sendNotification("workspace/didChangeConfiguration", {
    settings: {
      "vscode-languageclient": {
        trace: { server: "verbose" },
      },
    },
  });
  client.onNotification(
    "custom/evaluationReport",
    (
      ...payload: {
        LineIndex: number;
        EvalResult: string;
        Diagnostics: vscode.Diagnostic[];
      }[]
    ) => {
      const editor = vscode.window.activeTextEditor;
      if (!editor) {
        return;
      }

      const decorationOptions: vscode.DecorationOptions[] = payload.map(
        (evaluation) => {
          const line = editor.document.lineAt(evaluation.LineIndex);
          return {
            range: new vscode.Range(
              evaluation.LineIndex,
              line.range.end.character,
              evaluation.LineIndex,
              line.range.end.character,
            ),
            renderOptions: {
              after: {
                contentText: evaluation.EvalResult,
              },
            },
          };
        },
      );

      editor.setDecorations(evalDecorationType, decorationOptions);
    },
  );
}

export async function deactivate() {
  await client?.stop();
}
