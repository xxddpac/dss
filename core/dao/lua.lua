package.cpath = package.cpath .. ';C:/Users/administrator/AppData/Roaming/JetBrains/GoLand2022.2/plugins/EmmyLua/debugger/emmy/windows/x86/?.dll'
local dbg = require('emmy_core')
dbg.tcpConnect('localhost', 9966)

function test (n)
    if n == 0 then
        return "lua"
    else
        local resp = "golang"
        return resp
    end
end
print(test(1))