# data-stuff

## Price data
This is just some ways of handling price data that I've encountered and needed
for specific purposes - thought I'd share it in case anyone found it useful.

### The Basics

The first call should always be to one of the Fill... functions. This initiates
the maps and fills in the whole set of data.

Subsequent calls can use the Update... functions. These will update the price 
data and keep track of the latest update.
