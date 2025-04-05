class TaskLane extends HTMLElement {
	constructor () {
		super()

		this.collapseLabel = this.querySelector('task-lane-header label')
		this.toggle = this.collapseLabel.querySelector("input")
	}

	connectedCallback () {
		// on startup check session storage for lane state
		const storedState = window.sessionStorage.getItem(this.getAttribute('id'))

		if (storedState === 'true') {
			this.toggle.setAttribute('checked', true)
		}

		this.collapseLabel.addEventListener('click', (e) => {
			e.preventDefault()
			this.toggle.toggleAttribute('checked')
			window.sessionStorage.setItem(this.getAttribute('id'), this.toggle.checked)
		})
	}
}

if (!customElements.get('task-lane')) {
	customElements.define('task-lane', TaskLane)
}

