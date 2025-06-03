class TaskLane extends HTMLElement {
	constructor () {
		super()

		this.body = this.querySelector('task-lane-body')
		this.collapseLabel = this.querySelector('task-lane-header label')
		this.toggle = this.collapseLabel.querySelector("input")
	}

	connectedCallback () {
		// Maintain collapsed state
		const storedState = window.sessionStorage.getItem(this.getAttribute('id'))

		if (storedState === 'true') {
			this.toggle.setAttribute('checked', true)
		}

		this.collapseLabel.addEventListener('click', (e) => {
			e.preventDefault()
			this.toggle.toggleAttribute('checked')
			window.sessionStorage.setItem(this.getAttribute('id'), this.toggle.checked)
		})

		// Maintain scroll position
		this.body.addEventListener('scroll', (e) => {
			window.sessionStorage.setItem(this.getAttribute('id'), e.target.scrollTop)
		})

		this.body.scrollTop = window.sessionStorage.getItem(this.getAttribute('id'))

		this.addEventListener('drop', this)
		this.addEventListener('dragover', this)
		this.addEventListener('dragleave', this)
	}

	async handleEvent (e) {
		switch(e.type) {
			case 'drop':
				e.preventDefault()
				e.currentTarget.classList.remove('dragover')

				const dt = e.dataTransfer
				const files = dt.files

				if (files.length) {
					const fd = new FormData()
					fd.append('lane', this.getAttribute('data-name'))

					for (const file of files) {
						fd.append('files[]', file)
					}

					const response = await fetch("/attach-file", {
						method: "POST",
						body: fd,
					});

					window.location.href = response.url;
				}

				if (dt.getData("text") !== "") {
					//  could we post to an endpoint and have it render the task view with content filled out?
					const lane = this.getAttribute('data-name')
					const content = dt.getData("text")

					const form = document.createElement("form")
					form.setAttribute('method', 'POST')
					form.setAttribute('action', '/new/')

					const inputLane = document.createElement('input')
					inputLane.setAttribute('type', 'hidden')
					inputLane.setAttribute('name', 'lane')
					inputLane.setAttribute('value', lane)
					form.appendChild(inputLane)

					const inputContent = document.createElement('input')
					inputContent.setAttribute('type', 'hidden')
					inputContent.setAttribute('name', 'content')
					inputContent.setAttribute('value', content)
					form.appendChild(inputContent)

					document.body.appendChild(form);
					form.submit();
				}

				break
			case 'dragover':
				e.preventDefault()
				if (e.currentTarget.classList.contains('dragover')) return
				e.currentTarget.classList.add('dragover')
				break
			case 'dragleave':
				e.currentTarget.classList.remove('dragover')
				break
			default:
				console.info('no handler assigned to event')
		}
	}
}

if (!customElements.get('task-lane')) {
	customElements.define('task-lane', TaskLane)
}

