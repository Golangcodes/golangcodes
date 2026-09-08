/**
 * A deliberately small Go highlighter for read-only Go snippets.
 *
 * Pulling in a full highlighter for a handful of static blocks is not worth the
 * bytes. This tokenises the raw text and escapes as it goes, so no unescaped
 * input ever reaches innerHTML.
 */
(function () {
    "use strict";

    const KEYWORDS = new Set([
        "break", "case", "chan", "const", "continue", "default", "defer", "else",
        "fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
        "map", "package", "range", "return", "select", "struct", "switch", "type", "var",
    ]);

    const BUILTINS = new Set([
        "any", "bool", "byte", "comparable", "complex64", "complex128", "error",
        "float32", "float64", "int", "int8", "int16", "int32", "int64", "rune",
        "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
        "true", "false", "nil", "iota", "append", "cap", "clear", "close", "copy",
        "delete", "len", "make", "max", "min", "new", "panic", "print", "println", "recover",
    ]);

    // Comments and strings first: inside them, nothing else is a token.
    const TOKEN = new RegExp(
        [
            "(//[^\\n]*|/\\*[\\s\\S]*?\\*/)",              // 1: comment
            "(`[^`]*`|\"(?:[^\"\\\\\\n]|\\\\.)*\"|'(?:[^'\\\\\\n]|\\\\.)*')", // 2: string or rune
            "\\b(0[xXbBoO][0-9a-fA-F_]+|\\d[\\d_]*(?:\\.[\\d_]+)?)\\b",       // 3: number
            "([A-Za-z_][A-Za-z0-9_]*)(\\s*\\()",           // 4+5: call
            "\\b([A-Za-z_][A-Za-z0-9_]*)\\b",              // 6: identifier
        ].join("|"),
        "g"
    );

    function escape(s) {
        return s
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;");
    }

    function span(cls, text) {
        return '<span class="' + cls + '">' + escape(text) + "</span>";
    }

    function highlightGo(source) {
        let out = "";
        let last = 0;
        let m;

        TOKEN.lastIndex = 0;
        while ((m = TOKEN.exec(source)) !== null) {
            out += escape(source.slice(last, m.index));
            last = TOKEN.lastIndex;

            if (m[1]) {
                out += span("tok-com", m[1]);
            } else if (m[2]) {
                out += span("tok-str", m[2]);
            } else if (m[3]) {
                out += span("tok-num", m[3]);
            } else if (m[4]) {
                // A call: highlight the name, then let the paren fall through.
                out += KEYWORDS.has(m[4])
                    ? span("tok-kw", m[4])
                    : span(BUILTINS.has(m[4]) ? "tok-typ" : "tok-fn", m[4]);
                out += escape(m[5]);
            } else if (m[6]) {
                if (KEYWORDS.has(m[6])) out += span("tok-kw", m[6]);
                else if (BUILTINS.has(m[6])) out += span("tok-typ", m[6]);
                else out += escape(m[6]);
            }
        }

        return out + escape(source.slice(last));
    }

    window.GCHighlight = highlightGo;

    document.addEventListener("DOMContentLoaded", function () {
        document.querySelectorAll("pre[data-code] code").forEach(function (el) {
            el.innerHTML = highlightGo(el.textContent);
        });
    });
})();
