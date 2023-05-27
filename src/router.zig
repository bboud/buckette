const std = @import("std");
const mem = std.mem;
const http = std.http;
const fs = std.fs;

const print = std.debug.print;

//Route function signature
const RouteFnPtr = *const fn (response: *http.Server.Response) void;

pub const Router = struct {
    routes: std.StringHashMap(RouteFnPtr),

    pub fn init(allocator: mem.Allocator) Router {
        return .{ .routes = std.StringHashMap(RouteFnPtr).init(allocator) };
    }

    pub fn deinit(self: *Router) void {
        self.deinit();
    }

    pub fn addRoute(self: *Router, target: []const u8, r: RouteFnPtr) !void {
        try self.routes.put(target, r);
    }

    pub fn route(self: *Router, response: *http.Server.Response) !void {
        const request = response.request;
        const target = request.target;

        //Search on the target skipping the first '/'
        const i = mem.indexOf(u8, target[1..], "/") orelse target.len;

        //In the slice, we return the whole thing upto the index
        const slice = target[0..i];

        // Try to get the route and if not, load index as a fileserver. Failure to get index results in canned response
        const foundRoute: RouteFnPtr = self.routes.get(slice) orelse self.routes.get("/") orelse noDefault;

        foundRoute(response);
    }
};

fn noDefault(response: *http.Server.Response) void {
    const cwd = fs.cwd();
    const file = cwd.openFile("html/no_default.html", .{}) catch unreachable;
    const stat = file.stat() catch unreachable;
    const size = stat.size;

    var buffer: [1024]u8 = undefined;
    const read = file.readAll(&buffer) catch unreachable;

    response.status = http.Status.ok;
    response.transfer_encoding = .{ .content_length = size };
    response.headers.append("connection", "close") catch unreachable;
    response.do() catch unreachable;

    response.writeAll(buffer[0..read]) catch unreachable;
    response.finish() catch unreachable;
}
