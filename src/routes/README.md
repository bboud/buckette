## Note:
There always needs to be the default route of "/" however you can point it to whatever function you like. 

## Multiple routes
If you need more complex behavior, like you have a route that has a large number of sub-routes, you can define another router within!

```zig
    var r = router.Router.init(myAllocator.allocator());
    defer r.deinit();

    try r.addRoute("/mynewcoolroute", .{
        .method = http.Method.GET,
        .fileserver = true,
        .route = index,
    });
```