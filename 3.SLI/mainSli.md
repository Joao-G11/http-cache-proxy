The main SLIs that guarantee the reliability, performance and scalability are the following.

###• The latency of requests:
The latency is the time that it takes for a client to receive a reply to its request. When providing a service one wants to make sure that a given amount of requests is answered within a given latency. One way to measure it is the following:

Acceptable latency / All answered requests

###• Error rate
For a particular proxy configuration it is the ratio of errors over all requests. In order to detect errors one can check if the http reply status code is in the 500 range, for example. After that the following formula is applied:

Number of errors / All responses

###• Response rate/Availability:
Is quite similar to the error rate but instead of finding the error ratio in requests it measures the amount of requests that were attended. Big response rates imply high availability. It can be calculated with:

Number of answered requests / All requests

###• System throughput:
The number of requests the proxy can attend in a given time frame. The way to calculate the throughput is by making increasing amounts of requests over equal time frames. When the request amount increases but the number of server requests attended doesn’t it usually means that the proxy is already drowning on requests and that is the throughput of the proxy. When the proxy starts to stop attending some requests given the overload performance might even decrease a little.

Note: When testing it is also important to be aware that the throughput might drop because the machine running the client requests being overloaded before the proxy does.
