-- 流量控制脚本, 设置一个TTL = 1的HASH对象, 有2个entry, 分别是当前请求次数current和当前最大请求次数max;
-- 当current + 1 > max时则返回0, 当HASH对象不存在或者current + 1 <= max时返回1;
-- 返回0表示达到限流, 返回1表示OK

local HASH_KEY = "goate:ratelimiter"

local mapResult = redis.call("hgetall", HASH_KEY)
-- 判断hgetall是否为空
if nil == next(mapResult) then
    local max = ARGV[1]
    if nil == max then
        max = 100
    end

    redis.call("hmset", HASH_KEY, "current", 1, "max", max)
    redis.call("expire", HASH_KEY, 60)

    return 1
end

-- 将hgetall的返回对象转换成table(hashKey-val)
local limiterMap = {}
local nextkey
for i, v in ipairs(mapResult) do
    if i % 2 == 1 then
        nextkey = v
    else
        limiterMap[nextkey] = v
    end
end


local current = tonumber(limiterMap["current"])
local max = tonumber(limiterMap["max"])
if current + 1 > max then
    return 0
end

redis.call("hincrby", HASH_KEY, "current", 1)
return 1

