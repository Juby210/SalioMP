_waiting = true

function init()
  drawText(20, 20, 13, "SalioMP v0.2.1")
  writeAll()

  setPlayerX(-18)
  setPlayerY(-18)
  disablePlayer()
  waitingText = drawText(430, 100, 15, "Waiting for other player..")

  p2 = newObject('player.png', -18, -18)
  setUpdate(10)
end

function update()
  writeAll()

  local lines = lines_from("mp/p2.txt")
  if table.getn(lines) == 2 then
    if _waiting and tonumber(lines[1]) > -18 then
      _waiting = false
      enablePlayer()
      deleteText(waitingText)
      killPlayer()
    end
    setObjectX(p2, lines[1])
    setObjectY(p2, lines[2])
  end
end

function lines_from(file)
  lines = {}
  for line in io.lines(file) do
    lines[#lines + 1] = line
  end
  return lines
end

function writeAll()
  writeP1()
  writeLevel()
end

function writeP1()
  lastx = getPlayerX()
  lasty = getPlayerY()
  file = io.open("mp/p1.txt", "w")
  file:write(lastx .. "\n" .. lasty)
  file:close()
end

function writeLevel()
  level = getLevelNumber()
  file = io.open("mp/level.txt", "w")
  file:write(level)
  file:close()
end

function getKey(id)
end

function onDeath()
end
