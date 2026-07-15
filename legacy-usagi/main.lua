-- typeinc — incremental typing game
local words = require("words")

local WORD_SCALE = 1
local ERROR_FLASH_TIME = 0.35

-- Scene state: "menu" | "game"
local scene = "menu"

-- Menu state
local menu_options = { { label = "ESPANOL", lang = "es" }, { label = "ENGLISH", lang = "en" } }
local menu_index = 1

-- Game state
local word_list = nil
local word = ""
local pos = 0 -- letters typed correctly so far
local score = 0
local error_t = 0 -- red flash timer

-- Raw key -> letter lookup, built in _init
local keymap = {}

function _config()
  return {
    name = "typeinc",
    -- native 1080p render target: text draws 1:1 at fullscreen, no upscale blur
    game_width = 1920,
    game_height = 1080,
    pause_menu = false, -- letters like P and keys like Enter/Escape must reach the game
  }
end

function _init()
  for i = 0, 25 do
    local letter = string.char(97 + i)
    keymap[input["KEY_" .. letter:upper()]] = letter
  end
  -- v1.1.1 has no fullscreen _config option; the setting persists on disk,
  -- so only toggle when not already fullscreen
  if not usagi.is_fullscreen() then
    usagi.toggle_fullscreen()
  end
end

local function next_word()
  local candidate = word_list[math.random(#word_list)]
  while candidate == word and #word_list > 1 do
    candidate = word_list[math.random(#word_list)]
  end
  word = candidate
  pos = 0
end

local function start_game(lang)
  word_list = words[lang]
  score = 0
  error_t = 0
  word = ""
  next_word()
  scene = "game"
end

-- Menu ----------------------------------------------------------------------

local function update_menu(_dt)
  if input.key_pressed(input.KEY_UP) then
    menu_index = menu_index == 1 and #menu_options or menu_index - 1
  end
  if input.key_pressed(input.KEY_DOWN) then
    menu_index = menu_index == #menu_options and 1 or menu_index + 1
  end
  if input.key_pressed(input.KEY_ENTER) then
    start_game(menu_options[menu_index].lang)
  end
end

local function draw_menu()
  local title = "TYPEINC"
  local tw = usagi.measure_text(title)
  gfx.text_ex(title, (usagi.GAME_W - tw * 2) / 2, usagi.GAME_H * 0.28, 2, 0, gfx.COLOR_WHITE, 1.0)

  for i, option in ipairs(menu_options) do
    local selected = i == menu_index
    local label = selected and ("> " .. option.label) or option.label
    local w = usagi.measure_text(label)
    local y = usagi.GAME_H * 0.5 + (i - 1) * 48
    local color = selected and gfx.COLOR_WHITE or gfx.COLOR_LIGHT_GRAY
    local alpha = selected and 1.0 or 0.4
    gfx.text(label, (usagi.GAME_W - w) / 2, y, color, alpha)
  end

  local hint = "UP/DOWN + ENTER"
  local hw = usagi.measure_text(hint)
  gfx.text(hint, (usagi.GAME_W - hw) / 2, usagi.GAME_H - 80, gfx.COLOR_LIGHT_GRAY, 0.4)
end

-- Game ----------------------------------------------------------------------

local function update_game(dt)
  if input.key_pressed(input.KEY_ESCAPE) then
    scene = "menu"
    return
  end

  if error_t > 0 then
    error_t = error_t - dt
  end

  for key, letter in pairs(keymap) do
    if input.key_pressed(key) then
      if letter == word:sub(pos + 1, pos + 1) then
        pos = pos + 1
        if pos == #word then
          score = score + #word
          next_word()
        end
      else
        -- wrong key: flash red and send the cursor back to the start
        error_t = ERROR_FLASH_TIME
        pos = 0
      end
    end
  end
end

local function draw_game()
  local w, h = usagi.measure_text(word)
  local x = (usagi.GAME_W - w * WORD_SCALE) / 2
  local y = (usagi.GAME_H - h * WORD_SCALE) / 2

  if error_t > 0 then
    gfx.rect_fill(x - 12, y - 8, w * WORD_SCALE + 24, h * WORD_SCALE + 16, gfx.COLOR_RED, 0.6)
  end

  -- ghost word, then the correctly typed prefix drawn on top (monospace font
  -- keeps both perfectly aligned)
  gfx.text_ex(word, x, y, WORD_SCALE, 0, gfx.COLOR_LIGHT_GRAY, 0.35)
  if pos > 0 then
    gfx.text_ex(word:sub(1, pos), x, y, WORD_SCALE, 0, gfx.COLOR_WHITE, 1.0)
  end

  -- cursor: thin bar in front of the current letter
  local prefix_w = pos > 0 and usagi.measure_text(word:sub(1, pos)) or 0
  gfx.rect_fill(x + prefix_w * WORD_SCALE, y, 2, h * WORD_SCALE, gfx.COLOR_WHITE, 0.9)

  -- score, gold, top-right corner
  local score_text = tostring(score)
  local sw = usagi.measure_text(score_text)
  gfx.text(score_text, usagi.GAME_W - sw - 24, 16, gfx.COLOR_YELLOW, 1.0)
end

-- Loop ----------------------------------------------------------------------

function _update(dt)
  if scene == "menu" then
    update_menu(dt)
  else
    update_game(dt)
  end
end

function _draw(_dt)
  gfx.clear(gfx.COLOR_DARK_BLUE)
  if scene == "menu" then
    draw_menu()
  else
    draw_game()
  end
end
