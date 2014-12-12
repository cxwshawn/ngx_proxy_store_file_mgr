local redis = require "resty.redis"
local red = redis:new()

red:set_timeout(1000) -- 1 sec

                -- or connect to a unix domain socket file listened
                -- by a redis server:
                --     local ok, err = red:connect("unix:/path/to/redis.sock")

local ok, err = red:connect("127.0.0.1", 6379)
if not ok then
    ngx.log(ngx.ERR, "connect to redis failed.", err)
    return
end
                local res, err = red:get("cache")
                if  res == ngx.null then
                    ngx.log(ngx.ERR, "get cache from redis failed.", err)
                elseif res == "lock" then
                    ngx.exec("@proxyNoStore", ngx.var.args)
                end
                --ngx.say("lock name", res)
                local md5 = ngx.md5(ngx.var.uri)
                local ok, err = red:zadd("defset", ngx.now(), md5)
                if not ok then
                    ngx.log(ngx.ERR, "add to defset sorted set failed.", err)
                    return
                end
                local storepath = ngx.var.document_root..ngx.var.uri
                local ok, err = red:hset("defhash", md5, storepath)
                if not ok then
                    ngx.log(ngx.ERR, "set to hash failed.", err)
                    return
                end

                -- put it into the connection pool of size 100,
                -- with 10 seconds max idle time
                local ok, err = red:set_keepalive(10000, 100)
                if not ok then
                    ngx.log(ngx.ERR, "failed to set keepalive", err)
                    return
                end

                -- or just close the connection right away:
                -- local ok, err = red:close()
                -- if not ok then
                --     ngx.say("failed to close: ", err)
                --     return
                -- end