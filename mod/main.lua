--[[
    SalioMP
    Copyright (C) 2020 Juby210

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
--]]

json = require("mp/json")
waiting = false
firstWaiting = true
data = { level = 0, waiting = false }
lastx = -18
lasty = -18

function drawWatermark()
  drawText(20, 20, 13, "SalioMP v0.2.5")
end

function init()
  drawWatermark()
  writeAll()

  p2 = newObject('player.png', -18, -18)
  setUpdate(15)
end

function update()
  if waiting then
    local lines = lines_from("mp/data2.json")
    local data2 = json.decode(lines[1])
    if data2["level"] == getLevelNumber() or (data2["waiting"] and firstWaiting) then
      firstWaiting = false
      waiting = false
      enablePlayer()
      deleteText(waitingText)
      deleteText(waitingText2)
      killPlayer()
      drawWatermark()
    end
  else
    local lines = lines_from("mp/data2.json")
    local data2 = json.decode(lines[1])
    if data2["level"] < getLevelNumber() then
      waiting = true
      setPlayerX(-18)
      setPlayerY(-18)
      disablePlayer()
      waitingText = drawText(430, 100, 15, "Waiting for other player..")
      waitingText2 = drawText(361, 123, 12, "Salio doesn't support forwarding levels :v")
    end
  end

  writeAll()

  local lines = lines_from("mp/p2.txt")
  if table.getn(lines) == 2 then
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
  writeData()
end

function writeP1()
  x = getPlayerX()
  y = getPlayerY()
  if x ~= lastx or y ~= lasty then
    lastx = x
    lasty = y
    file = io.open("mp/p1.txt", "w")
    file:write(lastx .. "\n" .. lasty)
    file:close()
  end
end

function writeData()
  if data["level"] ~= getLevelNumber() or data["waiting"] ~= waiting then
    data["level"] = getLevelNumber()
    data["waiting"] = waiting
    file = io.open("mp/data.json", "w")
    file:write(json.encode(data))
    file:close()
  end
end

function getKey(id)
end

function onDeath()
end
