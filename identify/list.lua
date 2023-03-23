#!/usr/bin/env resty

--[[

This script takes the "schemas.json" file, as generated by "scrape2.lua" and then analyses
the properties of the schema's.

]]
local json_decode = require("cjson").decode
local readfile = require("pl.utils").readfile


local filename = "./schemas.json"
local data = assert(readfile(filename))
data = json_decode(data)

local function minlen(l, ...)
  local r = table.concat({...})
  if #r < l then
    r = r .. (" "):rep(l-#r)
  end
  return r
end

local function log(...)
  io.stdout:write(...)
end


-- table to track which entities can have nested ones in a deckfile
-- key = tablename, value = array of possible nested entities
local can_have_nested_entities = {}
local function add_nested_entity(table_having_nested_entity, name_of_nested_entity)
  local t = can_have_nested_entities[table_having_nested_entity]
  if not t then
    t = {}
    can_have_nested_entities[table_having_nested_entity] = t
  end
  t[#t+1] = name_of_nested_entity
end

-- count foreign references
for name, schema in pairs(data) do
  -- foreign relation ships
  local foreign_key_count = 0
  local foreign_required_count = 0
  local foreign_required_ref -- only a single value
  local foreign_optional_refs = {}
  for fieldname, field_definition in pairs(schema.fields) do
    if field_definition.type == "foreign" then
      foreign_key_count = foreign_key_count + 1
      if field_definition.required then
        foreign_required_count = foreign_required_count + 1
        foreign_required_ref = field_definition.reference
      else
        foreign_optional_refs[#foreign_optional_refs+1] = field_definition.reference
      end
    end
  end

  -- Deck file: a nested entity can only have 1 foreign relation ship specified, to the parent
  -- it is nested in. This means an entity can only be nested in a deck file if;
  --  * it has a foreign relationships to the parent (this is obvious)
  --  * of all the foreign relationships a maximum of 1 is required (if more would be required,
  --    then they all would have to be specified, and then it is no longer allowed to be nested
  --    according to the deck format)
  -- so if it can be nested then it can ONLY be nested in the entity referenced by the
  -- REQUIRED property. If there is no REQUIRED reference, then it can be in any of the references
  -- (as long as the other references are not set).
  local can_be_nested = (foreign_key_count > 0) and (foreign_required_count <= 1)
  if can_be_nested then
    if foreign_required_ref then
      -- there is a REQUIRED ref, so that's the only place where the entity can be nested
      -- log(tostring(can_be_nested), ": ", foreign_required_ref)
      add_nested_entity(foreign_required_ref, name)
    else
      -- no required ref, so it can be nested in ANY of the optional references
      -- log(tostring(can_be_nested), ": ", table.concat( foreign_optional_refs, ", "))
      for i, ref in ipairs(foreign_optional_refs) do
        add_nested_entity(ref, name)
      end
    end
  -- else
  --   log(tostring(can_be_nested))
  end
  -- log("\n")
end

-- print("title: ", require("pl.pretty").write(can_have_nested_entities))
-- os.exit()

-- sort output
local sorted = {}
for name in pairs(data) do
  sorted[#sorted+1] = name
end
table.sort(sorted)

print("title: ", require("pl.pretty").write(data.keyauth_credentials))
for i, name in ipairs(sorted) do
  local schema = data[name]

  -- name
  log(minlen(50,name))

  -- id field
  if schema.fields and
     schema.fields.id and
     schema.fields.id.uuid == true then
      if schema.fields.id.auto == true then
        log(minlen(10, "auto-id"))
      else
        log(minlen(10, "id"))
      end
  else
    log(minlen(10, "no-id"))
  end

  -- primary-keys
  log(minlen(60, table.concat(schema.primary_key or {}, ", ")))

  -- unique fields
  local uf = {}
  for fieldname, field_definition in pairs(schema.fields) do
    if field_definition.unique then
      uf[#uf + 1] = fieldname
    end
  end
  log(minlen(40, table.concat(uf, ", ")))

  -- cache-key
  log(minlen(50, table.concat(schema.cache_key or {}, ", ")))

  -- endpoint-key
  log(minlen(50, schema.endpoint_key or ""))

  -- nested entities
  log(minlen(50, table.concat(can_have_nested_entities[name] or {}, ", ")))


  log("\n")
end
