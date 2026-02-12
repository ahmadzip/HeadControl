const HC = {};
HC.Sidebar = {
    toggle() {
        const sidebar = document.getElementById('sidebar');
        const overlay = document.getElementById('sidebar-overlay');
        if (!sidebar) return;
        sidebar.classList.toggle('open');
        if (overlay) overlay.classList.toggle('open');
    },
    close() {
        const sidebar = document.getElementById('sidebar');
        const overlay = document.getElementById('sidebar-overlay');
        if (sidebar) sidebar.classList.remove('open');
        if (overlay) overlay.classList.remove('open');
    }
};


HC.Theme = {
    set(theme) {
        document.documentElement.setAttribute('data-theme', theme);
        localStorage.setItem('hc-theme', theme);

        const label = document.getElementById('theme-label');
        if (label) {
            label.textContent = theme.charAt(0).toUpperCase() + theme.slice(1);
        }

        document.querySelectorAll('.theme-option').forEach(btn => {
            btn.classList.toggle('active', btn.getAttribute('data-theme') === theme);
        });

        this.closeDropdown();
    },

    toggleDropdown() {
        const menu = document.getElementById('theme-menu');
        if (!menu) return;
        menu.classList.toggle('open');

        if (menu.classList.contains('open')) {
            setTimeout(() => {
                const handler = (e) => {
                    const dropdown = document.getElementById('theme-dropdown');
                    if (dropdown && !dropdown.contains(e.target)) {
                        this.closeDropdown();
                        document.removeEventListener('click', handler);
                    }
                };
                document.addEventListener('click', handler);
            }, 10);
        }
    },

    closeDropdown() {
        const menu = document.getElementById('theme-menu');
        if (menu) menu.classList.remove('open');
    },

    init() {
        const saved = localStorage.getItem('hc-theme');
        if (saved) {
            this.set(saved);
        } else {
            const current = document.documentElement.getAttribute('data-theme') || 'light';
            const label = document.getElementById('theme-label');
            if (label) {
                label.textContent = current.charAt(0).toUpperCase() + current.slice(1);
            }
            document.querySelectorAll('.theme-option').forEach(btn => {
                btn.classList.toggle('active', btn.getAttribute('data-theme') === current);
            });
        }
    }
};


HC.Modal = {
    open(id) {
        const modal = document.getElementById(id);
        if (!modal) return;
        modal.style.display = 'flex';
        if (typeof lucide !== 'undefined') {
            lucide.createIcons();
        }
        setTimeout(() => {
            const input = modal.querySelector('input:not([type="hidden"])');
            if (input) input.focus();
        }, 100);
    },

    close(id) {
        const modal = document.getElementById(id);
        if (!modal) return;
        modal.style.display = 'none';
        const form = modal.querySelector('form');
        if (form) form.reset();
    },
    openRenameUser(id, name) {
        const idEl = document.getElementById('rename-user-id');
        const nameEl = document.getElementById('rename-user-current');
        const inputEl = document.getElementById('rename-user-newname');
        if (idEl) idEl.value = id;
        if (nameEl) nameEl.textContent = name;
        if (inputEl) inputEl.value = name;
        this.open('rename-user-modal');
    },

    openDeleteUser(id, name) {
        const idEl = document.getElementById('delete-user-id');
        const nameEl = document.getElementById('delete-user-name');
        if (idEl) idEl.value = id;
        if (nameEl) nameEl.textContent = name;
        this.open('delete-user-modal');
    },

    openNodeDetail(id) {
        const content = document.getElementById('node-detail-content');
        if (content) {
            content.innerHTML = '<div style="text-align:center; padding:24px;"><span class="spinner"></span></div>';
        }
        this.open('node-detail-modal');

        if (typeof htmx !== 'undefined') {
            htmx.ajax('GET', '/nodes/detail?id=' + id, {
                target: '#node-detail-content',
                swap: 'innerHTML'
            });
        }
    },

    openRenameNode(id, name) {
        const idEl = document.getElementById('rename-node-id');
        const nameEl = document.getElementById('rename-node-current');
        const inputEl = document.getElementById('rename-node-newname');
        if (idEl) idEl.value = id;
        if (nameEl) nameEl.textContent = name;
        if (inputEl) inputEl.value = name;
        this.open('rename-node-modal');
    },

    openExpireNode(id, name) {
        const idEl = document.getElementById('expire-node-id');
        const nameEl = document.getElementById('expire-node-name');
        if (idEl) idEl.value = id;
        if (nameEl) nameEl.textContent = name;
        this.open('expire-node-modal');
    },

    openDeleteNode(id, name) {
        const idEl = document.getElementById('delete-node-id');
        const nameEl = document.getElementById('delete-node-name');
        if (idEl) idEl.value = id;
        if (nameEl) nameEl.textContent = name;
        this.open('delete-node-modal');
    }
};


HC.Toast = {
    init() {
        const container = document.getElementById('toast-container');
        if (!container) return;

        const observer = new MutationObserver(mutations => {
            mutations.forEach(mutation => {
                mutation.addedNodes.forEach(node => {
                    if (node.classList && node.classList.contains('toast')) {
                        setTimeout(() => {
                            node.style.opacity = '0';
                            node.style.transform = 'translateX(100%)';
                            setTimeout(() => node.remove(), 300);
                        }, 4000);
                    }
                });
            });
        });

        observer.observe(container, { childList: true });
    }
};


HC.refreshUsers = function () {
    if (typeof htmx !== 'undefined') {
        htmx.ajax('GET', '/users/table', { target: '.content', swap: 'innerHTML' });
    }
};

HC.refreshNodes = function () {
    if (typeof htmx !== 'undefined') {
        htmx.ajax('GET', '/nodes/table', { target: '.content', swap: 'innerHTML' });
    }
};

document.addEventListener('click', function (e) {
    if (e.target.classList.contains('modal-overlay')) {
        e.target.style.display = 'none';
        const form = e.target.querySelector('form');
        if (form) form.reset();
    }
});


document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') {
        document.querySelectorAll('.modal-overlay').forEach(modal => {
            if (modal.style.display !== 'none') {
                modal.style.display = 'none';
                const form = modal.querySelector('form');
                if (form) form.reset();
            }
        });
        HC.Theme.closeDropdown();
    }
});


document.addEventListener('click', function (e) {
    const link = e.target.closest('.nav-link');
    if (link && window.innerWidth <= 768) {
        HC.Sidebar.close();
    }
});


document.addEventListener('DOMContentLoaded', function () {
    HC.Theme.init();
    HC.Toast.init();
});
