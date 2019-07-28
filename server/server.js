let p1 = null
let p2 = null

const io = require('socket.io')()
io.on('connection', socket => {
  let p11 = false
  if(p1 == null) {
    p1 = socket
    p11 = true
  } else p2 = socket;

  socket.on('playermove', data => {
    console.log(`${p11 ? 'p1' : 'p2'}; ${data.x} | ${data.y}`)
    if(p11) {
      if(p2 != null) p2.emit('playermove', data)
    } else {
      if(p2 != null) p1.emit('playermove', data)
    }
  })

  socket.on('disconnect', () => {
    if(p11) {
      p1 = null
      if(p2 != null) p2.emit('playermove', { x: -18, y: -18 })
    } else {
      p2 = null
      if(p1 != null) p1.emit('playermove', { x: -18, y: -18 })
    }
  })
})
io.listen(2410)
