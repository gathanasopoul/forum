document.addEventListener('DOMContentLoaded', () => {
    const themeToggleButton = document.getElementById('theme-toggle-btn');
    
    // Function to get the saved theme preference or system default
    const getSavedTheme = () => {
        const savedTheme = localStorage.getItem('forum-theme');
        if (savedTheme) {
            return savedTheme;
        }
        
        // System preference default
        const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        return systemPrefersDark ? 'dark' : 'light';
    };

    // Apply selected theme
    const applyTheme = (theme) => {
        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('forum-theme', theme);
    };

    // Init theme
    const currentTheme = getSavedTheme();
    applyTheme(currentTheme);

    // Toggle theme button click listener
    if (themeToggleButton) {
        themeToggleButton.addEventListener('click', () => {
            const activeTheme = document.documentElement.getAttribute('data-theme');
            const newTheme = activeTheme === 'dark' ? 'light' : 'dark';
            applyTheme(newTheme);
        });
    }
});
