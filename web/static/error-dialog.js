/* global HTMLElement customElements */

class ErrorDialog extends HTMLElement {
  constructor () {
    super()

    this.dialog = this.querySelector('dialog')
    this.closeForm = this.querySelector('form:has(> [type="submit"])')
  }

  connectedCallback () {
    this.dialog.removeAttribute('open')
    this.dialog.showModal()

    document.addEventListener('mouseup', e => {
      if (e.target !== this.dialog) return

      this.dialog.close()
    })

    this.closeForm.addEventListener('click', e => {
      e.preventDefault()

      this.dialog.close()
    })
  }
}

if (!customElements.get('error-dialog')) {
  customElements.define('error-dialog', ErrorDialog)
}

