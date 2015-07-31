# Sump: Raspberry Pi-based Sump Monitor [![Build Status](https://travis-ci.org/wblakecaldwell/sump.svg?branch=master)](https://travis-ci.org/wblakecaldwell/sump)

TLDR: [Raspberry Pi](https://www.raspberrypi.org/products/model-a-plus/) + custom circuitry + [Go](http://golang.org) service in [Google AppEngine](https://cloud.google.com/appengine/docs) = [http://sump.blakecaldwell.net](http://sump.blakecaldwell.net)


Details
-------

This is the code for my DIY [sump](https://en.wikipedia.org/wiki/Sump) water level monitoring system. I have a [Raspberry Pi A+](https://www.raspberrypi.org/products/model-a-plus/) in my
basement with some custom circuitry to monitor the water level in my sump pit. Every 2 seconds,
the Raspberry Pi posts the current status to my server, which is hosted in [Google App Engine](https://cloud.google.com/appengine/docs). The server keeps
track of the levels, and displays them in a pretty chart, here: [http://sump.blakecaldwell.net](http://sump.blakecaldwell.net).
If the level gets above a certain threshold (the pump died), the server will email me.


Client
------

The client code is written in [Python](https://www.python.org/), running on a [Raspberry Pi A+](https://www.raspberrypi.org/products/model-a-plus/).

It reads values from an analog-to-digital converter, estimates the water level,
and sends it to my server. At some point, I'll post the schematics of my custom circuitry.


Server
------

The server is written in [Go](http://golang.org), and hosted in [Google App Engine](https://cloud.google.com/appengine/docs).

The server hosts the following endpoints:

- **Main page:** [http://sump.blakecaldwell.net/index.html](http://sump.blakecaldwell.net/index.html) - the pretty graph showing the water level over the past 2 hours
- **Data:** [http://sump.blakecaldwell.net/info](http://sump.blakecaldwell.net/info) - all of the JSON data that the main page uses to build the chart
- **Data input:** [http://sump.blakecaldwell.net/water-level]() - the page that the client posts to. The server expects that the client knows its secret

At this point, the server just keeps all of the readings in memory, in an ever-expanding [slice](http://blog.golang.org/go-slices-usage-and-internals).
I never intended to keep that around outside of testing, but after launching this to test it out, it seemed "fine". From time to time,
I noticed that Google restarted my service. Maybe this is the reason. No big deal. I'll update it at some point.


Building
--------

I can't imagine anyone out there is going to replicate this project, but if you do, you'll need
to create your own `server/config_values.go` from `server/config_values.go.template`.  This is the
best way I found to get settings into a Go service that's hosted in AppEngine. I would have
preferred using environment variables, listed in app.yaml, but Google doesn't support this for
Go services.
