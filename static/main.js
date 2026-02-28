document.addEventListener("DOMContentLoaded", () => {
    // 1. Gamification & Tracking State
    const currentPath = window.location.pathname;

    let userXP = window.StorageAPI.getUserXP();
    const scoreDisplay = document.getElementById('user-score');

    function updateXPDisplay(amount) {
        userXP = window.StorageAPI.addXP(amount);
        if (scoreDisplay) {
            scoreDisplay.textContent = userXP + ' XP';
            // Simple animation
            scoreDisplay.classList.add('text-green-500', 'scale-110');
            setTimeout(() => scoreDisplay.classList.remove('text-green-500', 'scale-110'), 300);
        }
    }

    if (scoreDisplay) scoreDisplay.textContent = userXP + ' XP';

    // 2. Track page visits & XP
    if (currentPath.startsWith('/go/')) {
        const decodedPath = decodeURIComponent(currentPath);
        window.StorageAPI.trackPageVisit(decodedPath,
            () => updateXPDisplay(50), // 50 XP for completing a new lesson!
            () => updateXPDisplay(5)   // 5 XP for reviewing an old lesson!
        );
    }

    // 2.5 Dynamic Page Rating 
    const ratingContainer = document.getElementById('quality-score-container');
    const starButtons = document.querySelectorAll('#star-rating button');
    const percentageUI = document.getElementById('score-percentage');

    if (ratingContainer && starButtons.length > 0) {
        const slug = ratingContainer.getAttribute('data-slug');
        const starRatingWrapper = document.getElementById('star-rating');

        // Load existing ratings from localStorage
        let userRatings = {};
        try {
            userRatings = JSON.parse(localStorage.getItem('user_ratings')) || {};
        } catch (e) {
            userRatings = {};
        }

        const existingScore = userRatings[slug];

        // Function to lock stars visually
        const lockStars = (score) => {
            starButtons.forEach(b => {
                const s = parseInt(b.getAttribute('data-score'), 10);
                if (s <= score) {
                    b.classList.add('text-yellow-400');
                    b.classList.remove('text-slate-300', 'dark:text-slate-600');
                } else {
                    b.classList.remove('text-yellow-400');
                    b.classList.add('text-slate-300', 'dark:text-slate-600');
                }
                b.disabled = true;
                b.classList.add('cursor-default');
            });
        };

        // If user already rated this page, lock the UI on load
        if (existingScore) {
            lockStars(existingScore);
            const rateText = starRatingWrapper.querySelector('span');
            if (rateText) rateText.textContent = 'Rated';
        }

        starButtons.forEach(btn => {
            // Uniform Hover Logic
            btn.addEventListener('mouseenter', () => {
                if (existingScore) return; // Don't hover if already rated
                const hoverScore = parseInt(btn.getAttribute('data-score'), 10);
                starButtons.forEach(b => {
                    const s = parseInt(b.getAttribute('data-score'), 10);
                    if (s <= hoverScore) {
                        b.classList.add('text-yellow-400');
                        b.classList.remove('text-slate-300', 'dark:text-slate-600');
                    } else {
                        b.classList.remove('text-yellow-400');
                        b.classList.add('text-slate-300', 'dark:text-slate-600');
                    }
                });
            });

            btn.addEventListener('click', async (e) => {
                if (existingScore) return; // Prevent clicking if already rated

                const score = parseInt(btn.getAttribute('data-score'), 10);

                // Visual feedback for clicking
                btn.style.transform = 'scale(1.2)';
                setTimeout(() => btn.style.transform = 'scale(1)', 150);

                // Lock stars after voting
                lockStars(score);
                const rateText = starRatingWrapper.querySelector('span');
                if (rateText) rateText.textContent = 'Rated';

                // Save to localStorage
                userRatings[slug] = score;
                localStorage.setItem('user_ratings', JSON.stringify(userRatings));

                try {
                    const response = await fetch('/api/rate', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ slug: slug, score: score })
                    });

                    if (response.ok) {
                        const data = await response.json();
                        // Update UI seamlessly
                        if (percentageUI) {
                            percentageUI.textContent = data.percentage + '%';
                            percentageUI.classList.add('text-green-500', 'scale-110');
                            setTimeout(() => percentageUI.classList.remove('text-green-500', 'scale-110'), 400);
                        }
                        updateXPDisplay(5); // Reward user for rating!
                    }
                } catch (error) {
                    console.error("Failed to submit rating:", error);
                }
            });
        });

        if (starRatingWrapper) {
            starRatingWrapper.addEventListener('mouseleave', () => {
                if (userRatings[slug]) return; // Keep locked state
                starButtons.forEach(b => {
                    if (!b.disabled) {
                        b.classList.remove('text-yellow-400');
                        b.classList.add('text-slate-300', 'dark:text-slate-600');
                    }
                });
            });
        }
    }

    // 3. Mark links in sidebar as completed with read counts
    const sidebarLinks = document.querySelectorAll('aside a');
    const readCounts = window.StorageAPI.getReadCounts();
    sidebarLinks.forEach(link => {
        // HTMX links or Standard Links
        const href = decodeURIComponent(link.getAttribute('href') || '');
        if (readCounts[href] > 0) {
            const count = readCounts[href];
            const badge = document.createElement('span');
            badge.innerHTML = `
                <div class="flex items-center gap-1.5 ml-auto shrink-0">
                    <span class="text-[10px] font-bold text-slate-400 bg-slate-200 dark:bg-slate-700 dark:text-slate-300 px-1.5 py-0.5 rounded" title="Times Read">${count}x</span>
                    <svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                </div>
            `;
            if (!link.classList.contains('flex')) {
                link.classList.add('flex', 'justify-between', 'items-center', 'w-full');
            }
            link.appendChild(badge);
        }
    });

    // Default Go code
    const defaultCode = `package main

import "fmt"

func main() {
	fmt.Println("Hello, Go!")
}
`;

    // Load any saved code for the current page
    let pageCode = window.StorageAPI.getCode(currentPath) || defaultCode;

    require(['vs/editor/editor.main'], function () {
        // Initialize Monaco Editor
        const isDark = document.documentElement.classList.contains('dark');
        window.editor = monaco.editor.create(document.getElementById('editor'), {
            value: pageCode,
            language: 'go',
            theme: isDark ? 'vs-dark' : 'vs',
            automaticLayout: true,
            minimap: { enabled: false },
            fontSize: 14,
            fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
            scrollBeyondLastLine: false,
            roundedSelection: false,
            padding: { top: 16 }
        });

        // Save progress to local storage when code changes
        window.editor.onDidChangeModelContent(() => {
            const code = window.editor.getValue();
            window.StorageAPI.saveCode(currentPath, code);

            // Live update the code snippet in the article if one is actively selected!
            if (window.currentEditingSnippetIndex !== undefined) {
                const index = window.currentEditingSnippetIndex;
                window.StorageAPI.saveSnippet(currentPath, index, code);

                const codeBlocks = document.querySelectorAll('article pre');
                if (codeBlocks[index]) {
                    const codeNode = codeBlocks[index].querySelector('code') || codeBlocks[index];
                    codeNode.textContent = code;
                    if (window.Prism) Prism.highlightElement(codeNode);
                }
            }
        });

        // Setup HTMX integration inside the require block since window.editor is set here
        const codeInput = document.getElementById('code-input');

        // Sync editor content to the hidden textarea before HTMX request
        document.body.addEventListener('htmx:configRequest', function (evt) {
            if (evt.detail.path === '/api/run') {
                codeInput.value = window.editor.getValue();
                evt.detail.parameters['code'] = window.editor.getValue();
            }
        });

        // Grant XP for practicing code
        document.body.addEventListener('htmx:afterRequest', function (evt) {
            if (evt.detail.path === '/api/run' && evt.detail.successful) {
                updateXPDisplay(10); // 10 XP for practicing code
            }
        });
    });

    // Make the pane resizable (simple implementation)
    const divider = document.querySelector('.cursor-row-resize');
    const editorContainer = document.getElementById('editor');
    const outputContainer = document.querySelector('.h-\\[40\\%\\]');
    let isResizing = false;

    divider.addEventListener('mousedown', (e) => {
        isResizing = true;
        document.body.style.cursor = 'row-resize';
    });

    document.addEventListener('mousemove', (e) => {
        if (!isResizing) return;

        const section = document.querySelector('section');
        const sectionRect = section.getBoundingClientRect();

        // Calculate new heights (accounting for action bar)
        let newEditorHeight = e.clientY - sectionRect.top - 56; // 56px is action bar height

        // Add constraints
        if (newEditorHeight < 100) newEditorHeight = 100;
        if (sectionRect.height - newEditorHeight - 56 < 100) newEditorHeight = sectionRect.height - 156;

        const percentage = (newEditorHeight / (sectionRect.height - 56)) * 100;

        editorContainer.style.height = `${percentage}%`;
        editorContainer.style.flex = 'none';

        outputContainer.style.height = `${100 - percentage}%`;
        outputContainer.style.flex = 'none';

        // Tell Monaco to resize
        window.editor.layout();
    });

    document.addEventListener('mouseup', () => {
        isResizing = false;
        document.body.style.cursor = 'default';
    });
});
