function init()
  file = io.open("p1.txt", "w")
  file:write(getPlayerX() .. "\n" .. getPlayerY())
  file:close()
  p2 = newObject('player.png', -18, -18)
  setUpdate(2) -- probably not working function
end

function update()
  file = io.open("p1.txt", "w")
  file:write(getPlayerX() .. "\n" .. getPlayerY())
  file:close()

  local lines = lines_from("p2.txt")
  setObjectX(p2, lines[1])
  setObjectY(p2, lines[2])
end

function getKey(id)
end

function onDeath()
end

function lines_from(file)
  lines = {}
  for line in io.lines(file) do
    lines[#lines + 1] = line
  end
  return lines
end
