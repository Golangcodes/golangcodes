// A central interface for all storage operations, making it easy to swap 
// localStorage for a real database (e.g., fetch API calls) in the future.
window.StorageAPI = {
    // ---- Theme Settings ----
    getTheme: function () {
        return localStorage.getItem('theme') || 'light';
    },
    setTheme: function (theme) {
        localStorage.setItem('theme', theme);
    },

    // ---- Gamification & XP ----
    getUserXP: function () {
        return parseInt(localStorage.getItem('userXP') || '0', 10);
    },
    addXP: function (amount) {
        let xp = this.getUserXP() + amount;
        localStorage.setItem('userXP', xp);
        return xp;
    },

    // ---- Read Tracking ----
    getReadCounts: function () {
        let counts = {};
        try {
            counts = JSON.parse(localStorage.getItem('readCounts')) || {};
        } catch (e) { }

        // Backward compatibility for existing users with visitedPages
        try {
            let visitedPages = JSON.parse(localStorage.getItem('visitedPages')) || [];
            let updated = false;
            visitedPages.forEach(p => {
                if (!counts[p]) {
                    counts[p] = 1;
                    updated = true;
                }
            });
            if (updated) {
                this.setReadCounts(counts);
                localStorage.removeItem('visitedPages');
            }
        } catch (e) { }

        return counts;
    },
    setReadCounts: function (counts) {
        localStorage.setItem('readCounts', JSON.stringify(counts));
    },
    getLastVisitedPage: function () {
        return localStorage.getItem('lastVisitedPage');
    },
    setLastVisitedPage: function (path) {
        localStorage.setItem('lastVisitedPage', path);
    },
    trackPageVisit: function (path, onNewVisit, onRevisit) {
        let lastVisited = this.getLastVisitedPage();
        let counts = this.getReadCounts();

        // Prevent refresh-spamming XP tracking
        if (lastVisited !== path) {
            if (!counts[path]) {
                counts[path] = 0;
                if (onNewVisit) onNewVisit();
            } else {
                if (onRevisit) onRevisit();
            }
            counts[path]++;
            this.setReadCounts(counts);
            this.setLastVisitedPage(path);
        }
    },

    // ---- Code Editor Content ----
    getCode: function (path) {
        return localStorage.getItem('code_' + path);
    },
    saveCode: function (path, code) {
        localStorage.setItem('code_' + path, code);
    },
    removeCode: function (path) {
        localStorage.removeItem('code_' + path);
    },

    // ---- Lesson Code Snippets Edits ----
    getSnippet: function (path, index) {
        return localStorage.getItem(`snippet_${path}_${index}`);
    },
    saveSnippet: function (path, index, code) {
        localStorage.setItem(`snippet_${path}_${index}`, code);
    },
    removeSnippet: function (path, index) {
        localStorage.removeItem(`snippet_${path}_${index}`);
    }
};
