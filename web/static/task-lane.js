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
	}
}

if (!customElements.get('task-lane')) {
	customElements.define('task-lane', TaskLane)
}

