# Manager for Xiaomi (wifi) devices

## Instead of preface

Manager is a driver module to serve requests and statuses from different hardware. This module manage hardware from xiaomi
ecosystem using miio protocol specification. All updates and commands send thru mqtt server. I assume different module 
(named hub) should control mqtt traffic and serve automation.

## Known problems

nothing works for now as assumed. manager was dropped for some time. may be forever. 

actually device registration works fine. key -uid waits mi account id. in case of defined new device will be avail in 
mi home application and works fine.
BUT deviceId will be changed after device come to online (reconnect to defined network). So this manager will be unusable.

## License and author

This project licensed under GNU GPLv3.

Author Eugene Chertikhin
