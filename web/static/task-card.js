class TaskCard extends HTMLElement {
	static observedAttributes = ["data-mnemonic"];

	constructor () {
		super()

		this.header = this.querySelector('task-header')
		this.body = this.querySelector('task-body')
		this.footer = this.querySelector('task-footer')
		this.actions = this.querySelector('.actions')
	}

	connectedCallback () {
		this.classList.add('enhanced')
		this.setAttribute('tabindex', 0)

		const viewAction = this.actions.querySelector('a[href*="/view/"]')
		this.addEventListener('click', () => window.location.href = viewAction)
		this.addEventListener('keydown', (e) => {
			if (e.key === 'Enter') {
				e.preventDefault()
				window.location.href = viewAction
			}

			if (e.key === 'Backspace') {
				e.preventDefault()
				window.location.href = this.actions.querySelector('a[href*="/delete/"]')
			}

			if (e.key === 'a') {
				e.preventDefault()
				window.location.href = this.actions.querySelector('a[href*="/archive/"]')
			}

			if (e.key === 'Escape') {
				e.preventDefault()
				this.blur()
			}
		})

		if (this.body.innerText.trim() == "") {
			this.body.remove()
		}

		if (this.footer.innerText.trim() == "") {
			this.footer.remove()
		}

		this.actions.querySelectorAll('a').forEach(action => {
			action.setAttribute('tabindex', -1)

			if (action.getAttribute('href').includes('/view/')) {
				action.setAttribute('aria-hidden', true)
				action.classList.add('hidden')
			}
		})
	}

	attributeChangedCallback(name, _oldValue, newValue) {
		if (name === 'data-mnemonic') {
			this.header.querySelectorAll('.mnemonic').forEach(el => el.remove())

			if (newValue !== null) {
				const menmonic = document.createElement('span')
				menmonic.classList.add('mnemonic')
				menmonic.textContent = newValue
				this.header.appendChild(menmonic)
			}
		}
	}
}

if (!customElements.get('task-card')) {
	customElements.define('task-card', TaskCard)
}
