import * as path from "path";
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

const decorationType = vscode.window.createTextEditorDecorationType({
  after: {
    color: "#637777",
    fontStyle: "italic",
    margin: "0 0 0 3em",
  },
});

const diagnosticCollection =
  vscode.languages.createDiagnosticCollection("puter");

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
      ...payloads: Array<{
        uri: string;
        interpretations: Array<{
          LineIndex: number;
          EvalResult: string;
          Diagnostics: vscode.Diagnostic[];
        }>;
      }>
    ) => {
      const editor = vscode.window.activeTextEditor;
      if (!editor) {
        return;
      }

      if (
        payloads.length === 0 ||
        payloads.every((p) => p.interpretations.length === 0)
      ) {
        editor.setDecorations(decorationType, []);
        return;
      }

      for (const payload of payloads) {
        const decorationOptions: vscode.DecorationOptions[] =
          payload.interpretations.map((evaluation) => {
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
          });

        editor.setDecorations(decorationType, decorationOptions);

        const diagnostics = payload.interpretations.flatMap((p) => {
          return p.Diagnostics;
        });
        diagnosticCollection.set(vscode.Uri.parse(payload.uri), diagnostics);
      }
    },
  );

  const allLanguages = await vscode.languages.getLanguages();
  const rule: vscode.LanguageConfiguration = {
    onEnterRules: [
      {
        beforeText: /^\s*\/\/\s*\|.*$/,
        action: {
          indentAction: vscode.IndentAction.None,
          appendText: "// | ",
        },
      },
      {
        beforeText: /^\s*#\s*\|.*$/,
        action: {
          indentAction: vscode.IndentAction.None,
          appendText: "# | ",
        },
      },
    ],
  };
  const disposables = allLanguages.map((lang) =>
    vscode.languages.setLanguageConfiguration(lang, rule),
  );
  context.subscriptions.push(...disposables);
}

export async function deactivate() {
  await client?.stop();
}
