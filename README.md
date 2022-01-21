# Manager for Xiaomi (wifi) devices

## Instead of preface

Manager is a driver module to serve requests and statuses from different hardware. This module serve hardware from xiaomi
ecosystem using miio protocol specification. All updates and commands send thru mqtt server. I assume different module 
(named hub) should control mqtt traffic and serve automation.

## How to Xiaomi MIIO (and this manager) works (in short)

....

## Supported devices

* Probably any type of xiaomi wifi device.

## Command line paraments

$ ./manager_xiaomi -?
flag provided but not defined: -?
Usage of ./manager_yeelight:
  -clientid string
    client id for mqtt server (default "xiaomi-1")
  -debug
    print debuging hex dumps
  -keepalive int
    keepalive timeout for mqtt server (default 30)
  -login string
    login string for mqtt server
  -mqtt string
    mqtt server address (default "127.0.0.1:1883")
  -pass string
    password string for mqtt server

# Device registration

Manager can adopt new device from factory default state. If device was previously linked it should be reset as described
in documentation. For lamps usually need to turn on and off five times with 1 second delay. Other device has reset pin
or need press button for 5 seconds. Reset light should start fast flashing (yellow). After device was reset computer should
be connected to device wireless network *device-name-v1_mibt8217* without password and run manager binary with follow parameters:

$ ./manager_xiaomi -reg -sid MYNET -key MyNetPass [-clientid string] [-login mqtt-login] [-pass mqtt-pass] -mqtt 192.168.1.1:1883

After registration complete computer will be disconnected from device network (and probably connect to you home network),
from device manager receive *device_id* and *token*. This data is very important to next communication, so it will be printed 
to console and try to save to mqtt server as retain message.

## Format mqtt stored object

Each device found in the network stored in mqtt as retain message in follow format:

yeelight/78ab4e4 = {"id":78ab4e4,"model":"ceiling24","token":"##########", name":"My ceiling light","ver":5}

Ip address of device doen't stored because address can change by DHCP server in the network. Manager_yeelight doesn't
requre to use static address for device and address changes serve correctly.

When device found in the network mqtt record will be extended with support commands tag (support) and status variables.
Status variables set may (and actually should) be different on different types of device. (Lamps do not inform about
temperature, but thermometers does).

Extended example of device records in mqtt:

yeelight/78ab4e4 = {"id":78ab4e4,"token":"############", "ip":"1.1.1.1","model":"ceiling24","name":"My ceiling light","ver":5,
"support":"set_scene set_ct_abx adjust_ct set_bright set_name adjust_bright set_default toggle cron_get start_cf set_adjust get_prop set_power cron_add cron_del stop_cf",
"power":"false","bright":100,"mode":2,"temp":4000,"rgb":0,"hue":0,"sat":0}

## The commands from mqtt server to manage device

Command format: { "cmd":"xxxx", "value1":"aaa", "value2":"bbbb", "value3":"ccccc"}

Commands are dynamically formatted and processed. Different devices has different set of commands. Tag *support* in device
record define which command are actually supported by this device. Hub or process who format command should know set of
parameters for each command.

## Known problems

* only registration was done yet.
* mqtt client has several (or much more) issues. i.e. reconnect procedure eat too much cpu time
* All commands/parameters is pass to lamp correctly, but not all types of possible response correctly converted to mqtt device update, i.e. music flow and etc.

## License and author

This project licensed under GNU GPLv3.

Author Eugene Chertikhin
