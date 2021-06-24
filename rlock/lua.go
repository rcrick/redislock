package rlock

const (
	luaMutexLock = `
local ttl = tonumber(ARGV[2])
local setRS = redis.call("Set", KEYS[1], ARGV[1], "PX", ttl, "NX")

local rs = {}
rs[1] = 1
rs[2] = redis.call("Get", KEYS[1])
if (setRS == nil or setRS == false) then
    rs[1] = 0
end

return rs

`

	luaKvDelIfExists = `
if redis.call("Get", KEYS[1]) == ARGV[1] then
    return redis.call("Del", KEYS[1])
else
    return 0
end

`
)
