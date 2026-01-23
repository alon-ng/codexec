import { Decoration, type DecorationSet, EditorView, EditorState, StateField } from '@uiw/react-codemirror';
import { loadLanguage } from '@uiw/codemirror-extensions-langs';
import LANGUAGE_MAP from "~/utils/codeLang";


export function getCodeMirrorExtensions(language: string, readOnlyLines: number[]) {
    const langName = LANGUAGE_MAP[language];
    const langExt = loadLanguage(langName);

    if (!langExt) {
        return [];
    }

    const readOnlyLineDecoration = Decoration.line({
        class: "cm-readonly-line",
    });

    const readOnlyFilterExtension = EditorState.transactionFilter.of((tr) => {
        if (!tr.docChanged) return tr;

        let isBlocked = false;

        tr.changes.iterChanges((fromA, toA) => {
            const lineStart = tr.startState.doc.lineAt(fromA).number;
            const lineEnd = tr.startState.doc.lineAt(toA).number;

            for (let i = lineStart; i <= lineEnd; i++) {
                if (readOnlyLines.includes(i)) {
                    isBlocked = true;
                    break;
                }
            }
        });

        return isBlocked ? [] : tr;
    });

    const readOnlyLinesField = StateField.define<DecorationSet>({
        create(state) {
            const decorations = [];
            for (const lineNum of readOnlyLines) {
                if (lineNum <= state.doc.lines) {
                    const line = state.doc.line(lineNum);
                    decorations.push(readOnlyLineDecoration.range(line.from));
                }
            }
            return Decoration.set(decorations);
        },
        update(decorations, tr) {
            // Recreate decorations if document changes
            if (tr.docChanged) {
                const newDecorations = [];
                for (const lineNum of readOnlyLines) {
                    if (lineNum <= tr.newDoc.lines) {
                        const line = tr.newDoc.line(lineNum);
                        newDecorations.push(readOnlyLineDecoration.range(line.from));
                    }
                }
                return Decoration.set(newDecorations);
            }
            return decorations.map(tr.changes);
        },
        provide(field) {
            return EditorView.decorations.from(field);
        },
    });

    return [
        readOnlyLinesField,
        readOnlyFilterExtension,
        langExt,
        EditorView.theme({
            ".cm-readonly-line": {
                opacity: "0.4",
            },
        }),
    ];
}