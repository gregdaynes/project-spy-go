/* global HTMLElement customElements */

class InfoDialog extends HTMLElement {
  constructor () {
    super()

    this.dialog = this.querySelector('dialog')
    this.closeForm = this.querySelector('form:has(> [type="submit"])')
    this.interactor = this.querySelector('a#show-info')
    this.interactorIcon = this.interactor.querySelector('svg')
  }

  connectedCallback () {

    if (this.dialog.hasAttribute('open')) {
      this.dialog.removeAttribute('open')
      this.dialog.showModal()
    }

    this.replaceWithButton(this.interactor)

    this.eventHotkey = document.addEventListener('keydown', e => {
      if (!(e.ctrlKey && (e.key === '?' || e.key === '/'))) return

      this.dialog.showModal()
    })

    document.addEventListener('mouseup', e => {
      if (e.target !== this.dialog) return

      this.dialog.close()
    })

    document.addEventListener('keydown', e => {
      if (e.ctrlKey && e.key === 'w') {
        e.preventDefault()

        this.dialog.close()
      }
    })

    this.closeForm.addEventListener('click', e => {
      e.preventDefault()

      this.dialog.close()
    })
  }

  disconnectedCallback () {
    document.removeEventListener(this.eventHotkey)
  }

  replaceWithButton (interactor) {
    interactor.classList.add('hidden')
		interactor.setAttribute('aria-hidden', true)
		interactor.setAttribute('tabindex', -1)

    const button = document.createElement('button')
    button.setAttribute('id', 'show-info')
    button.classList.add('icon')
    button.appendChild(this.interactorIcon)
    button.addEventListener('click', () => {
      this.dialog.showModal()
    })

    this.appendChild(button)
  }
}

if (!customElements.get('info-dialog')) {
  customElements.define('info-dialog', InfoDialog)
}

