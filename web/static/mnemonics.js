document.addEventListener('DOMContentLoaded', () => {
	let cardLookup = {}
	let easyMotion = false
	let easyMotionTimeout = null
	let easyMotionKey = []

	document.addEventListener('keydown', (e) => {
		if (e.ctrlKey && e.key === ' ') {
			document.activeElement.blur()

			clearTimeout(easyMotionTimeout)
			easyMotion = false
			easyMotionKey = []
			cardLookup = {}

			const cards = document.querySelectorAll('task-card')
			const alphabet = 'abcdefghijklmnopqrstuvwxyz'.split('')

			for (let i = 0; i < cards.length; i++) {
				let letter = alphabet[i % 26]

				if (i >= 26) {
					let j = Math.floor(i / 26) - 1
					letter = alphabet[j] + letter
				}

				cardLookup[letter] = cards[i]
				cards[i].setAttribute("data-mnemonic", letter)
			}

			easyMotion = true
			return
		}

		if (e.key === 'Escape') {
			easyMotion = false
			clearTimeout(easyMotionTimeout)
			easyMotionKey = []
			Object.values(cardLookup).forEach(card => card.removeAttribute('data-mnemonic'))
			cardLookup = {}
			return
		}

		if (easyMotion) {
			easyMotionTimeout = setTimeout(() => {
				if (easyMotionKey.length) {
					const card = cardLookup[easyMotionKey.join('')]

					if (card) {
						card.focus()
					}
				}

				easyMotionTimeout = null
				easyMotion = false
				easyMotionKey = []

				Object.values(cardLookup).forEach(card => card.removeAttribute('data-mnemonic'))
				cardLookup = {}
			}, 250)

			if (['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'].includes(e.key)) {
				easyMotionKey.push(e.key)
				return
			}
		}
	})
})
