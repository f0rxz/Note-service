document.addEventListener('DOMContentLoaded', function() {
    fetchNotes();
});

function fetchNotes() {
    fetch('/api/notes?page=0')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(notes => {
            if (!notes || notes.length === 0) {
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
                    <button onclick="deleteNote(${note.id})">Delete</button>
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
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: `title=${title}&content=${content}` // Don't include an ID, let the server handle it
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.text;
    })
    .then(note => {
        fetchNotes();
        hideCreateNoteForm();
    })
    .catch(error => console.error('Error:', error));
}

function deleteNote(id) {
    fetch(`/api/notes/${id}`, { method: 'DELETE' })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            fetchNotes();
        })
        .catch(error => console.error('Error:', error));
}
