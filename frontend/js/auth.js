document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    const logoutBtn = document.getElementById('logout-btn');
    const errorMessage = document.getElementById('error-message');

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const login = document.getElementById('login').value;
            const password = document.getElementById('password').value;

            try {
                const response = await fetch(`http://localhost:8080/login`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ login, password })
                });

                if (!response.ok) {
                    const error = await response.json();
                    throw new Error(error.error || 'Ошибка авторизации');
                }

                const { token, role } = await response.json();
                localStorage.setItem('authToken', token);
                localStorage.setItem('userRole', role);

                window.location.href = 'dashboard.html';
            } catch (error) {
                errorMessage.textContent = error.message;
                errorMessage.style.display = 'block';
            }
        });
    }

    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            localStorage.removeItem('authToken');
            localStorage.removeItem('userRole');
            window.location.href = 'index.html';
        });
    }

    if (window.location.pathname.includes('dashboard.html')) {
        const token = localStorage.getItem('authToken');
        if (!token) {
            window.location.href = 'index.html';
        }
    }
});
