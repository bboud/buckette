const std = @import("std");
const heap = std.heap;
const http = std.http;
const mem = std.mem;

const print = std.debug.print;

const router = @import("router.zig");

fn index(response: *http.Server.Response) void {
    const body = "Index!";

    //TODO: Handle the errors lol
    response.transfer_encoding.content_length = body.len;
    response.headers.append("content-type", "text/plain") catch return;
    response.do() catch return;

    _ = response.writeAll(body) catch return;
    response.finish() catch return;
}

pub fn main() !void {
    var aAllocator = heap.ArenaAllocator.init(heap.page_allocator);
    defer aAllocator.deinit();

    var server = http.Server.init(aAllocator.allocator(), .{});
    defer server.deinit();

    var r = router.Router.init(aAllocator.allocator());
    defer r.deinit();

    try r.addRoute("/", http.Method.GET, index);

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

        try r.route(&response);
    }
}
