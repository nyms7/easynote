# easynote
## params
### -p
server port, default 9600
### -t
admin token, if not specified, it will be randomly generated
## nginx conf
```
location /  {
    proxy_pass http://127.0.0.1:9600;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "Upgrade";
    proxy_set_header Host $http_host;
    proxy_set_header Remote_Addr $remote_addr;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```
