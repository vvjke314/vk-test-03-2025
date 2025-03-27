#!/usr/bin/env tarantool
-- Configure database
box.cfg{
    listen = '0.0.0.0:3301',
}

-- Custom function to check key existence
box.schema.func.create('key_check', {
    body = [[ 
    function(key)
        return box.space.vault:get( { key } )
    end 
    ]]
})


-- Initing database
box.once("init", function()
    box.schema.space.create('vault')
    box.space.vault:format({
        { name = 'key', type = 'string' },
        { name = 'value', type = 'string' }
    })
    box.space.vault:create_index('primary',
        { parts = { 'key' } })

    box.schema.user.create('go-api', { password = os.getenv('TT_PASS') or 'password' })
    box.schema.user.grant('go-api', 'read,write', 'space', 'vault')
    box.schema.user.grant('go-api', 'execute', 'function', 'key_check')
 end)

