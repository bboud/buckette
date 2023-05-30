const std = @import("std");
const heap = std.heap;
const http = std.http;
const mem = std.mem;

const print = std.debug.print;

const router = @import("router.zig");
const setup = @import("routes/setup.zig").setup;

pub fn main() !void {
    var aAllocator = heap.ArenaAllocator.init(heap.page_allocator);
    defer aAllocator.deinit();

    var server = http.Server.init(aAllocator.allocator(), .{});
    defer server.deinit();

    // This FBA is used as the allocator for a StringHashMap to be fast for route lookups
    var buffer: [512]u8 = undefined;
    var fbaAllocator = heap.FixedBufferAllocator.init(&buffer);
    var r = router.Router.init(fbaAllocator.allocator());
    defer r.deinit();
    try setup(&r);

    const address = std.net.Address.parseIp4("127.0.0.1", 8080) catch unreachable;
    try server.listen(address);

    print("Listening Port: {}\n", .{address.getPort()});

    while (true) {
        var response = server.accept(.{
            .allocator = aAllocator.allocator(),
        }) catch |err| switch (err) {
            mem.Allocator.Error.OutOfMemory => break,
            else => {
                print("Unable to accept connection: {}", .{err});
                continue;
            },
        };
        defer response.deinit();
        defer _ = response.reset();
        try response.wait();

        // This route take an allocator for use in serving files and other things
        r.route(&response, aAllocator.allocator());
    }
}
