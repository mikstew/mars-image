# mars-image

## Endpoints
 * /image
    * Returns a URL for the given lat and long
    * Required query parameters:
        * lat (-90 to 90)
        * long (-180 to 180)
 * /get-metrics
    * Returns the cache miss and cache hit metrics
 * /reset-metrics
    * Resets the cache metrics