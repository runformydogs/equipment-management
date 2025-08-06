const API_BASE_URL = 'http://localhost:8080';

async function fetchWithAuth(url, options = {}) {
    const token = localStorage.getItem('authToken');

    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}${url}`, {
        ...options,
        headers
    });

    if (response.status === 204) {
        return null;
    }

    if (!response.ok) {
        let errorMsg = 'Ошибка запроса';

        try {
            const errorBody = await response.json();
            if (errorBody.error) {
                errorMsg = errorBody.error;
            }
        } catch (e) {
            const text = await response.text();
            if (text) {
                errorMsg = text;
            }
        }

        throw new Error(errorMsg);
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
        return response.json();
    }

    return response.text();
}

window.api = {
    getFullTree: () => fetchWithAuth('/network-nodes/tree'),

    getAllNodes: () => fetchWithAuth('/network-nodes'),
    getNode: (id) => fetchWithAuth(`/network-nodes/${id}`),
    createNode: (data) => fetchWithAuth('/network-nodes', {
        method: 'POST',
        body: JSON.stringify(data)
    }),
    updateNode: (id, data) => fetchWithAuth(`/network-nodes/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    }),
    deleteNode: (id) => fetchWithAuth(`/network-nodes/${id}`, {
        method: 'DELETE'
    }),

    getAllDevices: () => fetchWithAuth('/devices'),
    getDevice: (id) => fetchWithAuth(`/devices/${id}`),
    createDevice: (data) => fetchWithAuth('/devices', {
        method: 'POST',
        body: JSON.stringify(data)
    }),
    updateDevice: (id, data) => fetchWithAuth(`/devices/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    }),
    deleteDevice: (id) => fetchWithAuth(`/devices/${id}`, {
        method: 'DELETE'
    })
};