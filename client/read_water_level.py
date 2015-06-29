#! /usr/bin/env python
#
# Python Code running on Raspberry Pi to communicate with MCP3008 A/D converter to estimate the sump pit water level.
#
# Readings are submitted to the server at sump.blakecaldwell.net/water-level, along with the secret/password.
#

import spidev
import os
from time import sleep

is_debug_mode = False

server_secret = os.environ['SUMP_SECRET']
if len(server_secret) == 0:
    print('Missing environment variable SUMP_SECRET')
    exit -1


#Establish SPI Connection on Bus 0, Device 0
spi = spidev.SpiDev()
spi.open(0,0)


# Get and return the water level
def check_water_level(channel):

    #Check valid channel
    if((channel > 7) or (channel < 0)):
        return -1

    #Perform SPI transaction and store returned bits in 'r'
    r = spi.xfer([1, (8+channel)<<4, 0])
    
    #Filter data bits from returned bits
    adcout = ((r[1]&3) << 8) + r[2]
    percent = int(round(adcout/10.23))

    # Convert to inches
    depth = float(820 - adcout)/12.2 + 1
    
    if is_debug_mode:
        # Print 0-1023 and percentage
        print('ADC Output: {0:4d}     Percentage: {1:3d}%    Depth: {2:3.2f} inch'.format(adcout, percent, depth))
        print('{0:3.2f} inches'.format(depth))
        
    return depth
    

# Submit the water level to the server
def submit_water_level(level):
    import requests
    
    url = 'http://sump.blakecaldwell.net/water-level'
    payload = {
        'secret': server_secret,
        'level': level,
    }
    headers = {'content-type': 'application/json'}

    print("Posting %s" % str(level))	
    response = requests.post(url, data='{{ "secret":"{secret}", "level":{level} }}'.format(secret=server_secret, level=level), headers=headers)
    print("Posted.")


# Loop forever!
counter = 0
sum = 0
while True:
    try:
        # Every 2 seconds, send the average reading to the server
        if counter == 20:
            level_str = "{0:3.2f}".format(sum/20)
            print("%s inches" % level_str)
            submit_water_level(level_str)
            sum = 0
            counter = 0

        counter = counter + 1
        sum = sum + check_water_level(0)
        sleep(0.1)
    
    except KeyboardInterrupt:
        exit(0)
    
    except:
        pass
