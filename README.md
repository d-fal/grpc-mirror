# Foreword

With the gowth of **Microservice** architecture, we are witnessing the purpose-separation within whole the services in an neatly-designed ecosystem. For example, one may 
need to separate logging layer from server applications. GRPC, on anotherhand, takes a good chance to subside currently wide used HTTP/1.1 in the coming years. 
Migrating from mature HTTP/1.1 to faster and more secure grpc needs to 
pack all the good things we had within it and deploy it in our new destination. Some very wanted features, tools and assets that we already have in HTTP/1.1 
are tuned load-balancer (LB) like [haproxy](http://www.haproxy.org/). There is a [valid proposal](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md) 
for such service within GRPC as well. Powerful Nginx has recently opened the door for grpc and we expect a reliable proxy model from this.

In proxy model in which a world-exposed frontend takes care of requests and bruit them over its intended backends, we need a more handy and configurable 
load-balancer.


### GRPC-Mirror

This project, GRPC mirror, is intended to replicate grpc requests and dispatch them to their intended backend. At the beginnong, **grpc-mirror** was not intended to be used as LB since it was designed to be a dispatcher. Later, we used the proxy model but we didn't use round-robin or etc. However, it is open to be used to cover this as well.

In the following, the whole picture of the project is presented:
```
// GRPCEndpoint grpc endpoint
type GRPCEndpoint map[*Connection]interface{}

// Connection is the container that keeps connection
type Connection struct {
	Conn    *grpc.ClientConn
	Filters []DispatchOptions
	// Contains determines if a message with specified type should have a type or not.
	// This one is a little tricky. So here is the example:
	// {MessageType: 1, Contains: true}
	// The above example is equivalent to:
	// Select all the messages where MessageType=1
	// another example:
	// {MessageType: 3, Contains: false}
	// Select * messages where MessageType is not 3
	// To pass all the messages without any kinds of filters, just let it be nil
	Contains *bool
}

// DispatchOptions dispatch options
/*
									           
								           
                                    Server1 (Recv types {1},{2,RedirectResponse})
                                  /
 ______  Msg type: 3 ___________ /________ Mirror1 (Recv types {1},{2},{3}) KAFKA
|Client| --->-->--> |Dispatcher |/
|______|  ----<---<-|___________|\________ Mirror2 (Recv type {1} only) ELASTICSEARCH
                                 \
		Rsp from Reciver2 \
                                   \
                                     Server2 (Recv types {2},{3,RedirectResponse})
				                

*/
  type DispatchOptions struct {
    MessageType int32

    // RedirectResponse is the hack to redirect responses
    // If Contains in Connection block was true, then it looks into this flag
    // If RedirectResponse was true, the dispatcher sends back te recipent
    // response to the client that has made the request
    RedirectResponse bool
  }

  // Lookup looks up into filterd types and prevents message dispatch if filter is specified
  // This method searchs inside the filtered types and send back the result of search
  // in accord with specified filter
  // Assume we have a message with type:2 and a registered recipient with
  // the following Dispatch options
  // Option 1: {Type:1  RedirectResponse: false}
  // Option 2: {Type:2  RedirectResponse: true}
  // .....
  // Since we have message type=2, Option 2 asserts that we should block the message that
  // has this type. So, it will be blocked
  // If a message was not specified in any of
  func (c *Connection) Lookup(filter int32) (bool, bool) {
    if c.Contains != nil {
      for _, f := range c.Filters {
        if f.MessageType == filter {
          return *c.Contains, f.RedirectResponse
        }
      }
      return !*c.Contains, false
    }
    return true, false
  }


```

I appreciate your PRs and new ideas.
