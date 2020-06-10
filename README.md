# mars-image

## Endpoints
 * /image
    * Returns a URL for the given lat and long
    * Required query parameters:
        * lat (-90 to 90)
        * long (-180 to 180)
 * /get-metrics
    * Returns the cache miss, cache hit, cache eviction, and their associated average execution times.
 * /reset-metrics
    * Resets the cache metrics