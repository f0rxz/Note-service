document.addEventListener('DOMContentLoaded', function() {
	fetchNotes();
});

function fetchNotes() {
	fetch('/api/notes')
		.then(response => {
			if (!response.ok) {
				throw new Error('Network response was not ok');
			}
			return response.json();
		})
		.then(notes => {
			if (!notes) {
				throw new Error('No notes found');
			}
			const noteList = document.getElementById('note-list');
			noteList.innerHTML = '';
			notes.forEach(note => {
				const noteItem = document.createElement('div');
				noteItem.className = 'note-item';
				noteItem.innerHTML = `
					<h3>${note.title}</h3>
					<p>${note.content}</p>
					<button onclick="deleteNote('${note.id}')">Delete</button>
				`;
				noteList.appendChild(noteItem);
			});
		})
		.catch(error => {
			console.error('Fetch error:', error);
			const noteList = document.getElementById('note-list');
			noteList.innerHTML = `<p>Error fetching notes: ${error.message}</p>`;
		});
}

function showCreateNoteForm() {
	document.getElementById('note-form').style.display = 'block';
}

function hideCreateNoteForm() {
	document.getElementById('note-form').style.display = 'none';
}

function createNote() {
	const title = document.getElementById('note-title').value;
	const content = document.getElementById('note-content').value;

	fetch('/api/notes', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ id: Date.now().toString(), title, content }) // Include an ID
	})
	.then(response => {
		if (!response.ok) {
			throw new Error('Network response was not ok');
		}
		return response.json();
	})
	.then(note => {
		fetchNotes();
		hideCreateNoteForm();
	})
	.catch(error => console.error('Error:', error));
}

function deleteNote(id) {
	fetch(`/api/notes/${id}`, { method: 'DELETE' })
		.then(() => fetchNotes())
		.catch(error => console.error('Error:', error));
}
