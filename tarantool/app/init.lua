#!/usr/bin/env tarantool
-- Configure database
box.cfg{
    listen = '0.0.0.0:3301',
    log_level = 6,
}


box.once("init", function()
    box.schema.space.create('vault')
    box.space.vault:format({
        { name = 'key', type = 'string' },
        { name = 'value', type = 'string' }
    })
    box.space.vault:create_index('primary',
        { parts = { 'key' } })

    box.schema.user.create('go-api', { password = os.getenv('TRNTLPASS') or 'password' })
    box.schema.user.grant('go-api', 'read,write', 'space', 'vault')
 end)

