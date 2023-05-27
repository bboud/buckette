const std = @import("std");
const mem = std.mem;
const http = std.http;

const print = std.debug.print;

//Route function signature
const RouteFnPtr = *const fn (response: *http.Server.Response) void;

pub const Route = struct {
    method: http.Method,
    target: []const u8,
    route: RouteFnPtr,
};

pub const Router = struct {
    routes: std.StringHashMap(Route),

    pub fn init(allocator: mem.Allocator) Router {
        return .{ .routes = std.StringHashMap(Route).init(allocator) };
    }

    pub fn deinit(self: *Router) void {
        self.deinit();
    }

    pub fn addRoute(self: *Router, target: []const u8, method: http.Method, fnPtr: RouteFnPtr) !void {
        try self.routes.put(target, .{
            .method = method,
            .target = target,
            .route = fnPtr,
        });
    }

    pub fn route(self: *Router, response: *http.Server.Response) !void {
        const index = mem.indexOf(u8, response.request.target[1..], "/") orelse response.request.target.len;
        const slice = response.request.target[1..index];

        const funcToCall = self.routes.get(slice) orelse CannedResponses.failed404;
        funcToCall();
    }
};

//TODO: Handle errors correctly
pub const CannedResponses = struct {
    // Returns generic 200 with no data
    pub fn default(response: *http.Server.Response) void {
        response.status = http.Status.ok;
        response.headers.append("connection", "close") catch return;
        response.do() catch return;
        response.finish() catch return;
    }

    pub fn failed404(response: *http.Server.Response) void {
        response.status = http.Status.not_found;
        response.headers.append("connection", "close") catch return;
        response.do() catch return;
        response.finish() catch return;
    }
};
