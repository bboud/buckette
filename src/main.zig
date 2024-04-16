const std = @import("std");
const heap = std.heap;
const http = std.http;
const mem = std.mem;

const print = std.debug.print;

const router = @import("router.zig");
const index = @import("fileserver.zig").fileserver;
const upload = @import("upload.zig").upload;

pub fn main() !void {
    var gpa = heap.GeneralPurposeAllocator(.{}){};
    defer {
        const check = gpa.deinit();
        switch (check) {
            heap.Check.ok => print("No mem lead detected", .{}),
            heap.Check.leak => print("Mem leak!", .{}),
        }
    }

    var gpaAlloc = gpa.allocator();
    var server = http.Server.init(gpaAlloc, .{});
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
            .allocator = gpaAlloc,
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
        // Errors are handling inside each route handler
        r.route(&response, gpaAlloc);
    }
}

fn setup(r: *router.Router) !void {
    try r.addRoute("/", index);
    try r.addRoute("/upld", upload);
}
