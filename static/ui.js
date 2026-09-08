/**
 * Component behaviours. Vanilla, no build step, no framework.
 *
 * Everything is wired from data- attributes on DOMContentLoaded, so dropping a
 * component into any template is enough — there is nothing to initialise by
 * hand. The public surface is `window.UI`:
 *
 *     UI.toast("Saved")                     transient notification
 *     UI.copy(text)                         clipboard, with a fallback
 *     UI.modal.open(id) / UI.modal.close()  dialogs
 *     UI.editor(id)                         the Editor for a ui/editor block
 *     UI.run(id) / UI.share(id)             compile and share a buffer
 *
 * Browser alert/confirm/prompt are deliberately never used: they block the
 * page and cannot be styled.
 */
(function () {
    "use strict";

    var UI = {};
    var editors = new Map();

    function $(sel, root) { return (root || document).querySelector(sel); }
    function $$(sel, root) { return Array.prototype.slice.call((root || document).querySelectorAll(sel)); }

    /* ----------------------------------------------------------- toasts */

    var TOAST_SKIN = {
        info: "border-sky-500/30 bg-sky-950/90 text-sky-200",
        success: "border-emerald-500/30 bg-emerald-950/90 text-emerald-200",
        warn: "border-amber-500/30 bg-amber-950/90 text-amber-200",
        error: "border-red-500/30 bg-red-950/90 text-red-200",
    };

    UI.toast = function (message, opts) {
        opts = opts || {};
        var host = document.getElementById("ui-toasts");
        if (!host) return null;

        var el = document.createElement("div");
        el.className =
            "ui-toast flex items-start gap-3 border rounded-xl px-4 py-3 text-sm shadow-xl backdrop-blur max-w-sm " +
            (TOAST_SKIN[opts.variant] || TOAST_SKIN.info);
        el.textContent = message;

        // Errors interrupt; everything else waits its turn in the live region.
        el.setAttribute("role", opts.variant === "error" ? "alert" : "status");
        host.appendChild(el);

        var timeout = opts.timeout === undefined ? 3200 : opts.timeout;
        if (timeout > 0) {
            setTimeout(function () {
                el.setAttribute("data-leaving", "");
                setTimeout(function () { el.remove(); }, 160);
            }, timeout);
        }
        return el;
    };

    /* -------------------------------------------------------- clipboard */

    UI.copy = function (text) {
        if (navigator.clipboard && window.isSecureContext) {
            return navigator.clipboard.writeText(text);
        }
        // execCommand is deprecated but is the only path on plain http://.
        return new Promise(function (resolve, reject) {
            var ta = document.createElement("textarea");
            ta.value = text;
            ta.setAttribute("readonly", "");
            ta.style.cssText = "position:fixed;top:-9999px;opacity:0";
            document.body.appendChild(ta);
            ta.select();
            var ok = false;
            try { ok = document.execCommand("copy"); } catch (e) { ok = false; }
            ta.remove();
            ok ? resolve() : reject(new Error("copy failed"));
        });
    };

    function copyWithFeedback(text, label) {
        UI.copy(text).then(
            function () { UI.toast(label || "Copied", { variant: "success", timeout: 1600 }); },
            function () { UI.toast("Could not copy — select the text and copy manually", { variant: "error" }); }
        );
    }

    /* ----------------------------------------------------------- editor */

    function Editor(root) {
        this.root = root;
        this.id = root.getAttribute("data-ui-editor-id") || "editor";
        this.input = $("[data-ui-editor-input]", root);
        this.gutter = $("[data-ui-editor-gutter]", root);
        this.original = this.input ? this.input.value : "";
        this.bind();
        this.renderGutter();
    }

    Editor.prototype.value = function (v) {
        if (v === undefined) return this.input ? this.input.value : "";
        this.input.value = v;
        this.renderGutter();
        return v;
    };

    Editor.prototype.reset = function () {
        this.value(this.original);
        UI.toast("Reset to the original code", { timeout: 1600 });
    };

    Editor.prototype.renderGutter = function () {
        if (!this.gutter || !this.input) return;
        var lines = this.input.value.split("\n").length;
        var out = "";
        for (var i = 1; i <= lines; i++) out += i + "\n";
        this.gutter.textContent = out;
        this.gutter.scrollTop = this.input.scrollTop;
    };

    // Replace the selection and keep the caret where the typist expects it.
    Editor.prototype.splice = function (text, selStart, selEnd) {
        var el = this.input;
        var start = el.selectionStart;
        var end = el.selectionEnd;
        el.setRangeText(text, start, end, "end");
        if (selStart !== undefined) el.selectionStart = start + selStart;
        if (selEnd !== undefined) el.selectionEnd = start + selEnd;
        this.renderGutter();
    };

    Editor.prototype.indentSelection = function (outdent) {
        var el = this.input;
        var value = el.value;
        var from = value.lastIndexOf("\n", el.selectionStart - 1) + 1;
        var to = el.selectionEnd;
        var block = value.slice(from, to);

        var updated = outdent
            ? block.replace(/^(\t| {1,4})/gm, "")
            : block.replace(/^/gm, "\t");

        el.setRangeText(updated, from, to, "select");
        this.renderGutter();
    };

    Editor.prototype.bind = function () {
        var self = this;
        if (!this.input) return;

        this.input.addEventListener("input", function () { self.renderGutter(); });
        this.input.addEventListener("scroll", function () {
            if (self.gutter) self.gutter.scrollTop = self.input.scrollTop;
        });

        this.input.addEventListener("keydown", function (e) {
            // Cmd/Ctrl+Enter runs, matching the playground.
            if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) {
                e.preventDefault();
                UI.run(self.id);
                return;
            }

            if (e.key === "Tab") {
                e.preventDefault();
                if (self.input.selectionStart !== self.input.selectionEnd || e.shiftKey) {
                    self.indentSelection(e.shiftKey);
                } else {
                    self.splice("\t");
                }
                return;
            }

            // Enter keeps the current indent, and adds one inside a new block.
            if (e.key === "Enter" && !e.shiftKey) {
                var value = self.input.value;
                var lineStart = value.lastIndexOf("\n", self.input.selectionStart - 1) + 1;
                var line = value.slice(lineStart, self.input.selectionStart);
                var indent = (line.match(/^[\t ]*/) || [""])[0];
                if (/[{(\[]\s*$/.test(line)) indent += "\t";
                if (indent) {
                    e.preventDefault();
                    self.splice("\n" + indent);
                }
            }
        });
    };

    UI.editor = function (id) { return editors.get(id) || null; };

    function editorFrom(el) {
        var root = el.closest("[data-ui-editor]");
        if (root) return editors.get(root.getAttribute("data-ui-editor-id"));
        return null;
    }

    // The nearest thing that holds code: an editor, or any block with a source.
    function sourceFrom(el) {
        var ed = editorFrom(el);
        if (ed) return ed.value();
        var block = el.closest("[data-ui-share-block], figure, [data-ui-editor]");
        var src = block && (block.querySelector("[data-ui-copy-source]") || block.querySelector("[data-ui-editor-input]"));
        return src ? (src.value !== undefined ? src.value : src.textContent) : "";
    }

    /* ------------------------------------------------------- run / share */

    var PLAYGROUND_COMPILE = "https://play.golang.org/compile";
    var PLAYGROUND_SHARE = "https://play.golang.org/share";
    var PLAYGROUND_LINK = "https://go.dev/play/p/";

    // The playground answers {Errors, Events:[{Message, Kind}]}. Compile and
    // vet failures arrive in Errors; a program that ran reports through Events.
    function playgroundText(data) {
        if (data.Errors) return { output: data.Errors, error: true };

        var out = "";
        var stderr = false;
        (data.Events || []).forEach(function (e) {
            out += e.Message || "";
            if (e.Kind === "stderr") stderr = true;
        });
        return { output: out || "Program exited successfully with no output.", error: stderr };
    }

    function outputFor(id) {
        return $('[data-ui-output-for="' + id + '"]') || null;
    }

    UI.run = function (id) {
        var ed = editors.get(id);
        if (!ed) return Promise.resolve();

        var pane = outputFor(id);
        var body = pane && $("[data-ui-output-body]", pane);
        var status = pane && $("[data-ui-output-status]", pane);
        var button = $("[data-ui-run]", ed.root);

        if (pane) pane.classList.remove("hidden");
        if (body) body.textContent = "Running…";
        if (status) status.innerHTML = "";
        if (button) button.disabled = true;

        var form = new URLSearchParams();
        form.set("version", "2");
        form.set("withVet", "true");
        form.set("body", ed.value());

        return fetch(PLAYGROUND_COMPILE, {
            method: "POST",
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            body: form.toString(),
        })
            .then(function (r) { return r.json(); })
            .then(function (data) {
                var result = playgroundText(data);
                if (body) body.textContent = result.output;
                if (status) {
                    status.className = result.error
                        ? "font-mono text-xs text-red-400"
                        : "font-mono text-xs text-emerald-400";
                    status.textContent = result.error ? "error" : "ok";
                }
            })
            .catch(function () {
                if (body) body.textContent = "Could not reach the Go playground.";
                if (status) {
                    status.className = "font-mono text-xs text-red-400";
                    status.textContent = "error";
                }
            })
            .finally(function () { if (button) button.disabled = false; });
    };

    UI.share = function (idOrCode) {
        var code = editors.has(idOrCode) ? editors.get(idOrCode).value() : idOrCode;
        if (!code) return Promise.resolve();

        return fetch(PLAYGROUND_SHARE, {
            method: "POST",
            headers: { "Content-Type": "text/plain; charset=utf-8" },
            body: code,
        })
            .then(function (r) {
                if (!r.ok) throw new Error("share rejected");
                return r.text();
            })
            .then(function (id) {
                id = (id || "").trim();
                if (!id) throw new Error("no id");
                // The playground hosts the snippet, so the link outlives this site.
                copyWithFeedback(PLAYGROUND_LINK + id, "Playground link copied");
            })
            .catch(function () {
                UI.toast("Could not create a share link", { variant: "error" });
            });
    };

    /* ------------------------------------------------------------- tabs */

    function initTabs(root) {
        var tabs = $$('[role="tab"]', root);
        if (!tabs.length) return;

        function select(tab) {
            tabs.forEach(function (t) {
                var on = t === tab;
                t.setAttribute("aria-selected", on ? "true" : "false");
                t.tabIndex = on ? 0 : -1;
                var panel = document.getElementById(t.getAttribute("aria-controls"));
                if (panel) panel.hidden = !on;
            });
        }

        tabs.forEach(function (tab, i) {
            tab.addEventListener("click", function () { select(tab); });
            tab.addEventListener("keydown", function (e) {
                var next = null;
                if (e.key === "ArrowRight") next = tabs[(i + 1) % tabs.length];
                else if (e.key === "ArrowLeft") next = tabs[(i - 1 + tabs.length) % tabs.length];
                else if (e.key === "Home") next = tabs[0];
                else if (e.key === "End") next = tabs[tabs.length - 1];
                if (next) { e.preventDefault(); select(next); next.focus(); }
            });
        });
    }

    /* ------------------------------------------------------------ modal */

    var openModal = null;
    var lastFocus = null;

    UI.modal = {
        open: function (id) {
            var el = document.getElementById(id);
            if (!el) return;
            lastFocus = document.activeElement;
            el.classList.remove("hidden");
            el.classList.add("flex");
            openModal = el;
            var focusable = $('button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])', el);
            if (focusable) focusable.focus();
        },
        close: function () {
            if (!openModal) return;
            openModal.classList.add("hidden");
            openModal.classList.remove("flex");
            openModal = null;
            if (lastFocus && lastFocus.focus) lastFocus.focus();
        },
    };

    /* ------------------------------------------------------------- init */

    function init() {
        $$("[data-ui-editor]").forEach(function (root) {
            var ed = new Editor(root);
            editors.set(ed.id, ed);
        });

        // Share blocks are editors without a UI; give them a buffer to read.
        $$("[data-ui-tabs]").forEach(initTabs);

        if (window.GCHighlight) {
            $$("[data-ui-highlight]").forEach(function (el) {
                el.innerHTML = window.GCHighlight(el.textContent);
            });
        }

        document.addEventListener("click", function (e) {
            var el = e.target.closest ? e.target.closest("[data-ui-copy], [data-ui-run], [data-ui-share], [data-ui-editor-reset], [data-ui-dismiss], [data-ui-modal-close], [data-ui-modal-open]") : null;
            if (!el) return;

            if (el.hasAttribute("data-ui-copy")) {
                copyWithFeedback(sourceFrom(el), "Code copied");
            } else if (el.hasAttribute("data-ui-run")) {
                var runEd = editorFrom(el);
                if (runEd) UI.run(runEd.id);
            } else if (el.hasAttribute("data-ui-share")) {
                var shareEd = editorFrom(el);
                UI.share(shareEd ? shareEd.id : sourceFrom(el));
            } else if (el.hasAttribute("data-ui-editor-reset")) {
                var resetEd = editorFrom(el);
                if (resetEd) resetEd.reset();
            } else if (el.hasAttribute("data-ui-dismiss")) {
                var target = el.closest("[data-ui-dismiss-target]");
                if (target) target.remove();
            } else if (el.hasAttribute("data-ui-modal-close")) {
                UI.modal.close();
            } else if (el.hasAttribute("data-ui-modal-open")) {
                UI.modal.open(el.getAttribute("data-ui-modal-open"));
            }
        });

        document.addEventListener("keydown", function (e) {
            if (e.key === "Escape" && openModal) UI.modal.close();
        });
    }

    window.UI = UI;

    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", init);
    } else {
        init();
    }
})();
