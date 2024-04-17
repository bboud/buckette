const std = @import("std");
const heap = std.heap;
const http = std.http;
const mem = std.mem;

const config = @import("config.zig");

//TODO: Graceful shutdown
pub fn main() !void {
    // Primary allocator for the while server.
    var gpa = heap.GeneralPurposeAllocator(.{}){};
    defer {
        const check = gpa.deinit();
        switch (check) {
            heap.Check.ok => std.log.debug("No mem lead detected!", .{}),
            heap.Check.leak => std.log.debug("Mem leak!", .{}),
        }
    }
    var gpaAlloc = gpa.allocator();

    // HTTP server initialization
    var server = http.Server.init(gpaAlloc, .{});
    defer server.deinit();

    // Router initialization
    var buffer: [1024]u8 = undefined;
    var fba = heap.FixedBufferAllocator.init(&buffer);
    var router = @import("router.zig").Router.init(fba.allocator());
    defer router.deinit();
    try @import("routes.zig").setup(&router);

    // Bind do the address and listen
    const address = std.net.Address.parseIp4(config.ADDRESS, config.PORT) catch unreachable;
    try server.listen(address);
    std.log.debug("Listening Port: {}\n", .{address.getPort()});

    // Primary connection processing loop.
    // 1) Accept the connection
    // 2) Wait for the request to be completed
    // 3) route the request to the correct handler
    while (true) {
        var response = server.accept(.{
            .allocator = gpaAlloc,
        }) catch |err| {
            std.log.err("Unable to accept connection: {}", .{err});
            continue;
        };

        defer response.deinit();
        defer _ = response.reset();

        response.wait() catch |err| {
            std.log.err("Unable to complete wait for request: {}", .{err});
            continue;
        };

        const target = response.request.target;
        // Errors are handling inside each route handler
        router.route(target, &response, gpaAlloc);
    }
}
