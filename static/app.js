(() => {
	const state = { limit: 20, offset: 0, status: '' };
	const els = {
		form: document.getElementById('create-form'),
		title: document.getElementById('title'),
		description: document.getElementById('description'),
		due: document.getElementById('due'),
		status: document.getElementById('status'),
		filterStatus: document.getElementById('filter-status'),
		tasks: document.getElementById('tasks'),
		prev: document.getElementById('prev'),
		next: document.getElementById('next'),
		pageInfo: document.getElementById('page-info'),
	};

	async function fetchTasks() {
		const params = new URLSearchParams({ limit: String(state.limit), offset: String(state.offset) });
		if (state.status) params.set('status', state.status);
		const res = await fetch(`/api/tasks?${params.toString()}`);
		if (!res.ok) throw new Error('Failed fetching tasks');
		const json = await res.json();
		return json.data || [];
	}

	function renderTasks(tasks) {
		els.tasks.innerHTML = '';
		if (!Array.isArray(tasks) || tasks.length === 0) {
			els.tasks.innerHTML = '<li class="task">No tasks</li>';
			els.prev.disabled = state.offset === 0;
			els.pageInfo.textContent = `Page ${Math.floor(state.offset / state.limit) + 1}`;
			return;
		}
		tasks.forEach(task => {
			const li = document.createElement('li');
			li.className = 'task';
			const title = document.createElement('div');
			title.innerHTML = `<strong>${escapeHtml(task.title)}</strong><div class="meta">#${task.id} • ${escapeHtml(task.status)}${task.due_date ? ' • due ' + new Date(task.due_date).toLocaleDateString() : ''}</div><div>${escapeHtml(task.description || '')}</div>`;
			const setStatus = document.createElement('button');
			setStatus.textContent = task.status === 'completed' ? 'Mark pending' : 'Mark completed';
			setStatus.addEventListener('click', () => updateTask(task.id, { status: task.status === 'completed' ? 'pending' : 'completed' }));
			const del = document.createElement('button');
			del.textContent = 'Delete';
			del.addEventListener('click', () => deleteTask(task.id));
			const actions = document.createElement('div');
			actions.className = 'task-actions';
			actions.append(setStatus, del);
			li.append(title, actions);
			els.tasks.appendChild(li);
		});
		els.prev.disabled = state.offset === 0;
		els.pageInfo.textContent = `Page ${Math.floor(state.offset / state.limit) + 1}`;
	}

	function escapeHtml(s) {
		return String(s).replace(/[&<>"]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
	}

	async function createTask(e) {
		e.preventDefault();
		const payload = {
			title: els.title.value.trim(),
			description: els.description.value.trim(),
			status: els.status.value,
		};
		if (els.due.value) payload.due_date = new Date(els.due.value).toISOString();
		const res = await fetch('/api/tasks', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
		if (!res.ok) { alert('Failed to create'); return; }
		els.form.reset();
		await refresh();
	}

	async function updateTask(id, patch) {
		const res = await fetch(`/api/tasks/${id}`, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(patch) });
		if (!res.ok) { alert('Failed to update'); return; }
		await refresh();
	}

	async function deleteTask(id) {
		if (!confirm('Delete task?')) return;
		const res = await fetch(`/api/tasks/${id}`, { method: 'DELETE' });
		if (!res.ok) { alert('Failed to delete'); return; }
		await refresh();
	}

	async function refresh() {
		try {
			const tasks = await fetchTasks();
			renderTasks(tasks);
		} catch (e) {
			console.error(e);
			alert('Error loading tasks');
		}
	}

	els.form.addEventListener('submit', createTask);
	els.filterStatus.addEventListener('change', () => { state.status = els.filterStatus.value; state.offset = 0; refresh(); });
	els.prev.addEventListener('click', () => { state.offset = Math.max(0, state.offset - state.limit); refresh(); });
	els.next.addEventListener('click', () => { state.offset += state.limit; refresh(); });

	refresh();
})();