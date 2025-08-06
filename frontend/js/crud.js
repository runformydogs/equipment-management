document.addEventListener('DOMContentLoaded', () => {
    const userRole = localStorage.getItem('userRole');

    if (userRole !== 'admin') {
        document.querySelectorAll('.btn:not(.logout-btn)').forEach(btn => {
            btn.disabled = true;
            btn.title = 'Требуются права администратора';
        });
        return;
    }

    initNodeSelects();
    initDeviceSelects();

    document.getElementById('add-node-form').addEventListener('submit', addNode);
    document.getElementById('edit-node-form').addEventListener('submit', updateNode);
    document.getElementById('delete-node-form').addEventListener('submit', deleteNode);

    document.getElementById('add-device-form').addEventListener('submit', addDevice);
    document.getElementById('edit-device-form').addEventListener('submit', updateDevice);
    document.getElementById('delete-device-form').addEventListener('submit', deleteDevice);

    document.getElementById('edit-node-id').addEventListener('change', loadNodeDetails);
    document.getElementById('edit-device-id').addEventListener('change', loadDeviceDetails);
});

async function initNodeSelects() {
    try {
        const nodes = await api.getAllNodes();

        const rootOption = {
            id: null,
            name: "(Корневой узел)"
        };

        populateSelect('node-parent', [rootOption, ...nodes]);
        populateSelect('edit-node-id', nodes);
        populateSelect('edit-node-parent', [rootOption, ...nodes]);
        populateSelect('delete-node-id', nodes);
        populateSelect('device-node', [rootOption, ...nodes]);
        populateSelect('edit-device-node', [rootOption, ...nodes]);
    } catch (error) {
        console.error('Error loading nodes:', error);
        alert('Не удалось загрузить список узлов');
    }
}

async function initDeviceSelects() {
    try {
        const devices = await api.getAllDevices();

        const deviceOptions = devices.map(device => ({
            id: device.id,
            name: `${device.type}: ${device.model} (${device.serial})`
        }));

        populateSelect('edit-device-id', deviceOptions);
        populateSelect('delete-device-id', deviceOptions);
    } catch (error) {
        console.error('Error loading devices:', error);
        alert('Не удалось загрузить список устройств');
    }
}

function populateSelect(selectId, items) {
    const select = document.getElementById(selectId);
    if (!select) return;

    const currentValue = select.value;

    const placeholder = select.options[0]?.text.includes('--') ? select.options[0] : null;
    select.innerHTML = '';
    if (placeholder) select.appendChild(placeholder);

    items.forEach(item => {
        const option = document.createElement('option');
        option.value = item.id;
        option.textContent = item.name || `Устройство ${item.id}`;
        select.appendChild(option);
    });

    if (currentValue && Array.from(select.options).some(opt => opt.value === currentValue)) {
        select.value = currentValue;
    }
}

async function loadNodeDetails() {
    const nodeId = document.getElementById('edit-node-id').value;
    if (!nodeId) return;

    try {
        const node = await api.getNode(nodeId);
        document.getElementById('edit-node-name').value = node.name;
        document.getElementById('edit-node-description').value = node.description || '';
        document.getElementById('edit-node-parent').value = node.parent_id || '';
    } catch (error) {
        console.error('Error loading node details:', error);
        alert('Не удалось загрузить данные узла');
    }
}

async function loadDeviceDetails() {
    const deviceId = document.getElementById('edit-device-id').value;
    if (!deviceId) return;

    try {
        const device = await api.getDevice(deviceId);
        if (!device) {
            throw new Error('Устройство не найдено');
        }

        const setValue = (id, value) => {
            const el = document.getElementById(id);
            if (el) el.value = value || '';
        };

        setValue('edit-device-type', device.type);
        setValue('edit-device-vendor', device.vendor);
        setValue('edit-device-model', device.model);
        setValue('edit-device-serial', device.serial);
        setValue('edit-device-location', device.location);
        setValue('edit-device-node', device.network_node_id);
    } catch (error) {
        console.error('Error loading device details:', error);
        alert('Не удалось загрузить данные устройства: ' + error.message);
    }
}

async function addNode(e) {
    e.preventDefault();

    const nodeData = {
        name: document.getElementById('node-name').value,
        description: document.getElementById('node-description').value,
        parent_id: document.getElementById('node-parent').value || null
    };

    if (nodeData.parent_id) {
        nodeData.parent_id = parseInt(nodeData.parent_id);
    }

    try {
        await api.createNode(nodeData);
        alert('Узел успешно добавлен');
        initNodeSelects();
        refreshTree();
        e.target.reset();
    } catch (error) {
        console.error('Error adding node:', error);
        alert('Ошибка при добавлении узла: ' + error.message);
    }
}

async function updateNode(e) {
    e.preventDefault();

    const nodeId = document.getElementById('edit-node-id').value;
    if (!nodeId) {
        alert('Выберите узел для изменения');
        return;
    }

    const nodeData = {
        name: document.getElementById('edit-node-name').value,
        description: document.getElementById('edit-node-description').value,
        parent_id: document.getElementById('edit-node-parent').value || null
    };

    if (nodeData.parent_id) {
        nodeData.parent_id = parseInt(nodeData.parent_id);
    }

    try {
        await api.updateNode(nodeId, nodeData);
        alert('Узел успешно обновлен');
        initNodeSelects();
        refreshTree();
        document.getElementById('edit-node-id').value = '';
        e.target.reset();
    } catch (error) {
        console.error('Error updating node:', error);
        alert('Ошибка при обновлении узла: ' + error.message);
    }
}

async function deleteNode(e) {
    e.preventDefault();

    const nodeId = document.getElementById('delete-node-id').value;
    if (!nodeId) {
        alert('Выберите узел для удаления');
        return;
    }

    if (!confirm('Вы уверены, что хотите удалить этот узел?')) return;

    try {
        await api.deleteNode(nodeId);
        alert('Узел успешно удален');
        initNodeSelects();
        refreshTree();
        e.target.reset();
    } catch (error) {
        console.error('Error deleting node:', error);
        alert('Ошибка при удалении узла: ' + error.message);
    }
}

async function addDevice(e) {
    e.preventDefault();

    const deviceData = {
        type: document.getElementById('device-type').value,
        vendor: document.getElementById('device-vendor').value,
        model: document.getElementById('device-model').value,
        serial: document.getElementById('device-serial').value,
        location: document.getElementById('device-location').value,
        network_node_id: document.getElementById('device-node').value || null
    };

    if (deviceData.network_node_id) {
        deviceData.network_node_id = parseInt(deviceData.network_node_id);
    }

    try {
        await api.createDevice(deviceData);
        alert('Устройство успешно добавлено');
        initDeviceSelects();
        refreshTree();
        e.target.reset();
    } catch (error) {
        console.error('Error adding device:', error);
        alert('Ошибка при добавлении устройства: ' + error.message);
    }
}

async function updateDevice(e) {
    e.preventDefault();

    const deviceId = document.getElementById('edit-device-id').value;
    if (!deviceId) {
        alert('Выберите устройство для изменения');
        return;
    }

    const deviceData = {
        type: document.getElementById('edit-device-type').value,
        vendor: document.getElementById('edit-device-vendor').value,
        model: document.getElementById('edit-device-model').value,
        serial: document.getElementById('edit-device-serial').value,
        location: document.getElementById('edit-device-location').value,
        network_node_id: document.getElementById('edit-device-node').value || null
    };

    if (deviceData.network_node_id) {
        deviceData.network_node_id = parseInt(deviceData.network_node_id);
    }

    try {
        await api.updateDevice(deviceId, deviceData);
        alert('Устройство успешно обновлено');
        initDeviceSelects();
        refreshTree();
        document.getElementById('edit-device-id').value = '';
        e.target.reset();
    } catch (error) {
        console.error('Error updating device:', error);
        alert('Ошибка при обновлении устройства: ' + error.message);
    }
}

async function deleteDevice(e) {
    e.preventDefault();

    const deviceId = document.getElementById('delete-device-id').value;
    if (!deviceId) {
        alert('Выберите устройство для удаления');
        return;
    }

    if (!confirm('Вы уверены, что хотите удалить это устройство?')) return;

    try {
        await api.deleteDevice(deviceId);
        alert('Устройство успешно удалено');
        initDeviceSelects();
        refreshTree();
        e.target.reset();
    } catch (error) {
        console.error('Error deleting device:', error);
        alert('Ошибка при удалении устройства: ' + error.message);
    }
}

async function refreshTree() {
    try {
        const data = await api.getFullTree();
        currentTreeData = data.tree;

        const treeContainer = document.getElementById('tree-container');
        treeContainer.innerHTML = '';
        window.renderTree(currentTreeData);
    } catch (error) {
        console.error('Error refreshing tree:', error);
    }
}