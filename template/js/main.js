import '../css/main.css'

import Turbolinks from 'turbolinks'
import lazy from './lazy.js'

document.addEventListener('turbolinks:load', () => {
  lazy.observe()
  import('./app.js')
})

Turbolinks.start()
